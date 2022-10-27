package service

import (
	"context"
	"database/sql"
	. "github.com/glennliao/apijson-go/apijson/app/ind/internal/model"
	"github.com/glennliao/apijson-go/apijson/internal/dao"
	"github.com/glennliao/apijson-go/apijson/internal/do"
	"github.com/glennliao/apijson-go/apijson/util"
	"github.com/gogf/gf/v2/database/gdb"
	"strconv"
)

type requestService struct{}

var _request = requestService{}

func Request() *requestService {
	return &_request
}

func (s *requestService) List(ctx context.Context, in RequestListIn) (list []RequestListResult, total int, err error) {
	var (
		m, c = dao.Request.MC(ctx)
	)

	if total, err = m.Count(); total == 0 || err != nil {
		return
	}

	err = m.Page(in.PageNum, in.PageSize).
		OrderDesc(c.CreatedAt).
		Scan(&list)
	return
}

func (s *requestService) Get(ctx context.Context, in RequestGetIn) (row *RequestGetOut, err error) {
	var (
		d = dao.Request
	)
	err = d.WithRk(ctx, in.Id).Scan(&row)
	return
}

func (s *requestService) Add(ctx context.Context, in RequestAddIn) (out RequestAddOut, err error) {
	var (
		d = dao.Request
	)

	if err := util.NotExist(d.M(ctx).Where(do.Request{
		Version: in.Version,
		Method:  in.Method,
		Tag:     in.Tag,
	}), "数据已存在"); err != nil {
		return out, err
	}

	err = d.Transaction(ctx, func(ctx context.Context, tx *gdb.TX) error {
		ret, err := d.M(ctx).Insert(in)

		id, _ := ret.LastInsertId()
		out.Id = strconv.FormatInt(id, 10)
		return err
	})
	return
}

func (s *requestService) Update(ctx context.Context, in RequestUpdateIn) (ret sql.Result, err error) {
	var (
		d = dao.Request
	)

	err = d.Transaction(ctx, func(ctx context.Context, tx *gdb.TX) error {
		ret, err = d.WithRk(ctx, in.Id).Update(in)
		return err
	})
	return
}

func (s *requestService) Delete(ctx context.Context, ids []string) (ret sql.Result, err error) {
	var (
		d = dao.Request
	)
	err = d.Transaction(ctx, func(ctx context.Context, tx *gdb.TX) error {
		ret, err = d.WithRk(ctx, ids...).Delete()
		return err
	})
	return
}

// #d65be58308362f015cc2e27179c7bd0684296bfc:021026:52
