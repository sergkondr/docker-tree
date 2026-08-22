package docker

import (
	"testing"

	"github.com/docker/docker/api/types/image"
)

func Test_imageMatches(t *testing.T) {
	img := image.Summary{
		ID:          "sha256:0123456789abcdef",
		RepoTags:    []string{"alpine:latest"},
		RepoDigests: []string{"alpine@sha256:abcdef"},
	}

	tests := []struct {
		name    string
		imageID string
		want    bool
	}{
		{name: "full image ID", imageID: "sha256:0123456789abcdef", want: true},
		{name: "short image ID", imageID: "sha256:012345", want: true},
		{name: "tag", imageID: "alpine:latest", want: true},
		{name: "digest", imageID: "alpine@sha256:abcdef", want: true},
		{name: "unknown image", imageID: "alpine:3.20", want: false},
		{name: "empty image", imageID: "", want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := imageMatches(img, tt.imageID); got != tt.want {
				t.Errorf("imageMatches() = %v, want %v", got, tt.want)
			}
		})
	}
}
