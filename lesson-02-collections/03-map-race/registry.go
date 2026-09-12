package registry

// UnsafeRegistry 故意没有同步保护，用来观察 go test -race 的报告。
type UnsafeRegistry struct{ values map[string]int }

func New() *UnsafeRegistry                           { return &UnsafeRegistry{values: make(map[string]int)} }
func (r *UnsafeRegistry) Set(key string, value int)  { r.values[key] = value }
func (r *UnsafeRegistry) Get(key string) (int, bool) { value, ok := r.values[key]; return value, ok }
