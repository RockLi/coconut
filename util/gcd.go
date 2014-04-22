package util

func Gcd(nums ...uint) uint {
	if len(nums) < 2 {
		return 0
	}

	x := gcd(nums[0], nums[1])

	for _, num := range nums[2:] {
		x = gcd(x, num)
	}

	return x
}

func gcd(a, b uint) uint {
	if a == 0 || b == 0 {
		return 0
	}
	for b != 0 {
		a, b = b, a%b
	}

	return a
}
