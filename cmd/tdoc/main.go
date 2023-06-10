package main

import (
	_ "embed"
	"errors"
	"fmt"
	"log"

	"github.com/charmbracelet/glamour"
	"github.com/spf13/cobra"

	"github.com/zrcoder/tdoc"
	"github.com/zrcoder/tdoc/docmgr"
	"github.com/zrcoder/tdoc/view"
)

var (
	//go:embed help.md
	helpInfo         string
	renderedHelpInfo string
	rootCmd          = &cobra.Command{}
	sortBy           string
)

func main() {
	rootCmd.Run = run
	renderedHelpInfo = renderedMarkdown(helpInfo)
	rootCmd.SetHelpTemplate(renderedHelpInfo)
	rootCmd.PersistentFlags().StringVarP(&sortBy, "sort", "s", "", "sort the docs by title/time, if not set, will keep the original order")

	err := rootCmd.Execute()
	if err != nil {
		fmt.Println(view.ErrStyle.Copy().Render(err.Error()))
	}
}

func run(cmd *cobra.Command, args []string) {
	dir := "."
	if len(args) > 0 {
		dir = args[0]
	}
	mgr, err := docmgr.NewWithDirectory(dir)
	if err != nil {
		printHelpInfo(err)
		return
	}
	if sortBy != "" {
		if sortBy != "time" && sortBy != "title" {
			printHelpInfo(errors.New("only supported sort by time/title"))
		}
		if sortBy == "time" {
			mgr.Sort(docmgr.ByModTime)
		} else {
			mgr.Sort(docmgr.ByTitle)
		}
	}
	err = tdoc.Run(mgr.Docs())
	if err != nil {
		printHelpInfo(err)
	}
}

func printHelpInfo(err error) {
	fmt.Println(renderedHelpInfo)
	fmt.Println(view.ErrStyle.Copy().Render(err.Error()))
}

func renderedMarkdown(ori string) string {
	md, err := glamour.Render(ori, "auto")
	if err != nil {
		log.Fatal(err)
	}
	return md
}
