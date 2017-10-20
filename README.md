# Seafile SDK For Golang
Golang版本的Seafile Web API库

# 当前支持的接口
- [ ] 基础接口（Basic）
  - [x] Token获取（AuthToken）
  - [x] Ping
  - [x] 认证Ping（Auth Ping）
- [x] 账户（Account）
  - [x] 获取账户信息
  - [x] 获取服务器信息
- [ ] 资料库
  - [x] 获取资料库列表
  - [x] 获取资料库上传链接
  - [x] 获取资料库更新链接
- [ ] 文件
  - [x] 上传文件
  - [x] 下载文件
  - [x] 更新文件
  - [x] 删除文件
  - [x] 重命名文件
- [ ] 目录
  - [x] 获取目录内容
  - [x] 创建目录
  - [x] 删除目录

# TBD
由于目前Seafile官方的文档并不完善，尤其是错误处理方面。有时候用HTTP状态吗、有时候用字符串、有时候用非固定的JSON字符串。

所以这部分并没有很好的办法来处理。

只有后期建议官方完善后再做处理，或者遇到问题后做相应的处理。
