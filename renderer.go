package protyle

import (
	"github.com/88250/lute/ast"
	"github.com/88250/lute/parse"
	"github.com/88250/lute/render"
	"github.com/88250/lute/util"
)

type SYJSONRenderer struct {
	*render.BaseRenderer
}

func NewSYJSONRenderer(tree *parse.Tree, options *render.Options) render.Renderer {
	ret := &SYJSONRenderer{render.NewBaseRenderer(tree, options)}
	ret.DefaultRendererFunc = ret.renderNode
	return ret
}

func (r *SYJSONRenderer) renderNode(node *ast.Node, entering bool) ast.WalkStatus {
	if ast.NodeKramdownBlockIAL == node.Type {
		// TODO: 某些情况还是有 IAL 块，需要确认移除
		return ast.WalkContinue
	}

	if entering {
		if nil != node.Previous {
			r.WriteString(",")
		}
		node.Data, node.TypeStr = util.BytesToStr(node.Tokens), node.Type.String()
		node.Properties = parse.IAL2Map(node.KramdownIAL)
		data, err := MarshalJSON(node)
		node.Data, node.TypeStr = "", ""
		node.Properties = nil
		if nil != err {
			panic("marshal node to json failed: " + err.Error())
			return ast.WalkStop
		}
		n := util.BytesToStr(data)
		n = n[:len(n)-1] // 去掉结尾的 }
		r.WriteString(n)
		if nil != node.FirstChild {
			r.WriteString(",\"Children\":[")
		} else {
			r.WriteString("}")
		}
	} else {
		if nil != node.FirstChild {
			r.WriteByte(']')
			r.WriteString("}")
		}
	}
	return ast.WalkContinue
}
