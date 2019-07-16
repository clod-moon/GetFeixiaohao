package main

import (
	"log"
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"sync"
	"time"
	"github.com/clod-moon/goconf"
	"strconv"
)

var (

	wg sync.WaitGroup

	DB *gorm.DB

	username string = "root"
	password string = "root"
	dbName   string = "chun"
	host     string = "192.168.0.102"
	port     int    = 3306

	biAmount  string = "100"
	updatePreSecond time.Duration = 5

	tcpSreverPort int = 9088
	httpSreverPort int = 9080
)

func init() {

	readConf()

	DB.SingularTable(true)
	gorm.DefaultTableNameHandler = func(db *gorm.DB, defaultTableName string) string {
		return  defaultTableName
	}
	if !DB.HasTable(&Market{}) {
		err := DB.Set("gorm:table_options", "ENGINE=InnoDB DEFAULT CHARSET=utf8").CreateTable(&Market{}).Error
		if err != nil {
			panic(err)
		}
	}
}

func readConf(){
	conf := iniconf.InitConfig("./conf.ini")

	username = conf.GetValue("mysql","username")

	password = conf.GetValue("mysql","password")
	dbName= conf.GetValue("mysql","dbName")
	host= conf.GetValue("mysql","host")
	port ,err :=strconv.Atoi(conf.GetValue("mysql","port"))
	if err != nil{
		log.Fatal("mysql port conf err!")
	}
	biAmount = conf.GetValue("market","pageSize")
	second ,err := strconv.Atoi(conf.GetValue("market","updatePreSecond"))
	if err != nil{
		log.Fatal("UpdatePreSecond conf err!")
	}
	updatePreSecond = time.Duration(second)

	tcpSreverPort,err = strconv.Atoi(conf.GetValue("tcpServer","port"))
	if err != nil{
		log.Fatal("tcpServer port conf err!")
	}
	httpSreverPort,err = strconv.Atoi(conf.GetValue("httpServer","port"))
	if err != nil{
		log.Fatal("HttpServer port conf err!")
	}

	mysqlstr := fmt.Sprintf("%s:%s@(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local", username, password, host, port, dbName)
	DB, err = gorm.Open("mysql", mysqlstr)
	if err != nil {
		log.Fatalf(" gorm.Open.err: %v", err)
	}
}

const (
	teststr = `{"code": 200,"msg": "ok","data": [{"name": "BTC","current_price": "100","change_percent": "100","tee":"safaf"}]}`
)


func main() {
	// 设置路由，如果访问/，则调用index方法
	wg.Add(1)
	go UpdateMarket()

	wg.Add(1)
	go ServerSocket()

	HttpServer()

	wg.Wait()

}
