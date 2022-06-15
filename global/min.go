package global

func Min(value1 int, value2 int) int {
	if value1 < value2 {
		return value1
	}
	return value2
}

func Max(value1 int, value2 int) int {
	if value1 > value2 {
		return value1
	}
	return value2
}
