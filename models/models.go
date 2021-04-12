package models

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"go-gin-example/pkg/setting"
	"log"
	"time"
)

var db *gorm.DB

type Model struct {
	ID int	`gorm:"primary_key" json:"id"`
	CreatedOn	int 	`json:"created_on"`
	ModifiedOn 	int 	`json:"modified_on"`
}

func Setup() {
	var err error
	db, err = gorm.Open(setting.DatabaseSetting.Type, fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True&loc=Local",
		setting.DatabaseSetting.User,
		setting.DatabaseSetting.Password,
		setting.DatabaseSetting.Host,
		setting.DatabaseSetting.Name))

	if err != nil {
		log.Fatalf("models.Setup err: %v", err)
	}

	gorm.DefaultTableNameHandler = func(db *gorm.DB, defaultTableName string) string {
		return setting.DatabaseSetting.TablePrefix + defaultTableName
	}

	db.SingularTable(true)
	db.Callback().Create().Replace("gorm:update_time_stamp", updateTimeStampForCreateCallback)
	db.Callback().Update().Replace("gorm:update_time_stamp", updateTimeStampForUpdateCallback)

	db.DB().SetMaxIdleConns(10)
	db.DB().SetMaxOpenConns(100)
}

func CloseDB()  {
	defer db.Close()
}

func updateTimeStampForCreateCallback(scope *gorm.Scope){
	if !scope.HasError(){
		nowTime := time.Now().Unix()
		if createTimeField, ok := scope.FieldByName("CreatedOn"); ok{
			// 如何字段为空则给该字段设置nowTime时间戳
			if createTimeField.IsBlank{
				createTimeField.Set(nowTime)
			}
		}

		if modifyTimeField, ok := scope.FieldByName("ModifiedOn"); ok{
			if modifyTimeField.IsBlank{
				modifyTimeField.Set(nowTime)
			}
		}
	}
}

func updateTimeStampForUpdateCallback(scope *gorm.Scope){
	if !scope.HasError(){
		nowTime := time.Now().Unix()
		if _, ok := scope.Get("gorm:update_column"); ok{
			scope.SetColumn("ModifiedOn", nowTime)
		}
	}
}
