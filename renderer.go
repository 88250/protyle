// Protyle - 解析渲染思源笔记文档数据
// Copyright (c) 2021-present, b3log.org
//
// Lute is licensed under Mulan PSL v2.
// You can use this software according to the terms and conditions of the Mulan PSL v2.
// You may obtain a copy of Mulan PSL v2 at:
//         http://license.coscl.org.cn/MulanPSL2
// THIS SOFTWARE IS PROVIDED ON AN "AS IS" BASIS, WITHOUT WARRANTIES OF ANY KIND, EITHER EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT, MERCHANTABILITY OR FIT FOR A PARTICULAR PURPOSE.
// See the Mulan PSL v2 for more details.

package protyle

import (
	"github.com/88250/lute/ast"
	"github.com/88250/lute/parse"
	"github.com/88250/lute/render"
	"github.com/88250/lute/util"
)

type JSONRenderer struct {
	*render.BaseRenderer
}

func NewJSONRenderer(tree *parse.Tree, options *render.Options) render.Renderer {
	ret := &JSONRenderer{render.NewBaseRenderer(tree, options)}
	ret.DefaultRendererFunc = ret.renderNode
	return ret
}

func (r *JSONRenderer) renderNode(node *ast.Node, entering bool) ast.WalkStatus {
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
