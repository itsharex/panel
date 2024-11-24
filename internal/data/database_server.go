package data

import (
	"errors"
	"fmt"
	"slices"
	"strings"

	"github.com/TheTNB/panel/internal/app"
	"github.com/TheTNB/panel/internal/biz"
	"github.com/TheTNB/panel/internal/http/request"
	"github.com/TheTNB/panel/pkg/db"
)

type databaseServerRepo struct{}

func NewDatabaseServerRepo() biz.DatabaseServerRepo {
	return &databaseServerRepo{}
}

func (d databaseServerRepo) Count() (int64, error) {
	var count int64
	if err := app.Orm.Model(&biz.DatabaseServer{}).Count(&count).Error; err != nil {
		return 0, err
	}

	return count, nil
}

func (d databaseServerRepo) List(page, limit uint) ([]*biz.DatabaseServer, int64, error) {
	var databaseServer []*biz.DatabaseServer
	var total int64
	err := app.Orm.Model(&biz.DatabaseServer{}).Order("id desc").Count(&total).Offset(int((page - 1) * limit)).Limit(int(limit)).Find(&databaseServer).Error
	return databaseServer, total, err
}

func (d databaseServerRepo) Get(id uint) (*biz.DatabaseServer, error) {
	databaseServer := new(biz.DatabaseServer)
	if err := app.Orm.Where("id = ?", id).First(databaseServer).Error; err != nil {
		return nil, err
	}

	return databaseServer, nil
}

func (d databaseServerRepo) Create(req *request.DatabaseServerCreate) error {
	switch biz.DatabaseType(req.Type) {
	case biz.DatabaseTypeMysql:
		if _, err := db.NewMySQL(req.Username, req.Password, fmt.Sprintf("%s:%d", req.Host, req.Port)); err != nil {
			return err
		}
	case biz.DatabaseTypePostgresql:
		if _, err := db.NewPostgres(req.Username, req.Password, req.Host, req.Port, ""); err != nil {
			return err
		}
	case biz.DatabaseTypeRedis:
		if _, err := db.NewRedis(req.Username, req.Password, fmt.Sprintf("%s:%d", req.Host, req.Port)); err != nil {
			return err
		}

	}

	databaseServer := &biz.DatabaseServer{
		Name:     req.Name,
		Type:     biz.DatabaseType(req.Type),
		Host:     req.Host,
		Port:     req.Port,
		Username: req.Username,
		Password: req.Password,
		Remark:   req.Remark,
	}

	return app.Orm.Create(databaseServer).Error
}

func (d databaseServerRepo) Update(req *request.DatabaseServerUpdate) error {
	switch biz.DatabaseType(req.Type) {
	case biz.DatabaseTypeMysql:
		if _, err := db.NewMySQL(req.Username, req.Password, fmt.Sprintf("%s:%d", req.Host, req.Port)); err != nil {
			return err
		}
	case biz.DatabaseTypePostgresql:
		if _, err := db.NewPostgres(req.Username, req.Password, req.Host, req.Port, ""); err != nil {
			return err
		}
	case biz.DatabaseTypeRedis:
		if _, err := db.NewRedis(req.Username, req.Password, fmt.Sprintf("%s:%d", req.Host, req.Port)); err != nil {
			return err
		}

	}

	return app.Orm.Model(&biz.DatabaseServer{}).Where("id = ?", req.ID).Select("*").Updates(&biz.DatabaseServer{
		Name:     req.Name,
		Type:     biz.DatabaseType(req.Type),
		Host:     req.Host,
		Port:     req.Port,
		Username: req.Username,
		Password: req.Password,
		Remark:   req.Remark,
	}).Error
}

func (d databaseServerRepo) Delete(id uint) error {
	ds := new(biz.DatabaseServer)
	if err := app.Orm.Where("id = ?", id).First(ds).Error; err != nil {
		return err
	}

	if slices.Contains([]string{"local_mysql", "local_postgresql", "local_redis"}, ds.Name) && !app.IsCli {
		return errors.New("can't delete " + ds.Name + ", if you must delete it, please uninstall " + strings.TrimPrefix(ds.Name, "local_"))
	}

	return app.Orm.Delete(&biz.DatabaseServer{}, id).Error
}
