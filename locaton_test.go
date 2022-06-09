package WebGuard

import (
	"fmt"
	"testing"
)

func _Test(t *testing.T) {
	fmt.Println(LoopUpLocation([]string{"上海", "郑州", "Beijing", "局域网"}, "127.0.0.1"))
}
