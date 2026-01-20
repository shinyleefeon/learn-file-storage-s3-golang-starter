package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os/exec"
)

// getVideoAspectRatio takes a filepath to a video file and returns the aspect ratio
// as "16:9" for landscape videos, "9:16" for portrait videos, or "other" for everything else
func getVideoAspectRatio(filepath string) (string, error) {
	cmd := exec.Command("ffprobe", "-v", "error", "-print_format", "json", "-show_entries", "stream=width,height", filepath)
	output := &bytes.Buffer{}
	cmd.Stdout = output
	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf("failed to run ffprobe: %w", err)
	}
	
	var ffprobeOutput struct {
		Streams []struct {
			Width  int `json:"width"`
			Height int `json:"height"`
		} `json:"streams"`
	}

	err = json.Unmarshal(output.Bytes(), &ffprobeOutput)
	if err != nil {
		return "", fmt.Errorf("failed to parse ffprobe output: %w", err)
	}

	if len(ffprobeOutput.Streams) == 0 {
		return "", fmt.Errorf("no video streams found in file: %s", filepath)
	}

	width := ffprobeOutput.Streams[0].Width
	height := ffprobeOutput.Streams[0].Height

	if height == 0 {
		return "", fmt.Errorf("invalid video height: %d", height)
	}

	aspectRatio := float64(width) / float64(height)
	
	// Define standard ratios and tolerance
	const (
		ratio16_9 = 16.0 / 9.0  // ~1.778
		ratio9_16 = 9.0 / 16.0  // ~0.5625
		tolerance = 0.1
	)
	
	// Check if aspect ratio is close to 16:9
	if aspectRatio >= ratio16_9-tolerance && aspectRatio <= ratio16_9+tolerance {
		return "16:9", nil
	}
	
	// Check if aspect ratio is close to 9:16
	if aspectRatio >= ratio9_16-tolerance && aspectRatio <= ratio9_16+tolerance {
		return "9:16", nil
	}
	
	return "other", nil
}

func processVideoForFastStart(filepath string) (string, error) {
	outputPath := filepath + ".processing"
	cmd := exec.Command("ffmpeg", "-i", filepath, "-movflags", "faststart", "-f", "mp4", outputPath)
	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf("failed to process video for fast start: %w", err)
	}
	return outputPath, nil
}

/*func generatePresignedURL(s3Client *s3.Client, bucket, key string, expireTime time.Duration) (string, error) {
	psClient := s3.NewPresignClient(s3Client)
	req, err := psClient.PresignGetObject(context.Background(), &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	}, s3.WithPresignExpires(expireTime))
	if err != nil {
		return "", fmt.Errorf("failed to generate presigned URL: %w", err)
	}
	return req.URL, nil
}*/