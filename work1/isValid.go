func isValid(s string) bool {
    stack := []rune{}
	mapping := map[rune]rune{
		')': '(',
		']': '[',
		'}': '{',
	}
	for _, str := range s {
		if str == '(' || str == '[' || str == '{' {
			stack = append(stack, str)
		} else if str == ')' || str == ']' || str == '}' {
			if len(stack) == 0 {
				return false
			}
			top := stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			if mapping[str] != top {
				return false
			}
		}
	}
	return len(stack) == 0
}