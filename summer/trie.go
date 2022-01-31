package summer

import (
	"errors"
	"fmt"
	"log"
	"path"
	"strings"
)

// Trie 路由树结构
type Trie struct {
	root *node
}

func NewTrie() *Trie {
	return &Trie{root: &node{parentPath: "/"}}
}

type node struct {
	isLeafNode bool                // 叶节点标志，如果是叶节点则说明该节点是一个uri
	parentPath string              // 该节点的所有父节点路径段，例如：'/user/:id', segment='name'
	segment    string              // uri中的路径段，例如：'user'属于'/user/:id/name'中的一个segment
	handlers   []ControllerHandler // 控制器链(中间件、业务逻辑控制器)
	handler    ControllerHandler   // 控制器，只有该节点是叶节点时才存在
	children   []*node             // 子节点
}

func newNode() *node {
	return &node{children: []*node{}}
}

// isWildSegment 用于判断该路径段是否为通配符路径段
func isWildSegment(segment string) bool {
	return strings.HasPrefix(segment, ":")
}

func (n *node) searchChildren(segment string) []*node {
	if len(n.children) == 0 {
		return nil
	}
	// 如果是通配符路径段，则所有子节点都作为匹配节点
	if isWildSegment(segment) {
		return n.children
	}

	// 如果不是通配符路径段，则查找匹配的子节点
	nodes := make([]*node, 0, len(n.children))
	for _, child := range n.children {
		// 如果子节点的路径段是通配符，或子节点路径段不是通配符，路径段文本完全匹配，则视为一个匹配的节点
		if isWildSegment(child.segment) || child.segment == segment {
			nodes = append(nodes, child)
		}
	}
	return nodes
}

func (n *node) matchNode(uri string) *node {
	segments := strings.SplitN(uri, "/", 2)
	// 取得第一个url路径段，用于匹配下一个节点
	segment := segments[0]
	if !isWildSegment(segment) {
		segment = strings.ToLower(segment)
	}

	// 匹配下一节点
	children := n.searchChildren(segment)
	if len(children) == 0 {
		return nil
	}

	// case 1: 如果uri只有一个segment，则查找其中符合路由规则(segment相等 或 通配符segment)的叶子节点
	if len(segments) == 1 {
		for _, child := range children {
			if child.isLeafNode && (child.segment == segment || isWildSegment(child.segment)) {
				return child
			}
		}
		// 没有找到则说明无匹配的节点，终止查询
		return nil
	}
	// case 2: 如果uri有2个及以上的segment，则说明需要继续递归查找下一层匹配的子节点
	for _, child := range children {
		matchedNode := child.matchNode(segments[1])
		if matchedNode != nil {
			return matchedNode
		}
	}
	return nil
}

func (t *Trie) AddRouter(uri string, handler ControllerHandler) error {
	parent := t.root
	n := parent.matchNode(uri)
	if n != nil {
		return errors.New(fmt.Sprintf("route '%s' exist: leaf-node parent path '%s'", uri, n.parentPath))
	}

	segments := strings.Split(uri, "/")
	for index, segment := range segments {
		// 对所有非通配符路径段都转为小写
		if !isWildSegment(segment) {
			segment = strings.ToLower(segment)
		}
		// 查找是否存在合适的节点
		var childNode *node
		children := parent.searchChildren(segment)
		if len(children) > 0 {
			for _, child := range children {
				if child.segment == segment {
					childNode = child
					break
				}
			}
		}
		if childNode == nil {
			childNode = &node{segment: segment, parentPath: path.Join(parent.parentPath, parent.segment)}
			if index == len(segments)-1 {
				childNode.isLeafNode = true
				childNode.handler = handler
			}
			parent.children = append(parent.children, childNode)
		}
		parent = childNode
	}
	return nil
}

func (t *Trie) FindHandler(uri string) ControllerHandler {
	log.Printf("start finding handler for uri '%s'", uri)
	n := t.root.matchNode(uri)
	if n == nil {
		log.Printf("handler not found for uri '%s'", uri)
		return nil
	}
	log.Printf("matched handler for uri '%s'", uri)
	return n.handler
}
