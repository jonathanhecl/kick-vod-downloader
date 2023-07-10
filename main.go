package main

import (
	"fmt"
	"github.com/jonathanhecl/gotimeleft"
	"os"
	"os/exec"
)

var (
	version = "0.1.6"
)

func main() {
	fmt.Println("Kick VOD Downloader v" + version)
	fmt.Println()

	timeleft := gotimeleft.Init(5)

	fmt.Println("This tool is for educational purpose only.")
	fmt.Println("Do not use this tool to download videos without permission from the owner.")

	_, err := exec.LookPath("ffmpeg")
	if err != nil {
		fmt.Println()
		fmt.Println("WARNING! FFmpeg is not detected. Please install FFmpeg first. https://ffmpeg.org/download.html")
		fmt.Println("	FFmpeg is required to convert video to MP4 format. If you don't want to convert video, you can skip this warning.")
	}

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

	timeleft.Step(1)

	metadata, err := getMetadataFromKickURL(videoID)
	if err != nil {
		fmt.Println(err)
		return
	}

	timeleft.Step(2)

	fmt.Println()
	fmt.Println("Title:", metadata.Livestream.SessionTitle)
	fmt.Println()

	err = downloadSegments(metadata)
	if err != nil {
		fmt.Println(err)
		return
	}

	timeleft.Step(3)

	err = mergeSegments(metadata.Livestream.Slug)
	if err != nil {
		fmt.Println(err)
		return
	}

	timeleft.Step(4)

	err = convertVideo(metadata.Livestream.Slug)
	if err != nil {
		fmt.Println(err)
		return
	}

	timeleft.Step(5)

	fmt.Println("Done in", timeleft.GetTimeSpent())
}
