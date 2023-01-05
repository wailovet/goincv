package goincv

import (
	"bytes"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
)

type Ast struct {
	file         string
	src          string
	doc          string
	types        string
	tkfs         *token.FileSet
	ast          interface{}
	CursorOffset int
	RawOffset    int
	StartPos     int
	EndPos       int
}

func NewAst(file, src string) *Ast {
	tokenfs := token.NewFileSet()
	ast := Ast{
		file: file,
		src:  src,
		tkfs: tokenfs,
	}
	return &ast
}

func (a *Ast) Parse() error {
	f, err := parser.ParseFile(a.tkfs, a.file, a.src, parser.ParseComments)
	if err != nil {
		return err
	}
	a.ast = f
	return nil
}

func (a *Ast) ReSynthesisSrc() string {
	dst := bytes.NewBuffer(nil)
	err := format.Node(dst, a.tkfs, a.ast)
	if err != nil {
		return ""
	}
	return dst.String()
}

func (a *Ast) ToGoSrcCursorBefore(cp int) string {
	return string([]rune(a.src)[:cp])
}

func (a *Ast) ToGoSrcCursorAfter(cp int) string {
	return string([]rune(a.src)[cp:])
}

func (a *Ast) Types() string {
	return a.types
}
func (a *Ast) Raw() interface{} {
	return a.ast
}

func (a *Ast) GetAllBlock() (ret []*Ast) {
	switch v := a.ast.(type) {
	case *ast.File:
		for _, i := range v.Decls {
			pos := int(i.Pos())
			start := len([]rune(string(a.src[:pos-1])))

			e := int(i.End())
			end := len([]rune(string(a.src[:e-1])))

			switch block := i.(type) {
			case *ast.FuncDecl:

				item := &Ast{
					file:      a.file + ":func",
					tkfs:      a.tkfs,
					ast:       block,
					types:     "func",
					doc:       block.Doc.Text(),
					RawOffset: int(i.Pos()),
					StartPos:  start,
					EndPos:    end,
					src:       string([]rune(a.src)[start:end]),
				}

				ret = append(ret, item)

			case *ast.GenDecl:
				item := &Ast{
					file:      a.file + ":" + block.Tok.String(),
					tkfs:      a.tkfs,
					ast:       block,
					types:     block.Tok.String(),
					doc:       block.Doc.Text(),
					RawOffset: int(i.Pos()),
					StartPos:  start,
					EndPos:    end,
					src:       string([]rune(a.src)[start:end]),
				}
				ret = append(ret, item)

			}

		}
		return ret
	}
	return
}

func (a *Ast) GetBlockFormOffset(offset int) (ret []*Ast) {
	switch v := a.ast.(type) {
	case *ast.File:
		for _, i := range v.Decls {
			pos := int(i.Pos())
			start := len([]rune(string(a.src[:pos-1])))

			e := int(i.End())
			end := len([]rune(string(a.src[:e-1])))

			if offset > int(start) && offset < int(end) {
				switch block := i.(type) {
				case *ast.FuncDecl:
					item := &Ast{
						file:         a.file + ":func",
						tkfs:         a.tkfs,
						ast:          block,
						types:        "func",
						doc:          block.Doc.Text(),
						CursorOffset: offset - start,
						RawOffset:    int(i.Pos()),
						StartPos:     start,
						EndPos:       end,
						src:          string([]rune(a.src)[start:end]),
					}
					ret = append(ret, item)

				case *ast.GenDecl:
					item := &Ast{
						file:         a.file + ":" + block.Tok.String(),
						tkfs:         a.tkfs,
						ast:          block,
						types:        block.Tok.String(),
						doc:          block.Doc.Text(),
						CursorOffset: offset - start,
						RawOffset:    int(i.Pos()),
						StartPos:     start,
						EndPos:       end,
						src:          string([]rune(a.src)[start:end]),
					}
					ret = append(ret, item)

				default:
					item := &Ast{
						file:         a.file + ":unknown",
						tkfs:         a.tkfs,
						ast:          block,
						CursorOffset: offset - start - 1,
						RawOffset:    int(i.Pos()),
						StartPos:     start,
						EndPos:       end,
						src:          string([]rune(a.src)[start:end]),
					}
					ret = append(ret, item)
				}
			}

		}
		return ret
	}
	return
}
