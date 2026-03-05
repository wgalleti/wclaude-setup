package main

import (
	"fmt"
	"os"
	"strings"

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

	fmt.Println(tui.Title.Render("Setup Projeto"))
	fmt.Printf("  Diretorio: %s\n", tui.Info.Render(dir))

	if stack == setup.StackUnknown {
		fmt.Println(tui.Info.Render("  Stack nao detectada automaticamente"))

		var selected string
		form := huh.NewSelect[string]().
			Title("Selecione a stack:").
			Options(
				huh.NewOption("Django / DRF", "django"),
				huh.NewOption("Vue.js / PrimeVue", "vue"),
				huh.NewOption("Flutter / Dart", "flutter"),
				huh.NewOption("Go", "go"),
			).
			Value(&selected)

		if err := huh.NewForm(huh.NewGroup(form)).Run(); err != nil {
			return
		}
		stack = setup.Stack(selected)
	} else {
		fmt.Printf("  Stack detectada: %s\n", tui.Success.Render(stack.Label()))
	}

	if err := setup.InstallProjectCLAUDE(dir, stack); err != nil {
		if strings.Contains(err.Error(), "ja existe") {
			fmt.Println(tui.Info.Render("  " + err.Error()))

			var doMerge bool
			confirm := huh.NewConfirm().
				Title("Deseja fazer merge com o template?").
				Value(&doMerge)

			if err := huh.NewForm(huh.NewGroup(confirm)).Run(); err != nil {
				return
			}
			if doMerge {
				runMergeForStack(dir, stack)
			}
		} else {
			fmt.Println(tui.Error.Render("  Erro: " + err.Error()))
		}
		return
	}

	fmt.Println(tui.Success.Render("  Projeto configurado! Edite o CLAUDE.md com os dados do projeto."))
}
