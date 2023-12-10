package main

import "testing"

// 编写一个猫会叫的测试
func TestCat(t *testing.T) {
	got := Cat()
	if got != "qaq" {
		t.Errorf("Cat say %s,excepted %s", got, "qaq")
	}
}
