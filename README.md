# wclaude-setup

CLI interativa para configurar o Claude Code em projetos. Detecta a stack do projeto, gerencia MCPs, gera prompts eficientes e faz merge inteligente de CLAUDE.md via API Anthropic.

## Instalacao

```bash
make install
```

O binario e instalado em `~/.local/bin/wclaude-setup`. Certifique-se que esse diretorio esta no `PATH`.

## Uso

```bash
# Menu interativo
wclaude-setup

# Comandos diretos
wclaude-setup init          # Setup global (~/.claude)
wclaude-setup project       # Setup CLAUDE.md no projeto atual
wclaude-setup mcp           # Gerenciar MCPs
wclaude-setup mcp list      # Listar MCPs instalados
wclaude-setup prompt        # Gerar prompt eficiente
wclaude-setup merge         # Merge CLAUDE.md via API Anthropic
wclaude-setup config set-key  # Configurar API key
wclaude-setup config show     # Ver configuracao
```

## Estrutura

```
cmd/wclaude-setup/       Entrypoint e comandos CLI
internal/
  config/                Configuracao persistente (~/.wclaude-setup/)
  mcp/                   Gerenciamento de MCP servers
  merge/                 Merge via API Anthropic
  prompt/                Gerador de prompts eficientes
  setup/                 Deteccao de stack e instalacao
  templates/             Templates CLAUDE.md embeddados
  tui/                   Estilos do terminal
docs/                    Documentacao de referencia
```

## Stacks suportadas

| Stack | Deteccao automatica |
|-------|-------------------|
| Django / DRF | `manage.py` |
| Vue.js / PrimeVue | `package.json` + `vite.config.*` |
| Flutter / Dart | `pubspec.yaml` |
| Go | `go.mod` |

## Requisitos

- Go 1.22+
- Claude Code CLI (`claude`) no PATH
- Node.js/npx (para MCPs baseados em npm)
