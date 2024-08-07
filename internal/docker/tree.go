package docker

import (
	"archive/tar"
	"fmt"
	"sort"
	"strings"
)

const (
	empty  = ""
	space  = "    "
	branch = "│   "
	middle = "├── "
	last   = "└── "
	link   = " -> "

	whiteoutFilePrefix = ".wh."
	whiteoutDirPrefix  = ".wh..wh..opq"
)

type fileTreeNode struct {
	Name     string
	Symlink  string
	IsDir    bool
	Children []*fileTreeNode
}

type getStringOpts struct {
	showLinks bool
	depth     int
}

func (n *fileTreeNode) getString(prefix string, opts getStringOpts, isFirst, isLast bool) string {
	opts.depth--
	if opts.depth == -1 {
		return empty
	}

	passPrefix := prefix
	currentPrefix := empty

	if !isFirst {
		if isLast {
			currentPrefix = last
			passPrefix += space
		} else {
			passPrefix += branch
			currentPrefix = middle
		}
	}

	name := n.Name
	if n.IsDir && n.Name != "/" {
		name += "/"
	}

	result := fmt.Sprintf("%s%s%s\n", prefix, currentPrefix, name)
	if opts.showLinks && n.Symlink != "" {
		result = fmt.Sprintf("%s%s%s%s%s\n", prefix, currentPrefix, name, link, n.Symlink)
	}

	for i, child := range n.Children {
		result += child.getString(passPrefix, opts, false, i == len(n.Children)-1)
	}

	return result
}

func (n *fileTreeNode) addChild(file *tar.Header) {
	pathDirs := strings.Split(file.Name, "/")
	for _, dir := range pathDirs {
		if dir == "" {
			continue
		}

		childIndex := n.findChild(dir)
		if childIndex != -1 {
			n = n.Children[childIndex]
			continue
		}

		child := &fileTreeNode{
			Name:    dir,
			IsDir:   file.Typeflag == tar.TypeDir,
			Symlink: file.Linkname,
		}

		n.Children = append(n.Children, child)
	}
}

func (n *fileTreeNode) findChild(name string) int {
	for i, child := range n.Children {
		if child.Name == name {
			return i
		}
	}

	return -1
}

func (n *fileTreeNode) findNode(path string) *fileTreeNode {
	pathDirs := strings.Split(path, "/")
	for _, dir := range pathDirs {
		if dir == "" {
			continue
		}

		childIndex := n.findChild(dir)
		if childIndex == -1 {
			return nil
		}
		n = n.Children[childIndex]
	}

	return n
}

func mergeFileTrees(original, updated *fileTreeNode) (*fileTreeNode, error) {
	if original == updated || updated == nil {
		return original, nil
	}
	if original == nil {
		return updated, nil
	}

	merged := &fileTreeNode{
		Name:     original.Name,
		Symlink:  "",
		IsDir:    original.IsDir,
		Children: original.Children,
	}

	for _, updatedChild := range updated.Children {
		// to avoid "/./" in tree for some images
		if updatedChild.Name == "." {
			continue
		}

		if strings.HasPrefix(updatedChild.Name, whiteoutFilePrefix) {
			updatedChild.Name = strings.TrimPrefix(updatedChild.Name, whiteoutFilePrefix)
			if err := original.deleteNode(updatedChild); err != nil {
				return nil, fmt.Errorf("error deleting file %s: %w", updatedChild.Name, err)
			}
			continue
		}

		childIndex := merged.findChild(updatedChild.Name)
		if childIndex == -1 {
			merged.Children = append(merged.Children, updatedChild)
			sort.Slice(merged.Children, func(i, j int) bool {
				return merged.Children[i].Name < merged.Children[j].Name
			})
		} else {
			_, err := mergeFileTrees(merged.Children[childIndex], updatedChild)
			if err != nil {
				return nil, err
			}
		}
	}

	return merged, nil
}

func (n *fileTreeNode) deleteNode(node *fileTreeNode) error {
	childIndex := n.findChild(node.Name)
	if childIndex != -1 {
		n.Children = append(n.Children[:childIndex], n.Children[childIndex+1:]...)
	}

	return nil
}
