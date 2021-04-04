package MyTimeParse

import (
	"fmt"
	"testing"
)

func TestTimeParse(t *testing.T) {
	s := TimeParse("13:45")
	fmt.Println(s)
}
