package database

import (
	"fmt"

	"gorm.io/gorm"
)

type TxFunc func(tx *gorm.DB) error

// 在db事务内执行一个/多个函数，返回可能出现的错误
// 【注意】请勿在事务中做以下操作，回滚时会导致一致性问题：
// 1. 对本连接之外的其它数据库连接，进行写操作
// 2. 写 redis / 内存缓存
// 3. 调接口修改远端数据
// 4. 修改全局变量
// 5. 其它类似行为
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
