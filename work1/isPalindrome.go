func isPalindrome(x int) bool {
   if x < 0 || (x%10 == 0 && x != 0) {
        return false
    }

    reverted := 0
    for x > reverted {
        reverted = reverted*10 + x%10
        x /= 10
    }

    // 偶数位时 x == reverted，奇数位时 x == reverted/10
    return x == reverted || x == reverted/10
}