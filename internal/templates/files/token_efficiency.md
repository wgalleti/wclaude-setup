# Guia de Eficiência de Tokens — Claude Code

## Princípio Central

> Contexto preciso → menos tokens gastos → respostas melhores

Claude Code lê o contexto inteiro a cada mensagem. Quanto mais limpo e relevante for o contexto, menos tokens ele consome para chegar na resposta certa.

---

## Comandos Essenciais de Contexto

```bash
/clear          # Limpa histórico — USE entre tarefas não relacionadas
/compact        # Resume o histórico — USE em sessões longas antes de nova tarefa
/cost           # Mostra tokens gastos na sessão atual
/memory         # Gerencia o que o Claude "lembra" entre sessões
```

**Regra prática:**
- Nova feature? → `/clear` primeiro
- Refatorando algo diferente? → `/clear`  
- Sessão virou um monstro de 50 mensagens? → `/compact` antes de continuar

---

## Padrões de Prompt Eficientes

### ❌ Prompt ruim (desperdiça tokens)
```
Oi Claude, tudo bem? Preciso da sua ajuda com uma coisa. 
Estou desenvolvendo uma API em Django e queria que você 
me ajudasse a criar um serializer para o model User que 
eu tenho no meu projeto. Pode me ajudar com isso?
```

### ✅ Prompt bom
```
Cria serializer para apps/users/models.py:User
Campos: id, email, full_name, role, created_at
read_only: id, created_at
Validação: email único na criação
```

---

## Templates de Prompt por Tarefa

### Nova Feature (Django)
```
Feature: [nome]
Model: apps/[app]/models.py:[Model]
Precisa:
- serializer com campos [x, y, z]
- viewset com actions [list, create, retrieve]  
- permissão: [IsAuthenticated | IsOwner | IsAdmin]
- teste de integração básico
```

### Bug Fix
```
Bug: [descrição do comportamento errado]
Arquivo: [path:linha]
Esperado: [o que deveria acontecer]
Atual: [o que está acontecendo]
[cole apenas o trecho relevante do código]
```

### Refactor
```
Refatorar: [arquivo:função]
Objetivo: [ex: extrair lógica de cálculo para service]
Manter: [o que não pode mudar — ex: interface pública]
[cole apenas o código alvo]
```

### Code Review
```
Review: [arquivo ou PR]
Focar em: [performance | segurança | padrões do projeto | tudo]
Ignorar: [estilo, formatação — já tem linter]
```

### Flutter Widget
```
Widget: [NomeWidget]
Propósito: [o que exibe/faz]
Props: [lista de parâmetros]
State: [stateless | riverpod provider | local]
Design ref: [descrição ou componente PrimeVue equivalente]
```

### Go Handler
```
Handler: [Método] /api/[resource]
Input: [body schema ou query params]
Output: [response schema]
Erros: [404 se não encontrar, 403 se não autorizado]
Service já existe: [sim/não]
```

---

## Técnicas de Economia de Tokens

### 1. Referencie, não cole
```
# Em vez de colar o model inteiro:
"Baseado em apps/users/models.py:User, cria..."

# Claude lê o arquivo diretamente — não precisa do conteúdo no chat
```

### 2. Scope preciso
```
# Ruim — Claude vai ler tudo
"Revisa o projeto e me diz o que melhorar"

# Bom — Claude foca no que importa
"Revisa apps/billing/services.py:StripeService para thread safety"
```

### 3. Separe tarefas grandes
```
# Ruim — tudo de uma vez
"Cria o model, serializer, view, urls, testes e documentação da API de produtos"

# Bom — incremental
1. "Cria model Product em apps/products/models.py"
2. [review] → "Cria serializer baseado em apps/products/models.py:Product"
3. [review] → "Cria ViewSet para Product com ações list e create"
4. [review] → "Testes de integração para GET /api/products/"
```

### 4. Use o CLAUDE.md do projeto
O `CLAUDE.md` no root do projeto é lido automaticamente. Quanto mais contexto útil ele tiver, menos você precisa repetir em cada prompt.

### 5. Slash commands para evitar retrabalho
```bash
/add [arquivo]     # Adiciona arquivo ao contexto sem colar no chat
/remove [arquivo]  # Remove do contexto o que não é mais relevante
```

---

## Configuração de Permissões (reduz interrupções)

Edite `~/.claude/settings.json` para auto-aprovar operações rotineiras:

```json
{
  "permissions": {
    "allow": [
      "Bash(git diff:*)",
      "Bash(git log:*)",
      "Bash(git status:*)",
      "Bash(python manage.py *)",
      "Bash(uv run python manage.py *)",
      "Bash(pytest *)",
      "Bash(uv run pytest *)",
      "Bash(dart analyze)",
      "Bash(flutter analyze)",
      "Bash(go test *)",
      "Bash(go vet *)"
    ]
  }
}
```

Isso evita que Claude pergunte confirmação em cada `git status` ou `pytest`.

---

## Fluxo de Trabalho Recomendado

```
1. Abrir projeto → Claude Code lê CLAUDE.md automaticamente
2. /clear se vindo de outra sessão
3. Contexto do trabalho atual (1-2 frases)
4. Tarefa específica com template acima
5. Review incremental — não aprove código que não leu
6. /compact quando histórico ficar grande
7. /clear ao mudar de feature/contexto
```

---

## Métricas de Eficiência

Você está usando bem os tokens quando:
- [ ] Cada prompt tem menos de 150 palavras
- [ ] Claude não faz perguntas de clarificação (contexto estava claro)
- [ ] A primeira resposta já está ~80% certa
- [ ] Você usa `/clear` pelo menos uma vez por hora de trabalho
- [ ] O `CLAUDE.md` do projeto está atualizado com as convenções atuais
