package context

import (
	. "github.com/leyle/ginbase/consolelog"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Ds struct {
	Db *gorm.DB
}

func NewDs(dbFile string) *Ds {
	conn, err := gorm.Open(sqlite.Open(dbFile), &gorm.Config{})
	if err != nil {
		Logger.Errorf("", "Open Sqlite3 failed, %s", err.Error())
		panic(err)
	}

	db := &Ds{
		Db: conn,
	}
	return db
}

func (d *Ds) Close() {
	sqlDb, _ := d.Db.DB()
	err := sqlDb.Close()
	if err != nil {
		Logger.Errorf("", "close sqlite failed, %s", err.Error())
	}
}
