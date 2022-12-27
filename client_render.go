package mermaid

import (
	_ "embed"
	"html/template"

	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/util"
)

const defaultMermaidJS = "https://cdn.jsdelivr.net/npm/mermaid/dist/mermaid.min.js"

//go:generate curl -O -L https://cdn.jsdelivr.net/npm/mermaid/dist/mermaid.min.js
//go:embed mermaid.min.js
var mermaidJSEmbedded []byte

// ClientRenderer renders Mermaid diagrams as HTML,
// to be rendered into images client side.
//
// It operates by installing a <script> tag into the document
// that renders the Mermaid diagrams client-side.
type ClientRenderer struct {
	// URL of Mermaid Javascript to be included in the page.
	//
	// Defaults to the latest version available on cdn.jsdelivr.net.
	// "<embed>" uses the embedded source file
	MermaidJS string
}

// RegisterFuncs registers the renderer for Mermaid blocks with the provided
// Goldmark Registerer.
func (r *ClientRenderer) RegisterFuncs(reg renderer.NodeRendererFuncRegisterer) {
	reg.Register(Kind, r.Render)
	reg.Register(ScriptKind, r.RenderScript)
}

// Render renders mermaid.Block nodes.
func (*ClientRenderer) Render(w util.BufWriter, src []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	n := node.(*Block)
	if entering {
		w.WriteString(`<div class="mermaid">`)
		lines := n.Lines()
		for i := 0; i < lines.Len(); i++ {
			line := lines.At(i)
			template.HTMLEscape(w, line.Value(src))
		}
	} else {
		w.WriteString("</div>")
	}
	return ast.WalkContinue, nil
}

// RenderScript renders mermaid.ScriptBlock nodes.
func (r *ClientRenderer) RenderScript(w util.BufWriter, src []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	_ = node.(*ScriptBlock) // sanity check
	if entering {
		mermaidJS := r.MermaidJS
		if mermaidJS == "<embed>" {
			w.WriteString(`<script>`)
			w.Write(mermaidJSEmbedded)
			w.WriteString(`</script>`)
		} else {
			if len(mermaidJS) == 0 {
				mermaidJS = defaultMermaidJS
			}
			w.WriteString(`<script src="`)
			w.WriteString(mermaidJS)
			w.WriteString(`"></script>`)
		}
	} else {
		w.WriteString("<script>mermaid.initialize({startOnLoad: true});</script>")
	}

	return ast.WalkContinue, nil
}
