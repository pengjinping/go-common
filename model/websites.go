package model

import (
	"context"
)

type Websites struct {
	BaseModel
	Id   uint   `gorm:"primarykey"`
	Name string `json:"name" gorm:"comment:name"`
	UUID string `json:"uuid" gorm:"comment:uuid"`
}

func (w *Websites) getByUUID(ctx context.Context, uuid string) {
	// db := w.getConn(ctx)
}
