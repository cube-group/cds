package models

import(
	"sys/core"
	"time"
)

//fcds.FLog表模型
type FLog struct {
	ID            uint `gorm:"AUTO_INCREMENT"`
	Uid           uint
	Username      string `gorm:"type:varchar(50)"`
	Route         string `gorm:"type:varchar(100)"`
	CreateTime    time.Time `gorm:"datetime"`
	content       string `gorm:"type:varchar(100)"`
}

func NewFLog() *FLog {
    return new(FLog)
}

func Create(log *FLog) error{
    db, err := core.Mysql()
    if err != nil {
    	return err
    }
    defer core.PoolRollBack(db)
    db.Save(log)
    return nil
}