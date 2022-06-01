package tree

import "container/list"

type Node struct {
	Val         int
	Left, Right *Node
}

// BuildTree 输入一个切片 ：[3,9,20,0,0,15,7]
func BuildTree(l []int) (root *Node) {
	length := len(l)
	if length == 0 {
		return root
	}

	var nodes = make([]*Node, length)
	root = &Node{
		Val: l[0],
	}
	nodes[0] = root
	//循环输入的数组切片，依次判断每一个结点的左右结点是否存在并创建
	for i := 0; i < length; i++ {
		currentNode := nodes[i]

		if currentNode == nil {
			continue
		}

		leftIndex := 2*i + 1
		if leftIndex < length && l[leftIndex] != 0 {
			currentNode.Left = &Node{
				Val: l[leftIndex],
			}
			nodes[leftIndex] = currentNode.Left
		}

		rightIndex := 2*i + 2
		if rightIndex < length && l[rightIndex] != 0 {
			currentNode.Right = &Node{
				Val: l[rightIndex],
			}
			nodes[rightIndex] = currentNode.Right
		}
	}

	return root
}

// MaxDepth 计算树的深度
func MaxDepth(root *Node) int {
	return len(dfsRecursion(root, 0, [][]int{}))
}

// 深度优先
// [1,[2,3],[4,5,6]]
func dfsRecursion(node *Node, level int, nodes [][]int) [][]int {
	if node == nil {
		return nodes
	}

	// 判断切片长度是否满足要求
	if level < len(nodes) {
		nodes[level] = append(nodes[level], node.Val)
	} else {
		nodes = append(nodes, []int{node.Val})
	}
	nodes = dfsRecursion(node.Left, level+1, nodes)
	nodes = dfsRecursion(node.Right, level+1, nodes)

	return nodes
}

// ConvertToArr 树转化成数组
func ConvertToArr(root *Node) []int {
	if root == nil {
		return []int{}
	}

	result := indexRecursion(root, 1, 0, []int{})
	// 删除末尾的0
	for i := len(result) - 1; i >= 0; i-- {
		if result[i] == 0 {
			result = result[:i]
		} else {
			break
		}
	}

	return result
}

// 索引循环
func indexRecursion(node *Node, level int, index int, result []int) []int {
	//深度为level的树最多有2^level-1个节点，容量不够时扩容依据
	if len(result) < (1<<level - 1) {
		newArr := make([]int, 1<<level-1)
		copy(newArr, result)
		result = newArr
	}
	result[index] = node.Val

	if node.Left != nil {
		result = indexRecursion(node.Left, level+1, 2*index+1, result)
	}

	if node.Right != nil {
		result = indexRecursion(node.Right, level+1, 2*index+2, result)
	}

	return result
}

// BFS 宽度优先搜索
func BFS(root *Node) (result [][]int) {
	if root == nil {
		return result
	}

	queue := list.New()
	queue.PushFront(root)

	for queue.Len() > 0 {
		var currentLevel []int
		listLength := queue.Len()
		for i := 0; i < listLength; i++ {
			node := queue.Remove(queue.Back()).(*Node)
			currentLevel = append(currentLevel, node.Val)
			if node.Left != nil {
				queue.PushFront(node.Left)
			}
			if node.Right != nil {
				queue.PushFront(node.Right)
			}
		}
		result = append(result, currentLevel)
	}

	return result
}
