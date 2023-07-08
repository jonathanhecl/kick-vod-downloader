package main

import (
	"encoding/json"
	"fmt"
	"regexp"
	"sync"
	"time"

	"github.com/gocolly/colly"
)

type KickResponse struct {
	Source       string `json:"source"`
	SessionTitle string `json:"session_title"`
	Livestream   struct {
		Slug string `json:"slug"`
	} `json:"livestream"`
}

func extractVideoID(url string) string {
	re := regexp.MustCompile(`^https://kick.com/video/([a-zA-Z0-9-]+$)`)
	matches := re.FindStringSubmatch(url)
	if len(matches) < 2 {
		return ""
	}
	return matches[1]
}

func getMetadataFromKickURL(videoID string) (KickResponse, error) {
	url := fmt.Sprintf("https://kick.com/api/v1/video/%s", videoID)

	var wait sync.WaitGroup
	var res = KickResponse{}

	c := colly.NewCollector()
	c.SetRequestTimeout(120 * time.Second)

	c.OnError(func(r *colly.Response, err error) {
		fmt.Println("Request URL:", r.Request.URL, "failed with error:", err)
		wait.Done()
	})

	c.OnResponse(func(r *colly.Response) {
		//fmt.Println("Response received", r.StatusCode)
		//fmt.Println("Response body", string(r.Body))

		err := json.Unmarshal(r.Body, &res)
		if err != nil {
			return
		}

		wait.Done()
	})

	c.OnRequest(func(r *colly.Request) {
		r.Headers.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/86.0.4240.75 Safari/537.36")
		r.Headers.Set("Cache-Control", "no-cache")
		r.Headers.Set("Accept", "*/*")
		r.Headers.Set("Connection", "Keep-Alive")
		r.Headers.Set("Accept-Encoding", "gzip")
		r.Headers.Set("Host", "kick.com")
		r.Headers.Set("Referer", "https://kick.com/")
	})

	wait.Add(1)
	err := c.Visit(url)
	if err != nil {
		return KickResponse{}, err
	}

	wait.Wait()

	return res, nil
}
