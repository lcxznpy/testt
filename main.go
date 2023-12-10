package main

import (
	"fmt"
	"test/animal"
)

func main() {
	var password = "qaqaqqaqq"
	fmt.Println("password:%s", password)
	got := animal.Cat()
	fmt.Println(got)
}
