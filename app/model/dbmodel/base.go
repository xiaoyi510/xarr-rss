package dbmodel

import (
	"context"
	"encoding/json"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"time"
)

type Model struct {
	Id        int       `gorm:"primarykey;auto_increment;" json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type JsonStringArray []string

func (loc JsonStringArray) GormDataType() string {
	return "json"
}

func (loc JsonStringArray) GormValue(ctx context.Context, db *gorm.DB) clause.Expr {
	marshal, err := json.Marshal(loc)
	if err != nil {
		panic("GG")
		return clause.Expr{}
	}

	return clause.Expr{
		SQL:  "?",
		Vars: []interface{}{string(marshal)},
	}
}

// Scan 方法实现了 sql.Scanner 接口
func (loc *JsonStringArray) Scan(v interface{}) error {
	// Scan a value into struct from database driver
	err := json.Unmarshal([]byte(v.(string)), loc)
	if err != nil {
		return err
	}
	return nil
}
