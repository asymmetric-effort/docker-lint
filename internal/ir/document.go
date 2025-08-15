// file: internal/ir/document.go
// (c) 2025 Asymmetric Effort, LLC. scaldwell@asymmetric-effort.com
package ir

import (
	"strings"

	"github.com/moby/buildkit/frontend/dockerfile/parser"
)

// Document is a normalized representation of a Dockerfile.
//
// Document retains stage information extracted from the Dockerfile AST.
type Document struct {
	Filepath string
	Stages   []*Stage
	AST      *parser.Node
}

// Stage represents a single FROM instruction.
//
// Stage records the source image, optional name, and AST node for positioning.
type Stage struct {
	Index int
	Name  string
	From  string
	Node  *parser.Node
}

// BuildDocument converts an AST into a Document.
//
// BuildDocument iterates the AST, collecting FROM instructions as stages.
func BuildDocument(path string, ast *parser.Node) (*Document, error) {
	doc := &Document{Filepath: path, AST: ast}
	idx := 0
	for _, n := range ast.Children {
		if strings.EqualFold(n.Value, "from") {
			from := ""
			name := ""
			if n.Next != nil {
				from = n.Next.Value
				for tok := n.Next.Next; tok != nil; tok = tok.Next {
					if strings.EqualFold(tok.Value, "as") && tok.Next != nil {
						name = tok.Next.Value
						break
					}
				}
			}
			doc.Stages = append(doc.Stages, &Stage{Index: idx, Name: name, From: from, Node: n})
			idx++
		}
	}
	return doc, nil
}
