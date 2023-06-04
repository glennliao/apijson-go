package tables

import (
	"github.com/glennliao/table-sync/tablesync"
	"time"
)

type Access struct {
	tablesync.TableMeta `tableName:"_access" charset:"utf8mb4" comment:"权限配置"`
	Id                  uint32         `ddl:"primaryKey"`
	Debug               int8           `ddl:"not null;default:0;comment:是否调试表,开发环境可用"`
	Name                string         `ddl:"size:32;not null;comment:实际表名"`
	Alias               string         `ddl:"size:32;uniqueIndex;comment:表别名,外部调用"`
	Get                 string         `ddl:"size:128;not null;default:'LOGIN,OWNER,ADMIN';comment:允许get的权限角色列表"`
	Head                string         `ddl:"size:128;not null;default:'LOGIN,OWNER,ADMIN';comment:允许head的权限角色列表"`
	Gets                string         `ddl:"size:128;not null;default:'LOGIN,OWNER,ADMIN';comment:允许gets的权限角色列表"`
	Heads               string         `ddl:"size:128;not null;default:'LOGIN,OWNER,ADMIN';comment:允许heads的权限角色列表"`
	Post                string         `ddl:"size:128;not null;default:'LOGIN,OWNER,ADMIN';comment:允许post的权限角色列表"`
	Put                 string         `ddl:"size:128;not null;default:'LOGIN,OWNER,ADMIN';comment:允许put的权限角色列表"`
	Delete              string         `ddl:"size:128;not null;default:'LOGIN,OWNER,ADMIN';comment:允许delete的权限角色列表"`
	CreatedAt           *time.Time     `ddl:"notnull;default:CURRENT_TIMESTAMP;comment:创建时间"`
	Detail              string         `ddl:"size:512;"`
	RowKey              string         `ddl:"size:32;comment:(逻辑)主键字段名,联合主键使用,分割"`
	FieldsGet           map[string]any `ddl:"type:json;comment:get查询时字段配置"`
	RowKeyGen           string         `ddl:"comment:rowKey生成策略"`
	Executor            string         `ddl:"size:32;comment:执行器"`
}

type Request struct {
	tablesync.TableMeta `tableName:"_request" charset:"utf8mb4" comment:"请求参数校验配置"`
	Id                  uint32         `ddl:"primaryKey"`
	Debug               int8           `ddl:"not null;default:0;comment:是否调试,开发环境可用"`
	Tag                 string         `ddl:"not null;size:32;not null;comment:标签名(表别名)"`
	Version             string         `ddl:"not null;size:8;comment:版本号"`
	Method              string         `ddl:"not null;size:6;comment:请求方式"`
	Structure           map[string]any `ddl:"not null;type:json;comment:请求结构"`
	Detail              string         `ddl:"size:512;comment:描述说明"`
	CreatedAt           *time.Time     `ddl:"NOT NULL;comment:创建时间"`
	ExecQueue           string         `ddl:"size:512;comment:节点执行队列,使用,分割  请求结构确定的,不用每次计算依赖关系"`
	Executor            map[string]any `ddl:"type:json;comment:节点执行器,格式为Tag:executor;Tag2:executor 未配置为default"`
}

type Function struct {
	tablesync.TableMeta `tableName:"_function" charset:"utf8mb4" comment:"远程函数(暂未使用)"`
	Id                  uint32     `ddl:"primaryKey"`
	Debug               int8       `ddl:"not null;default:0;comment:是否调试,开发环境可用"`
	UserId              string     `ddl:"not null;comment:管理员id"`
	Name                string     `ddl:"size:50;comment:方法名"`
	Arguments           string     `ddl:"size:100;comment:参数类型列表"`
	Demo                string     `ddl:"size:256;comment:参数示例"`
	Type                string     `ddl:"size:16;comment:返回值类型"`
	Tag                 string     `ddl:"not null;size:32;not null;comment:标签名(表别名)"`
	Version             string     `ddl:"not null;size:8;comment:版本号"`
	Method              string     `ddl:"not null;size:5;comment:请求方式"`
	Detail              string     `ddl:"size:512;comment:描述说明"`
	CreatedAt           *time.Time `ddl:"NOT NULL;comment:创建时间"`
	Back                string     `ddl:"size:128;comment:返回值示例"`
}
