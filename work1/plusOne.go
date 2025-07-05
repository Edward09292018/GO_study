package main

import "fmt"

func plusOne(digits []int) []int {
	//n:=len(digits)
	for i := len(digits)-1; i >=0 ; i-- {
		digits[i]++
		fmt.Println(digits[i])
		if digits[i]==10 {
			digits[i]=0
		}else {
			return digits
		}
	}
	return append([]int{1},digits...)
}

func main() {
	fmt.Println(plusOne([]int{1, 2, 3}))
	fmt.Println(plusOne([]int{4, 3, 2, 1}))
	fmt.Println(plusOne([]int{9}))
}
