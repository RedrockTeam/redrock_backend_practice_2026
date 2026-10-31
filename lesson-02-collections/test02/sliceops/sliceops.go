package sliceops

// Clone 返回 values 的一份独立副本。
// 修改副本中的元素，不应该影响原来的 values。
func Clone(values []int) []int {
	// TODO: 用 make + copy，或者 append([]int(nil), values...)。
	// 直接 return values 只是复制了切片头，底层数组仍然是同一份。
	return values
}

// DoubleInPlace 把 values 中每个元素都乘以 2，并且原地修改 values。
func DoubleInPlace(values []int) {
	// TODO: 用下标访问 values[i]。
	// for _, value := range values 拿到的 value 是元素副本，改它不会改切片。
	for _, value := range values {
		value *= 2
	}
}

// AppendOne 返回在 values 末尾追加 x 之后的新切片。
// 注意：不能丢掉 append 的返回值。
func AppendOne(values []int, x int) []int {
	// TODO: return append(values, x)。
	// append 可能扩容并返回新的切片头，不能只写 append(...) 然后 return values。
	return values
}

// Sum 返回 values 中所有元素的和；空切片返回 0。
func Sum(values []int) int {
	total := 0
	for _, value := range values {
		total += value
	}
	return total
}
