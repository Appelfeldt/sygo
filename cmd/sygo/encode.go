package sygo

import (
	"fmt"
	"os"
	"path"
	"regexp"

	st "github.com/Appelfeldt/sygo/pkg/steganography"
	"github.com/spf13/cobra"
)

var encodeCmd = &cobra.Command{
	Use:   "encode",
	Short: "Encodes data into an image",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		outpath, _ := cmd.Flags().GetString("output")
		if ext := path.Ext(outpath); ext == "" {
			outpath += ".png"
		}

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
			OutputPath:     outpath,
			DataString:     args[1],
			Channels:       channels,
			BitsPerChannel: bpc,
		}

		st.Encode(params)
	},
}

func init() {
	rootCmd.AddCommand(encodeCmd)
	encodeCmd.PersistentFlags().String("output", "encoded.png", "Output filepath")
	encodeCmd.PersistentFlags().Int("bits-per-channel", 1, "How many least-significant-bits to use per channel for data encoding")
	encodeCmd.PersistentFlags().String("channels", "rgb", "Which color channels to use for data encoding. Examples: 'rgba', 'rba' or 'gb'")
}
