package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"

	cloudflarebp "github.com/DaRealFreak/cloudflare-bp-go"
)

type KickMetadataResponse struct {
	Source     string `json:"source"`
	Livestream struct {
		SessionTitle string `json:"session_title"`
		Slug         string `json:"slug"`
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

func getMetadataFromKickURL(videoID string) (KickMetadataResponse, error) {
	var res = KickMetadataResponse{}
	url := fmt.Sprintf("https://kick.com/api/v1/video/%s", videoID)

	client := &http.Client{}
	client.Transport = cloudflarebp.AddCloudFlareByPass(client.Transport)

	rep, err := client.Get(url)
	if err != nil {
		return KickMetadataResponse{}, err
	}

	body, err := io.ReadAll(rep.Body)
	rep.Body.Close()
	if err != nil {
		return KickMetadataResponse{}, err
	}

	err = json.Unmarshal(body, &res)
	if err != nil {
		return KickMetadataResponse{}, err
	}

	return res, nil
}
