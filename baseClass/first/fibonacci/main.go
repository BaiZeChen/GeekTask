package main

import "fmt"

func main() {

	var a = fibonacci(4)

	fmt.Println(a)

}

func fibonacci(num int) int {
	if num < 0 {
		return -1
	}
	if num < 2 {
		return 1
	}

	return fibonacci(num-1) + fibonacci(num-2)

}
