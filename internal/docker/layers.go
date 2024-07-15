package docker

import (
	"archive/tar"
	"context"
	"fmt"
	"io"
	"slices"
	"strings"

	"github.com/docker/cli/cli/command"
	"github.com/docker/docker/api/types/image"
)

type layer struct {
	ID       string
	FileTree *fileTreeNode
}

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

func getLayersOrderedArrFromImage(cli command.Cli, imageID string) ([]string, error) {
	imageInspect, _, err := cli.Client().ImageInspectWithRaw(context.Background(), imageID)
	if err != nil {
		return nil, fmt.Errorf("can't inspect image: %w", err)
	}

	layersOrderedArr := make([]string, len(imageInspect.RootFS.Layers))
	for i, layer := range imageInspect.RootFS.Layers {
		layersOrderedArr[i] = fmt.Sprintf("%s", strings.Split(layer, ":")[1])
	}

	return layersOrderedArr, nil
}

func readLayers(cli command.Cli, imageID string, layersOrderedArr []string) (map[string]layer, error) {
	imageReader, err := cli.Client().ImageSave(context.Background(), []string{imageID})
	if err != nil {
		return nil, fmt.Errorf("error saving image: %v", err)
	}
	defer imageReader.Close()

	tarReader := tar.NewReader(imageReader)
	filesInImage := make(map[string]layer, len(layersOrderedArr))
	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		}

		if err != nil {
			return nil, fmt.Errorf("error reading tar header: %w", err)
		}

		fileName := header.Name
		if header.Typeflag == tar.TypeReg { //|| header.Typeflag == tar.TypeSymlink {
			fName := strings.TrimPrefix(fileName, "blobs/sha256/")
			if slices.Contains(layersOrderedArr, fName) {
				layerReader := tar.NewReader(tarReader)
				files, err := getFileTreeFromLayer(layerReader)
				if err != nil {
					return nil, fmt.Errorf("error getting files from layer: %w", err)
				}

				filesInImage[fName] = layer{
					ID:       fileName,
					FileTree: files,
				}
			}
		}
	}

	return filesInImage, nil
}

func getFileTreeFromLayer(layerReader *tar.Reader) (*fileTreeNode, error) {
	fileTree := &fileTreeNode{"/", true, make([]*fileTreeNode, 0)}

	for {
		header, err := layerReader.Next()
		if err == io.EOF {
			break
		}

		if err != nil {
			return nil, fmt.Errorf("error reading tar header: %w", err)
		}

		fileTree.addChild(header)
	}

	return fileTree, nil
}
