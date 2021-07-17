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
	"bytes"
	"strings"

	"github.com/88250/gulu"
	"github.com/88250/lute"
	"github.com/88250/lute/ast"
	"github.com/88250/lute/html"
	"github.com/88250/lute/parse"
)

func ParseJSON(luteEngine *lute.Lute, jsonData []byte) (ret *parse.Tree, needFix bool, err error) {
	root := &ast.Node{}
	err = gulu.JSON.UnmarshalJSON(jsonData, root)
	if nil != err {
		return
	}

	ret = &parse.Tree{Name: "", Root: &ast.Node{Type: ast.NodeDocument, ID: root.ID}, Context: &parse.Context{ParseOption: luteEngine.ParseOptions}}
	ret.Root.KramdownIAL = parse.Map2IAL(root.Properties)
	ret.Context.Tip = ret.Root
	if nil == root.Children {
		return
	}

	idMap := map[string]bool{}
	for _, child := range root.Children {
		genTreeByJSON(child, ret, &idMap, &needFix)
	}
	ret.ID = ret.Root.ID
	return
}

func genTreeByJSON(node *ast.Node, tree *parse.Tree, idMap *map[string]bool, needFix *bool) {
	node.Tokens, node.Type = gulu.Str.ToBytes(node.Data), ast.Str2NodeType(node.TypeStr)
	node.Data, node.TypeStr = "", ""
	node.KramdownIAL = parse.Map2IAL(node.Properties)
	node.Properties = nil

	// 历史数据订正
	if ast.NodeList == node.Type {
		if 1 > len(node.Children) {
			*needFix = true
			return // 忽略空列表
		}
	} else if ast.NodeSuperBlock == node.Type && (0 == len(node.Children) || (3 == len(node.Children) && "NodeSuperBlockCloseMarker" == node.Children[2].TypeStr)) {
		*needFix = true
		return // 忽略空超级块
	} else if ast.NodeMathBlock == node.Type {
		if 1 > len(node.Children) {
			*needFix = true
			return // 忽略空公式
		}
	}
	fixLegacyData(node, idMap, needFix)

	tree.Context.Tip.AppendChild(node)
	tree.Context.Tip = node
	defer tree.Context.ParentTip()
	if nil == node.Children {
		return
	}
	for _, child := range node.Children {
		genTreeByJSON(child, tree, idMap, needFix)
	}
	node.Children = nil
}

func fixLegacyData(node *ast.Node, idMap *map[string]bool, needFix *bool) {
	if node.IsBlock() && "" == node.ID {
		node.ID = ast.NewNodeID()
		node.SetIALAttr("id", node.ID)
		*needFix = true
	}
	if "" != node.ID {
		if _, ok := (*idMap)[node.ID]; ok {
			node.ID = ast.NewNodeID()
			node.SetIALAttr("id", node.ID)
			*needFix = true
		}
		(*idMap)[node.ID] = true
	}

	if ast.NodeIFrame == node.Type && bytes.Contains(node.Tokens, gulu.Str.ToBytes("iframe-content")) {
		start := bytes.Index(node.Tokens, gulu.Str.ToBytes("<iframe"))
		end := bytes.Index(node.Tokens, gulu.Str.ToBytes("</iframe>"))
		node.Tokens = node.Tokens[start : end+9]
		*needFix = true
	}

	if ast.NodeList == node.Type && nil != node.ListData && 3 != node.ListData.Typ && 0 < len(node.Children) &&
		nil != node.Children[0].ListData && 3 == node.Children[0].ListData.Typ {
		node.ListData.Typ = 3
		*needFix = true
	}

	if ast.NodeMark == node.Type && 3 == len(node.Children) && "NodeText" == node.Children[1].TypeStr {
		if strings.HasPrefix(node.Children[1].Data, " ") || strings.HasSuffix(node.Children[1].Data, " ") {
			node.Children[1].Data = strings.TrimSpace(node.Children[1].Data)
			*needFix = true
		}
	}

	if ast.NodeInlineMathContent == node.Type {
		if bytes.Contains(node.Tokens, []byte("&gt;")) || bytes.Contains(node.Tokens, []byte("&lt;")) || bytes.Contains(node.Tokens, []byte("&amp;")) || bytes.Contains(node.Tokens, []byte("&quot;")) {
			node.Tokens = html.UnescapeHTML(node.Tokens)
			*needFix = true
		}
	}

	if ast.NodeHeading == node.Type && 6 < node.HeadingLevel {
		node.HeadingLevel = 6
		*needFix = true
	}

	if ast.NodeLinkDest == node.Type && bytes.HasPrefix(node.Tokens, []byte("assets/")) && bytes.HasSuffix(node.Tokens, []byte(" ")) {
		node.Tokens = bytes.TrimSpace(node.Tokens)
		*needFix = true
	}
}
