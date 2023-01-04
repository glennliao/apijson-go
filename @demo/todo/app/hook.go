package app

import (
	"github.com/glennliao/apijson-go/action"
)

func init() {
	action.RegHook(action.Hook{
		For:        "*",
		BeforeExec: nil,
		AfterExec:  nil,
		BeforeDo:   nil,
		AfterDo:    nil,
	})

	//action.RegHook(action.Hook{
	//	For: "*",
	//	BeforeExec: func(ctx context.Context, n *action.Node, method string) error {
	//		g.Log().Debug(ctx, " iam BeforeExec", n.TableName, method)
	//		return nil
	//	},
	//	AfterExec: func(ctx context.Context, n *action.Node, method string) error {
	//		g.Log().Debug(ctx, " iam AfterExec", n.TableName, method)
	//		return nil
	//	},
	//	BeforeDo: func(ctx context.Context, n *action.Node, method string) error {
	//		g.Log().Debug(ctx, " iam BeforeDo", n.TableName, method)
	//		return nil
	//	},
	//	AfterDo: func(ctx context.Context, n *action.Node, method string) error {
	//		g.Log().Debug(ctx, " iam AfterDo", n.TableName, method)
	//		return nil
	//	},
	//})
	//
	//action.RegHook(action.Hook{
	//	For: "t_todox",
	//	BeforeExec: func(ctx context.Context, n *action.Node, method string) error {
	//		g.Log().Debug(ctx, " iam BeforeExec For todo", n.TableName, method)
	//		return nil
	//	},
	//	AfterExec: func(ctx context.Context, n *action.Node, method string) error {
	//		g.Log().Debug(ctx, " iam AfterExec For todo", n.TableName, method)
	//		return nil
	//	},
	//	BeforeDo: func(ctx context.Context, n *action.Node, method string) error {
	//		g.Log().Debug(ctx, " iam BeforeDo For todo", n.TableName, method)
	//		return nil
	//	},
	//	AfterDo: func(ctx context.Context, n *action.Node, method string) error {
	//		g.Log().Debug(ctx, " iam AfterDo For todo", n.TableName, method)
	//		return nil
	//	},
	//})
}
