package main

import (
	"fmt"
	"os"
)

func convertVideo(slug string) {
	mergeFile := fmt.Sprintf("%s/%s.ts", slug, slug)
	if _, err := os.Stat(mergeFile); err != nil {
		panic(err)
	}

	convertFile := fmt.Sprintf("%s.mp4", slug)
	if _, err := os.Stat(convertFile); err == nil {
		err := os.Remove(convertFile)
		if err != nil {
			panic(err)
		}
	}

	fmt.Println("- Converting ", slug, " to mp4")

}
