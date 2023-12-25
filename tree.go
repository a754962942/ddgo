package ddgo

import "strings"

/*路由前缀书实现*/

type treeNode struct {
	name       string
	children   []*treeNode
	routerName string
	isEnd      bool
}

// Put path: /user/get/:id
func (t *treeNode) Put(path string) {
	root := t
	strs := strings.Split(path, "/")
	for index, name := range strs {
		if index == 0 {
			continue
		}
		isMatch := false
		for _, child := range t.children {
			if child.name == name {
				isMatch = true
				t = child
				break
			}
		}
		if !isMatch {
			isEnd := false
			if index == len(strs)-1 {
				isEnd = true
			}
			t2 := &treeNode{
				isEnd:    isEnd,
				name:     name,
				children: make([]*treeNode, 0),
			}
			t.children = append(t.children, t2)
			t = t2
		}
	}
	t = root

}

// get path: /user/get/1
func (t *treeNode) Get(path string) *treeNode {
	strs := strings.Split(path, "/")
	var routerName string

	for index, str := range strs {
		isMatch := false
		for _, child := range t.children {
			if child.name == str ||
				child.name == "*" ||
				strings.Contains(child.name, ":") {
				isMatch = true
				routerName += "/" + child.name
				child.routerName = routerName
				t = child
				if index == len(strs)-1 {
					return child
				}
				break
			}
		}
		if !isMatch {
			for _, child := range t.children {
				if child.name == "**" {
					routerName += "/" + child.name
					child.routerName = routerName
					return child
				}
			}
		}
	}
	return nil
}
