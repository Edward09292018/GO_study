package main

func PlusOne(digits []int) []int {
	//n:=len(digits)
	for i := len(digits) - 1; i >= 0; i-- {
		digits[i]++
		if digits[i] == 10 {
			digits[i] = 0
		} else {
			return digits
		}
	}
	return append([]int{1}, digits...)
}
