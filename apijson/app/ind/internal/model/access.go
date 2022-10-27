package model

type (
	AccessListIn struct {
		ReqList
	}

	AccessListResult struct {
		Id     string   `json:"id" des:""`
		Debug  int8     `json:"debug" des:"是否为调试表，只允许在开发环境使用，测试和线上环境禁用"`
		Name   string   `json:"name" des:"实际表名，例如 apijson_user"`
		Alias  string   `json:"alias" des:"外部调用的表别名，例如 User"`
		Get    []string `json:"get" des:"允许 get 的角色列表"`
		Head   []string `json:"head" des:"允许 head 的角色列表"`
		Gets   []string `json:"gets" des:"允许 gets 的角色列表"`
		Heads  []string `json:"heads" des:"允许 heads 的角色列表"`
		Post   []string `json:"post" des:"允许 post 的角色列表"`
		Put    []string `json:"put" des:"允许 put 的角色列表"`
		Delete []string `json:"delete" des:"允许 delete 的角色列表"`
		Detail string   `json:"detail" des:""`
	}
)

type AccessGetOut struct {
	Id     string   `json:"id" des:""`
	Debug  int8     `json:"debug" des:"是否为调试表，只允许在开发环境使用，测试和线上环境禁用"`
	Name   string   `json:"name" des:"实际表名，例如 apijson_user"`
	Alias  string   `json:"alias" des:"外部调用的表别名，例如 User"`
	Get    []string `json:"get" des:"允许 get 的角色列表"`
	Head   []string `json:"head" des:"允许 head 的角色列表"`
	Gets   []string `json:"gets" des:"允许 gets 的角色列表"`
	Heads  []string `json:"heads" des:"允许 heads 的角色列表"`
	Post   []string `json:"post" des:"允许 post 的角色列表"`
	Put    []string `json:"put" des:"允许 put 的角色列表"`
	Delete []string `json:"delete" des:"允许 delete 的角色列表"`
	Detail string   `json:"detail" des:""`
}

type (
	AccessAddIn struct {
		Debug  *int8     `des:"是否为调试表，只允许在开发环境使用，测试和线上环境禁用"`
		Name   *string   `des:"实际表名，例如 apijson_user"`
		Alias  *string   `des:"外部调用的表别名，例如 User"`
		Get    *[]string `des:"允许 get 的角色列表"`
		Head   *[]string `des:"允许 head 的角色列表"`
		Gets   *[]string `des:"允许 gets 的角色列表"`
		Heads  *[]string `des:"允许 heads 的角色列表"`
		Post   *[]string `des:"允许 post 的角色列表"`
		Put    *[]string `des:"允许 put 的角色列表"`
		Delete *[]string `des:"允许 delete 的角色列表"`
		Detail *string   `des:""`
	}

	AccessAddOut struct {
		Id string
	}
)

type AccessUpdateIn struct {
	Id    string `des:""`
	Debug *int8  `des:"是否为调试表，只允许在开发环境使用，测试和线上环境禁用"`
	//Name   *string   `des:"实际表名，例如 apijson_user"`
	Alias  *string   `des:"外部调用的表别名，例如 User"`
	Get    *[]string `des:"允许 get 的角色列表"`
	Head   *[]string `des:"允许 head 的角色列表"`
	Gets   *[]string `des:"允许 gets 的角色列表"`
	Heads  *[]string `des:"允许 heads 的角色列表"`
	Post   *[]string `des:"允许 post 的角色列表"`
	Put    *[]string `des:"允许 put 的角色列表"`
	Delete *[]string `des:"允许 delete 的角色列表"`
	Detail *string   `des:""`
}

// #035351d91781135326f3411607a61e6a1fb29c81:021026:52
