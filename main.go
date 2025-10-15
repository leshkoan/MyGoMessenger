package main

import "fmt"

func main() {
	a := false
	var s string
	if a {
		s = "Да"
	} else {
		s = "Нет"
	}
	fmt.Printf("Hello, World!%s\n", s)
}
