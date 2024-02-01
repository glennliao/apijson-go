package action

import (
	"context"
	"strings"

	"github.com/glennliao/apijson-go/consts"
	"github.com/glennliao/apijson-go/util"

	"github.com/glennliao/apijson-go/config"
	"github.com/glennliao/apijson-go/model"
)

type Node struct {
	ctx context.Context
	// key为原始请求的key
	Key        string
	AccessName string
	req        model.Map
	Ret        model.Map

	Action *Action
	Role   string

	Data            model.Map
	Where           model.Map
	AccessCondition *config.ConditionRet

	TableName string
	RowKey    string
	structure *config.Structure
	executor  string

	keyNode map[string]*Node
}

func newNode(a *Action, key string, req model.Map, structure *config.Structure) *Node {
	n := &Node{
		ctx:        a.ctx,
		Action:     a,
		keyNode:    a.keyNode,
		Key:        key,
		AccessName: key,
		req:        req,
		structure:  structure,
		executor:   a.tagRequestConfig.Executor[key],
		Data:       make(model.Map),
		Where:      make(model.Map),
	}

	if strings.HasSuffix(key, consts.ListKeySuffix) {
		n.AccessName = util.RemoveSuffix(key, consts.ListKeySuffix)
	}

	return n
}
