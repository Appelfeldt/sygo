package sygo

import (
	"fmt"

	st "github.com/Appelfeldt/steganography/pkg/steganography"
	"github.com/spf13/cobra"
)

var extractCmd = &cobra.Command{
	Use:   "extract",
	Short: "Extracts data embedded in an image",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		res := st.Extract(args[0])
		fmt.Printf("%s\n", string(res))
	},
}

func init() {
	rootCmd.AddCommand(extractCmd)
}
