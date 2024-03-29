package cmd

import (
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/schnoddelbotz/albutim/lib"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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

		album := &lib.Album{
			SubTitle:         "all the fun pics!",
			RootPath:         albumRoot,
			AlbutimVersion:   AppVersion,
			CreatedAt:        time.Now().String(),
			Title:            viper.GetString("title"),
			NoScaledPreviews: viper.GetBool("no-scaled-previews"),
			NoScaledThumbs:   viper.GetBool("no-scaled-thumbs"),
			NoCacheScaled:    viper.GetBool("no-cache-scaled")}

		var e error
		album.Data, e = lib.ScanDir(albumRoot, album)
		if e != nil {
			log.Fatalf("Cannot scan '%s': %s", albumRoot, e)
		}

		lib.Serve(*album, httpPort)
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
	serveCmd.Flags().StringVar(&httpPort, "port", "3000", "HTTP port to serve on")

	serveCmd.PersistentFlags().BoolP("no-cache-scaled", "n", false, "do not cache scaled thumbs/previews")
	_ = viper.BindPFlag("no-cache-scaled", serveCmd.PersistentFlags().Lookup("no-cache-scaled"))
}
