package main

import (
	"fmt"
	"os"

	"github.com/skondrashov/docker-tree/internal/docker"

	"github.com/docker/cli/cli-plugins/manager"
	"github.com/docker/cli/cli-plugins/plugin"
	"github.com/docker/cli/cli/command"
	"github.com/spf13/cobra"
)

var (
	version string
)

func main() {
	plugin.Run(
		func(dockerCli command.Cli) *cobra.Command {
			var (
				quiet     bool
				showLinks bool
			)

			cmd := &cobra.Command{
				Use:   "tree NAME[:TAG|@DIGEST] [DIRECTORY]",
				Short: "Display the directory tree of a Docker image",
				Long: `This command displays the directory tree of a Docker image, similar to the 'tree' command. 
Provide the image name and an optional tag or digest to view the file structure within the image. 
You can also specify a directory to see the file tree relative to this directory.`,
				Run: func(cmd *cobra.Command, args []string) {
					if len(args) == 0 || len(args) > 2 {
						fmt.Fprintln(dockerCli.Out(), cmd.UsageString())
						os.Exit(0)
					}

					imageID := args[0]
					if imageID == "" {
						fmt.Fprintln(dockerCli.Out(), "no imageID provided")
						os.Exit(1)
					}

					treeRoot := "/"
					if len(args) == 2 && args[1] != "" {
						treeRoot = args[1]
					}

					treeStrings, err := docker.GetImageTree(docker.GetTreeOpts{
						Cli:       dockerCli,
						ImageID:   imageID,
						Quiet:     quiet,
						ShowLinks: showLinks,
						TreeRoot:  treeRoot,
					})
					if err != nil {
						fmt.Fprintf(dockerCli.Err(), "can't get image tree: %s\n", err)
						os.Exit(1)
					}

					fmt.Fprintf(dockerCli.Out(), "%s", treeStrings)
				},
			}

			flags := cmd.Flags()
			flags.BoolVarP(&quiet, "quiet", "q", false, "Suppress verbose output")
			flags.BoolVarP(&showLinks, "links", "l", false, "Show symlinks destination")

			cmd.AddCommand()
			return cmd
		},

		manager.Metadata{
			SchemaVersion:    "0.1.0",
			Vendor:           "Sergei Kondrashov",
			Version:          version,
			ShortDescription: "Docker image tree",
			URL:              "https://github.com/sergkondr/docker-tree",
		})
}
