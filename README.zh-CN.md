# apijson-go [WIP]
基于 go + goframe 实现的 apijson

> 暂仍处于前期开发探索中, 请仅使用在 个人探索项目

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
- [ ] 可用的权限方案

## 查询实现
1. 根据json构造节点树, 并检查节点结构(不符合直接返回)
2. parse 节点树内容, 并分析关联关系(不要求json的key顺序, 因为go的原生map不支持顺序遍历)
3. 从依赖关系中逐步fetch数据
4. 构造响应数据


# 列表查询限制

[//]: # (1. page,count 最大值)
- []下只能有一个主查询表 (不依赖于列表中其他表)


# 开发指南

1. go >= 1.18
2. 创建mysql数据库
3. 导入test.sql文件
4. 修改配置文件config.yaml中数据库连接
5. 运行go run main.go
6. 查看测试test.http


# 感谢
- [GoFrame](https://gitee.com/johng/gf)
- [APIJSON](https://gitee.com/Tencent/APIJSON)

# 参考链接
1. [详细的说明文档.md](https://github.com/Tencent/APIJSON/blob/master/%E8%AF%A6%E7%BB%86%E7%9A%84%E8%AF%B4%E6%98%8E%E6%96%87%E6%A1%A3.md)
2. [最新规范文档](https://github.com/Tencent/APIJSON/blob/master/Document.md)
3. [todo demo doc](https://github.com/jerrylususu/apijson_todo_demo/blob/master/FULLTEXT.md)
4. [如何实现其它语言的APIJSON？](https://github.com/Tencent/APIJSON/issues/38)