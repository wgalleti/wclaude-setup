package main

import (
	"fmt"
	"os"

	"github.com/charmbracelet/huh"
	"github.com/spf13/cobra"

	"github.com/wgalleti/wclaude-setup/internal/config"
	"github.com/wgalleti/wclaude-setup/internal/merge"
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
			initGlobalCLAUDE()
		case "settings":
			initSettings()
		case "mcps":
			runMCPInstall()
		}
	}

	tui.WaitForEnter()
}

func initGlobalCLAUDE() {
	tui.LogStep("Verificando CLAUDE.md global...")

	if setup.GlobalCLAUDEExists() {
		tui.LogWarn("CLAUDE.md ja existe em " + setup.GlobalCLAUDEPath())

		var action string
		form := huh.NewSelect[string]().
			Title("CLAUDE.md ja existe. O que fazer?").
			Options(
				huh.NewOption("Merge inteligente (via API Anthropic)", "merge"),
				huh.NewOption("Sobrescrever com template padrao", "overwrite"),
				huh.NewOption("Pular", "skip"),
			).
			Value(&action)

		if err := huh.NewForm(huh.NewGroup(form)).Run(); err != nil {
			return
		}

		switch action {
		case "merge":
			initGlobalCLAUDEMerge()
			return
		case "overwrite":
			tui.LogStep("Sobrescrevendo CLAUDE.md...")
			res, err := setup.InstallGlobalCLAUDE()
			if err != nil {
				tui.LogError(err.Error())
				return
			}
			tui.LogSuccess("Backup criado em " + res.BackupPath)
			tui.LogSuccess("CLAUDE.md sobrescrito em " + res.DestPath)
		case "skip":
			tui.LogInfo("CLAUDE.md mantido sem alteracoes")
		}
		return
	}

	tui.LogStep("Instalando CLAUDE.md global...")
	res, err := setup.InstallGlobalCLAUDE()
	if err != nil {
		tui.LogError(err.Error())
		return
	}
	tui.LogSuccess("CLAUDE.md instalado em " + res.DestPath)
}

func initGlobalCLAUDEMerge() {
	tui.LogStep("Lendo CLAUDE.md existente...")
	existing, err := setup.ReadGlobalCLAUDE()
	if err != nil {
		tui.LogError("Erro ao ler existente: " + err.Error())
		return
	}
	tui.LogSuccess(fmt.Sprintf("Existente lido (%d bytes)", len(existing)))

	tui.LogStep("Carregando template global...")
	template, err := setup.GlobalCLAUDETemplate()
	if err != nil {
		tui.LogError("Erro ao ler template: " + err.Error())
		return
	}
	tui.LogSuccess(fmt.Sprintf("Template carregado (%d bytes)", len(template)))

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

		if err := huh.NewForm(huh.NewGroup(field)).Run(); err != nil || apiKey == "" {
			tui.LogError("API key obrigatoria para merge")
			return
		}

		cfg.AnthropicAPIKey = apiKey

		var saveKey bool
		confirm := huh.NewConfirm().
			Title("Salvar API key para uso futuro?").
			Value(&saveKey)

		if err := huh.NewForm(huh.NewGroup(confirm)).Run(); err == nil && saveKey {
			if err := cfg.Save(); err != nil {
				tui.LogError("Erro ao salvar config: " + err.Error())
			} else {
				tui.LogSuccess("API key salva")
			}
		}
	}

	tui.LogStep("Enviando para API Anthropic (model: " + cfg.DefaultModel + ")...")

	result, err := merge.MergeFiles(existing, template, cfg)
	if err != nil {
		tui.LogError("Falha no merge: " + err.Error())
		return
	}
	tui.LogSuccess(fmt.Sprintf("Merge concluido (%d bytes)", len(result)))

	// Backup antes de escrever
	backup := setup.GlobalCLAUDEPath() + ".bak"
	tui.LogStep("Criando backup...")
	if err := os.WriteFile(backup, []byte(existing), 0o644); err != nil {
		tui.LogError("Erro ao criar backup: " + err.Error())
		return
	}
	tui.LogSuccess("Backup em " + backup)

	tui.LogStep("Escrevendo CLAUDE.md atualizado...")
	if err := setup.WriteGlobalCLAUDE(result); err != nil {
		tui.LogError("Erro ao escrever: " + err.Error())
		return
	}
	tui.LogSuccess("CLAUDE.md global atualizado com merge!")
}

func initSettings() {
	tui.LogStep("Verificando settings.json...")

	if setup.SettingsExists() {
		tui.LogWarn("settings.json ja existe em " + setup.SettingsPath())

		var action string
		form := huh.NewSelect[string]().
			Title("settings.json ja existe. O que fazer?").
			Options(
				huh.NewOption("Mesclar permissoes (adicionar novas sem remover existentes)", "merge"),
				huh.NewOption("Sobrescrever com permissoes padrao", "overwrite"),
				huh.NewOption("Pular", "skip"),
			).
			Value(&action)

		if err := huh.NewForm(huh.NewGroup(form)).Run(); err != nil {
			return
		}

		switch action {
		case "merge":
			tui.LogStep("Lendo settings.json existente...")
			existing, err := setup.ReadSettings()
			if err != nil {
				tui.LogError("Erro ao ler settings.json: " + err.Error())
				return
			}
			tui.LogInfo(fmt.Sprintf("Permissoes existentes: %d", len(existing.Permissions.Allow)))

			merged := setup.MergeSettings(existing)
			tui.LogInfo(fmt.Sprintf("Permissoes apos merge: %d", len(merged.Permissions.Allow)))

			if len(merged.Permissions.Allow) == len(existing.Permissions.Allow) {
				tui.LogSuccess("Nenhuma permissao nova para adicionar, ja esta atualizado")
				return
			}

			added := len(merged.Permissions.Allow) - len(existing.Permissions.Allow)
			tui.LogStep(fmt.Sprintf("Adicionando %d permissao(oes) nova(s)...", added))

			if err := setup.WriteSettings(merged); err != nil {
				tui.LogError("Erro ao escrever settings.json: " + err.Error())
				return
			}
			tui.LogSuccess("settings.json atualizado com novas permissoes")

		case "overwrite":
			tui.LogStep("Sobrescrevendo settings.json...")
			os.Remove(setup.SettingsPath())
			res, err := setup.InstallSettings()
			if err != nil {
				tui.LogError(err.Error())
				return
			}
			tui.LogSuccess("settings.json sobrescrito em " + res.DestPath)

		case "skip":
			tui.LogInfo("settings.json mantido sem alteracoes")
		}
		return
	}

	tui.LogStep("Criando settings.json...")
	res, err := setup.InstallSettings()
	if err != nil {
		tui.LogError(err.Error())
		return
	}
	tui.LogSuccess("settings.json criado em " + res.DestPath)
}
