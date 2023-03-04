package config

type Config struct {
	Access           *Access
	MaxTreeWidth     int
	MaxTreeDeep      int
	rowKeyGenFuncMap map[string]RowKeyGenFuncHandler
}

func New() *Config {
	a := &Config{}
	a.Access = NewAccess()

	a.MaxTreeWidth = 5
	a.MaxTreeDeep = 5

	a.rowKeyGenFuncMap = make(map[string]RowKeyGenFuncHandler)

	return a
}
