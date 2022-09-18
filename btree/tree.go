package btree

import (
	"math"
	"sort"
)

type position int

const (
	left  = position(-1)
	none  = position(0)
	right = position(+1)
)

type Int int64

func (a Int) Less(b Item) bool {
	return a < b.(Int)
}

type Item interface {
	Less(than Item) bool
}

// NewBtree 创建m阶B树
func NewBtree(m int) *BTree {
	if m < 3 {
		panic("m >= 3")
	}
	var bt = &BTree{}
	bt.root = createNode(bt.m)
	bt.m = m
	bt.minItems = int(math.Ceil(float64(m)/2) - 1)
	bt.maxItems = m - 1
	return bt
}

func createNode(m int) *node {
	var node = &node{}
	node.items = make([]Item, m+1)
	node.items = node.items[0:0]

	return node
}

type items []Item

func (s items) find(key Item) (index int, exist bool) {
	i := sort.Search(len(s), func(i int) bool {
		return key.Less(s[i])
	})
	if i > 0 && !s[i-1].Less(key) {
		return i - 1, true
	}
	return i, false
}

type children []*node

type node struct {
	items    items    //关键字
	children children //叶子节点
	parent   *node
}

// 获取key
func (n *node) get(key Item) Item {
	i, found := n.items.find(key)
	if found {
		return n.items[i]
	} else if len(n.children) > 0 {
		return n.children[i].get(key)
	}
	return nil
}

// 获取key
func (n *node) getKey(key Item) (*node, int) {
	i, found := n.items.find(key)
	if found {
		return n, i
	} else if len(n.children) > 0 {
		return n.children[i].getKey(key)
	}
	return nil, 0
}

func (n *node) insert(bt *BTree, key Item) {
	i := sort.Search(len(n.items), func(i int) bool {
		return key.Less(n.items[i])
	})

	newItems := append(n.items[0:i], append(items{key}, n.items[i:]...)...)
	n.items = newItems

	bt.split(n)
}

// 获取key节点的左右子节点
func (n *node) getChildNodes(index int) (l, r *node) {
	if len(n.children) == 0 {
		return nil, nil
	}
	return n.children[index], n.children[index+1]
}

// 获取key节点的左右兄弟节点
func (n *node) getBrotherNodes() (l, r *node) {
	if n.parent == nil {
		return nil, nil
	}

	index := 0
	for i, pn := range n.parent.children {
		if pn == n {
			index = i
			break
		}
	}

	if index == 0 {
		return nil, n.parent.children[1]
	} else {
		if len(n.parent.children) == index+1 {
			return n.parent.children[index-1], nil
		} else {
			return n.parent.children[index-1], n.parent.children[index+1]
		}
	}
}

// 获取后继节点
func (n *node) getSuccessorNode(index int) *node {
	ln := n.children[index+1]
	if len(ln.children) == 0 {
		return ln
	}
	for nod := ln.children[0]; nod != nil; ln = ln.children[0] {
	}
	return ln
}

func (n *node) delete(index int) Item {
	key := n.items[index]
	if index == 0 {
		n.items = append(items{}, n.items[1:]...)
	} else {
		n.items = append(n.items[0:index], n.items[index+1:]...)
	}
	return key
}

// 移除一个子节点
//TODO 下标越界可能
func (n *node) deleteChild(index int) {
	if index == 0 {
		n.children = append(children{}, n.children[1:]...)
	} else {
		n.children = append(n.children[0:index], n.children[index+1:]...)
	}
}

//前后插入一个子节点
func (n *node) addChild(move *node, p position) {
	if p == left {
		n.children = append(children{move}, n.children...)
	}

	if p == right {
		n.children = append(n.children, move)
	}

	move.parent = n
}

// 获取node在父节点中的index
func (n *node) getIndexInParent() int {
	if n.parent == nil {
		panic("n is root node")
	}
	index := 0
	for i, pn := range n.parent.children {
		if pn == n {
			index = i
			break
		}
	}
	return index
}

type BTree struct {
	root     *node
	m        int //阶
	minItems int //除根节点外单个节点最少包含的元素数量
	maxItems int //单个节点最多包含的元素数量
}

// Get 查找
func (bt *BTree) Get(key Item) Item {
	return bt.root.get(key)
}

// InsertMultiple 批量插入
func (bt *BTree) InsertMultiple(keys []int) {
	for _, v := range keys {
		bt.Insert(Int(v))
	}
}

// Insert 插入数据
func (bt *BTree) Insert(key Item) {
	if bt.Get(key) != nil {
		return
	}
	node := getRightNode(bt.root, key)
	node.insert(bt, key)
}

// Delete 删除key
func (bt *BTree) Delete(key Item) bool {
	n, i := bt.root.getKey(key)
	if n == nil {
		return false
	}

	if len(n.children) == 0 && n.parent == nil { //只有一个根节点
		n.delete(i)
		return true
	}

	if len(n.children) == 0 && n.parent != nil && len(n.items) > bt.minItems { //叶子节点 & 节点数量充足
		n.delete(i)
		return true
	}

	if len(n.children) == 0 && n.parent != nil && len(n.items) <= bt.minItems { //叶子节点、节点不足
		n.delete(i)
		bt.reBalance(n)
		return true
	}

	if len(n.children) != 0 { // 非叶子节点
		sn := n.getSuccessorNode(i)
		n.items[i] = sn.items[0] //后继替换
		sn.delete(0)             //后继key删除
		bt.reBalance(sn)
		return true
	}

	return false
}

//合并两个节点
func (bt *BTree) mergeNodes(l, r *node) {
	ix := l.getIndexInParent()
	k := l.parent.delete(ix)
	rChildren := r.children

	l.items = append(l.items, k)
	l.items = append(l.items, r.items...)

	l.parent.deleteChild(ix + 1)

	// 右子分支合并到左分支中
	if len(l.children) > 0 {
		for _, v := range rChildren {
			v.parent = l
		}
		l.children = append(l.children, rChildren...)
	}

	if l.parent.parent == nil && len(l.parent.items) == 0 { //根节点
		l.parent.children = make(children, 0)
		l.parent = nil
		bt.root = l
	} else {
		//递归向根验证
		bt.reBalance(l.parent)
	}
}

//再平衡
func (bt *BTree) reBalance(n *node) {

	if n.parent == nil {
		return
	} //根节点

	if len(n.items) >= bt.minItems {
		return
	}

	//判断兄弟节点是否可借
	var bn *node = nil
	p := none
	l, r := n.getBrotherNodes()
	if l != nil && len(l.items) > bt.minItems {
		bn = l
		p = left
	}

	if bn == nil && r != nil && len(r.items) > bt.minItems {
		bn = r
		p = right
	}

	if bn != nil { //可借
		var key, pk Item
		ix := bn.getIndexInParent()
		if p == left {
			key = bn.delete(len(bn.items) - 1)
			pk = bn.parent.items[ix]
			bn.parent.items[ix] = key
			n.items = append(items{pk}, n.items...)

			//借完需要处理子节点的归属问题
			if len(bn.children) > 0 {
				moveChild := bn.children[len(bn.children)-1]
				bn.deleteChild(len(bn.children) - 1)
				n.addChild(moveChild, left)
			}
		} else {
			key = bn.delete(0)
			pk = bn.parent.items[ix-1]
			bn.parent.items[ix-1] = key
			n.items = append(n.items, pk)

			//借完需要处理子节点的归属问题
			if len(bn.children) > 0 {
				moveChild := bn.children[0]
				bn.deleteChild(0)
				n.addChild(moveChild, right)
			}
		}

	} else { //不可借
		if l != nil {
			bt.mergeNodes(l, n)
		} else {
			bt.mergeNodes(n, r)
		}
	}
}

// 插入拆分
func (bt *BTree) split(n *node) {
	if len(n.items) <= bt.maxItems {
		return
	}
	middle := bt.m / 2
	l := make(items, middle)
	r := make(items, middle)
	copy(l, n.items[:middle])
	copy(r, n.items[middle+1:])
	middleItem := n.items[middle]
	n.items = l

	//右节点
	rNode := createNode(bt.m)
	rNode.items = r

	if n.parent != nil { //有父节点
		i := n.getIndexInParent()

		if i == 0 {
			n.parent.items = append(items{middleItem}, n.parent.items[:]...)
			n.parent.children = append(n.parent.children[0:1], append(children{rNode}, n.parent.children[1:]...)...)
		} else {
			n.parent.items = append(n.parent.items[0:i], append(items{middleItem}, n.parent.items[i:]...)...)
			n.parent.children = append(n.parent.children[0:i+1], append(children{rNode}, n.parent.children[i+1:]...)...)
		}

		rNode.parent = n.parent
		//递归处理父节点
		bt.split(n.parent)
	} else { //没有父节点
		newRoot := createNode(bt.m)
		newRoot.items = append(newRoot.items, middleItem)
		newRoot.children = append(newRoot.children, n, rNode)
		n.parent = newRoot
		bt.root = newRoot
		rNode.parent = bt.root

		if len(n.children) > 0 { //处理被拆分节点的子节点
			rChildren := make(children, bt.m-middle)
			copy(rChildren, n.children[middle+1:])
			n.children = n.children[0 : middle+1]
			rNode.children = rChildren
			for _, v := range rNode.children {
				v.parent = rNode
			}
		}
	}
}

// 获取待插入的叶子节点
func getRightNode(root *node, key Item) *node {
	//根节点
	if len(root.children) == 0 {
		return root
	}

	for index, item := range root.items {
		if key.Less(item) {
			return getRightNode(root.children[index], key)
		}
	}
	return getRightNode(root.children[len(root.items)], key)
}
