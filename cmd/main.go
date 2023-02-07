package main

import (
	_ "embed"
	"errors"
	"fmt"
	"log"

	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"

	"github.com/zrcoder/tdoc"
)

var (
	//go:embed help.md
	helpInfo string

	renderedHelpInfo string
	rootCmd          = &cobra.Command{}
	sortBy           string
	errStyle         lipgloss.Style
)

func main() {
	rootCmd.Run = run
	renderedHelpInfo = renderedMarkdown(helpInfo)
	rootCmd.SetHelpTemplate(renderedHelpInfo)
	rootCmd.PersistentFlags().StringVarP(&sortBy, "sort", "s", "", "sort the docs by title/time, if not set, will keep the original order")
	errStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#f00"))

	err := rootCmd.Execute()
	if err != nil {
		fmt.Println(errStyle.Render(err.Error()))
	}
}

func run(cmd *cobra.Command, args []string) {
	dir := "."
	if len(args) > 0 {
		dir = args[0]
	}
	docs, err := tdoc.ParseFromDir(dir)
	if err != nil {
		printHelpInfo(err)
		return
	}
	if sortBy != "" {
		if sortBy != "time" && sortBy != "title" {
			printHelpInfo(errors.New("only supported sort by time/title"))
		}
		if sortBy == "time" {
			tdoc.Sort(docs, tdoc.ByModTime)
		} else {
			tdoc.Sort(docs, tdoc.ByTitle)
		}
	}
	err = tdoc.Run(docs)
	if err != nil {
		printHelpInfo(err)
	}
}

func printHelpInfo(err error) {
	fmt.Println(renderedHelpInfo)
	fmt.Println(errStyle.Render(err.Error()))
}

func renderedMarkdown(ori string) string {
	md, err := glamour.Render(ori, "auto")
	if err != nil {
		log.Fatal(err)
	}
	return md
}
