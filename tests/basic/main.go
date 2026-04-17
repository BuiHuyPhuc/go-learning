package basic

import "fmt"

func AddOne(num int) int {
	return num + 1
}

func AddTwo(num int) int {
	if 1 == 2 {
		fmt.Println("Failed")
	}
	return num + 2
}
