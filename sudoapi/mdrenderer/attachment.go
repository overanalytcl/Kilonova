package mdrenderer

import (
	"fmt"
	"html"
	"net/url"
	"strings"

	"github.com/KiloProjects/kilonova"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/text"
	"github.com/yuin/goldmark/util"
	"go.uber.org/zap"
)

// Attachments are of the form ~[name.xyz]

var _ goldmark.Extender = &attNode{}
var _ renderer.NodeRenderer = &attachmentRenderer{}
var _ parser.InlineParser = &attachmentParser{}

var attNodeKind = ast.NewNodeKind("attachment")

type attachmentParser struct{}

func (attachmentParser) Trigger() []byte {
	return []byte{'~'}
}

func (attachmentParser) Parse(parent ast.Node, block text.Reader, pc parser.Context) ast.Node {
	line, _ := block.PeekLine()
	if len(line) < 2 {
		return nil
	}
	if line[1] != '[' {
		return nil
	}
	i := 2
	for ; i < len(line); i++ {
		if line[i] == ']' {
			break
		}
	}
	if i >= len(line) || line[i] != ']' {
		return nil
	}
	block.Advance(i + 1)
	fileName := line[2:i]
	return &AttachmentNode{Filename: string(fileName)}
}

type attachmentRenderer struct{}

func (att *attachmentRenderer) RegisterFuncs(rd renderer.NodeRendererFuncRegisterer) {
	rd.Register(attNodeKind, att.renderAttachment)
}

func (att *attachmentRenderer) renderAttachment(writer util.BufWriter, source []byte, n ast.Node, entering bool) (ast.WalkStatus, error) {
	if !entering {
		return ast.WalkContinue, nil
	}

	align := "left"
	width := ""
	var inline bool

	node := n.(*AttachmentNode)
	parts := strings.Split(node.Filename, "|")
	name := parts[0]
	if len(parts) > 1 {
		for _, part := range parts {
			kv := strings.SplitN(part, "=", 2)
			if len(kv) == 2 {
				switch kv[0] {
				case "align":
					if align == "left" || align == "right" || align == "center" {
						align = kv[1]
					}
				case "width":
					width = kv[1]
				case "inline":
					inline = true
				}
			} else if len(kv) == 1 && kv[0] == "inline" {
				inline = true
			}
		}
	}
	ctx, ok := n.OwnerDocument().Meta()["ctx"].(*kilonova.RenderContext)
	var link string
	if !ok || ctx == nil || (ctx.Problem == nil && ctx.BlogPost == nil) {
		link = url.PathEscape(name)
	} else {
		if ctx.Problem != nil {
			link = fmt.Sprintf("/problems/%d/attachments/%s", ctx.Problem.ID, url.PathEscape(name))
		} else if ctx.BlogPost != nil {
			link = fmt.Sprintf("/posts/%s/attachments/%s", ctx.BlogPost.Slug, url.PathEscape(name))
		} else {
			zap.S().Warn("WTF")
		}
	}

	extra := ""
	if inline {
		extra += ` data-imginline="true" `
	}
	if width != "" {
		extra += ` style="width:` + html.EscapeString(width) + `" `
	}
	fmt.Fprintf(writer, `<img src="%s" data-imgalign="%s" %s></img>`, link, align, extra)
	return ast.WalkContinue, nil
}

type attNode struct{}

func (*attNode) Extend(md goldmark.Markdown) {
	md.Renderer().AddOptions(renderer.WithNodeRenderers(util.Prioritized(&attachmentRenderer{}, 900)))
	md.Parser().AddOptions(parser.WithInlineParsers(util.Prioritized(&attachmentParser{}, 900)))
}

type AttachmentNode struct {
	ast.BaseInline

	Filename string
}

func (yt *AttachmentNode) Dump(source []byte, level int) {
	ast.DumpHelper(yt, source, level, map[string]string{"filename": yt.Filename}, nil)
}

func (yt *AttachmentNode) Kind() ast.NodeKind {
	return attNodeKind
}
