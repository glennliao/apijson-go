## 1. 调用方式一  来自返回的字段 -> 因为默认simple字段会返回 -> 需要返回吗
```json
{
  "name": "hi",
  "aaa": "demo",
  "ref()": "sayHello(name,aaa)",
  "@a": 0,
  "ref2()": "ret(@a)",
  "User": {
    "pic()": "getPic(userId)" // 来自当前User的字段, 需要分析函数依赖的字段和依赖函数字段的节点
  }
}
```

## 2
```json
{
  "msg()": "sayHi",
  "msg2()": "sayHi()"
}
```