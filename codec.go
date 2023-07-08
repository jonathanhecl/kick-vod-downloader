package main

import (
	"fmt"
	"os"
	"os/exec"
)

func convertVideo(slug string) error {
	mergeFile := fmt.Sprintf("%s/%s.ts", slug, slug)
	if _, err := os.Stat(mergeFile); err != nil {
		return err
	}

	convertFile := fmt.Sprintf("%s.mp4", slug)
	if _, err := os.Stat(convertFile); err == nil {
		err := os.Remove(convertFile)
		if err != nil {
			return err
		}
	}

	fmt.Println("- Converting ", slug, " to mp4")
	fmt.Println("- Please wait...")

	cmd := exec.Command("ffmpeg", "-i", mergeFile, "-c", "copy", "-loglevel", "error", "-stats", convertFile)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		return err
	}

	if _, err := os.Stat(convertFile); err == nil {
		fmt.Println("- Converted ", slug, " to mp4!!")

		fmt.Println("- Removing folder ", slug)

		err := os.RemoveAll(slug)
		if err != nil {
			return err
		}
	}

	return nil

}
