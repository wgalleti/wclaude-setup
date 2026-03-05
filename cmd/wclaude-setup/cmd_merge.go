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
	tui.PrintHeader("Merge CLAUDE.md")

	dir, _ := os.Getwd()
	stack := setup.DetectStack(dir)

	tui.LogInfo("Diretorio: " + dir)

	if stack == setup.StackUnknown {
		tui.LogWarn("Stack nao detectada automaticamente")

		var selected string
		form := huh.NewSelect[string]().
			Title("Selecione a stack:").
			Options(
				huh.NewOption("Django / DRF", "django"),
				huh.NewOption("Vue.js / PrimeVue", "vue"),
				huh.NewOption("Flutter / Dart", "flutter"),
				huh.NewOption("Go", "go"),
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

	runMergeForStack(dir, stack)
	tui.WaitForEnter()
}

func runMergeForStack(dir string, stack setup.Stack) {
	existingPath := dir + "/CLAUDE.md"

	tui.LogStep("Lendo CLAUDE.md existente...")
	existing, err := os.ReadFile(existingPath)
	if err != nil {
		tui.LogError("CLAUDE.md nao encontrado em " + dir)
		return
	}
	tui.LogSuccess(fmt.Sprintf("CLAUDE.md lido (%d bytes)", len(existing)))

	tmplName := stack.TemplateName()
	if tmplName == "" {
		tui.LogError("Nenhum template disponivel para " + stack.Label())
		return
	}

	tui.LogStep("Carregando template " + tmplName + "...")
	tmplContent, err := templates.FS.ReadFile("files/" + tmplName)
	if err != nil {
		tui.LogError("Erro ao ler template: " + err.Error())
		return
	}
	tui.LogSuccess(fmt.Sprintf("Template carregado (%d bytes)", len(tmplContent)))

	cfg, err := config.Load()
	if err != nil {
		tui.LogError("Erro ao carregar config: " + err.Error())
		return
	}

	if cfg.AnthropicAPIKey == "" {
		cfg.AnthropicAPIKey = os.Getenv("ANTHROPIC_API_KEY")
	}

	if cfg.AnthropicAPIKey == "" {
		tui.LogWarn("API key nao configurada")

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
				tui.LogError("Erro ao salvar config: " + err.Error())
			} else {
				tui.LogSuccess("API key salva em " + config.DefaultConfigPath())
			}
		}
	}

	tui.LogStep("Enviando para API Anthropic (model: " + cfg.DefaultModel + ")...")

	result, err := merge.MergeFiles(string(existing), string(tmplContent), cfg)
	if err != nil {
		tui.LogError("Falha no merge: " + err.Error())
		return
	}
	tui.LogSuccess(fmt.Sprintf("Merge concluido (%d bytes gerados)", len(result)))

	backup := existingPath + ".bak"
	tui.LogStep("Criando backup em " + backup + "...")
	if err := os.WriteFile(backup, existing, 0o644); err != nil {
		tui.LogError("Erro ao criar backup: " + err.Error())
		return
	}
	tui.LogSuccess("Backup criado")

	tui.LogStep("Escrevendo CLAUDE.md atualizado...")
	if err := os.WriteFile(existingPath, []byte(result), 0o644); err != nil {
		tui.LogError("Erro ao escrever resultado: " + err.Error())
		return
	}

	tui.LogSuccess("CLAUDE.md atualizado com sucesso!")
}
