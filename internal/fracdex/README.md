# Fractional Indexing (src https://github.com/rocicorp/fracdex)

This is based on [Implementing Fractional Indexing
](https://observablehq.com/@dgreensp/implementing-fractional-indexing) by [David Greenspan
](https://github.com/dgreensp).

Fractional indexing is a technique to create an ordering that can be used for [Realtime Editing of Ordered Sequences](https://www.figma.com/blog/realtime-editing-of-ordered-sequences/).

This implementation includes variable-length integers, and the prepend/append optimization described in David's article.

This should be byte-for-byte compatible with https://github.com/rocicorp/fractional-indexing.

## Example

```go
package main

import (
	"fmt"

	"github.com/StevenWeathers/thunderdome-planning-poker/internal/fracdex"
)

func main() {
	first, _ := fracdex.KeyBetween("", "") // a0
	fmt.Println(first)

	// Insert after 1st
	second, _ := fracdex.KeyBetween(first, "") // "a1"
	fmt.Println(second)

	// Insert after 2nd
	third, _ := fracdex.KeyBetween(second, "") // "a2"
	fmt.Println(third)

	// Insert before 1st
	zeroth, _ := fracdex.KeyBetween("", first) // "Zz"
	fmt.Println(zeroth)

	// Insert in between 2nd and 3rd
	secondAndHalf, _ := fracdex.KeyBetween(second, third) // "a1V"
	fmt.Println(secondAndHalf)
}
```