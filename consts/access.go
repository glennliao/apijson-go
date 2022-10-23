package consts

const (
	UNKNOWN = "UNKNOWN" // 未登录用户
	LOGIN   = "LOGIN"   // 登录用户 (用于需要登录才能查看的公开资源)
	OWNER   = "OWNER"   // 用户 自己创建的数据
	ADMIN   = "ADMIN"   // 管理员
)

const (
	MethodGet    = "GET"
	MethodHead   = "HEAD"
	MethodPost   = "POST"
	MethodPut    = "PUT"
	MethodDelete = "DELETE"
)
