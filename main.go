package main

import "fmt"

func Cat() string {
	return "~~~qqq"
}

func main() {
	got := Cat()
	fmt.Println(got)
}
