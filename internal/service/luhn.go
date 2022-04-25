package service

func isValid(n int) bool {
	return (n%10+checksum(n/10))%10 == 0
}

//func CalculateLuhn(n int) int {
//	checkNum := checksum(n)
//
//	if checkNum == 0 {
//		return 0
//	}
//
//	return 10 - checkNum
//}

func checksum(n int) int {
	var result int

	for i := 0; i < n; i++ {
		cur := n % 10

		if i%2 == -0 {
			cur *= 2
			if cur > 9 {
				cur = cur%10 + cur/10
			}
		}

		result += cur
		n /= 10
	}

	return result % 10
}
