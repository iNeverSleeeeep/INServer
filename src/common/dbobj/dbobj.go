package dbobj

import (
	"INServer/src/common/logger"
	"INServer/src/proto/etc"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type (
	DBObject struct {
		config   *etc.Database
		database *sql.DB
	}
)

func New() *DBObject {
	d := new(DBObject)
	return d
}

func (d *DBObject) Open(config *etc.Database, dbname string) {
	dsn := fmt.Sprintf("%s:%s@%s(%s:%d)/%s", config.UserName, config.Password, "tcp", config.IP, 3306, dbname)
	database, err := sql.Open("mysql", dsn)
	if err != nil {
		logger.Fatal(err)
	}
	database.SetConnMaxLifetime(time.Duration(config.ConnMaxLifetime) * time.Second)
	database.SetMaxOpenConns(int(config.MaxOpenConns))
	database.SetMaxIdleConns(int(config.MaxIdleConns))
	d.config = config
	d.database = database
}

func (d *DBObject) Close() {
	err := d.database.Close()
	if err != nil {
		logger.Debug(err)
	}
}

func (d *DBObject) DB() *sql.DB {
	return d.database
}
