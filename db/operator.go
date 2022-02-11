package db

import (
	"github.com/latifrons/commongo/utilfuncs"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"time"
)

type DbOperator struct {
	Source string
	Db     *gorm.DB
}

func (d *DbOperator) InitDefault() {
	newLogger := logger.New(
		logrus.StandardLogger(),
		logger.Config{
			SlowThreshold:             time.Second, // Slow SQL threshold
			LogLevel:                  logger.Info, // Log level
			IgnoreRecordNotFoundError: true,        // Ignore ErrRecordNotFound berror for logger
			Colorful:                  true,        // Disable color
		},
	)

	db, err := gorm.Open(mysql.Open(d.Source), &gorm.Config{PrepareStmt: true, Logger: newLogger})
	//Db, err := gorm.Open(sqlite.Open(d.Source), &gorm.Config{PrepareStmt: true})
	if err != nil {
		logrus.WithError(err).Fatal("failed to connect to database")
	}
	d.Db = db

	sqlDB, err := d.Db.DB()

	// SetMaxIdleConns sets the maximum number of connections in the idle connection pool.
	sqlDB.SetMaxIdleConns(10)

	// SetMaxOpenConns sets the maximum number of open connections to the database.
	sqlDB.SetMaxOpenConns(100)

	// SetConnMaxLifetime sets the maximum amount of time a connection may be reused.
	sqlDB.SetConnMaxLifetime(time.Hour)

	for _, dbo := range []struct {
		Obj  interface{}
		Name string
	}{
		//{Obj: &MarketPlaceUser{}, Name: "MarketPlaceUser"},
	} {
		err = d.Db.AutoMigrate(dbo.Obj)
		utilfuncs.PanicIfError(err, "failed to migrate "+dbo.Name)
	}

	logrus.Info("Db inited")
}
