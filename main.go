package main


import (
	"log"
	"net/http"
	"io/ioutil"
	"encoding/json"
	"fmt"
)


type Market struct {
	Name string `json:"name"`
	Price float64 `json:"current_price"`
	Rose float64 `json:"change_percent"`
}


type FeixiaoData struct {
	Data []Market `json:"data"`
	Code int `json:"code"`
	Msg  string `json:"msg"`
}


const (
	teststr = `{"code": 200,"msg": "ok","data": [{"name": "BTC","current_price": "100","change_percent": "100","tee":"safaf"}]}`
)

func index(w http.ResponseWriter, r *http.Request) {
	// 往w里写入内容，就会在浏览器里输出

	w.Header().Set("Access-Control-Allow-Origin", "*")             //允许访问所有域
	w.Header().Add("Access-Control-Allow-Headers", "Content-Type") //header的类型
	w.Header().Set("content-type", "application/json")             //返回数据格式是json

	rQuery :=r.URL.Query()

	url :=fmt.Sprintf("https://dncapi.bqiapp.com/api/coin/web-coinrank?page=1&type=-1&pagesize=%s&webp=1",rQuery["get-size"][0])


	resp,err:= http.Get(url)
	if err != nil{

	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		// handle error
	}
	var feixiaoData FeixiaoData
	json.Unmarshal(body,&feixiaoData)

	ret ,err:=json.Marshal(feixiaoData)

	w.Write(ret)

}

func main() {
	// 设置路由，如果访问/，则调用index方法
	http.HandleFunc("/", index)

	// 启动web服务，监听9090端口
	err := http.ListenAndServe(":9090", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
