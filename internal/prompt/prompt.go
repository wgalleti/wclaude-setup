package prompt

import "fmt"

type PromptType string

const (
	PromptFeature   PromptType = "feature"
	PromptBugfix    PromptType = "bugfix"
	PromptRefactor  PromptType = "refactor"
	PromptReview    PromptType = "review"
	PromptWidget    PromptType = "widget"
	PromptHandler   PromptType = "handler"
	PromptComponent PromptType = "component"
	PromptPage      PromptType = "page"
)

type PromptData struct {
	Type PromptType

	// Feature fields
	FeatureName string
	Model       string
	Fields      string
	Actions     string
	Permission  string

	// Bugfix fields
	BugDescription string
	FilePath       string
	Expected       string
	Actual         string

	// Refactor fields
	RefactorTarget string
	Objective      string
	Preserve       string

	// Review fields
	ReviewTarget string
	FocusAreas   string

	// Widget fields (Flutter)
	WidgetName    string
	WidgetPurpose string
	WidgetProps   string
	WidgetState   string

	// Handler fields (Go)
	HTTPMethod   string
	Endpoint     string
	InputSchema  string
	OutputSchema string
	Errors       string

	// Component fields (React/Next)
	ComponentName    string
	ComponentPurpose string
	ComponentProps   string
	ComponentState   string
	ComponentLibs    string

	// Page fields (Next/Nuxt)
	PageRoute       string
	PageDataSource  string
	PageAuth        string
	PageComponents  string
}

func Generate(d PromptData) string {
	switch d.Type {
	case PromptFeature:
		return fmt.Sprintf(`Feature: %s
Model: %s
Precisa:
- serializer com campos %s
- viewset com actions %s
- permissao: %s
- teste de integracao basico`, d.FeatureName, d.Model, d.Fields, d.Actions, d.Permission)

	case PromptBugfix:
		return fmt.Sprintf(`Bug: %s
Arquivo: %s
Esperado: %s
Atual: %s`, d.BugDescription, d.FilePath, d.Expected, d.Actual)

	case PromptRefactor:
		return fmt.Sprintf(`Refatorar: %s
Objetivo: %s
Manter: %s`, d.RefactorTarget, d.Objective, d.Preserve)

	case PromptReview:
		return fmt.Sprintf(`Review: %s
Focar em: %s
Ignorar: estilo, formatacao (ja tem linter)`, d.ReviewTarget, d.FocusAreas)

	case PromptWidget:
		return fmt.Sprintf(`Widget: %s
Proposito: %s
Props: %s
State: %s`, d.WidgetName, d.WidgetPurpose, d.WidgetProps, d.WidgetState)

	case PromptHandler:
		return fmt.Sprintf(`Handler: %s %s
Input: %s
Output: %s
Erros: %s`, d.HTTPMethod, d.Endpoint, d.InputSchema, d.OutputSchema, d.Errors)

	case PromptComponent:
		return fmt.Sprintf(`Component: %s
Proposito: %s
Props: %s
State: %s
Libs: %s`, d.ComponentName, d.ComponentPurpose, d.ComponentProps, d.ComponentState, d.ComponentLibs)

	case PromptPage:
		return fmt.Sprintf(`Page: %s
Data source: %s
Auth: %s
Componentes: %s`, d.PageRoute, d.PageDataSource, d.PageAuth, d.PageComponents)

	default:
		return ""
	}
}
