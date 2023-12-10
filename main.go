package main

import (
	"fmt"
	"test/animal"
)

func main() {
	var password = "ed4b76e52f61a1fc8f3d997815aaafcf731a2f78"
	fmt.Printf("password:%s", password)
	got := animal.Cat()
	fmt.Println(got)
}
