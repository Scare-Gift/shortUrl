package db

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
)

var (
	Db  *gorm.DB
	err error
)

func init() {
	Db, err = gorm.Open(mysql.New(mysql.Config{
		DSN:                       "root:123456@tcp(localhost:3306)/shorturl?charset=utf8&parseTime=True&loc=Local",
		DefaultStringSize:         256,
		DisableDatetimePrecision:  true,
		DontSupportRenameIndex:    true,
		DontSupportRenameColumn:   true,
		SkipInitializeWithVersion: false,
	}), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}
	if Db.Error != nil {
		log.Fatal(Db.Error)
	}
}
