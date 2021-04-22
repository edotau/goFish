package algorithms

import(
    "testing"
    "fmt"
)

// Input: k = 2, k = 3, A = {0, 1)

func TestDeBruijn(t *testing.T) {
    k, n := byte(2), byte(3)
    db := DeBruijn(k, n)
    fmt.Printf("the debruijn sequencing of k=%b and n=%b is: %v\n", k, n, db)
}
