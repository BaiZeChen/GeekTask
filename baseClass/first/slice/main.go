package main

import "fmt"

func main() {
	s := []int{1, 2, 4, 7}
	// 结果应该是 5, 1, 2, 4, 7
	s = Add(s, 0, 5)
	fmt.Println(s)

	// 结果应该是5, 9, 1, 2, 4, 7
	s = Add(s, 1, 9)
	fmt.Println(s)

	// 结果应该是5, 9, 1, 2, 4, 7, 13
	s = Add(s, 6, 13)
	fmt.Println(s)

	// 结果应该是5, 9, 2, 4, 7, 13
	s = Delete(s, 2)
	fmt.Println(s)

	// 结果应该是9, 2, 4, 7, 13
	s = Delete(s, 0)
	fmt.Println(s)

	// 结果应该是9, 2, 4, 7
	s = Delete(s, 4)
	fmt.Println(s)

}

func Add(s []int, index int, value int) []int {
	//TODO
	if len(s) > index {
		newSlice := make([]int, len(s)+1)
		for i, k := 0, 0; i < len(s); i++ {
			if i == index {
				newSlice[k] = value
				newSlice[k+1] = s[i]
				k++
			} else {
				newSlice[k] = s[i]
			}
			k++
		}
		s = newSlice
	} else {
		newSlice := make([]int, index+1)
		for key, value := range s {
			newSlice[key] = value
		}
		newSlice[index] = value
		s = newSlice
	}
	return s
}

func Delete(s []int, index int) []int {
	// TODO
	var newSlice []int
	for key, value := range s {
		if key == index {
			continue
		}
		newSlice = append(newSlice, value)
	}
	s = newSlice
	return s
}
