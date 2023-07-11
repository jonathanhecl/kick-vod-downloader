package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"time"

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
	url := fmt.Sprintf("https://kick.com/api/v1/video/%s?%d", videoID, time.Now().Unix())

	client := &http.Client{}
	client.Transport = cloudflarebp.AddCloudFlareByPass(client.Transport)

	rep, err := client.Get(url)
	if err != nil {
		return KickMetadataResponse{}, errors.New("client.Get: " + err.Error())
	}

	body, err := io.ReadAll(rep.Body)
	defer rep.Body.Close()
	if err != nil {
		return KickMetadataResponse{}, errors.New("io.ReadAll: " + err.Error())
	}

	if rep.StatusCode == 403 { // Cloudflare
		return KickMetadataResponse{}, errors.New("StatusCode: Cloudflare bypass failed")
	}

	err = json.Unmarshal(body, &res)
	if err != nil {
		return KickMetadataResponse{}, errors.New("json.Unmarshal: " + err.Error())
	}

	return res, nil
}
