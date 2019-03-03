//go:generate esc -pkg lib -prefix "templates" -o lib/templates.go -private templates

package main

import "github.com/schnoddelbotz/albutim/cmd"

func main() {
	cmd.Execute()
}