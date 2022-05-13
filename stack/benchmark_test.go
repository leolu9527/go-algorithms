package stack

import (
	"fmt"
	"math/rand"
	"testing"
)

//go test -benchmem -bench . -run=none
func BenchmarkArrayStackParallel(b *testing.B) {
	var stack = new(ArrayStack)
	b.RunParallel(func(pb *testing.PB) {
		s := RandStringBytes(16)
		for pb.Next() {
			x := rand.Intn(1)
			if x == 0 {
				stack.Push(s)
			} else {
				_, err := stack.Pop()
				if err != nil {
					fmt.Println(err)
				}
			}
		}
	})
}

func BenchmarkLinkStackParallel(b *testing.B) {
	var stack = new(LinkStack)
	b.RunParallel(func(pb *testing.PB) {
		s := RandStringBytes(16)
		for pb.Next() {
			x := rand.Intn(1)
			if x == 0 {
				stack.Push(s)
			} else {
				_, err := stack.Pop()
				if err != nil {
					fmt.Println(err)
				}
			}
		}
	})
}
