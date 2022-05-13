package list_test

import (
	"fmt"
	"github.com/leolu9527/algorithm/list"
)

func Example() {
	l0 := new(list.Node)
	l0.Data = 1

	l1 := new(list.Node)
	l1.Data = 2

	l0.Add(l1)

	for root := l0; root != nil; root = root.Next {
		fmt.Println(root.Data)
	}

	// Output:
	// 1
	// 2
}
