package main

import (
	"test/animal"
	"testing"
)

// 编写一个猫会叫的测试
func TestCat(t *testing.T) {
	got := animal.Cat()
	if got != "~~~qqq" {
		t.Errorf("Cat say %s,excepted %s", got, "~~~qqq")
	}
}
