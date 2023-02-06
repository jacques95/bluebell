package mysql

import (
	"bluebell/settings"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

var db *sqlx.DB

func Init(cfg *settings.MySQLConfig) (err error) {
	dns := fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=true",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.DbName,
	)

	db, err := sqlx.Connect("mysql", dns)
	if err != nil {
		zap.L().Error("connect DB failed", zap.Error(err))
	}
	db.SetMaxOpenConns(cfg.MaxCons)
	db.SetMaxIdleConns(cfg.MaxIdles)
	return
}

func Close() {
	_ = db.Close()
}
