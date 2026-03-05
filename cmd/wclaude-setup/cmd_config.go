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
		},
	})

	cmd.AddCommand(&cobra.Command{
		Use:   "show",
		Short: "Mostrar configuracao atual",
		Run: func(cmd *cobra.Command, args []string) {
			cfg, err := config.Load()
			if err != nil {
				fmt.Println(tui.Error.Render(err.Error()))
				return
			}
			fmt.Println(tui.Title.Render("Configuracao"))
			keyDisplay := "(nao configurada)"
			if cfg.AnthropicAPIKey != "" {
				keyDisplay = cfg.AnthropicAPIKey[:8] + "..."
			}
			fmt.Printf("  API Key: %s\n", keyDisplay)
			fmt.Printf("  Model:   %s\n", cfg.DefaultModel)
			fmt.Printf("  Arquivo: %s\n", config.DefaultConfigPath())
		},
	})

	return cmd
}

func runConfig() {
	cfg, err := config.Load()
	if err != nil {
		fmt.Println(tui.Error.Render(err.Error()))
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
	}
	if model != "" {
		cfg.DefaultModel = model
	}

	if err := cfg.Save(); err != nil {
		fmt.Println(tui.Error.Render("Erro ao salvar: " + err.Error()))
		return
	}

	fmt.Println(tui.Success.Render("Configuracao salva em " + config.DefaultConfigPath()))
}
