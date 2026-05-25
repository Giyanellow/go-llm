package llm

type LLMClient interface {
	Run(prompt string) ([]Message, error)
	GetModel() string
}
