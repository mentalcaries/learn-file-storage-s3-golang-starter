package main

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/bootdotdev/learn-file-storage-s3-golang-starter/internal/database"
)

func generatePresignedURL(s3Client *s3.Client, bucket, key string, expireTime time.Duration) (string, error) {
	psClient := s3.NewPresignClient(s3Client)

	req, err := psClient.PresignGetObject(context.Background(), &s3.GetObjectInput{
		Bucket: &bucket,
		Key:    &key,
	}, s3.WithPresignExpires(expireTime))

	if err != nil {
		return "", fmt.Errorf("could not create request: %s", err)
	}

	return req.URL, nil
}

func (cfg *apiConfig) dbVideoToSignedVideo(video database.Video) (database.Video, error) {
	values := strings.Split(*video.VideoURL, ",")
	if len(values) < 2 {
		return database.Video{}, fmt.Errorf("invalid video format")
	}
	bucket := values[0]
	key := values[1]

	signedVideo, err := generatePresignedURL(cfg.s3Client, bucket, key, time.Hour)
	if err != nil {
		return video, fmt.Errorf("could not create signed video: %v", err)
	}

	*video.VideoURL = signedVideo
	return video, nil
}
