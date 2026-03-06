package main

import (
	"fmt"
	"os"

	"github.com/charmbracelet/huh"
	"github.com/spf13/cobra"

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

	cfg, err := ensureAPIKey()
	if cfg == nil {
		return
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
