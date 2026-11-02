# test02:把静态资源挂到 URL 上

对应课件:第六课。

## 任务

授课人员会在 `cat/` 目录下放一张图片 `cat.jpg`。编辑 `cat/main.go`,注册 `GET /cat`,让浏览器访问 `http://localhost:8080/cat` 时直接显示这张图片。

提示:gin 提供了 `StaticFile` 方法,一行就能把磁盘上的文件绑到一个 URL 上,文件的 MIME 类型(这里是 `image/jpeg`)由扩展名自动识别。服务器从**包目录**启动,所以文件路径写 `./cat.jpg`。

## 自己验证

```bash
cd lesson-06-web/test02/cat
go run .                       # 监听 http://localhost:8080
```

浏览器打开 `http://localhost:8080/cat`,应该看到那张猫图。

## 评分点

| 测试 | 检查内容 |
|---|---|
| `TestCat` | `/cat` 返回 200、`Content-Type` 是 `image/jpeg`,响应体与 `cat.jpg` 逐字节一致 |
