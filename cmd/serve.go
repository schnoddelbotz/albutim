package cmd

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/schnoddelbotz/albutim/lib"

	"github.com/spf13/cobra"
)

var httpPort string

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "run the built-in webserver to serve your album",
	Long:  `Just try it`,
	Run: func(cmd *cobra.Command, args []string) {
		if _, err := os.Stat(albumRoot); os.IsNotExist(err) {
			fmt.Printf("Album root directory '%s' does not exist!\n", albumRoot)
			os.Exit(1)
		}

		albumRoot = filepath.Clean(albumRoot)
		var err error
		album := new(lib.Album)
		album.Data, err = lib.ScanDir(albumRoot)
		album.Title = albumTitle
		if err != nil {
			log.Fatalf("Cannot scan '%s': %s", albumRoot, err)
		}

		lib.Serve(*album, httpPort)
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
	serveCmd.Flags().StringVar(&httpPort, "port", "3000", "HTTP port to serve on")
}
