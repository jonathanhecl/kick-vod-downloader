package main

import "fmt"

var (
	version = "0.0.1"
)

func main() {
	fmt.Println("Kick VOD Downloader v" + version)
	fmt.Println()

	var url string
	fmt.Print("Input URL video: ")
	_, err := fmt.Scanln(&url)
	if err != nil {
		return
	}

	var videoID string
	videoID = extractVideoID(url)
	if len(videoID) == 0 {
		fmt.Println("URL is not valid")
		return
	}

	metadata, err := getMetadataFromKickURL(videoID)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Title:", metadata.SessionTitle)
	fmt.Println("M3U8 URL:", metadata.Source)
	fmt.Println("Slug:", metadata.Livestream.Slug)

	// TODO: Download m3u8 file
}
