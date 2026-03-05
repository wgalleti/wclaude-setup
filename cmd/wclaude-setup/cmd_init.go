package main

import (
	"fmt"

	"github.com/charmbracelet/huh"
	"github.com/spf13/cobra"

	"github.com/wgalleti/wclaude-setup/internal/setup"
	"github.com/wgalleti/wclaude-setup/internal/tui"
)

func newInitCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "init",
		Short: "Setup global do Claude Code (~/.claude)",
		Run: func(cmd *cobra.Command, args []string) {
			runInit()
		},
	}
}

func runInit() {
	fmt.Println(tui.Title.Render("Setup Global"))

	var tasks []string
	form := huh.NewMultiSelect[string]().
		Title("Selecione o que instalar:").
		Options(
			huh.NewOption("CLAUDE.md global", "claude"),
			huh.NewOption("settings.json (permissoes auto)", "settings"),
			huh.NewOption("MCPs essenciais", "mcps"),
		).
		Value(&tasks)

	if err := huh.NewForm(huh.NewGroup(form)).Run(); err != nil {
		return
	}

	for _, t := range tasks {
		switch t {
		case "claude":
			if err := setup.InstallGlobalCLAUDE(); err != nil {
				fmt.Println(tui.Error.Render("  Erro: " + err.Error()))
			} else {
				fmt.Println(tui.Success.Render("  CLAUDE.md global instalado"))
			}
		case "settings":
			if err := setup.InstallSettings(); err != nil {
				fmt.Println(tui.Error.Render("  Erro: " + err.Error()))
			} else {
				fmt.Println(tui.Success.Render("  settings.json configurado"))
			}
		case "mcps":
			runMCPInstall()
		}
	}
}
