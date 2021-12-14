package database

import (
	"fmt"

	"gorm.io/gorm"
)

type TxFunc func(db *gorm.DB) error

// 在db事务内执行一个/多个函数
func Tx(db *gorm.DB, funcs ...TxFunc) (err error) {
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			err = fmt.Errorf("%v", err)
		}
	}()
	for _, f := range funcs {
		err = f(tx)
		if err != nil {
			tx.Rollback()
			return
		}
	}
	err = tx.Commit().Error
	return
}
