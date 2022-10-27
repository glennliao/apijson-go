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
- [x] 分页返回total@

- [x] 可用的权限方案
  - [x] get只有access中定义的才能访问
  - [x] 非get操作则必须与request指定一致才可请求
  - [x] 基于角色控制

- [ ] 远程函数
- [ ] 错误提示
- [ ] 查询节点 自定义查询数据


## 查询实现
1. 根据json构造节点树, 并检查节点结构(不符合直接返回)
2. parse 节点树内容, 并分析关联关系(不要求json的key顺序, 因为go的原生map不支持顺序遍历)
3. 从依赖关系中逐步fetch数据
4. 构造响应数据


# 列表查询限制

[//]: # (1. page,count 最大值)
- `[]`下只能有一个主查询表 (不依赖于列表中其他表)
- 由于是应用内拼接结果处理1+n的问题， 当前以下写法total并不能获取到[]下单个节点的Todo[]的total值,也不能对Todo[]分页,所以勿使用于拔出萝卜带出泥的场景
```json
{
	"[]":{
		"User":{

		},
		"Todo[]":{
			"user_id@":"/User/user_id"
		},
		"total@":"/Todo[]/total"
	}
}
```


# 权限控制方案
## 表
- get 操作只能访问access中定义的表
- 其他操作需和request中定义的结构一致方可访问

## 行
- access 中定义各操作的角色权限, 角色有: 未登录用户/登录用户(OWNER/ADMIN/ 其他自定义角色)
- 提供自定义函数 针对不同表不同角色添加 过滤条件

## 列
- access_ext 中定义各操作的in/out字段列表, 限制各操作字段只能是此处的子集

# 开发指南

1. go >= 1.18
2. 创建mysql数据库
3. 导入cmd/test.sql文件
4. 从cmd/config.yaml.example 改成 cmd/config.yaml, 然后修改配置文件cmd/config.yaml中数据库连接
5. 在cmd目录运行go run main.go
6. 查看测试 test.http / z_test.go

# 感谢
- [GoFrame](https://gitee.com/johng/gf)
- [APIJSON](https://gitee.com/Tencent/APIJSON)
- [tiangao/apijson-go](https://gitee.com/tiangao/apijson-go)

# 参考链接
1. [详细的说明文档.md](https://github.com/Tencent/APIJSON/blob/master/%E8%AF%A6%E7%BB%86%E7%9A%84%E8%AF%B4%E6%98%8E%E6%96%87%E6%A1%A3.md)
2. [最新规范文档](https://github.com/Tencent/APIJSON/blob/master/Document.md)
3. [todo demo doc](https://github.com/jerrylususu/apijson_todo_demo/blob/master/FULLTEXT.md)
4. [如何实现其它语言的APIJSON？](https://github.com/Tencent/APIJSON/issues/38)