package query

import (
	"time"

	"github.com/gogf/gf/v2/frame/g"
)

func (n *Node) fetch() {
	defer func() {
		n.finish = true
		n.endAt = time.Now()
		if n.queryContext.PrintProcessLog {
			g.Log().Debugf(n.ctx, "【node】(%s) <fetch-endAt> ", n.Path)
		}
	}()

	if n.queryContext.PrintProcessLog {
		g.Log().Debugf(n.ctx, "【node】(%s) <fetch> hasFinish: 【%v】", n.Path, n.finish)
	}

	if n.finish {
		g.Log().Error(n.ctx, "再次执行", n.Path)
		return
	}

	if n.err != nil {
		return
	}

	n.nodeHandler.fetch()
}
