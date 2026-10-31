package missing

// Verify 检查 got 是不是 values 逐个平方的结果，返回有多少个位置不对。
//
//   - 第 i 个位置不对：got[i] != values[i] * values[i]；
//   - got 比 values 短时，缺的那些位置都算不对；
//   - got 比 values 长时，多出来的部分不管。
//
// 这是你的“量尺”：桥接测试会用它数出 Collect 每次少了多少个结果。
func Verify(values, got []int) int {
	// TODO
	return 0
}
