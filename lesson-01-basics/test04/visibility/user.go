package visibility

// User 用来练习结构体字段、方法和包可见性。
// 类型名 User 大写，所以其他包可以使用这个类型。
type User struct {
	Name string // 大写字段：其他包可以直接读写
	age  int    // 小写字段：只有 visibility 包内可以访问
}

// NewUser 创建一个 User。
// 因为 age 是小写字段，其他包无法自己填它，只能通过这个构造函数。
func NewUser(name string, age int) *User {
	// TODO: 同时保存传入的 name 和 age。
	return &User{}
}

// Age 把小写字段 age 以只读的方式暴露给其他包。
func (u *User) Age() int {
	// TODO: 返回小写字段 age。
	return 0
}

// Grow 是其他包唯一能用来增加年龄的入口。
func (u *User) Grow() {
	// TODO: 调用包内的小写方法 grow，不要在这里直接写 u.age++。
}

// grow 只在 visibility 包内可见。
func (u *User) grow() {
	// TODO: 把 age 增加 1。
}

// Stage 返回用户所处的阶段。
// 注意：它调用的 stage 函数定义在同一个包的另一个文件 stage.go 里。
// Go 的小写标识符是「包内可见」，不是 C 的 static「文件内可见」。
func (u *User) Stage() string {
	// TODO: 用 stage 函数把 age 换算成阶段字符串。
	return ""
}
