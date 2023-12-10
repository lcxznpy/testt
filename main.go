package main

import (
	"fmt"
	"test/animal"
)

func main() {
	var password = "qaqaqqaqq"
	fmt.Printf("password:%s", password)
	got := animal.Cat()
	fmt.Println(got)
}
