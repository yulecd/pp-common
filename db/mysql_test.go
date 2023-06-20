package db

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/yulecd/pp-common/plog"
)

func TestConn(t *testing.T) {
	plog.Init("error")
	platoMysql, err := InitMysqlClient(MysqlConf{
		DataBase:        "pay_server",
		Addr:            "10.0.2.26:3306",
		User:            "zhangbing",
		Password:        "rDFSmIXD6M",
		MaxIdleConns:    5,
		MaxOpenConns:    2,
		ConnMaxLifeTime: 100000,
		ConnTimeOut:     time.Second * 1,
		WriteTimeOut:    time.Second * 1,
		ReadTimeOut:     time.Second * 1,
		LogMode:         true,
		LogLevel:        Info,
		SlowThreshold:   time.Second,
	})

	if err != nil {
		fmt.Println("error select:", err.Error())
	}

	ctx := context.Background()
	ctx = context.WithValue(ctx, "logId", "1111111122222333333444445555")

	type aaa struct {
		MerchantId int64  `json:"merchant_id"`
		ApiKey     string `json:"api_key"`
	}

	var list []aaa

	if platoMysql != nil {
		if err = platoMysql.WithContext(ctx).Table("api_key").Where("merchant_id = ?", 2).Find(&list).Error; err != nil {

		}
		//fmt.Println("error select:", err.Error())
		fmt.Println(list)
	} else {
		fmt.Println("empty")
	}

	time.Sleep(time.Second * 5)
}
