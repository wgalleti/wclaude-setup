package main

import (
	"fmt"
	"os"

	"github.com/charmbracelet/huh"
	"github.com/spf13/cobra"

	"github.com/wgalleti/wclaude-setup/internal/tui"
)

var version = "0.1.0"

func main() {
	rootCmd := &cobra.Command{
		Use:     "wclaude-setup",
		Short:   "CLI para configurar Claude Code em projetos",
		Version: version,
		Run:     runInteractive,
	}

	rootCmd.AddCommand(
		newInitCmd(),
		newProjectCmd(),
		newMCPCmd(),
		newPromptCmd(),
		newMergeCmd(),
		newConfigCmd(),
	)

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func runInteractive(cmd *cobra.Command, args []string) {
	fmt.Println(tui.Title.Render("wclaude-setup"))

	var action string
	form := huh.NewSelect[string]().
		Title("O que deseja fazer?").
		Options(
			huh.NewOption("Setup global (~/.claude)", "init"),
			huh.NewOption("Setup projeto (diretorio atual)", "project"),
			huh.NewOption("Gerenciar MCPs", "mcp"),
			huh.NewOption("Gerar prompt eficiente", "prompt"),
			huh.NewOption("Merge CLAUDE.md (existente + template)", "merge"),
			huh.NewOption("Configurar API key", "config"),
			huh.NewOption("Sair", "exit"),
		).
		Value(&action)

	if err := huh.NewForm(huh.NewGroup(form)).Run(); err != nil {
		return
	}

	switch action {
	case "init":
		runInit()
	case "project":
		runProject()
	case "mcp":
		runMCP()
	case "prompt":
		runPrompt()
	case "merge":
		runMerge()
	case "config":
		runConfig()
	}
}
