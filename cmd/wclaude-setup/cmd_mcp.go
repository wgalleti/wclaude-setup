package main

import (
	"fmt"
	"os"

	"github.com/charmbracelet/huh"
	"github.com/spf13/cobra"

	"github.com/wgalleti/wclaude-setup/internal/mcp"
	"github.com/wgalleti/wclaude-setup/internal/tui"
)

func newMCPCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "mcp",
		Short: "Gerenciar MCPs do Claude Code",
		Run: func(cmd *cobra.Command, args []string) {
			runMCP()
		},
	}

	cmd.AddCommand(&cobra.Command{
		Use:   "list",
		Short: "Listar MCPs instalados",
		Run: func(cmd *cobra.Command, args []string) {
			runMCPList()
		},
	})

	return cmd
}

func runMCP() {
	for {
		tui.PrintHeader("Gerenciar MCPs")

		var action string
		form := huh.NewSelect[string]().
			Title("MCPs").
			Options(
				huh.NewOption("Instalar MCPs", "install"),
				huh.NewOption("Remover MCP", "remove"),
				huh.NewOption("Listar instalados", "list"),
				huh.NewOption("<< Voltar", "back"),
			).
			Value(&action)

		if err := huh.NewForm(huh.NewGroup(form)).Run(); err != nil || action == "back" {
			return
		}

		switch action {
		case "install":
			runMCPInstall()
		case "remove":
			runMCPRemove()
		case "list":
			runMCPList()
		}

		tui.WaitForEnter()
	}
}

func runMCPList() {
	tui.LogStep("Consultando MCPs instalados...")

	installed, err := mcp.ListInstalled()
	if err != nil {
		tui.LogError(err.Error())
		return
	}

	if len(installed) == 0 {
		tui.LogWarn("Nenhum MCP instalado")
		return
	}

	tui.LogSuccess(fmt.Sprintf("%d MCP(s) encontrado(s):", len(installed)))
	for _, name := range installed {
		fmt.Printf("    - %s\n", name)
	}
}

func runMCPInstall() {
	tui.LogStep("Consultando MCPs instalados...")
	installed, _ := mcp.ListInstalled()

	var options []huh.Option[string]
	for _, m := range mcp.AvailableMCPs {
		label := m.Name + " — " + m.Description
		if mcp.IsInstalled(m.Name, installed) {
			label += " (instalado)"
		}
		options = append(options, huh.NewOption(label, m.Name))
	}

	var selected []string
	form := huh.NewMultiSelect[string]().
		Title("Selecione os MCPs para instalar:").
		Options(options...).
		Value(&selected)

	if err := huh.NewForm(huh.NewGroup(form)).Run(); err != nil {
		return
	}

	if len(selected) == 0 {
		tui.LogWarn("Nenhum MCP selecionado")
		return
	}

	for _, name := range selected {
		if mcp.IsInstalled(name, installed) {
			tui.LogWarn(name + " ja esta instalado, pulando")
			continue
		}

		for _, m := range mcp.AvailableMCPs {
			if m.Name != name {
				continue
			}

			if len(m.EnvVars) > 0 {
				for k, v := range m.EnvVars {
					if v != "" {
						continue
					}
					envVal := os.Getenv(k)
					if envVal != "" {
						m.EnvVars[k] = envVal
						tui.LogInfo(fmt.Sprintf("%s carregado do ambiente", k))
						continue
					}
					var input string
					field := huh.NewInput().
						Title(fmt.Sprintf("Valor para %s:", k)).
						Value(&input)

					if err := huh.NewForm(huh.NewGroup(field)).Run(); err != nil {
						continue
					}
					m.EnvVars[k] = input
				}
			}

			tui.LogStep("Instalando " + name + "...")
			if err := mcp.Install(m); err != nil {
				tui.LogError("Falha ao instalar " + name + ": " + err.Error())
			} else {
				tui.LogSuccess(name + " instalado com sucesso")
			}
			break
		}
	}
}

func runMCPRemove() {
	tui.LogStep("Consultando MCPs instalados...")
	installed, err := mcp.ListInstalled()
	if err != nil {
		tui.LogError(err.Error())
		return
	}

	if len(installed) == 0 {
		tui.LogWarn("Nenhum MCP instalado para remover")
		return
	}

	var options []huh.Option[string]
	options = append(options, huh.NewOption("<< Voltar", "back"))
	for _, name := range installed {
		options = append(options, huh.NewOption(name, name))
	}

	var selected string
	form := huh.NewSelect[string]().
		Title("Selecione o MCP para remover:").
		Options(options...).
		Value(&selected)

	if err := huh.NewForm(huh.NewGroup(form)).Run(); err != nil || selected == "back" {
		return
	}

	tui.LogStep("Removendo " + selected + "...")
	if err := mcp.Remove(selected); err != nil {
		tui.LogError("Falha ao remover " + selected + ": " + err.Error())
	} else {
		tui.LogSuccess(selected + " removido com sucesso")
	}
}
