package db

import (
	"fmt"
	"sync"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"github.com/nhatminhk63j/uetvoting/config"
)

var (
	db *gorm.DB
)

func newDatabase() *gorm.DB {
	var once sync.Once
	once.Do(func() {
		cfg := config.LoadMysqlConfig()
		_db, err := gorm.Open(mysql.Open(cfg.ToURI()), &gorm.Config{})
		if err != nil {
			panic(fmt.Errorf("cannot connect database: %+v", err))
		}
		db = _db
	})
	return db
}
