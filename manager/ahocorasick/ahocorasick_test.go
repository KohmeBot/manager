package ahocorasick

import (
	"fmt"
	"testing"
)

func TestMatcher_Contains(t *testing.T) {
	m := NewStringMatcher([]string{"okok", "nb"})
	fmt.Println(m.Contains([]byte("ok123oknbok")))

}
