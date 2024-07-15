package docker

import (
	"fmt"

	"github.com/docker/cli/cli/command"
)

type GetTreeOpts struct {
	Cli      command.Cli
	ImageID  string
	Quiet    bool
	TreeRoot string
}

func GetImageTree(opts GetTreeOpts) (string, error) {
	imageExists, err := checkImageExists(opts.Cli, opts.ImageID)
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

	layersOrderedArr, err := getLayersOrderedArrFromImage(opts.Cli, opts.ImageID)
	if err != nil {
		return "", fmt.Errorf("can't get layersOrderedArr: %w", err)
	}

	layerMap, err := readLayers(opts.Cli, opts.ImageID, layersOrderedArr)
	if err != nil {
		return "", fmt.Errorf("can't read layer: %w", err)
	}

	originalLayer := layerMap[layersOrderedArr[0]].FileTree
	for i := 1; i <= len(layersOrderedArr)-1; i++ {
		updatedLayer := layerMap[layersOrderedArr[i]].FileTree
		originalLayer, _ = mergeFileTrees(originalLayer, updatedLayer)
	}

	node := originalLayer.findNode(opts.TreeRoot)
	if node == nil {
		return "", fmt.Errorf("tree is no such path in the image: %s", opts.TreeRoot)
	}

	return node.String(), nil
}
