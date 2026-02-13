package llm

import "context"

type LLM interface {
	RefineText(ctx context.Context, input string) (string, error)
}
