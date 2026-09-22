package common // TODO: 这个包只负责用户名校验，请使用更明确的包名

func Valid(name string) bool {
	return name != ""
}
