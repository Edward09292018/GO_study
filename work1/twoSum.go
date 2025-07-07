package main

func twoSum(nums []int, target int) []int {
	numMap := make(map[int]int)
	for i, num := range nums {
		cal := target - num
		if j, e := numMap[cal]; e {
			return []int{j, i}
		}
		numMap[num] = i
	}
	return []int{}
}
