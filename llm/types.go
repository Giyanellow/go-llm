package llm

type Message struct {
	Role    string `json:"role"`
	Content string `json:"message"`
	Type    string `json:"type"`
}
