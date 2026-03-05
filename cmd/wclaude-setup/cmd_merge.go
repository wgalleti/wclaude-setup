package main

import (
	"fmt"
	"os"

	"github.com/charmbracelet/huh"
	"github.com/spf13/cobra"

	"github.com/wgalleti/wclaude-setup/internal/config"
	"github.com/wgalleti/wclaude-setup/internal/merge"
	"github.com/wgalleti/wclaude-setup/internal/setup"
	"github.com/wgalleti/wclaude-setup/internal/templates"
	"github.com/wgalleti/wclaude-setup/internal/tui"
)

func newMergeCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "merge",
		Short: "Merge inteligente do CLAUDE.md existente com template (via API Anthropic)",
		Run: func(cmd *cobra.Command, args []string) {
			runMerge()
		},
	}
}

func runMerge() {
	dir, _ := os.Getwd()
	stack := setup.DetectStack(dir)

	if stack == setup.StackUnknown {
		var selected string
		form := huh.NewSelect[string]().
			Title("Stack nao detectada. Selecione:").
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
	}

	runMergeForStack(dir, stack)
}

func runMergeForStack(dir string, stack setup.Stack) {
	fmt.Println(tui.Title.Render("Merge CLAUDE.md"))

	existingPath := dir + "/CLAUDE.md"
	existing, err := os.ReadFile(existingPath)
	if err != nil {
		fmt.Println(tui.Error.Render("CLAUDE.md nao encontrado no diretorio atual"))
		return
	}

	tmplName := stack.TemplateName()
	if tmplName == "" {
		fmt.Println(tui.Error.Render("Nenhum template para essa stack"))
		return
	}

	tmplContent, err := templates.FS.ReadFile("files/" + tmplName)
	if err != nil {
		fmt.Println(tui.Error.Render("Erro ao ler template: " + err.Error()))
		return
	}

	cfg, err := config.Load()
	if err != nil {
		fmt.Println(tui.Error.Render("Erro ao carregar config: " + err.Error()))
		return
	}

	if cfg.AnthropicAPIKey == "" {
		cfg.AnthropicAPIKey = os.Getenv("ANTHROPIC_API_KEY")
	}

	if cfg.AnthropicAPIKey == "" {
		var apiKey string
		field := huh.NewInput().
			Title("ANTHROPIC_API_KEY (necessaria para merge):").
			EchoMode(huh.EchoModePassword).
			Value(&apiKey)

		if err := huh.NewForm(huh.NewGroup(field)).Run(); err != nil {
			return
		}

		cfg.AnthropicAPIKey = apiKey

		var saveKey bool
		confirm := huh.NewConfirm().
			Title("Salvar API key para uso futuro?").
			Value(&saveKey)

		if err := huh.NewForm(huh.NewGroup(confirm)).Run(); err != nil {
			return
		}

		if saveKey {
			if err := cfg.Save(); err != nil {
				fmt.Println(tui.Error.Render("Erro ao salvar config: " + err.Error()))
			}
		}
	}

	fmt.Println(tui.Info.Render("  Enviando para API Anthropic..."))

	result, err := merge.MergeFiles(string(existing), string(tmplContent), cfg)
	if err != nil {
		fmt.Println(tui.Error.Render("  Erro no merge: " + err.Error()))
		return
	}

	backup := existingPath + ".bak"
	if err := os.WriteFile(backup, existing, 0o644); err != nil {
		fmt.Println(tui.Error.Render("  Erro ao criar backup: " + err.Error()))
		return
	}
	fmt.Printf("  Backup: %s\n", backup)

	if err := os.WriteFile(existingPath, []byte(result), 0o644); err != nil {
		fmt.Println(tui.Error.Render("  Erro ao escrever resultado: " + err.Error()))
		return
	}

	fmt.Println(tui.Success.Render("  CLAUDE.md atualizado com sucesso!"))
}
