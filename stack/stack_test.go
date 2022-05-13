package stack

import (
	"math/rand"
	"sync"
	"testing"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

const stringLen = 16

const stackSize = 1000

func RandStringBytes(n int) string {
	l := len(letterBytes)
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(l)]
	}
	return string(b)
}

// test stack Push by interface
func testStackPush(t *testing.T, stack Stack) {
	for i := 0; i < stackSize; i++ {
		s := RandStringBytes(stringLen)
		stack.Push(s)
		p, err := stack.Peek()
		if err != nil || p != s {
			t.Errorf("push:%s != peek: %s", s, p)
		}
	}
}

// test stack Pop by interface
func testStackPop(t *testing.T, stack Stack) {
	var storage [stackSize]string
	for i := 0; i < stackSize; i++ {
		s := RandStringBytes(stringLen)
		storage[i] = s
		stack.Push(s)
		p, err := stack.Peek()
		if err != nil || p != s {
			t.Errorf("push:%s, peek: %s", s, p)
		}
	}

	for i := 0; i < stackSize; i++ {
		p, err := stack.Pop()
		if err != nil {
			s := storage[stackSize-1-i]
			if p != s {
				t.Errorf("storage:%s != pop: %s", s, p)
			}
		}
	}
}

// test stack On Goroutines by interface
func testStackOnGoroutine(t *testing.T, stack Stack) {
	IsTest = false
	//IsTest = true
	wg := sync.WaitGroup{}
	for i := 0; i < stackSize; i++ {
		s := RandStringBytes(stringLen)
		wg.Add(1)
		go func() {
			stack.Push(s)
		}()
	}

	go func() {
		for {
			_, err := stack.Pop()
			if err == nil {
				wg.Done()
			}
		}
	}()

	wg.Wait()
	IsTest = false
}

func TestArrayStackPush(t *testing.T) {
	stack := new(ArrayStack)
	testStackPush(t, stack)
}

func TestArrayStackPop(t *testing.T) {
	stack := new(ArrayStack)
	testStackPop(t, stack)
}

func TestArrayStackOnGoroutine(t *testing.T) {
	stack := new(ArrayStack)
	testStackOnGoroutine(t, stack)
}

func TestLinkStackPush(t *testing.T) {
	stack := new(LinkStack)
	testStackPush(t, stack)
}

func TestLinkStackPop(t *testing.T) {
	stack := new(LinkStack)
	testStackPop(t, stack)
}

func TestLinkStackOnGoroutine(t *testing.T) {
	stack := new(LinkStack)
	testStackOnGoroutine(t, stack)
}
