package cmd

import (
	"crypto/sha256"
	"fmt"
	"outline-cli/internal/cli"
	"outline-cli/internal/models"
	syncp "outline-cli/internal/sync"

	"github.com/spf13/cobra"
)

var diffCmd = &cobra.Command{
	Use:   "diff <collection-name> <folder-path>",
	Short: "Compare local folder against remote collection",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		collectionName := args[0]
		folderPath := args[1]

		client, err := getClient()
		if err != nil {
			return err
		}

		localTree, err := syncp.BuildLocalTree(folderPath)
		if err != nil {
			return fmt.Errorf("scan folder: %w", err)
		}

		coll, err := client.Collections.FindByName(getContext(), collectionName)
		if err != nil {
			return fmt.Errorf("find collection: %w", err)
		}
		if coll == nil {
			cli.Output.Info("Collection '%s' does not exist remotely. All local files would be created.", collectionName)
			printLocalTree(localTree, 0, "+")
			return nil
		}

		remoteNodes, err := client.Collections.Documents(getContext(), coll.ID)
		if err != nil {
			return fmt.Errorf("get remote tree: %w", err)
		}

		stats := diffNodes(localTree.Children, remoteNodes, 0)

		fmt.Printf("\nSummary: %d new, %d modified, %d unchanged, %d remote-only\n",
			stats.added, stats.modified, stats.unchanged, stats.remoteOnly)
		return nil
	},
}

type diffStats struct {
	added      int
	modified   int
	unchanged  int
	remoteOnly int
}

func diffNodes(localNodes []*syncp.LocalNode, remoteNodes []models.NavigationNode, depth int) diffStats {
	stats := diffStats{}
	indent := ""
	for i := 0; i < depth; i++ {
		indent += "  "
	}

	remoteTitles := map[string]models.NavigationNode{}
	for _, n := range remoteNodes {
		remoteTitles[n.Title] = n
	}

	localTitles := map[string]bool{}
	for _, n := range localNodes {
		localTitles[n.Title] = true
		remote, exists := remoteTitles[n.Title]

		if n.IsDir {
			var childRemote []models.NavigationNode
			if exists {
				childRemote = remote.Children
			}
			fmt.Printf("%s  %s/\n", indent, n.Title)
			childStats := diffNodes(n.Children, childRemote, depth+1)
			stats.added += childStats.added
			stats.modified += childStats.modified
			stats.unchanged += childStats.unchanged
			stats.remoteOnly += childStats.remoteOnly
		} else if exists {
			fmt.Printf("%s~ %s (exists, would update)\n", indent, n.Title)
			stats.modified++
		} else {
			fmt.Printf("%s+ %s (new)\n", indent, n.Title)
			stats.added++
		}
	}

	for _, n := range remoteNodes {
		if !localTitles[n.Title] {
			fmt.Printf("%s- %s (remote only)\n", indent, n.Title)
			stats.remoteOnly++
		}
	}

	return stats
}

func printLocalTree(node *syncp.LocalNode, depth int, prefix string) {
	indent := ""
	for i := 0; i < depth; i++ {
		indent += "  "
	}
	if node.IsDir {
		fmt.Printf("%s%s %s/\n", indent, prefix, node.Title)
	} else {
		fmt.Printf("%s%s %s\n", indent, prefix, node.Title)
	}
	for _, child := range node.Children {
		printLocalTree(child, depth+1, prefix)
	}
}

func contentHash(content string) string {
	h := sha256.Sum256([]byte(content))
	return fmt.Sprintf("%x", h[:8])
}
