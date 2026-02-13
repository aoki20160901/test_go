package llm

import "context"

type OpusClient struct {
	apiKey string
}

func NewOpusClient(apiKey string) *OpusClient {
	return &OpusClient{
		apiKey: apiKey,
	}
}

func (c *OpusClient) RefineText(
	ctx context.Context,
	input string,
) (string, error) {

	return "整形済み: " + input, nil
}
