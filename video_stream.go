package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os/exec"
)

type ffprobeOut struct {
	Streams []struct {
		Width  int `json:"width"`
		Height int `json:"height"`
	} `json:"streams"`
}

func getVideoAspectRatio(filepath string) (string, error) {
	cmd := exec.Command("ffprobe", "-v", "error", "-print_format", "json", "-show_streams", filepath)

	var out bytes.Buffer
	cmd.Stdout = &out

	err := cmd.Run()
	var errBuf bytes.Buffer
	cmd.Stderr = &errBuf
	if err != nil {
		return "", fmt.Errorf("could not get streams; %v", err)
	}

	streamData := ffprobeOut{}

	err = json.Unmarshal(out.Bytes(), &streamData)
	if err != nil {
		return "", fmt.Errorf("could not unmarshal data: %v", err)
	}

	aspectRatio := calculateAspectRatio(streamData.Streams[0].Height, streamData.Streams[0].Width)
	return aspectRatio, nil
}

func calculateAspectRatio(height, width int) string {
	ratio := float64(width) / float64(height)
	if ratio > 1.6 && ratio <= 1.8 {
		return "16:9"
	}
	if ratio >= 0.5 && ratio < 0.65 {
		return "9:16"
	} else {

		return "other"
	}

}
