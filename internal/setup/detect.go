package setup

import (
	"os"
	"path/filepath"
)

type Stack string

const (
	StackDjango   Stack = "django"
	StackVue      Stack = "vue"
	StackNuxt     Stack = "nuxt"
	StackReact    Stack = "react"
	StackNext     Stack = "next"
	StackFlutter  Stack = "flutter"
	StackGo       Stack = "go"
	StackSupabase Stack = "supabase"
	StackUnknown  Stack = "unknown"
)

func DetectStack(dir string) Stack {
	// Non-JS stacks first (unambiguous)
	if fileExists(filepath.Join(dir, "manage.py")) {
		return StackDjango
	}
	if fileExists(filepath.Join(dir, "pubspec.yaml")) {
		return StackFlutter
	}
	if fileExists(filepath.Join(dir, "go.mod")) {
		return StackGo
	}

	// Supabase: has supabase config
	if fileExists(filepath.Join(dir, "supabase", "config.toml")) {
		// Supabase can coexist with frontend — detect frontend first, fallback to supabase
		if frontend := detectJSStack(dir); frontend != StackUnknown {
			return frontend
		}
		return StackSupabase
	}

	// JS/TS stacks: need deeper inspection
	if fileExists(filepath.Join(dir, "package.json")) {
		return detectJSStack(dir)
	}

	return StackUnknown
}

func detectJSStack(dir string) Stack {
	// Nuxt: nuxt.config.ts or .nuxt directory
	if fileExists(filepath.Join(dir, "nuxt.config.ts")) ||
		fileExists(filepath.Join(dir, "nuxt.config.js")) ||
		dirExists(filepath.Join(dir, ".nuxt")) {
		return StackNuxt
	}

	// Next.js: next.config.* or .next directory
	if fileExists(filepath.Join(dir, "next.config.ts")) ||
		fileExists(filepath.Join(dir, "next.config.js")) ||
		fileExists(filepath.Join(dir, "next.config.mjs")) ||
		dirExists(filepath.Join(dir, ".next")) ||
		dirExists(filepath.Join(dir, "app")) && fileExists(filepath.Join(dir, "app", "layout.tsx")) {
		return StackNext
	}

	// Vue: vite.config + Vue-specific dirs
	if fileExists(filepath.Join(dir, "vite.config.ts")) ||
		fileExists(filepath.Join(dir, "vite.config.js")) {
		if dirExists(filepath.Join(dir, "src", "composables")) ||
			dirExists(filepath.Join(dir, "src", "stores")) ||
			fileExists(filepath.Join(dir, "src", "App.vue")) {
			return StackVue
		}
	}

	// React: src/App.tsx or src/main.tsx (Vite React)
	if fileExists(filepath.Join(dir, "src", "App.tsx")) ||
		fileExists(filepath.Join(dir, "src", "App.jsx")) ||
		fileExists(filepath.Join(dir, "src", "main.tsx")) ||
		fileExists(filepath.Join(dir, "src", "index.tsx")) {
		return StackReact
	}

	return StackUnknown
}

func (s Stack) TemplateName() string {
	switch s {
	case StackDjango:
		return "django.md"
	case StackVue:
		return "vue.md"
	case StackNuxt:
		return "nuxt.md"
	case StackReact:
		return "react.md"
	case StackNext:
		return "next.md"
	case StackFlutter:
		return "flutter.md"
	case StackGo:
		return "go.md"
	case StackSupabase:
		return "supabase.md"
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
	case StackNuxt:
		return "Nuxt / Vue.js"
	case StackReact:
		return "React / Vite"
	case StackNext:
		return "Next.js / React"
	case StackFlutter:
		return "Flutter / Dart"
	case StackGo:
		return "Go"
	case StackSupabase:
		return "Supabase"
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
