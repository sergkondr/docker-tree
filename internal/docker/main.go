package docker

import (
	"context"
	"fmt"

	"github.com/docker/cli/cli/command"
)

type GetTreeOpts struct {
	Cli command.Cli

	ImageID  string
	Quiet    bool
	TreeRoot string
}

func GetImageTree(opts GetTreeOpts) (string, error) {
	ctx := context.Background()

	imageExists, err := checkImageExists(ctx, opts.Cli, opts.ImageID)
	if err != nil {
		return "", fmt.Errorf("can't check if image exists: %w", err)
	}

	if !imageExists {
		args := []string{"pull", opts.ImageID}
		if opts.Quiet {
			args = append(args, "--quiet")
		}
		envs := []string{"DOCKER_CLI_HINTS=false"}
		if err = runDockerCliCommand(args, envs); err != nil {
			return "", fmt.Errorf("error inspecting image: %w", err)
		}
	}

	if !opts.Quiet {
		fmt.Fprintf(opts.Cli.Out(), "precessing image: %s\n", opts.ImageID)
	}

	imageReader, err := opts.Cli.Client().ImageSave(ctx, []string{opts.ImageID})
	if err != nil {
		return "", fmt.Errorf("error saving image: %v", err)
	}
	defer imageReader.Close()

	layersOrderedArr, err := getLayersOrderedArrFromImage(imageReader)
	if err != nil {
		return "", fmt.Errorf("can't get layersOrderedArr: %w", err)
	}

	originalLayer := layersOrderedArr[0].FileTree
	for i := 1; i <= len(layersOrderedArr)-1; i++ {
		updatedLayer := layersOrderedArr[i].FileTree
		originalLayer, err = mergeFileTrees(originalLayer, updatedLayer)
		if err != nil {
			return "", fmt.Errorf("can't merge layers: %w", err)
		}
	}

	node := originalLayer.findNode(opts.TreeRoot)
	if node == nil {
		return "", fmt.Errorf("there is no such path in the image: %s", opts.TreeRoot)
	}

	return node.String(), nil
}
