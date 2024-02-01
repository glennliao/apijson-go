package action

import (
	"context"

	"github.com/samber/lo"

	"github.com/glennliao/apijson-go/config"
	"github.com/glennliao/apijson-go/consts"
	"github.com/glennliao/apijson-go/model"
	"github.com/gogf/gf/v2/util/gconv"
)

// CheckTag
// check tag is valid
func CheckTag(req model.Map, method string, requestCfg *config.ActionConfig) (*config.RequestConfig, error) {
	_tag, ok := req[consts.Tag]
	if !ok {
		return nil, consts.ErrNoTag
	}

	tag := gconv.String(_tag)
	version := req[consts.Version]

	request, err := requestCfg.GetRequest(tag, method, gconv.String(version))
	if err != nil {
		return nil, err
	}

	return request, nil
}

// CheckRoleAccess
// check role access
func CheckRoleAccess(ctx context.Context, node *Node, accessRoles []string) error {
	roleFunc := node.Action.ActionConfig.DefaultRoleFunc()

	role, err := roleFunc(ctx, config.RoleReq{
		AccessName: node.TableName,
		Method:     node.Action.Method,
		NodeRole:   node.Role,
	})
	if err != nil {
		return err
	}

	if role == consts.DENY {
		return consts.NewDenyErr(node.Key, node.Role)
	}

	node.Role = role

	if !lo.Contains(accessRoles, role) {
		return consts.NewNoAccessErr(node.Key, node.Role)
	}

	return nil
}

func CheckReqStructure(req model.Map, structure *config.Structure) error {
	// must
	for _, key := range structure.Must {
		if _, exists := req[key]; !exists {
			return consts.NewStructureKeyNoFoundErr(key)
		}
	}

	// refuse
	if len(structure.Refuse) > 0 && structure.Refuse[0] == "!" {
		if len(structure.Must) == 0 {
			return consts.NewValidStructureErr("must set 'MUST' when 'REFUSE' is !")
		}
		for key := range req {
			if !lo.Contains(structure.Must, key) {
				return consts.NewValidStructureErr("can't has " + key)
			}
		}
	} else {
		for _, key := range structure.Refuse {
			if _, exists := req[key]; exists {
				return consts.NewValidStructureErr("refuse " + key)
			}
		}
	}

	return nil
}
