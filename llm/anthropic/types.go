package anthropic

type anthropicRequestBody struct {
	Model     string             `json:"model"`
	MaxTokens int                `json:"max_tokens,omitempty"`
	Input     []anthropicMessage `json:"messages"`
}

type anthropicResponseBody struct {
	Model       string                 `json:"model"`
	ID          string                 `json:"id"`
	MessageType string                 `json:"type"`
	Role        string                 `json:"role"`
	Content     []anthropicContentBody `json:"content"`
}

type anthropicContentBody struct {
	MessageType string `json:"type"`
	Text        string `json:"text"`
}

type anthropicMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}
