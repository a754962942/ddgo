package ddgo

import (
	"fmt"
	"testing"
)

func TestTreeNode(t *testing.T) {
	node := &treeNode{
		name:     "/",
		children: make([]*treeNode, 0),
	}
	node.Put("/user/get/:id")
	node.Put("/user/create/hello")
	node.Put("/user/create/aaa")
	node.Put("/users/create/bbb")
	get := node.Get("/user/get/1")
	fmt.Println(get)
	get = node.Get("/user/create/hello")
	fmt.Println(get)
	get = node.Get("/user/create/aaa")
	fmt.Println(get)
	get = node.Get("/users/create/bbb")
	fmt.Println(get)
	get = node.Get("/users/create/ccc")
	fmt.Println(get)
}
