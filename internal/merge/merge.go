package merge

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/wgalleti/wclaude-setup/internal/config"
)

type message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type request struct {
	Model     string    `json:"model"`
	MaxTokens int       `json:"max_tokens"`
	Messages  []message `json:"messages"`
	System    string    `json:"system"`
}

type contentBlock struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

type response struct {
	Content []contentBlock `json:"content"`
	Error   *struct {
		Message string `json:"message"`
	} `json:"error"`
}

const systemPrompt = `Voce e um especialista em configuracao do Claude Code.
Sua tarefa e fazer merge de dois arquivos CLAUDE.md:
1. O arquivo EXISTENTE do projeto (que tem contexto especifico do projeto)
2. O arquivo TEMPLATE (que tem as melhores praticas e convencoes)

Regras:
- PRESERVE todo o conteudo especifico do projeto existente (nomes, dominio, models, gotchas, etc)
- ADICIONE secoes do template que nao existem no arquivo atual
- ATUALIZE secoes existentes apenas se o template tiver informacao mais completa
- NUNCA remova informacao especifica do projeto
- Mantenha o formato Markdown limpo e organizado
- Responda APENAS com o conteudo do novo CLAUDE.md, sem explicacoes`

func MergeFiles(existing, template string, cfg *config.Config) (string, error) {
	if cfg.AnthropicAPIKey == "" {
		return "", fmt.Errorf("ANTHROPIC_API_KEY nao configurada. Use: claude-setup config --api-key <key>")
	}

	userMsg := fmt.Sprintf("## ARQUIVO EXISTENTE DO PROJETO:\n\n%s\n\n---\n\n## TEMPLATE COM MELHORES PRATICAS:\n\n%s", existing, template)

	reqBody := request{
		Model:     cfg.DefaultModel,
		MaxTokens: 8192,
		System:    systemPrompt,
		Messages: []message{
			{Role: "user", Content: userMsg},
		},
	}

	body, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("erro ao serializar request: %w", err)
	}

	req, err := http.NewRequest("POST", "https://api.anthropic.com/v1/messages", bytes.NewReader(body))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", cfg.AnthropicAPIKey)
	req.Header.Set("anthropic-version", "2023-06-01")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("erro na chamada API: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var apiResp response
	if err := json.Unmarshal(respBody, &apiResp); err != nil {
		return "", fmt.Errorf("erro ao parsear resposta: %w", err)
	}

	if apiResp.Error != nil {
		return "", fmt.Errorf("API error: %s", apiResp.Error.Message)
	}

	if len(apiResp.Content) == 0 {
		return "", fmt.Errorf("resposta vazia da API")
	}

	return apiResp.Content[0].Text, nil
}
