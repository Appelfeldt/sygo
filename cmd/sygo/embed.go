package sygo

import (
	st "github.com/Appelfeldt/steganography/pkg/steganography"
	"github.com/spf13/cobra"
)

var embedCmd = &cobra.Command{
	Use:   "embed",
	Short: "Embeds data in an image",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		outpath, _ := cmd.Flags().GetString("out")

		if outpath != "" {
			st.Embed(args[0], args[1], outpath)
		} else {
			st.Embed(args[0], args[1], "embedded.png")
		}
	},
}

func init() {
	rootCmd.AddCommand(embedCmd)
	embedCmd.PersistentFlags().String("out", "", "Output filepath")
}
