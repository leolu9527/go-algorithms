package btree_test

import (
	"fmt"
	"github.com/leolu9527/algorithm/btree"
)

func Example() {
	bt := btree.NewBtree(5)
	bt.Insert(btree.Int(4))
	bt.InsertMultiple([]int{20, 44, 89, 96, 25, 30, 33, 60, 75, 81, 85, 110, 120, 101, 150, 158, 130, 135, 138})
	fmt.Println(bt.Delete(btree.Int(44)))
	fmt.Println(bt.Delete(btree.Int(138)))
	fmt.Println(bt.Delete(btree.Int(99)))

	// Output:
	// true
	// true
	// false
}
