// Package list 单链表
package list

type Node struct {
	Data interface{}
	Next *Node
}

func (n *Node) Add(next *Node) {
	n.Next = next
}
