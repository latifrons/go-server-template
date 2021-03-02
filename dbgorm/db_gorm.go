package dbgorm

import (
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"time"
)
import "gorm.io/driver/sqlite"

type DbOperator struct {
	Source string
	db     *gorm.DB
}

func (d *DbOperator) InitDefault() {
	//db, err := gorm.Open(mysql.Open(d.Source), &gorm.Config{PrepareStmt: true})
	db, err := gorm.Open(sqlite.Open(d.Source), &gorm.Config{PrepareStmt: true})
	if err != nil {
		logrus.WithError(err).Fatal("failed to connect to database")
	}
	d.db = db

	sqlDB, err := db.DB()

	// SetMaxIdleConns sets the maximum number of connections in the idle connection pool.
	sqlDB.SetMaxIdleConns(10)

	// SetMaxOpenConns sets the maximum number of open connections to the database.
	sqlDB.SetMaxOpenConns(100)

	// SetConnMaxLifetime sets the maximum amount of time a connection may be reused.
	sqlDB.SetConnMaxLifetime(time.Hour)

	tables := []interface{}{
		//&DbContact{},
	}
	for _, table := range tables {
		err := db.AutoMigrate(table)
		if err != nil {
			logrus.WithError(err).Fatal("failed to migrate")
		}
	}
	logrus.Info("db inited")
}

func (d *DbOperator) Begin() {
	d.db.Begin()
}

func (d *DbOperator) Commit() {
	d.db.Commit()
}
func (d *DbOperator) Rollback() {
	d.db.Rollback()
}
