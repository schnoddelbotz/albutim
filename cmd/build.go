package cmd

import (
	"log"
	"os"
	"path/filepath"
	"runtime"

	"github.com/schnoddelbotz/albutim/lib"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var threads int

// buildCmd represents the build command
var buildCmd = &cobra.Command{
	Use:   "build",
	Short: "builds the album, suitable for static web-serving",
	Long: `The build command scans your original images root, retrieves EXIF meta
data from images and generates required index.html files and thumbnails. Usage example:

albutim --root my-images build --output my-album`,
	Run: func(cmd *cobra.Command, args []string) {
		if _, err := os.Stat(albumRoot); os.IsNotExist(err) {
			log.Fatalf("Album root directory '%s' does not exist!\n", albumRoot)
		}
		albumRoot = filepath.Clean(albumRoot)

		album := &lib.Album{
			SubTitle:         viper.GetString("tagline"),
			RootPath:         albumRoot,
			Title:            viper.GetString("title"),
			NoScaledPreviews: viper.GetBool("no-scaled-previews"),
			NoScaledThumbs:   viper.GetBool("no-scaled-thumbs"),
			ServeStatically:  true,
			AlbutimVersion:   AppVersion,
			NumThreads:       threads}

		var e error
		album.Data, e = lib.ScanDir(albumRoot, album)
		if e != nil {
			log.Fatalf("Cannot scan '%s': %s", albumRoot, e)
		}

		lib.BuildAlbum(*album)
	},
}

func init() {
	rootCmd.AddCommand(buildCmd)
	buildCmd.PersistentFlags().IntVar(&threads, "threads", runtime.NumCPU(), "threads for thumb generation")
}
