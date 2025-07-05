func removeDuplicates(nums []int) int {
    if len(nums) == 0 {
        return 0
    }

    // 指针 i 记录唯一元素的位置
    k := 1

    // 遍历数组，如果当前元素与前一个不同，则放入 nums[k]
    for i := 1; i < len(nums); i++ {
        if nums[i] != nums[k-1] {
            nums[k] = nums[i]
            k++
        }
    }

    return k
}