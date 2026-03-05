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
	var action string
	form := huh.NewSelect[string]().
		Title("MCPs:").
		Options(
			huh.NewOption("Instalar MCPs", "install"),
			huh.NewOption("Remover MCP", "remove"),
			huh.NewOption("Listar instalados", "list"),
		).
		Value(&action)

	if err := huh.NewForm(huh.NewGroup(form)).Run(); err != nil {
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
}

func runMCPList() {
	installed, err := mcp.ListInstalled()
	if err != nil {
		fmt.Println(tui.Error.Render(err.Error()))
		return
	}
	if len(installed) == 0 {
		fmt.Println(tui.Info.Render("Nenhum MCP instalado"))
		return
	}
	fmt.Println(tui.Title.Render("MCPs instalados:"))
	for _, name := range installed {
		fmt.Printf("  - %s\n", name)
	}
}

func runMCPInstall() {
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

	for _, name := range selected {
		if mcp.IsInstalled(name, installed) {
			fmt.Printf("  %s ja instalado\n", tui.Info.Render(name))
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

			fmt.Printf("  Instalando %s...\n", name)
			if err := mcp.Install(m); err != nil {
				fmt.Println(tui.Error.Render("  Erro: " + err.Error()))
			} else {
				fmt.Println(tui.Success.Render("  " + name + " instalado"))
			}
			break
		}
	}
}

func runMCPRemove() {
	installed, err := mcp.ListInstalled()
	if err != nil {
		fmt.Println(tui.Error.Render(err.Error()))
		return
	}

	if len(installed) == 0 {
		fmt.Println(tui.Info.Render("Nenhum MCP instalado"))
		return
	}

	var options []huh.Option[string]
	for _, name := range installed {
		options = append(options, huh.NewOption(name, name))
	}

	var selected string
	form := huh.NewSelect[string]().
		Title("Selecione o MCP para remover:").
		Options(options...).
		Value(&selected)

	if err := huh.NewForm(huh.NewGroup(form)).Run(); err != nil {
		return
	}

	if err := mcp.Remove(selected); err != nil {
		fmt.Println(tui.Error.Render("Erro: " + err.Error()))
	} else {
		fmt.Println(tui.Success.Render(selected + " removido"))
	}
}
