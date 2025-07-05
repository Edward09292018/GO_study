func merge(intervals [][]int) [][]int {
    if len(intervals) == 0 {
        return [][]int{}
    }

    // 按照起始位置排序
    sort.Slice(intervals, func(i, j int) bool {
        return intervals[i][0] < intervals[j][0]
    })

    var merged [][]int
    for _, interval := range intervals {
        // 如果 merged 为空或者当前区间的起始位置大于 merged 中最后一个区间的结束位置，则直接添加
        if len(merged) == 0 || interval[0] > merged[len(merged)-1][1] {
            merged = append(merged, interval)
        } else {
            // 否则，合并区间
            merged[len(merged)-1][1] = max(interval[1], merged[len(merged)-1][1])
        }
    }

    return merged
}

// 辅助函数用于取最大值
func max(a, b int) int {
    if a > b {
        return a
    }
    return b
}