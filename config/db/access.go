package db

import (
	"github.com/glennliao/apijson-go/config"
	"github.com/glennliao/apijson-go/util"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/samber/lo"
	"net/http"
)

var accessMap = map[string]Access{}

type FieldsGetValue struct {
	In  map[string][]string
	Out map[string]string
}

type Access struct {
	Debug     int8
	Name      string
	Alias     string
	Get       []string
	Head      []string
	Gets      []string
	Heads     []string
	Post      []string
	Put       []string
	Delete    []string
	CreatedAt *gtime.Time
	Detail    string

	RowKeyGen string // 主键生成策略
	RowKey    string
	FieldsGet map[string]FieldsGetValue
	Executor  string
}

func (a *Access) GetFieldsGetOutByRole(role string) []string {
	var fieldsMap map[string]string

	if val, exists := a.FieldsGet[role]; exists {
		fieldsMap = val.Out
	} else {
		fieldsMap = a.FieldsGet["default"].Out
	}
	return lo.Keys(fieldsMap)
}

func (a *Access) GetFieldsGetInByRole(role string) map[string][]string {
	var inFieldsMap map[string][]string

	if val, exists := a.FieldsGet[role]; exists {
		inFieldsMap = val.In
	} else {
		inFieldsMap = a.FieldsGet["default"].In
	}

	return inFieldsMap
}

func loadAccessMap() {
	_accessMap := make(map[string]Access)

	var accessList []Access

	db := g.DB()

	err := db.Model(config.TableAccess).Scan(&accessList)
	if err != nil {
		panic(err)
	}

	type AccessExt struct {
		RowKey    string
		FieldsGet map[string]FieldsGetValue
	}

	for _, access := range accessList {
		name := access.Alias
		if name == "" {
			name = access.Name
		}
		_accessMap[name] = access
	}

	accessMap = _accessMap
}

func GetAccess(tableAlias string, accessVerify bool) (*Access, error) {
	tableAlias, _ = util.ParseNodeKey(tableAlias)
	access, ok := accessMap[tableAlias]

	if !ok {
		if accessVerify {
			return nil, gerror.Newf("access[%s]: 404", tableAlias)
		}
		return &Access{
			Debug: 0,
			Name:  tableAlias,
			Alias: tableAlias,
		}, nil
	}

	return &access, nil
}

func GetAccessRole(table string, method string) ([]string, string, error) {
	access, ok := accessMap[table]

	if !ok {
		return nil, "", gerror.Newf("access[%s]: 404", table)
	}

	switch method {
	case http.MethodGet:
		return access.Get, access.Name, nil
	case http.MethodHead:
		return access.Head, access.Name, nil
	case http.MethodPost:
		return access.Post, access.Name, nil
	case http.MethodPut:
		return access.Put, access.Name, nil
	case http.MethodDelete:
		return access.Delete, access.Name, nil
	}

	return []string{}, access.Name, nil
}

func Init() {

	initTable()

	Reload()
}

// Reload 重载刷新配置
func Reload() {
	loadAccessMap()
	loadRequestMap()
	loadTableMeta()
}

// initTable 暂时先这样吧
func initTable() {
	sql_access := "CREATE TABLE `_access` (\n  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,\n  `debug` tinyint(4) NOT NULL DEFAULT '0' COMMENT '是否为调试表，只允许在开发环境使用，测试和线上环境禁用',\n  `name` varchar(50) NOT NULL COMMENT '实际表名，例如 apijson_user',\n  `alias` varchar(20) DEFAULT NULL COMMENT '外部调用的表别名，例如 User',\n  `get` varchar(100) NOT NULL DEFAULT '[\"UNKNOWN\", \"LOGIN\", \"CONTACT\", \"CIRCLE\", \"OWNER\", \"ADMIN\"]' COMMENT '允许 get 的角色列表，例如 [\"LOGIN\", \"CONTACT\", \"CIRCLE\", \"OWNER\"]\\n用 JSON 类型不能设置默认值，反正权限对应的需求是明确的，也不需要自动转 JSONArray。\\nTODO: 直接 LOGIN,CONTACT,CIRCLE,OWNER 更简单，反正是开发内部用，不需要复杂查询。',\n  `head` varchar(100) NOT NULL DEFAULT '[\"UNKNOWN\", \"LOGIN\", \"CONTACT\", \"CIRCLE\", \"OWNER\", \"ADMIN\"]' COMMENT '允许 head 的角色列表，例如 [\"LOGIN\", \"CONTACT\", \"CIRCLE\", \"OWNER\"]',\n  `gets` varchar(100) NOT NULL DEFAULT '[\"LOGIN\", \"CONTACT\", \"CIRCLE\", \"OWNER\", \"ADMIN\"]' COMMENT '允许 gets 的角色列表，例如 [\"LOGIN\", \"CONTACT\", \"CIRCLE\", \"OWNER\"]',\n  `heads` varchar(100) NOT NULL DEFAULT '[\"LOGIN\", \"CONTACT\", \"CIRCLE\", \"OWNER\", \"ADMIN\"]' COMMENT '允许 heads 的角色列表，例如 [\"LOGIN\", \"CONTACT\", \"CIRCLE\", \"OWNER\"]',\n  `post` varchar(100) NOT NULL DEFAULT '[\"OWNER\", \"ADMIN\"]' COMMENT '允许 post 的角色列表，例如 [\"LOGIN\", \"CONTACT\", \"CIRCLE\", \"OWNER\"]',\n  `put` varchar(100) NOT NULL DEFAULT '[\"OWNER\", \"ADMIN\"]' COMMENT '允许 put 的角色列表，例如 [\"LOGIN\", \"CONTACT\", \"CIRCLE\", \"OWNER\"]',\n  `delete` varchar(100) NOT NULL DEFAULT '[\"OWNER\", \"ADMIN\"]' COMMENT '允许 delete 的角色列表，例如 [\"LOGIN\", \"CONTACT\", \"CIRCLE\", \"OWNER\"]',\n  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',\n  `detail` varchar(1000) DEFAULT NULL,\n  `row_key` varchar(32) DEFAULT NULL COMMENT '@ext 关联主键字段名,联合主键时使用,分割',\n  `fields_get` json DEFAULT NULL COMMENT '@ext get查询时字段配置',\n  `row_key_gen` varchar(255) DEFAULT NULL,\n  `executor` varchar(32) DEFAULT NULL COMMENT '执行器name',\n  PRIMARY KEY (`id`) USING BTREE,\n  UNIQUE KEY `alias_UNIQUE` (`alias`) USING BTREE\n) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8 ROW_FORMAT=DYNAMIC COMMENT='权限配置(必须)';"
	sql_request := "CREATE TABLE `_request` (\n  `id` int(10) NOT NULL AUTO_INCREMENT COMMENT '唯一标识',\n  `debug` tinyint(4) NOT NULL DEFAULT '0' COMMENT '是否为 DEBUG 调试数据，只允许在开发环境使用，测试和线上环境禁用：0-否，1-是。',\n  `version` tinyint(4) NOT NULL DEFAULT '1' COMMENT 'GET,HEAD可用任意结构访问任意开放内容，不需要这个字段。\\n其它的操作因为写入了结构和内容，所以都需要，按照不同的version选择对应的structure。\\n\\n自动化版本管理：\\nRequest JSON最外层可以传  “version”:Integer 。\\n1.未传或 <= 0，用最新版。 “@order”:”version-“\\n2.已传且 > 0，用version以上的可用版本的最低版本。 “@order”:”version+”, “version{}”:”>={version}”',\n  `method` varchar(10) DEFAULT 'GETS' COMMENT '只限于GET,HEAD外的操作方法。',\n  `tag` varchar(20) NOT NULL COMMENT '标签',\n  `structure` json NOT NULL COMMENT '结构。\\nTODO 里面的 PUT 改为 UPDATE，避免和请求 PUT 搞混。',\n  `detail` varchar(10000) DEFAULT NULL COMMENT '详细说明',\n  `created_at` datetime DEFAULT CURRENT_TIMESTAMP,\n  `exec_queue` varchar(255) DEFAULT NULL COMMENT '@ext 节点执行顺序 执行队列, 因为请求的结构是确定的, 所以固定住节点的执行顺序,不用每次计算',\n  `executor` json DEFAULT NULL COMMENT '执行器映射 格式为Tag:executor;Tag2:executor 未配置为default',\n  PRIMARY KEY (`id`) USING BTREE\n) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8 ROW_FORMAT=DYNAMIC COMMENT='请求参数校验配置(必须)。\\r\\n最好编辑完后删除主键，这样就是只读状态，不能随意更改。需要更改就重新加上主键。\\r\\n\\r\\n每次启动服务器时加载整个表到内存。\\r\\n这个表不可省略，model内注解的权限只是客户端能用的，其它可以保证即便服务端代码错误时也不会误删数据。';"
	sql_function := "CREATE TABLE `_function` (\n  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,\n  `debug` tinyint(4) NOT NULL DEFAULT '0' COMMENT '是否为 DEBUG 调试数据，只允许在开发环境使用，测试和线上环境禁用：0-否，1-是。',\n  `userId` bigint(20) NOT NULL COMMENT '管理员用户Id',\n  `name` varchar(50) NOT NULL COMMENT '方法名',\n  `arguments` varchar(100) DEFAULT NULL COMMENT '参数列表，每个参数的类型都是 String。\\n用 , 分割的字符串 比 [JSONArray] 更好，例如 array,item ，更直观，还方便拼接函数。',\n  `demo` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL COMMENT '可用的示例。\\nTODO 改成 call，和返回值示例 back 对应。',\n  `detail` varchar(1000) NOT NULL COMMENT '详细描述',\n  `type` varchar(50) NOT NULL DEFAULT 'Object' COMMENT '返回值类型。TODO RemoteFunction 校验 type 和 back',\n  `version` tinyint(4) NOT NULL DEFAULT '0' COMMENT '允许的最低版本号，只限于GET,HEAD外的操作方法。\\nTODO 使用 requestIdList 替代 version,tag,methods',\n  `tag` varchar(20) DEFAULT NULL COMMENT '允许的标签.\\nnull - 允许全部\\nTODO 使用 requestIdList 替代 version,tag,methods',\n  `methods` varchar(50) DEFAULT NULL COMMENT '允许的操作方法。\\nnull - 允许全部\\nTODO 使用 requestIdList 替代 version,tag,methods',\n  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',\n  `back` varchar(45) DEFAULT NULL COMMENT '返回值示例',\n  PRIMARY KEY (`id`) USING BTREE\n) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8 ROW_FORMAT=DYNAMIC COMMENT='远程函数。强制在启动时校验所有demo是否能正常运行通过';"

	ctx := gctx.New()

	tables, err := g.DB().Tables(ctx)
	if err != nil {
		panic(err)
	}
	if !lo.Contains(tables, "_access") {
		g.DB().Exec(ctx, sql_access)
	}
	if !lo.Contains(tables, "_request") {
		g.DB().Exec(ctx, sql_request)
	}
	if !lo.Contains(tables, "_function") {
		g.DB().Exec(ctx, sql_function)
	}

}
