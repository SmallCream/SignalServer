package main

import (
	_ "SignalServer/routers"
	_ "SignalServer/utils"
	"github.com/astaxie/beego"
)

//func init()  {
//	utils.CacheConnect()
//	b := utils.Cache.Get("k1")
//	fmt.Println(b)
//
//}
func main() {
	beego.Run()
}

