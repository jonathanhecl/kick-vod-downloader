package main

import (
	"bytes"
	"fmt"
	"github.com/jonathanhecl/gotimeleft"
	"io"
	"net/http"
	"net/url"
	"os"
	"path"
	"sync"

	cloudflarebp "github.com/DaRealFreak/cloudflare-bp-go"
	"github.com/grafov/m3u8"
)

func getURLContent(url string) (string, error) {
	client := &http.Client{}
	client.Transport = cloudflarebp.AddCloudFlareByPass(client.Transport)

	rep, err := client.Get(url)
	if err != nil {
		return "", err
	}

	body, err := io.ReadAll(rep.Body)
	rep.Body.Close()
	if err != nil {
		return "", err
	}

	return string(body), nil
}

func downloadUrlSegment(url string, pathDest string) error {
	filename := path.Base(url)
	filePath := path.Join(pathDest, filename)
	if _, err := os.Stat(filePath); err == nil {
		return nil
	}

	//fmt.Println("- Downloading", url)

	client := &http.Client{}
	client.Transport = cloudflarebp.AddCloudFlareByPass(client.Transport)

	rep, err := client.Get(url)
	if err != nil {
		return err
	}
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, rep.Body)
	if err != nil {
		return err
	}

	return nil
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

func downloadSegments(metadata KickMetadataResponse) error {
	fmt.Println("- Downloading M3U8 Master Playlist")

	content, err := getURLContent(metadata.Source)
	if err != nil {
		return err
	}

	masterpl := decodeMasterPlaylist(content)

	urlVideo := fmt.Sprintf("%s/%s", getBaseUrl(metadata.Source), masterpl.Variants[0].URI)

	fmt.Println("- Downloading M3U8 Video Playlist")

	contentVideo, err := getURLContent(urlVideo)
	if err != nil {
		return err
	}

	videopl := decodeMediaPlaylist(contentVideo)

	if _, err := os.Stat(metadata.Livestream.Slug); os.IsNotExist(err) {
		os.Mkdir(metadata.Livestream.Slug, 0755)
	}

	var maxConcurrentDownloads = 16
	var downloadCounter int = 0
	var wg sync.WaitGroup
	var sem = make(chan struct{}, maxConcurrentDownloads)

	var totalSegments int
	for _, segment := range videopl.Segments {
		if segment == nil {
			continue
		}
		totalSegments++
	}

	fmt.Printf("- Downloading %d segments\n", totalSegments)
	fmt.Println("- Please wait...")

	timeleft := gotimeleft.Init(totalSegments)

	for _, segment := range videopl.Segments {
		if segment == nil {
			continue
		}

		wg.Add(1)
		sem <- struct{}{}
		go func(segment *m3u8.MediaSegment) {
			defer func() {
				wg.Done()
				timeleft.Value(downloadCounter)

				if downloadCounter%9 == 0 {
					fmt.Printf("%s %s - %s\n", timeleft.GetProgressBar(30), timeleft.GetProgressValues(), timeleft.GetProgress(1))
				}
			}()

			urlSegment := fmt.Sprintf("%s/%s", getBaseUrl(urlVideo), segment.URI)
			err := downloadUrlSegment(urlSegment, metadata.Livestream.Slug)
			if err != nil {
				panic(err)
			}
			downloadCounter++
			<-sem
		}(segment)
	}
	wg.Wait()

	fmt.Printf("- Downloaded %d segments\n", totalSegments)

	return nil
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
