
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// buildCmd represents the build command
var buildCmd = &cobra.Command{
	Use:   "build",
	Short: "builds the album, suitable for static web-serving",
	Long: `The build command scans your original images root, retrieves EXIF meta
data from images and generates required index.html files and thumbnails. Usage example:

albutim --root my-images build --output my-album`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("build called")
	},
}

func init() {
	rootCmd.AddCommand(buildCmd)
}
