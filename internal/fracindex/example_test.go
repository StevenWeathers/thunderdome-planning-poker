package fracindex

import (
	"fmt"
)

func ExampleKeyBetween() {
	a := "a0"
	b := "a1"
	key, err := KeyBetween(&a, &b)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println("Key between a0 and a1:", *key)
	// Output: Key between a0 and a1: a0P
}
