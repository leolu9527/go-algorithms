package stack

import (
	"errors"
	"fmt"
	"sync"
)

type Stack interface {
	Push(s string)
	Pop() (string, error)
	Peek() (string, error)
}

// IsTest 用于测试
var IsTest = false

// ArrayStack 数组栈
type ArrayStack struct {
	array   []string   //切片作为底层
	maxSize int        //栈大小
	lock    sync.Mutex //并发安全锁
}

// Push 入栈
func (stack *ArrayStack) Push(s string) {
	stack.lock.Lock()
	defer stack.lock.Unlock()
	//插入数据最后
	stack.array = append(stack.array, s)
	if IsTest {
		fmt.Println("push:", s)
	}
	stack.maxSize++
}

// Pop 出栈
func (stack *ArrayStack) Pop() (string, error) {
	stack.lock.Lock()
	defer stack.lock.Unlock()
	if stack.maxSize == 0 {
		return "", errors.New("empty")
	}
	s := stack.array[stack.maxSize-1]

	//切片收缩
	destArray := make([]string, stack.maxSize-1, stack.maxSize-1)
	//这里使用copy函数，copy函数会丢弃原切片中最后一个元素
	copy(destArray, stack.array)
	stack.array = destArray

	//循环
	//newArray := make([]string, stack.maxSize-1, stack.maxSize-1)
	//for i := 0; i < stack.maxSize-1; i++ {
	//	newArray[i] = stack.array[i]
	//}
	//stack.array = newArray

	if IsTest {
		fmt.Println("pop:", s)
	}
	// 栈中元素数量-1
	stack.maxSize--
	return s, nil
}

// Peek 获取栈顶元素
func (stack *ArrayStack) Peek() (string, error) {
	if stack.maxSize == 0 {
		return "", errors.New("empty")
	}
	return stack.array[stack.maxSize-1], nil
}

// LinkStack 单链表栈
type LinkStack struct {
	root *LinkNode  //栈顶
	size int        //栈深
	lock sync.Mutex //并发锁
}

type LinkNode struct {
	value string
	Next  *LinkNode
}

// Push 入栈
func (stack *LinkStack) Push(s string) {
	stack.lock.Lock()
	defer stack.lock.Unlock()
	lastNode := stack.root
	node := new(LinkNode)
	node.value = s
	node.Next = lastNode
	stack.root = node
	stack.size++
}

func (stack *LinkStack) Pop() (string, error) {
	stack.lock.Lock()
	defer stack.lock.Unlock()
	if stack.root == nil {
		return "", errors.New("empty")
	}
	s := stack.root.value
	nextNode := stack.root.Next
	stack.root = nextNode
	stack.size--
	return s, nil
}

func (stack *LinkStack) Peek() (string, error) {
	stack.lock.Lock()
	defer stack.lock.Unlock()
	if stack.size == 0 {
		return "", errors.New("empty")
	}
	s := stack.root.value

	return s, nil
}
