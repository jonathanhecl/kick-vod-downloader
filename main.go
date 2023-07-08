package main

import (
	"fmt"
	"os"
)

var (
	version = "0.0.3"
)

func main() {
	fmt.Println("Kick VOD Downloader v" + version)
	fmt.Println()

	fmt.Println("This tool is for educational purpose only.")
	fmt.Println("Do not use this tool to download videos without permission from the owner.")

	var url string
	if len(os.Args) == 2 {
		url = os.Args[1]
	} else {
		fmt.Println()
		fmt.Print("Input URL video: ")
		_, err := fmt.Scanln(&url)
		if err != nil {
			return
		}
		fmt.Println()
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

	fmt.Println()
	fmt.Println("Title:", metadata.Livestream.SessionTitle)
	fmt.Println()

	DownloadM3U8(metadata)

	fmt.Println("Done")
}
