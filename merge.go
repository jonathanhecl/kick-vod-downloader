package main

import (
	"fmt"
	"github.com/jonathanhecl/gotimeleft"
	"io"
	"os"
	"time"
)

func mergeSegments(slug string) error {
	var totalSegments = 0
	for {
		filePath := fmt.Sprintf("%s/%d.ts", slug, totalSegments)
		if _, err := os.Stat(filePath); err != nil {
			break
		}
		totalSegments++
	}

	if totalSegments == 0 {
		return fmt.Errorf("no segments found")
	}

	mergeFile := fmt.Sprintf("%s/%s.ts", slug, slug)
	if _, err := os.Stat(mergeFile); err == nil {
		err := os.Remove(mergeFile)
		if err != nil {
			return err
		}
	}

	fmt.Println("- Merging ", slug)
	fmt.Println("- Please wait...")

	file, err := os.Create(mergeFile)
	if err != nil {
		return err
	}
	defer file.Close()

	timeleft := gotimeleft.Init(totalSegments)

	for i := 0; i < totalSegments; i++ {
		filePath := fmt.Sprintf("%s/%d.ts", slug, i)
		if _, err := os.Stat(filePath); err != nil {
			break
		}

		if i%49 == 0 {
			fmt.Printf("%s %s - %s\n", timeleft.GetProgressBar(30), timeleft.GetProgressValues(), timeleft.GetProgress(1))
		}

		timeleft.Step(1)

		segFile, err := os.Open(filePath)
		if err != nil {
			return err
		}

		_, err = io.Copy(file, segFile)
		if err != nil {
			return err
		}

		segFile.Close()
	}

	fmt.Println("- Merged")

	time.Sleep(800 * time.Millisecond)

	fmt.Println("- Removing segments")

	for n := 0; n < totalSegments; n++ {
		filePath := fmt.Sprintf("%s/%d.ts", slug, n)
		if _, err := os.Stat(filePath); err != nil {
			break
		}

		err := os.Remove(filePath)
		if err != nil {
			return err
		}
	}

	fmt.Println("- Segments removed")

	return nil
}
