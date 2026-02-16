package llm

import "context"

type Request struct {
	System string
	User   string
}

type Client interface {
	Generate(ctx context.Context, req Request) (string, error)
}
