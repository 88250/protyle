package protyle

import (
	"github.com/88250/lute"
	"github.com/88250/lute/ast"
	"github.com/88250/lute/parse"
	"github.com/88250/lute/render"
	"github.com/88250/lute/util"
	jsoniter "github.com/json-iterator/go"
)

func ParseSYJSON(luteEngine *lute.Lute, jsonData []byte) (ret *parse.Tree, err error) {
	root := &ast.Node{}
	err = UnmarshalJSON(jsonData, root)
	if nil != err {
		return nil, err
	}

	ret = &parse.Tree{Name: "", Root: &ast.Node{Type: ast.NodeDocument, ID: root.ID}, Context: &parse.Context{ParseOption: luteEngine.ParseOptions}}
	ret.Root.KramdownIAL = parse.Map2IAL(root.Properties)
	ret.Context.Tip = ret.Root
	if nil == root.Children {
		return
	}

	idMap := map[string]bool{}
	for _, child := range root.Children {
		genTreeByJSON(child, ret, &idMap)
	}
	ret.ID = ret.Root.ID
	return
}

func genTreeByJSON(node *ast.Node, tree *parse.Tree, idMap *map[string]bool) {
	node.Tokens, node.Type = util.StrToBytes(node.Data), ast.Str2NodeType(node.TypeStr)
	node.Data, node.TypeStr = "", ""
	node.KramdownIAL = parse.Map2IAL(node.Properties)
	node.Properties = nil

	// 历史数据订正
	if node.IsBlock() && "" == node.ID {
		node.ID = ast.NewNodeID()
		node.SetIALAttr("id", node.ID)
	}
	if "" != node.ID {
		if _, ok := (*idMap)[node.ID]; ok {
			node.ID = ast.NewNodeID()
			node.SetIALAttr("id", node.ID)
		}
		(*idMap)[node.ID] = true
	}

	tree.Context.Tip.AppendChild(node)
	tree.Context.Tip = node
	defer tree.Context.ParentTip()
	if nil == node.Children {
		return
	}
	for _, child := range node.Children {
		genTreeByJSON(child, tree, idMap)
	}
	node.Children = nil
}

func UnmarshalJSON(data []byte, v interface{}) error {
	json := jsoniter.ConfigCompatibleWithStandardLibrary
	return json.Unmarshal(data, v)
}

func MarshalJSON(v interface{}) ([]byte, error) {
	json := jsoniter.ConfigCompatibleWithStandardLibrary
	return json.Marshal(v)
}

func MarshalIndentJSON(v interface{}, prefix, indent string) ([]byte, error) {
	json := jsoniter.ConfigCompatibleWithStandardLibrary
	return json.MarshalIndent(v, prefix, indent)
}

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
