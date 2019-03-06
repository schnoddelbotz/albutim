package cmd

import (
	"log"
	"os"
	"path/filepath"

	"github.com/schnoddelbotz/albutim/lib"

	"github.com/spf13/cobra"
)

var httpPort string

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "run the built-in webserver to serve your album",
	Long:  `Just try it`,
	Run: func(cmd *cobra.Command, args []string) {
		if _, err := os.Stat(albumRoot); os.IsNotExist(err) {
			log.Fatalf("Album root directory '%s' does not exist!\n", albumRoot)
		}
		albumRoot = filepath.Clean(albumRoot)

		albumData, err := lib.ScanDir(albumRoot)
		if err != nil {
			log.Fatalf("Cannot scan '%s': %s", albumRoot, err)
		}

		album := &lib.Album{
			SubTitle: "all the fun pics!",
			RootPath: albumRoot,
			Title:    albumTitle,
			Data:     albumData}
		lib.Serve(*album, httpPort)
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
	serveCmd.Flags().StringVar(&httpPort, "port", "3000", "HTTP port to serve on")
}
