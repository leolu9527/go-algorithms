package stack

import (
	"fmt"
	"math/rand"
	"testing"
)

// Benchmark for ArrayStack
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

// Benchmark for LinkStack
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
