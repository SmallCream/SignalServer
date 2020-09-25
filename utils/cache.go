package utils

import (
	"github.com/astaxie/beego"
	_ "github.com/astaxie/beego/cache/redis"
	"github.com/go-redis/redis"
)
//var Cache cache.Cache
//func CacheConnect ()  {
//	collectionName := beego.AppConfig.String("cache.collectionName")
//	conn := beego.AppConfig.String("cache.conn")
//	dbNum := beego.AppConfig.String("cache.dbNum")
//	password := beego.AppConfig.String("cache.password")
//	config := orm.Params{
//		"key": collectionName,
//		"conn": conn,
//		"dbNum": dbNum,
//		"password": password,
//	}
//	configStr, err := json.Marshal(config)
//	if err != nil {
//		logs.Error("redis配置模型转换失败")
//		return
//	}
//	Cache,err = cache.NewCache("redis",string(configStr))
//	if err != nil {
//		fmt.Println(err)
//		logs.Error("redis初始化失败")
//		return
//	}
//}

var RedisCache *redis.Client
func CacheConnect (){
	Addr := beego.AppConfig.String("cache.conn")
	dbNum, _ := beego.AppConfig.Int("cache.dbNum")
	Password := beego.AppConfig.String("cache.password")
	RedisCache = redis.NewClient(&redis.Options{
		Addr:     Addr, // use default Addr
		Password: Password,               // no password set
		DB:       dbNum,                // use default DB
	})

}

func init()  {
	CacheConnect()
	//RedisCache.SAdd("aaa","aaaa")
	//RedisCache.SAdd("aaa","vvvv")
	//RedisCache.SAdd("aaa","bbbb")
	//a,_:=RedisCache.SMembers("aaa").Result()
	//fmt.Println(a)


}
