package sygo

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var BuildVersion string

var rootCmd = &cobra.Command{
	Use:     "sygo",
	Version: BuildVersion,
	Short:   "sygo - A basic steganography tool",
	Long:    "sygo is a steganography CLI tool\n\nIt can be used to encode data into an image and decode data encoded in an image",
	Run: func(cmd *cobra.Command, args []string) {

	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "An error occurred '%s'", err)
		os.Exit(1)
	}
}
