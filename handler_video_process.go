package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func processVideoForFastStart(filePath string)(string, error) {
  outputPath := filePath + ".processing"

  cmd := exec.Command("ffmpeg", "-i", filePath, "-c", "copy", "-movflags", "faststart", "-f", "mp4", outputPath)

  var stderr strings.Builder
  cmd.Stderr = &stderr

  err :=  cmd.Run()
  if err != nil {
    return "", fmt.Errorf("could not process: %v | %s", err, stderr.String())
  }

  fileInfo, err := os.Stat(outputPath)
	if err != nil {
		return "", fmt.Errorf("could not stat processed file: %v", err)
	}
	if fileInfo.Size() == 0 {
		return "", fmt.Errorf("processed file is empty")
	}
  return outputPath, nil
}