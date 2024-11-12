package enum

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"os"
	"regexp"
	"strings"

	"github.com/pkg/errors"
)

// RenderParams 模板渲染需要的参数
type RenderParams struct {
	PKG       string
	Backquote string
	Enums     *Set
	Marshal   bool
	Stringer  bool
	ValueMap  bool
}

// NewRenderParams 创建新的模板参数
func NewRenderParams() *RenderParams {
	return &RenderParams{
		Enums:     NewEnumSet(),
		Stringer:  true,
		ValueMap:  true,
		Backquote: "`",
	}
}

// Set 枚举集合
type Set struct {
	Items []*Enum
}

// NewEnumSet todo
func NewEnumSet() *Set {
	return &Set{
		Items: []*Enum{},
	}
}

// Length 长度
func (s *Set) Length() int {
	return len(s.Items)
}

// Get 获取一个枚举类
func (s *Set) Get(name string) *Enum {
	for _, e := range s.Items {
		if e.Name == name {
			return e
		}
	}

	enum := NewEnum(name)
	s.Add(enum)
	return enum
}

// Add 添加一个类型
func (s *Set) Add(i *Enum) {
	s.Items = append(s.Items, i)
}

// Enum 枚举类型
type Enum struct {
	Name  string
	Doc   string
	Items []*Item
}

// NewEnum todo
func NewEnum(name string) *Enum {
	return &Enum{
		Name:  name,
		Items: []*Item{},
	}
}

// Add todo
func (e *Enum) Add(i *Item) {
	e.Items = append(e.Items, i)
}

// Item 枚举项
// 通过获取注释里()中的内容作为自定义参数
type Item struct {
	Name string // 枚举的名称, 对应常量的名称
	Doc  string // 文档, 常量的文档
}

// NewItem todo
func NewItem(name, doc string) *Item {
	return &Item{
		Name: name,
		Doc:  doc,
	}
}

var (
	itemParamRe = regexp.MustCompile(`.*?\((.*?)\).*`)
)

// Show 枚举项显示
func (i *Item) Show() string {
	matchs := itemParamRe.FindStringSubmatch(i.Doc)
	if len(matchs) < 2 {
		return i.defaultShow()
	}

	return matchs[1]
}

func (i *Item) defaultShow() string {
	return strings.ToLower(i.Name)
}

// 解析代码源文件，获取常量和类型
func (g *Generater) parse() (*RenderParams, error) {
	fset := token.NewFileSet() // 用于记录位置信息
	//create a AST logical tree
	f, err := parser.ParseFile(fset, g.getFile(), nil, parser.ParseComments)
	if err != nil {
		return nil, fmt.Errorf("parse file error, %s", err)
	}

	RenderParams := NewRenderParams()
	RenderParams.PKG = f.Name.Name
	// 遍历AST节点
	for _, decl := range f.Decls {
		switch d := decl.(type) {
		case *ast.GenDecl:
			switch d.Tok {
			case token.CONST:
				for _, spec := range d.Specs {
					vs, _ := spec.(*ast.ValueSpec)

					ident := vs.Names[0] // 提取常量名称
					doc := vs.Doc.Text() // 提取注释文档
					var enum *Enum
					vst, _ := vs.Type.(*ast.Ident)
					if vst != nil {
						enum = RenderParams.Enums.Get(vst.Name)
						enum.Add(NewItem(ident.Name, doc))
					}

				}
			}
		}
	}
	// ast.Print(fset, f) 	// 打印语法树
	return RenderParams, nil
}

func (g *Generater) getFile() string {
	if g.file != "" {
		return g.file
	}
	return os.Getenv("GOFILE")
}

func (g *Generater) gen(params *RenderParams) ([]byte, error) {
	buf := bytes.NewBufferString("")
	t, err := g.t.Parse(tmp)
	if err != nil {
		return nil, errors.Wrapf(err, "template init err")
	}

	err = t.Execute(buf, params)
	if err != nil {
		return nil, errors.Wrapf(err, "template data err")
	}
	return format.Source(buf.Bytes())
}
