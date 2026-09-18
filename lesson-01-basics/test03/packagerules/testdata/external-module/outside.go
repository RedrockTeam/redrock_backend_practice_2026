package outside

import _ "example.com/hello/internal/secret" // TODO: 独立 module 不能导入另一个 module 的 internal 包
