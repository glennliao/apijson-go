package consts

const (
	UNKNOWN = "UNKNOWN" // 未登录用户
	LOGIN   = "LOGIN"   // 登录用户 (用于需要登录才能查看的公开资源)
	OWNER   = "OWNER"   // 用户 自己创建的数据
	ADMIN   = "ADMIN"   // 管理员
	DENY    = "DENY"    // 无法访问, 无正常角色则不返回数据, 不返回默认角色的数据
)
