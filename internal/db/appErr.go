package db

import (
	"log"
	"time"
)

type AppErr struct {
	ID        int       `gorm:"primaryKey;autoIncrement"`
	Msg       string    `gorm:"not null"`
	CreatedAt time.Time `gorm:"not null"`
}

func InsertErr(msg string) {
	e := AppErr{
		Msg:       msg,
		CreatedAt: time.Now(),
	}

	result := _db.Create(&e)
	if result.Error != nil {
		log.Fatalf("db cannot insert Err: %s, Err time %s", result.Error.Error(), e.CreatedAt)
	}
}

func DeleteErrAfter(t time.Time) {
	result := _db.Delete(&AppErr{}, "created_at >= ?", t)
	if result.Error != nil {
		log.Fatalf("db cannot delete Err: %s", result.Error.Error())
	}
}

func GetErrAfter(t time.Time) []AppErr {
	var errs []AppErr
	result := _db.Where("created_at >= ?", t).Find(&errs)
	if result.Error != nil {
		log.Fatalf("db cannot get Err: %s", result.Error.Error())
	}
	return errs
}
