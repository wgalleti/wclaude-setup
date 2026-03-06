package main

import (
	"os"

	"github.com/charmbracelet/huh"
	"github.com/spf13/cobra"

	"github.com/wgalleti/wclaude-setup/internal/setup"
	"github.com/wgalleti/wclaude-setup/internal/tui"
)

func newProjectCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "project",
		Short: "Setup CLAUDE.md para o projeto no diretorio atual",
		Run: func(cmd *cobra.Command, args []string) {
			runProject()
		},
	}
}

func runProject() {
	dir, _ := os.Getwd()
	stack := setup.DetectStack(dir)

	tui.PrintHeader("Setup Projeto")
	tui.LogInfo("Diretorio: " + dir)

	if stack == setup.StackUnknown {
		tui.LogWarn("Stack nao detectada automaticamente")

		var selected string
		form := huh.NewSelect[string]().
			Title("Selecione a stack:").
			Options(
				huh.NewOption("Django / DRF", "django"),
				huh.NewOption("Vue.js / PrimeVue", "vue"),
				huh.NewOption("Nuxt / Vue.js", "nuxt"),
				huh.NewOption("React / Vite", "react"),
				huh.NewOption("Next.js / React", "next"),
				huh.NewOption("Flutter / Dart", "flutter"),
				huh.NewOption("Go", "go"),
				huh.NewOption("Supabase", "supabase"),
				huh.NewOption("<< Voltar", "back"),
			).
			Value(&selected)

		if err := huh.NewForm(huh.NewGroup(form)).Run(); err != nil || selected == "back" {
			return
		}
		stack = setup.Stack(selected)
	} else {
		tui.LogSuccess("Stack detectada: " + stack.Label())
	}

	tui.LogStep("Verificando CLAUDE.md...")

	res, err := setup.InstallProjectCLAUDE(dir, stack)
	if err != nil {
		tui.LogError(err.Error())
		tui.WaitForEnter()
		return
	}

	if res.Existed {
		tui.LogWarn("CLAUDE.md ja existe em " + res.DestPath)

		var doMerge bool
		confirm := huh.NewConfirm().
			Title("Deseja fazer merge com o template?").
			Value(&doMerge)

		if err := huh.NewForm(huh.NewGroup(confirm)).Run(); err != nil {
			return
		}
		if doMerge {
			runMergeForStack(dir, stack)
		} else {
			tui.LogInfo("CLAUDE.md mantido sem alteracoes")
		}
		tui.WaitForEnter()
		return
	}

	tui.LogSuccess("CLAUDE.md criado em " + res.DestPath)
	tui.LogInfo("Edite o arquivo com os dados especificos do projeto")
	tui.WaitForEnter()
}
