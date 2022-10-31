# apijson-go [WIP]
基于 go + goframe 实现的 apijson

> 暂仍处于前期开发探索中, 请仅使用在 个人探索项目


# 快速体验
<a href="https://gitpod.io/#https://github.com/glennliao/apijson-go"  target="_blank"><img src="https://gitpod.io/button/open-in-gitpod.svg" /> </a>

创建后 执行 demo/todo/todo/tests 下的 *_test.go 访问测试

# 功能实现

- [x] 单表查询、单表数组查询
- [x] 双表一对一关联查询、数组关联查询
- [x] 双表一对多关联查询、数组关联查询
- [x] @column, @order, @group, page, count
- [x] 单表单条新增
- [x] 单表单条修改
- [x] 单表单条、批量删除
- [x] Request表的tag校验
  - [x] MUST
  - [x] REFUSE
- [x] 分页返回total@

- [x] 可用的权限方案
  - [x] get只有access中定义的才能访问
  - [x] 非get操作则必须与request指定一致才可请求
  - [x] 基于角色控制
- [ ] 远程函数
- [ ] 错误提示
- [ ] 查询节点 自定义查询数据
- [ ] 字段限制
- [ ] 请求结构复杂度限制

## 文档参考
1. [Get开放查询](./doc/query.md)
2. [非开放请求](./doc/action.md)
3. [权限控制](./doc/access.md)


# 开发指南

1. go >= 1.18
2. 创建mysql数据库
3. 导入demo/todo/todo/todo.sql文件
4. demo/todo/config.yaml.example 改成 demo/todo/config.yaml, 然后修改配置文件 config.yaml 中数据库连接
5. 在demo/todo目录运行go run main.go
6. 查看测试 demo/todo/todo/tests



# 感谢
- [GoFrame](https://gitee.com/johng/gf)
- [APIJSON](https://gitee.com/Tencent/APIJSON)
- [tiangao/apijson-go](https://gitee.com/tiangao/apijson-go)

# 参考链接
1. [详细的说明文档.md](https://github.com/Tencent/APIJSON/blob/master/%E8%AF%A6%E7%BB%86%E7%9A%84%E8%AF%B4%E6%98%8E%E6%96%87%E6%A1%A3.md)
2. [最新规范文档](https://github.com/Tencent/APIJSON/blob/master/Document.md)
3. [todo demo doc](https://github.com/jerrylususu/apijson_todo_demo/blob/master/FULLTEXT.md)
4. [如何实现其它语言的APIJSON？](https://github.com/Tencent/APIJSON/issues/38)