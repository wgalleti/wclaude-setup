package main

import (
	"fmt"

	"github.com/charmbracelet/huh"
	"github.com/spf13/cobra"

	"github.com/wgalleti/wclaude-setup/internal/config"
	"github.com/wgalleti/wclaude-setup/internal/tui"
)

func newConfigCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "config",
		Short: "Configurar wclaude-setup",
	}

	cmd.AddCommand(&cobra.Command{
		Use:   "set-key",
		Short: "Configurar API key da Anthropic",
		Run: func(cmd *cobra.Command, args []string) {
			runConfig()
			tui.WaitForEnter()
		},
	})

	cmd.AddCommand(&cobra.Command{
		Use:   "show",
		Short: "Mostrar configuracao atual",
		Run: func(cmd *cobra.Command, args []string) {
			runConfigShow()
		},
	})

	return cmd
}

func runConfig() {
	for {
		tui.PrintHeader("Configuracao")

		var action string
		form := huh.NewSelect[string]().
			Title("Configuracao").
			Options(
				huh.NewOption("Configurar API key e model", "set"),
				huh.NewOption("Mostrar configuracao atual", "show"),
				huh.NewOption("<< Voltar", "back"),
			).
			Value(&action)

		if err := huh.NewForm(huh.NewGroup(form)).Run(); err != nil || action == "back" {
			return
		}

		switch action {
		case "set":
			runConfigSet()
		case "show":
			runConfigShow()
		}

		tui.WaitForEnter()
	}
}

func runConfigShow() {
	tui.LogStep("Carregando configuracao...")

	cfg, err := config.Load()
	if err != nil {
		tui.LogError(err.Error())
		return
	}

	keyDisplay := "(nao configurada)"
	if cfg.AnthropicAPIKey != "" {
		keyDisplay = cfg.AnthropicAPIKey[:8] + "..."
	}

	tui.LogSuccess("Configuracao carregada:")
	fmt.Printf("    API Key:  %s\n", keyDisplay)
	fmt.Printf("    Model:    %s\n", cfg.DefaultModel)
	fmt.Printf("    Arquivo:  %s\n", config.DefaultConfigPath())
}

func runConfigSet() {
	cfg, err := config.Load()
	if err != nil {
		tui.LogError(err.Error())
		return
	}

	var apiKey string
	var model string

	if err := huh.NewForm(huh.NewGroup(
		huh.NewInput().
			Title("ANTHROPIC_API_KEY:").
			EchoMode(huh.EchoModePassword).
			Value(&apiKey),
		huh.NewInput().
			Title("Model padrao (ex: claude-sonnet-4-6):").
			Value(&model).
			Placeholder(cfg.DefaultModel),
	)).Run(); err != nil {
		return
	}

	if apiKey != "" {
		cfg.AnthropicAPIKey = apiKey
		tui.LogInfo("API key atualizada")
	}
	if model != "" {
		cfg.DefaultModel = model
		tui.LogInfo("Model atualizado para " + model)
	}

	if apiKey == "" && model == "" {
		tui.LogWarn("Nenhuma alteracao feita")
		return
	}

	tui.LogStep("Salvando configuracao...")
	if err := cfg.Save(); err != nil {
		tui.LogError("Erro ao salvar: " + err.Error())
		return
	}
	tui.LogSuccess("Configuracao salva em " + config.DefaultConfigPath())
}
