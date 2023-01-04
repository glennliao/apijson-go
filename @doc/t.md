# 书签模型
可以为纯 ind 项目

## 个人
可查看、编辑、新增、删除自己的  -> userId == self
可创建分组 -> 

## 群组成员
可查看其他人分享
可分享、取消分享

## 群组管理员
可创建、编辑、删除分类
> 删除前检查如何处理

可强制取消别人的分享

## 群创建人
删除群+群管理员功能

# 一般业务系统模型

## 用户
查看发布的资源
查看自己的内容
查看自己区域（上级： 企业/行政区域）的内容
业务流程: 提交内容、考试、查看进度

> 复杂业务独立接口, 切勿拿着锤子看什么都是钉子

## 企业/门店
查看自己下面用户的数据
查看上级给自己的数据

## 区域/总管理员
查看下级数据
创建、修改、删除自己的数据

# todo 任务 模型
可为纯ind

## 个人
新增、发布、删除、完成 自己发布给自己的todo
无法删除非自己创建的
完成别人发布的时候 需要消息通知对方(消息通知, 不属于框架, 使用远程函数(hook)完成这个过程)


## 协助者
查看、完成别人分配给自己的todo


## 发布者/企业/区域管理员
发布、查看别人的todo


access_ext
    - ? idGen // id生成策略
    - ? deleteCheck // 删除检查

request_ext
    - tag
    - version
    - method
    - hooks
        - after 成功后续  清理缓存、发送消息等 
        - fail  失败后续 
        - before  操作前处理 (数据处理(function/js)) + 特殊权限校验

function
    - debug
    - name (唯一)
    - arguments
    - demo
    - detail
    - type
    - version 允许操作的最低版本
    - tag(允许的操作)
    - methods 允许的操作   // 使用 requestIdList 替代 version,tag, methods
    - back  返回值示例  ,  系统启动时校验demo是否正常通过(demo 是不是得多个)


# 功能点

## delete - 删除

- 权限检查  -- condition Where 可以返回sql
    - 所有者
    - 可删除权限(管理员、上级)


