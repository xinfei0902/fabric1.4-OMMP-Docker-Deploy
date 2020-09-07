package objectdefine

import (
	"fmt"
	"testing"
)

func Test_a(t *testing.T) {
	a := make([]TaskType, 5)

	a[3].ID = "aa"

	fmt.Println(a)
}
