package sygo

import (
	"fmt"
	"os"
	"regexp"

	st "github.com/Appelfeldt/steganography/pkg/steganography"
	"github.com/spf13/cobra"
)

var extractCmd = &cobra.Command{
	Use:   "extract",
	Short: "Extracts data embedded in an image",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		bpc, _ := cmd.Flags().GetInt("bits-per-channel")
		if bpc < 1 || bpc > 8 {
			fmt.Fprintf(os.Stderr, "Invalid bits-per-channel value: %d", bpc)
			os.Exit(1)
		}

		channels, _ := cmd.Flags().GetString("channels")
		match, _ := regexp.MatchString("^(r?g?b?a?)$", channels)
		if length := len(channels); length < 1 || length > 4 || !match {
			fmt.Fprintf(os.Stderr, "Invalid channels value: %s", channels)
			os.Exit(1)
		}

		params := st.WorkParams{
			InputPath:      args[0],
			Channels:       channels,
			BitsPerChannel: bpc,
		}

		res := st.Extract(params)

		fmt.Printf("%s\n", string(res))
	},
}

func init() {
	rootCmd.AddCommand(extractCmd)
	extractCmd.PersistentFlags().Int("bits-per-channel", 1, "Amount of bits used for embedding data per pixel, per channel. Expects value 1-8")
	extractCmd.PersistentFlags().String("channels", "rgb", "Which color channels to use for embedding. Expects strings such as 'rgba', 'rba' or 'gb'")
}
