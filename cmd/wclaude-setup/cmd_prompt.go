package main

import (
	"fmt"

	"github.com/charmbracelet/huh"
	"github.com/spf13/cobra"

	"github.com/wgalleti/wclaude-setup/internal/prompt"
	"github.com/wgalleti/wclaude-setup/internal/tui"
)

func newPromptCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "prompt",
		Short: "Gerar prompt eficiente interativamente",
		Run: func(cmd *cobra.Command, args []string) {
			runPrompt()
		},
	}
}

func runPrompt() {
	for {
		tui.PrintHeader("Gerador de Prompts")

		var promptType string
		form := huh.NewSelect[string]().
			Title("Tipo de tarefa:").
			Options(
				huh.NewOption("Nova Feature (Django)", "feature"),
				huh.NewOption("Bug Fix", "bugfix"),
				huh.NewOption("Refactor", "refactor"),
				huh.NewOption("Code Review", "review"),
				huh.NewOption("React/Next Component", "component"),
				huh.NewOption("Page (Next/Nuxt)", "page"),
				huh.NewOption("Flutter Widget", "widget"),
				huh.NewOption("Go Handler", "handler"),
				huh.NewOption("<< Voltar", "back"),
			).
			Value(&promptType)

		if err := huh.NewForm(huh.NewGroup(form)).Run(); err != nil || promptType == "back" {
			return
		}

		var data prompt.PromptData
		data.Type = prompt.PromptType(promptType)

		cancelled := false

		switch data.Type {
		case prompt.PromptFeature:
			if err := huh.NewForm(huh.NewGroup(
				huh.NewInput().Title("Nome da feature:").Value(&data.FeatureName),
				huh.NewInput().Title("Model (ex: apps/users/models.py:User):").Value(&data.Model),
				huh.NewInput().Title("Campos (ex: id, email, name):").Value(&data.Fields),
				huh.NewInput().Title("Actions (ex: list, create, retrieve):").Value(&data.Actions),
				huh.NewInput().Title("Permissao (ex: IsAuthenticated):").Value(&data.Permission),
			)).Run(); err != nil {
				cancelled = true
			}

		case prompt.PromptBugfix:
			if err := huh.NewForm(huh.NewGroup(
				huh.NewInput().Title("Descricao do bug:").Value(&data.BugDescription),
				huh.NewInput().Title("Arquivo:linha:").Value(&data.FilePath),
				huh.NewInput().Title("Comportamento esperado:").Value(&data.Expected),
				huh.NewInput().Title("Comportamento atual:").Value(&data.Actual),
			)).Run(); err != nil {
				cancelled = true
			}

		case prompt.PromptRefactor:
			if err := huh.NewForm(huh.NewGroup(
				huh.NewInput().Title("Alvo (arquivo:funcao):").Value(&data.RefactorTarget),
				huh.NewInput().Title("Objetivo:").Value(&data.Objective),
				huh.NewInput().Title("O que manter (interface publica, etc):").Value(&data.Preserve),
			)).Run(); err != nil {
				cancelled = true
			}

		case prompt.PromptReview:
			if err := huh.NewForm(huh.NewGroup(
				huh.NewInput().Title("Arquivo ou PR:").Value(&data.ReviewTarget),
				huh.NewInput().Title("Focar em (performance, seguranca, padroes):").Value(&data.FocusAreas),
			)).Run(); err != nil {
				cancelled = true
			}

		case prompt.PromptWidget:
			if err := huh.NewForm(huh.NewGroup(
				huh.NewInput().Title("Nome do widget:").Value(&data.WidgetName),
				huh.NewInput().Title("Proposito:").Value(&data.WidgetPurpose),
				huh.NewInput().Title("Props:").Value(&data.WidgetProps),
				huh.NewInput().Title("State (stateless, riverpod, local):").Value(&data.WidgetState),
			)).Run(); err != nil {
				cancelled = true
			}

		case prompt.PromptHandler:
			if err := huh.NewForm(huh.NewGroup(
				huh.NewInput().Title("Metodo HTTP (GET, POST, etc):").Value(&data.HTTPMethod),
				huh.NewInput().Title("Endpoint (ex: /api/users):").Value(&data.Endpoint),
				huh.NewInput().Title("Input schema:").Value(&data.InputSchema),
				huh.NewInput().Title("Output schema:").Value(&data.OutputSchema),
				huh.NewInput().Title("Erros esperados:").Value(&data.Errors),
			)).Run(); err != nil {
				cancelled = true
			}

		case prompt.PromptComponent:
			if err := huh.NewForm(huh.NewGroup(
				huh.NewInput().Title("Nome do componente:").Value(&data.ComponentName),
				huh.NewInput().Title("Proposito:").Value(&data.ComponentPurpose),
				huh.NewInput().Title("Props (ex: user, onEdit, isLoading):").Value(&data.ComponentProps),
				huh.NewInput().Title("State (local, zustand, tanstack-query):").Value(&data.ComponentState),
				huh.NewInput().Title("Libs (ex: framer-motion, tanstack-table):").Value(&data.ComponentLibs),
			)).Run(); err != nil {
				cancelled = true
			}

		case prompt.PromptPage:
			if err := huh.NewForm(huh.NewGroup(
				huh.NewInput().Title("Rota (ex: /users/[id]):").Value(&data.PageRoute),
				huh.NewInput().Title("Data source (server fetch, useFetch, useQuery):").Value(&data.PageDataSource),
				huh.NewInput().Title("Auth (public, authenticated, admin):").Value(&data.PageAuth),
				huh.NewInput().Title("Componentes usados:").Value(&data.PageComponents),
			)).Run(); err != nil {
				cancelled = true
			}
		}

		if cancelled {
			tui.LogWarn("Cancelado")
			continue
		}

		result := prompt.Generate(data)

		tui.LogSuccess("Prompt gerado:")
		tui.PrintSeparator()
		fmt.Println()
		fmt.Println(result)
		fmt.Println()
		tui.PrintSeparator()

		tui.WaitForEnter()
	}
}
