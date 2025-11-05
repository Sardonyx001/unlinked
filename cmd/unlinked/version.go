package main

import (
	"fmt"

	"github.com/sardonyx001/unlinked/internal/version"
	"github.com/spf13/cobra"
)

var (
	shortVersion bool
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print version information",
	Long:  `Print detailed version information including build time, git commit, and platform.`,
	Run: func(cmd *cobra.Command, args []string) {
		info := version.Get()
		if shortVersion {
			fmt.Println(info.Short())
		} else {
			fmt.Println(info.String())
		}
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
	versionCmd.Flags().BoolVarP(&shortVersion, "short", "s", false, "print short version")
}
