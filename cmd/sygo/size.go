package sygo

import (
	"fmt"

	st "github.com/Appelfeldt/steganography/pkg/steganography"
	"github.com/spf13/cobra"
)

var sizeCmd = &cobra.Command{
	Use:   "size",
	Short: "Calculates a strings size when embedded",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		bits := st.Size(args[0])
		fmt.Printf("%d bits\n", bits)
	},
}

func init() {
	rootCmd.AddCommand(sizeCmd)
}
