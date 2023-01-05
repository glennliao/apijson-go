# apijson-go [下水阶段]
基于 go + goframe 实现的 apijson

> ~~暂仍处于前期开发探索中, 请仅使用在 个人探索项目~~

> 目前处于 【下水阶段】, 欢迎测试、issue、建议、pr

[RoadMap 阶段规划](./@doc/roadmap.md)

# 快速体验
<a href="https://gitpod.io/#https://github.com/glennliao/apijson-go"  target="_blank"><img src="https://gitpod.io/button/open-in-gitpod.svg" /> </a>

创建后 执行 @demo/todo/tests 下的 *_test.go 访问测试


# 使用指南
暂参考demo目录下的todo

## 文档参考
1. [Get开放查询](./@doc/query.md)
2. [非开放请求](./@doc/action.md)
3. [权限控制](./@doc/access.md)


# 开发指南
1. go >= 1.18
2. 创建 mysql 数据库
3. 导入 demo/todo/doc/todo.sql文件
4. demo/todo/config.yaml.example 改成 demo/todo/config.yaml, 然后修改配置文件 config.yaml 中数据库连接
5. 在demo/todo 目录运行 go run main.go 或者 查看测试 demo/todo/tests




# 感谢
- [GoFrame](https://github.com/gogf/gf)
- [APIJSON](https://github.com/Tencent/APIJSON)

# 参考链接
1. [详细的说明文档.md](https://github.com/Tencent/APIJSON/blob/master/%E8%AF%A6%E7%BB%86%E7%9A%84%E8%AF%B4%E6%98%8E%E6%96%87%E6%A1%A3.md)
2. [最新规范文档](https://github.com/Tencent/APIJSON/blob/master/Document.md)
3. [todo demo doc](https://github.com/jerrylususu/apijson_todo_demo/blob/master/FULLTEXT.md)
4. [如何实现其它语言的APIJSON？](https://github.com/Tencent/APIJSON/issues/38)