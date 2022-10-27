package service

import (
	"context"
	"database/sql"
	. "github.com/glennliao/apijson-go/apijson/app/ind/internal/model"
	"github.com/glennliao/apijson-go/apijson/internal/dao"
	"github.com/gogf/gf/v2/database/gdb"
	"strconv"
)

type accessService struct{}

var _access = accessService{}

func Access() *accessService {
	return &_access
}

func (s *accessService) List(ctx context.Context, in AccessListIn) (list []AccessListResult, total int, err error) {
	var (
		m, c = dao.Access.MC(ctx)
	)

	if total, err = m.Count(); total == 0 || err != nil {
		return
	}

	err = m.Page(in.PageNum, in.PageSize).
		OrderDesc(c.CreatedAt).
		Scan(&list)
	return
}

func (s *accessService) Get(ctx context.Context, id string) (row *AccessGetOut, err error) {
	var (
		d = dao.Access
	)
	err = d.WithRk(ctx, id).Scan(&row)
	return
}

func (s *accessService) Add(ctx context.Context, in AccessAddIn) (out AccessAddOut, err error) {
	var (
		d = dao.Access
	)
	err = d.Transaction(ctx, func(ctx context.Context, tx *gdb.TX) error {
		ret, err := d.M(ctx).Insert(in)

		id, _ := ret.LastInsertId()
		out.Id = strconv.FormatInt(id, 10)
		return err
	})
	return
}

func (s *accessService) Update(ctx context.Context, in AccessUpdateIn) (ret sql.Result, err error) {
	var (
		d = dao.Access
	)
	err = d.Transaction(ctx, func(ctx context.Context, tx *gdb.TX) error {
		ret, err = d.WithRk(ctx, in.Id).
			//Where(do.Access{Name: in.Name}).
			Update(in)
		return err
	})
	return
}

func (s *accessService) Delete(ctx context.Context, ids []string) (ret sql.Result, err error) {
	var (
		d = dao.Access
	)
	err = d.Transaction(ctx, func(ctx context.Context, tx *gdb.TX) error {
		ret, err = d.WithRk(ctx, ids...).Delete()
		return err
	})
	return
}

// #624550add5a789a56be127b25688f07d581d4f67:021026:52
