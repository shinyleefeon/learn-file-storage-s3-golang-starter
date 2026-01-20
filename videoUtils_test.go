package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func TestGetVideoAspectRatio(t *testing.T) {
	// Check if ffprobe is available
	if _, err := exec.LookPath("ffprobe"); err != nil {
		t.Skip("ffprobe not found, skipping test")
	}

	// Create a temporary test video file (this would need an actual video file)
	// For now, we'll test error cases

	t.Run("returns error for non-existent file", func(t *testing.T) {
		_, err := getVideoAspectRatio("/path/to/nonexistent/video.mp4")
		if err == nil {
			t.Error("expected error for non-existent file, got nil")
		}
	})

	t.Run("returns error for empty filepath", func(t *testing.T) {
		_, err := getVideoAspectRatio("")
		if err == nil {
			t.Error("expected error for empty filepath, got nil")
		}
	})

	t.Run("returns error for non-video file", func(t *testing.T) {
		// Create a temporary text file
		tmpDir := t.TempDir()
		tmpFile := filepath.Join(tmpDir, "test.txt")
		err := os.WriteFile(tmpFile, []byte("not a video"), 0644)
		if err != nil {
			t.Fatalf("failed to create test file: %v", err)
		}

		_, err = getVideoAspectRatio(tmpFile)
		if err == nil {
			t.Error("expected error for non-video file, got nil")
		}
	})

	// If you have sample video files in the samples/ directory, you can test with real files
	t.Run("calculates aspect ratio for real video file", func(t *testing.T) {
		// Check if samples directory exists and has video files
		samples, err := filepath.Glob("samples/*.mp4")
		if err != nil || len(samples) == 0 {
			t.Skip("no sample video files found in samples/ directory")
		}

		ratio, err := getVideoAspectRatio(samples[0])
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		// Should return one of the valid values
		if ratio != "16:9" && ratio != "9:16" && ratio != "other" {
			t.Errorf("expected '16:9', '9:16', or 'other', got: %s", ratio)
		}
		
		t.Logf("Aspect ratio for %s: %s", samples[0], ratio)
	})
}
