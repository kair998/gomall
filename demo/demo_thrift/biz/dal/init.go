package dal

import (
	"github.com/kair998/gomall/demo/demo_thrift/biz/dal/mysql"
	"github.com/kair998/gomall/demo/demo_thrift/biz/dal/redis"
)

func Init() {
	redis.Init()
	mysql.Init()
}
