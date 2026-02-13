package pdf

import (
	"bytes"
)

type Generator interface {
	Generate(text string, image []byte) ([]byte, error)
}

type TemplateGenerator struct {
	templatePath string
}

func NewTemplateGenerator(path string) *TemplateGenerator {
	return &TemplateGenerator{
		templatePath: path,
	}
}

func (g *TemplateGenerator) Generate(
	text string,
	image []byte,
) ([]byte, error) {

	// 仮実装
	var buf bytes.Buffer
	buf.WriteString("PDF生成\n")
	buf.WriteString(text)

	return buf.Bytes(), nil
}
