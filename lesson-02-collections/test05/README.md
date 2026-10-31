# test05：从零实现一个本地 DNS 缓存

对应课件：第二课 · map（创建、增删改查、`v, ok := m[k]`、range 顺序随机、句柄传递）+ 课件 map 小练习“DNS 缓存”。
题型：**从零实现**。

## 背景

课件里的 DNS 缓存只有几行：`dns["www.baidu.com"] = "110.242.68.66"`。真实的解析器还要处理三件事：

1. 域名**不区分大小写**，`WWW.Baidu.com.` 和 `www.baidu.com` 是同一个名字（末尾的 `.` 表示根域）；
2. 经常要**反查**：这个 IP 上挂了哪些域名？
3. 正查和反查是两份索引，**改一边必须同步另一边**。

`dns/dns.go` 只给了类型、错误变量和方法签名，所有方法都返回零值。`Resolver` 的字段由你自己设计。

## 任务

| 方法 | 要求 |
|---|---|
| `Normalize(domain) string` | 去两端空白 → 转小写 → 去掉末尾一个 `.` |
| `New() *Resolver` | 空 Resolver |
| `(*Resolver) Add(domain, ip) error` | 规范化后记录；任一为空返回 `ErrInvalid` 且不修改状态；同域名覆盖旧 IP，反查跟着更新 |
| `(*Resolver) Lookup(domain) (string, bool)` | 查询前规范化；不存在 `"", false` |
| `(*Resolver) Remove(domain)` | 删除前规范化；反查同步删除；不存在不 panic |
| `(*Resolver) Len() int` | 域名数 |
| `(*Resolver) ReverseLookup(ip) []string` | 指向该 IP 的所有域名，字典序升序；IP 去两端空白；返回值是副本 |
| `(*Resolver) Domains() []string` | 所有域名，字典序升序 |

零值 `var r dns.Resolver` 也必须能直接用：第一次写入之前再 `make`（lazy 初始化）。

## 设计提示（不是答案）

先在纸上回答这三个问题，再动手：

1. 反查结果用什么存？`map[string][]string` 删除时要在切片里找、再删；`map[string]map[string]bool`（集合）删除是一步，但输出时要排序。两种都可以，各有代价；
2. `Add` 覆盖一个已存在的域名时，旧 IP 那边的记录谁来删？
3. 一个 IP 上的域名全部删光后，这个 IP 的键要不要留下？（留下不算错，但想清楚 `ReverseLookup` 返回什么）

其它：

- `strings.TrimSpace`、`strings.ToLower`、`strings.TrimSuffix`；
- range map 的顺序是随机的，有序输出必须 `sort.Strings`；
- 返回内部切片会让调用方绕过 Resolver 修改数据——这和 `test02` 的窗口问题、`test08` 的 `Snapshot` 是同一类问题；
- 对 `nil` map 读和 `delete` 都安全，写会 panic。

## 评分点

| 测试 | 检查内容 |
|---|---|
| `TestNormalize` | 大小写、空白、末尾 `.` |
| `TestAddAndLookupAreCaseInsensitive` | 写入和查询都规范化 |
| `TestAddRejectsEmptyInput` | 非法输入返回 `ErrInvalid`，状态不变 |
| `TestOverwriteUpdatesReverseIndex` | 覆盖后旧 IP 的反查记录消失 |
| `TestRemoveKeepsBothDirectionsInSync` | 删除同步两份索引 |
| `TestReverseLookupIsSortedAndACopy` | 有序；修改返回值不影响内部 |
| `TestDomainsSorted` | 有序；空 Resolver |
| `TestZeroValueResolver` | 零值可用 |

## 做完以后

你刚写的是一个“两份索引必须保持一致”的数据结构。它在单个执行流里正确，因为每次 `Add` / `Remove` 都会把两张表一起改完，中间没人能看到“改了一半”的状态。

如果两个 goroutine 同时 `Add` 呢？一个刚改完正查、还没改反查，另一个就来 `ReverseLookup`——这就是第四课 `test05` 账本里“总额和明细对不上”的问题。

另外，DNS 本身就是第五课计算机网络的内容：浏览器访问 `www.baidu.com` 之前，第一步就是把域名解析成 IP。
