package config

type Config struct {
	Access *Access

	Functions *Functions

	MaxTreeWidth int
	MaxTreeDeep  int

	rowKeyGenFuncMap map[string]RowKeyGenFuncHandler

	// dbFieldStyle 数据库字段命名风格 请求传递到数据库中
	DbFieldStyle FieldStyle

	// jsonFieldStyle 数据库返回的字段
	JsonFieldStyle FieldStyle

	DbMeta *DBMeta

	AccessList []AccessConfig // todo to access

	RequestConfig *RequestConfig
}

func New() *Config {
	a := &Config{}
	a.Access = NewAccess()

	a.MaxTreeWidth = 5
	a.MaxTreeDeep = 5

	a.rowKeyGenFuncMap = make(map[string]RowKeyGenFuncHandler)

	a.DbFieldStyle = CaseSnake
	a.JsonFieldStyle = CaseCamel

	a.Functions = &Functions{}

	return a
}
