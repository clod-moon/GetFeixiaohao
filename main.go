package main

import (
	"log"
	"net/http"
	"io/ioutil"
	"encoding/json"
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"time"
)

var (
	DB *gorm.DB

	username string = "root"
	password string = "root"
	dbName   string = "chun"
	host     string = "192.168.1.224"
	port     int    = 3306
)

func init() {
	var err error
	mysqlstr := fmt.Sprintf("%s:%s@(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local", username, password, host, port, dbName)
	DB, err = gorm.Open("mysql", mysqlstr)
	if err != nil {
		log.Fatalf(" gorm.Open.err: %v", err)
	}

	DB.SingularTable(true)
	gorm.DefaultTableNameHandler = func(db *gorm.DB, defaultTableName string) string {
		return "sp_" + defaultTableName
	}
	if !DB.HasTable(&Market{}) {
		err := DB.Set("gorm:table_options", "ENGINE=InnoDB DEFAULT CHARSET=utf8").CreateTable(&Market{}).Error
		if err != nil {
			panic(err)
		}
	}
}

type Market struct {
	Id         int       `gorm:"primary_key;type:int(11);AUTO_INCREMENT`
	Name       string    `json:"name" gorm:"type:varchar(30);not null"`
	Price      float64   `json:"current_price" gorm:"type:float;not null;"`
	Rose       float64   `json:"change_percent" gorm:"type:float;not null;"`
	CreateTime time.Time `gorm:"type:datetime;not null;"`
	UpdateTime time.Time `gorm:"type:datetime;not null;"`
}

type FeixiaoData struct {
	Data []Market `json:"data"`
	Code int      `json:"code"`
	Msg  string   `json:"msg"`
}

const (
	teststr = `{"code": 200,"msg": "ok","data": [{"name": "BTC","current_price": "100","change_percent": "100","tee":"safaf"}]}`
)

func index(w http.ResponseWriter, r *http.Request) {
	// 往w里写入内容，就会在浏览器里输出

	w.Header().Set("Access-Control-Allow-Origin", "*")             //允许访问所有域
	w.Header().Add("Access-Control-Allow-Headers", "Content-Type") //header的类型
	w.Header().Set("content-type", "application/json")             //返回数据格式是json

	//rQuery := r.URL.Query()
	//body := getMarket(rQuery["get-size"][0])

	var feixiaoData FeixiaoData

	//fmt.Println(DB.Find(&Market{}).Value)

	var list []Market

	rows,err:=DB.Debug().Find(&list).Rows()

	if err != nil{
		w.Write([]byte(err.Error()))
		return
	}

	var m Market
	rows.Scan(&m.Id,&m.Name,&m.Price,&m.Rose)
	fmt.Println(m)


	//json.Unmarshal(body, &feixiaoData)

	ret, _ := json.Marshal(feixiaoData)

	w.Write(ret)

}

func getMarket(pageSize string) []byte {
	url := fmt.Sprintf("https://dncapi.bqiapp.com/api/coin/web-coinrank?page=1&type=-1&pagesize=%s&webp=1", pageSize)
	resp, err := http.Get(url)
	if err != nil {

	}


	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	return body
}

func UpdateDB() {

	body := getMarket("100")

	var feixiaoData FeixiaoData
	json.Unmarshal(body, &feixiaoData)
	for i, v := range feixiaoData.Data {
		//db.Where(&User{Name: "jinzhu", Age: 20}).First(&user)
		var market Market
		v.Id = i+1
		DB.Where(&Market{Id: v.Id}).First(&market)
		//DB.Where(&v).First(&market)
		fmt.Println("market:",market.Name)
		if  len(market.Name) == 0 {
			v.CreateTime = time.Now()
			v.UpdateTime = time.Now()
			DB.Create(v)
		}else{
			DB.Model(&v).Updates(v)
		}
	}

}

func main() {
	// 设置路由，如果访问/，则调用index方法
	http.HandleFunc("/", index)

	// 启动web服务，监听9090端口
	err := http.ListenAndServe(":9089", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
