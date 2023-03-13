package consts

const (
	ListKeySuffix       = "[]"
	RefKeySuffix        = "@"
	FunctionsKeySuffix  = "()"
	FunctionOriReqParam = "$req"
)

const (
	RowKey = "rowKey"
	Raw    = "@raw"
)

const (
	Role  = "@role"
	Page  = "page"  // page num
	Count = "count" // page size // todo access中增加限制count,防止被恶意下载数据
	Query = "query"
)

const (
	OpLike   = "$"
	OpRegexp = "~"
	OpSub    = "-"
	OpPLus   = "+"
)

const (
	SqlLike   = "LIKE"
	SqlEqual  = "="
	SqlRegexp = "REGEXP"
)
