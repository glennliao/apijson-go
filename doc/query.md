
# Get开放查询接口

### 核心思路
json是可递归的, 可以将每一个json节点转换成Node对象, 最终构成节点树

例如
```json
{
  "[]": {
    "Todo": {
      "user_id@": "User/user_id"
    },
    "User": {
      "user_id@": "[]/Todo/user_id"
    }
  },
  "User": {
    "user_id": "10001"
  }
}
```
可以看成最外层是一个节点, 节点请求内容即为整个json, 他具有两个子节点, `[]`和`User`,
```json
{
  "Todo": {
    "user_id@": "User/user_id"
  },
  "User": {
    "user_id@": "[]/Todo/user_id"
  }
}
```
```json
{
    "user_id": "10001"
}
```
然后对于`[]`,又具有两个子节点, `Todo`、`User`

对于每一个Node, 分别有 `new->buildChild->parse->fetch->Result` 阶段

- new: 新建
- buildChild: 构建子节点
- parse: 解析当前节点的请求参数
- fetch: 获取值
- Result: 组装返回值

每一个节点还有Key和Path属性, Key 则为当前json节点中的Key,Path 则是该节点在整个json中的路径。 可以将Key看成为当前文件名, Path则为他的绝对路径

例如
```json
{
  "[]": {
    "Todo": {
      "user_id@": "User/user_id"
    },
    "User": {
      "user_id@": "[]/Todo/user_id"
    }
  }
}

```
中的User, Key为`User`, 路径则为 `[]/User`


## 查询流程
1. 创建一个`Query`
2. 将原始json请求生成为一个`rootNode`
3. 执行`buildChild`构建`rootNode`子节点. 
4. 执行`rootNode` 的 `parse` 解析请求, 并解析关联关系(不要求json的key顺序, 因为go的原生map不支持顺序遍历)
5. 分析节点树的依赖关系, 获取节点执行顺序, 并依次执行 `fetch`
6. 结果组装, 返回`rootNode`的`Result()` 完成本次查询



## 节点类型
节点根据内容划分为以下类型
- 查询节点: 该节点为实际查询数据库的节点, 其下面的内容可以看成是查询条件,不再往下构建子节点
- 引用节点: 该节点的值引用其他节点的值 (暂只为`total@`使用)
- 结构节点: 该节点仅为结构支撑 (例如: `[]`)


查询节点的判定:
- key 大写开头 (对应数据表)

引用节点判定:
- key 为 total@

其他则为结构节点


## 限制
1. `[]`节点下有且只有一个主查询表(不依赖兄弟节点的查询节点)
2. 由于是应用内拼接数据完成`n+1`的问题, 所以以下写法的total并不能获取到 (Todo[]是列表中主查询表User的副表)
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



## 待实现
- [ ] 限制page的最大值,count区间
- [ ] 分析节点树的复杂度, 限制最大复杂度
- [ ] 增加 replace节点的sqlexecutor, 使用自定义完成数据的获取 (例如实际存储时候使用同一个表保存不同数据, 实际需要根据用户id或者其他信息来完成， 或者数据来自别的数据源)
