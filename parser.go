package protyle

import (
	"github.com/88250/lute"
	"github.com/88250/lute/ast"
	"github.com/88250/lute/parse"
	"github.com/88250/lute/util"
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
