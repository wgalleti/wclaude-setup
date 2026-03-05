# MCP Setup — Claude Code
# wGalleti Dev Environment (macOS Tahoe, Apple Silicon)

## MCPs Instalados e Configurados

### 1. Dart & Flutter MCP (oficial)
**O que faz:** Analisa código, roda analyzer, gerencia pub.dev packages, hot reload, introspect widget tree

```bash
# Instalar (já configurado se Dart 3.9+ no PATH)
claude mcp add --transport stdio --scope user dart-flutter \
  -- dart mcp-server --force-roots-fallback
```

**Usar Flutter com MCP:**
```bash
# Sempre rodar o app assim para conectar ao MCP
flutter run -d macos --debug \
  --host-vmservice-port=9100 \
  --enable-vm-service \
  --disable-service-auth-codes
```

---

### 2. Context7 — Documentação up-to-date
**O que faz:** Injeta docs atuais de qualquer biblioteca no contexto. Resolve o problema de Claude sugerir APIs deprecadas.

```bash
claude mcp add --transport stdio --scope user context7 \
  -- npx -y @upstash/context7-mcp@latest
```

**Como usar nos prompts:**
```
"Como implementar paginação com DRF? use context7"
"Qual a sintaxe atual de Riverpod AsyncNotifier? use context7"
```

---

### 3. GitHub MCP (oficial)
**O que faz:** Lê/cria issues, PRs, branches, reviews direto do Claude Code

```bash
# Precisa de token: github.com/settings/tokens
# Permissions: repo, read:org, read:user
export GITHUB_TOKEN=ghp_xxxx

claude mcp add --transport stdio --scope user github \
  -e GITHUB_TOKEN=$GITHUB_TOKEN \
  -- npx -y @modelcontextprotocol/server-github
```

---

### 4. PostgreSQL MCP
**O que faz:** Queries diretas ao banco, explorar schema, debug de dados

```bash
# Por projeto (não global — cada projeto tem seu banco)
claude mcp add --transport stdio --scope project postgres \
  -- npx -y @modelcontextprotocol/server-postgres \
  postgresql://user:pass@localhost:5432/dbname
```

---

### 5. Playwright MCP — Browser automation
**O que faz:** Testa fluxos de UI, scraping, automação de browser

```bash
claude mcp add --transport stdio --scope user playwright \
  -- npx -y @playwright/mcp@latest
```

---

## Verificar Status

```bash
# Listar todos os MCPs
claude mcp list

# Dentro do Claude Code
/mcp
```

Output esperado:
```
dart-flutter: dart mcp-server --force-roots-fallback    ✓ Connected
context7:     npx @upstash/context7-mcp                 ✓ Connected  
github:       npx @modelcontextprotocol/server-github   ✓ Connected
playwright:   npx @playwright/mcp                       ✓ Connected
```

---

## Troubleshooting

### "spawn stdio ENOENT"
O Claude Code não encontra o executável. Use path absoluto:
```bash
# Verificar path do dart
which dart
# Ex: /Users/wgalleti/fvm/default/bin/dart

claude mcp add --transport stdio --scope user dart-flutter \
  -- /Users/wgalleti/fvm/default/bin/dart mcp-server --force-roots-fallback
```

### MCP configurado mas não responde
1. Reiniciar Claude Code completamente (não só /clear)
2. Verificar logs:
```bash
ls ~/Library/Caches/claude-cli-nodejs/*/mcp-logs-*/
tail -f ~/Library/Caches/claude-cli-nodejs/*/mcp-logs-dart-flutter/*.log
```

### Context7 não encontra pacote
Especificar o pacote explicitamente no prompt:
```
"Como usar django-filter 24.x? use context7 com libraryId django-filter"
```

---

## MCPs a Adicionar por Projeto

| Projeto    | MCPs adicionais sugeridos              |
|------------|----------------------------------------|
| workHard   | postgres (local), github               |
| DojoHub    | postgres (local), github               |
| Flutter    | dart-flutter (obrigatório)             |
| Go service | github, postgres                       |
