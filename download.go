package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"
	"sync"

	cloudflarebp "github.com/DaRealFreak/cloudflare-bp-go"
	"github.com/grafov/m3u8"
)

func getURLContent(url string) string {
	client := &http.Client{}
	client.Transport = cloudflarebp.AddCloudFlareByPass(client.Transport)

	rep, err := client.Get(url)
	if err != nil {
		log.Fatal(err)
	}

	body, err := io.ReadAll(rep.Body)
	rep.Body.Close()
	if err != nil {
		log.Fatal(err)
	}

	return string(body)
}

func downloadUrlSegment(url string, pathDest string) {
	filename := path.Base(url)
	filePath := path.Join(pathDest, filename)
	if _, err := os.Stat(filePath); err == nil {
		return
	}

	fmt.Println("Downloading", url)

	client := &http.Client{}
	client.Transport = cloudflarebp.AddCloudFlareByPass(client.Transport)

	rep, err := client.Get(url)
	if err != nil {
		panic(err)
	}
	file, err := os.Create(filePath)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	_, err = io.Copy(file, rep.Body)
	if err != nil {
		panic(err)
	}
}

func decodeMasterPlaylist(content string) *m3u8.MasterPlaylist {
	buf := bytes.NewBufferString(content)
	p, listType, err := m3u8.Decode(*buf, true)
	if err != nil {
		panic(err)
	}

	if listType != m3u8.MASTER {
		panic("Not a master playlist")
	}

	return p.(*m3u8.MasterPlaylist)
}

func decodeMediaPlaylist(content string) *m3u8.MediaPlaylist {
	buf := bytes.NewBufferString(content)
	p, listType, err := m3u8.Decode(*buf, true)
	if err != nil {
		panic(err)
	}

	if listType != m3u8.MEDIA {
		panic("Not a media playlist")
	}

	return p.(*m3u8.MediaPlaylist)
}

func downloadSegments(metadata KickMetadataResponse) {
	fmt.Println("- Downloading M3U8 Master Playlist")

	content := getURLContent(metadata.Source)
	masterpl := decodeMasterPlaylist(content)

	urlVideo := fmt.Sprintf("%s/%s", getBaseUrl(metadata.Source), masterpl.Variants[0].URI)

	fmt.Println("- Downloading M3U8 Video Playlist")

	contentVideo := getURLContent(urlVideo)
	videopl := decodeMediaPlaylist(contentVideo)

	if _, err := os.Stat(metadata.Livestream.Slug); os.IsNotExist(err) {
		os.Mkdir(metadata.Livestream.Slug, 0755)
	}

	var maxConcurrentDownloads = 16
	var wg sync.WaitGroup
	var sem = make(chan struct{}, maxConcurrentDownloads)

	//var totalSegments int
	//for _, segment := range videopl.Segments {
	//	if segment == nil {
	//		continue
	//	}
	//	totalSegments++
	//}

	for _, segment := range videopl.Segments {
		if segment == nil {
			continue
		}

		wg.Add(1)
		sem <- struct{}{}
		go func(segment *m3u8.MediaSegment) {
			defer func() {
				wg.Done()
			}()

			urlSegment := fmt.Sprintf("%s/%s", getBaseUrl(urlVideo), segment.URI)
			downloadUrlSegment(urlSegment, metadata.Livestream.Slug)
			<-sem
		}(segment)
	}
	wg.Wait()
}

func getBaseUrl(fullUrl string) string {
	u, err := url.Parse(fullUrl)
	if err != nil {
		panic(err)
	}
	finalPath := path.Dir(u.Path)
	baseUrl := fmt.Sprintf("%s://%s%s", u.Scheme, u.Host, finalPath)
	return baseUrl
}
