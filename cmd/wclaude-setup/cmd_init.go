package main

import (
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
	tui.PrintHeader("Setup Global")

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

	if len(tasks) == 0 {
		tui.LogWarn("Nenhuma opcao selecionada")
		tui.WaitForEnter()
		return
	}

	for _, t := range tasks {
		switch t {
		case "claude":
			tui.LogStep("Instalando CLAUDE.md global...")
			if err := setup.InstallGlobalCLAUDE(); err != nil {
				tui.LogError(err.Error())
			} else {
				tui.LogSuccess("CLAUDE.md global instalado")
			}
		case "settings":
			tui.LogStep("Configurando settings.json...")
			if err := setup.InstallSettings(); err != nil {
				tui.LogError(err.Error())
			} else {
				tui.LogSuccess("settings.json configurado")
			}
		case "mcps":
			runMCPInstall()
		}
	}

	tui.WaitForEnter()
}
