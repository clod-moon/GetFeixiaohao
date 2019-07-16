
package main

import (
	"fmt"
	"time"
	"io/ioutil"
	"net/http"
)



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

func getMarket(pageSize string) []byte {
	url := fmt.Sprintf("https://dncapi.bqiapp.com/api/coin/web-coinrank?page=1&type=-1&pagesize=%s&webp=1", pageSize)
	resp, err := http.Get(url)
	if err != nil {

	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	return body
}