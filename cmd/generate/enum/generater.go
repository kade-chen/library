package enum

import (
	"text/template"
)

// G Generater
var G = Generater{
	t: template.New("enum"),
}

// Generater 用于生成枚举的生成器
type Generater struct {
	t           *template.Template
	file        string
	Marshal     bool
	ProtobufExt bool
}

// Generate 生成文件
func (g *Generater) Generate(file string) ([]byte, error) {
	g.file = file
	RenderParams, err := g.parse()
	if err != nil {
		return nil, err
	}
	RenderParams.Marshal = g.Marshal
	if g.ProtobufExt {
		RenderParams.Stringer = false
		RenderParams.ValueMap = false
	}

	if RenderParams.Enums.Length() == 0 {
		return []byte{}, nil
	}

	return g.gen(RenderParams)
}
