package docker

import (
	"archive/tar"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/docker/cli/cli/command"
	"github.com/docker/docker/api/types/image"
)

type layer struct {
	ID       string
	FileTree *fileTreeNode
}

type manifestItem struct {
	Config   string   `json:"Config"`
	RepoTags []string `json:"RepoTags"`
	Layers   []string `json:"Layers"`
}

var (
	errNotATar = errors.New("not a tar archive")
)

func checkImageExists(cli command.Cli, imageID string) (bool, error) {
	images, err := cli.Client().ImageList(context.Background(), image.ListOptions{})
	if err != nil {
		return false, fmt.Errorf("can't list images: %w", err)
	}

	for _, img := range images {
		for _, t := range img.RepoTags {
			if t == imageID {
				return true, nil
			}
		}
	}

	return false, nil
}

func getLayersOrderedArrFromImage(imageReader io.ReadCloser) ([]layer, error) {
	tarReader := tar.NewReader(imageReader)
	manifest := make([]manifestItem, 0)
	filesInImage := make(map[string]layer, 0)

	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		}

		if err != nil {
			return nil, fmt.Errorf("error reading tar header: %w", err)
		}

		if header.Typeflag == tar.TypeReg {
			if header.Name == "manifest.json" {
				fileReader, err := io.ReadAll(tarReader)
				if err != nil {
					return nil, fmt.Errorf("error reading tar content: %w", err)
				}

				if err = json.Unmarshal(fileReader, &manifest); err != nil {
					return nil, fmt.Errorf("error unmarshalling manifest: %w", err)
				}
			} else if strings.HasSuffix(header.Name, "layer.tar") {
				layerReader := tar.NewReader(tarReader)
				files, err := getFileTreeFromLayer(layerReader)
				if err != nil {
					return nil, fmt.Errorf("error getting files from layer: %w", err)
				}

				filesInImage[header.Name] = layer{
					ID:       header.Name,
					FileTree: files,
				}
			} else if strings.HasPrefix(header.Name, "blobs/sha256/") {
				layerReader := tar.NewReader(tarReader)
				files, err := getFileTreeFromLayer(layerReader)
				if errors.Is(err, errNotATar) {
					continue
				}

				if err != nil {
					return nil, fmt.Errorf("error getting files from layer: %w", err)
				}

				filesInImage[header.Name] = layer{
					ID:       header.Name,
					FileTree: files,
				}
			}
		}
	}

	orderedLayersArr := make([]layer, len(manifest[0].Layers))
	for i, layer := range manifest[0].Layers {
		orderedLayersArr[i] = filesInImage[layer]
	}

	return orderedLayersArr, nil
}

func getFileTreeFromLayer(layerReader *tar.Reader) (*fileTreeNode, error) {
	fileTree := &fileTreeNode{"/", true, make([]*fileTreeNode, 0)}
	for {
		header, err := layerReader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, errNotATar
		}
		fileTree.addChild(header)
	}
	return fileTree, nil
}
