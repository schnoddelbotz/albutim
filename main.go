//go:generate esc -prefix assets/ -pkg lib -o lib/assets.go -private assets

package main

import "github.com/schnoddelbotz/albutim/cmd"

func main() {
	cmd.Execute()
}