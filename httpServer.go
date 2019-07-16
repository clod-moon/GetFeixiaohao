package main

import (
	"net/http"
	"fmt"
	"encoding/json"
	"log"
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

	rows,err:=DB.Find(&list).Rows()

	if err != nil{
		w.Write([]byte(err.Error()))
		return
	}

	for rows.Next() {
		var m Market
		rows.Scan(&m.Id,&m.Name,&m.Price,&m.Rose)
	}
	feixiaoData.Code = 200
	feixiaoData.Msg = "success"
	feixiaoData.Data = append(feixiaoData.Data,list...)
	//json.Unmarshal(body, &feixiaoData)

	ret, _ := json.Marshal(feixiaoData)

	w.Write(ret)

}

func HttpServer(){
	http.HandleFunc("/", index)

	// 启动web服务，监听9090端口

	err := http.ListenAndServe(fmt.Sprintf(":%d",httpSreverPort), nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
