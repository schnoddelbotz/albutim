package cmd

import (
	"fmt"
	"os"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var AppVersion string
var cfgFile string
var albumRoot string
var albumTitle string

var rootCmd = &cobra.Command{
	Version: AppVersion,
	Use:     "albutim",
	Short:   "albutim is yet another photo album generator and server",
	Long: `Provided with a image root folder, albutim can generate a HTML
photo album -- either to be served statically or by using the built-in server.

Further documentation: https://github.com/schnoddelbotz/albutim`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.albutim.yaml)")
	rootCmd.PersistentFlags().StringVar(&albumRoot, "root", "", "album/original images root path")
	//rootCmd.PersistentFlags().StringP("root", "r", "", "album/original images root path")
	rootCmd.PersistentFlags().StringVar(&albumTitle, "title", "Yet another timalbum", "album title")

	// FIXME no worky worky coming from config file...
	// https://stackoverflow.com/questions/41578264/config-file-with-cobra-and-viper ??
	// viper.BindPFlag("root", rootCmd.PersistentFlags().Lookup("root"))

	rootCmd.MarkPersistentFlagRequired("root")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".albutim" (without extension).
		viper.AddConfigPath(home)
		viper.AddConfigPath(".")
		viper.SetConfigName(".albutim")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
