package setup

import (
	"os"
	"path/filepath"
)

type Stack string

const (
	StackDjango  Stack = "django"
	StackVue     Stack = "vue"
	StackFlutter Stack = "flutter"
	StackGo      Stack = "go"
	StackUnknown Stack = "unknown"
)

func DetectStack(dir string) Stack {
	checks := []struct {
		file  string
		stack Stack
	}{
		{"manage.py", StackDjango},
		{"pubspec.yaml", StackFlutter},
		{"package.json", StackVue}, // will refine below
		{"go.mod", StackGo},
	}

	for _, c := range checks {
		if fileExists(filepath.Join(dir, c.file)) {
			if c.stack == StackVue {
				// Distinguish Vue from other JS projects
				if fileExists(filepath.Join(dir, "vite.config.ts")) ||
					fileExists(filepath.Join(dir, "vite.config.js")) ||
					dirExists(filepath.Join(dir, "src", "composables")) ||
					dirExists(filepath.Join(dir, "src", "stores")) {
					return StackVue
				}
			}
			return c.stack
		}
	}
	return StackUnknown
}

func (s Stack) TemplateName() string {
	switch s {
	case StackDjango:
		return "django.md"
	case StackVue:
		return "vue.md"
	case StackFlutter:
		return "flutter.md"
	case StackGo:
		return "go.md"
	default:
		return ""
	}
}

func (s Stack) Label() string {
	switch s {
	case StackDjango:
		return "Django / DRF"
	case StackVue:
		return "Vue.js / PrimeVue"
	case StackFlutter:
		return "Flutter / Dart"
	case StackGo:
		return "Go"
	default:
		return "Desconhecido"
	}
}

func fileExists(path string) bool {
	info, err := os.Stat(path)
	return err == nil && !info.IsDir()
}

func dirExists(path string) bool {
	info, err := os.Stat(path)
	return err == nil && info.IsDir()
}
