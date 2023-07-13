package Gout

import (
	"strings"
)

type methodTree struct {
	method string
	root   *node
}

type methodTrees []methodTree
type node struct {
	fullPart  string        // 完整路由
	handlers  HandlersChain //处理函数
	part      string        // 路由中的一部分
	children  []*node       // 子节点
	wildChild bool          // 判断子节点有通配符节点
	isWild    bool          // 判断当前节点是否为通配符节点
}

// 匹配节点
func (n *node) matchChild(part string) *node {
	for _, child := range n.children {
		if child.part == part {
			return child
		}
	}
	return nil
}

// 所有匹配成功的节点
func (n *node) matchChildren(part string) []*node {
	nodes := make([]*node, 0)
	for _, child := range n.children {
		if child.part == part || child.isWild {
			nodes = append(nodes, child)
		}
	}
	return nodes
}

// 插入节点
func (n *node) insert(path string, parts []string, height int, handlers HandlersChain) {
	if len(parts) == height {
		n.fullPart = path
		n.handlers = handlers
		return
	}

	part := parts[height]
	child := n.matchChild(part)
	if child == nil {
		if part[0] == ':' || part[0] == '*' {
			n.wildChild = true
		}
		child = &node{part: part, isWild: part[0] == ':' || part[0] == '*'}
		n.children = append(n.children, child)
	}

	child.insert(path, parts, height+1, handlers)
}

// 更新 node 的 children 字段
func (n *node) addChild(child *node) {
	if n.wildChild && len(n.children) > 0 {
		wildcardChild := n.children[len(n.children)-1]
		n.children = append(n.children[:len(n.children)-1], child, wildcardChild)
	} else {
		n.children = append(n.children, child)
	}
}

// 搜索节点
func (n *node) search(parts []string, height int) *node {
	if len(parts) == height || strings.HasPrefix(n.part, "*") {
		if n.fullPart == "" {
			return nil
		}
		return n
	}

	part := parts[height]
	children := n.matchChildren(part)

	for _, child := range children {
		result := child.search(parts, height+1)
		if result != nil {
			return result
		}
	}

	return nil
}
