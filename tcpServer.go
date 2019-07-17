package main

import (
	"log"
	"net"

	"fmt"
	"time"
	"encoding/json"
)

func UpdateDB(body *[]byte) {
	var feixiaoData FeixiaoData
	json.Unmarshal(*body, &feixiaoData)
	for i, v := range feixiaoData.Data {
		//db.Where(&User{Name: "jinzhu", Age: 20}).First(&user)
		var market Market
		v.Id = i + 1
		DB.Where(&Market{Id: v.Id}).First(&market)
		//DB.Where(&v).First(&market)
		//fmt.Println("market:", market.Name)
		if len(market.Name) == 0 {
			v.CreateTime = time.Now()
			v.UpdateTime = time.Now()
			DB.Create(v)
		} else {
			DB.Model(&v).Updates(v)
		}
	}

}

func connHandler(c net.Conn) {
	//1.conn是否有效
	if c == nil {
		log.Panic("无效的 socket 连接")
	}

	defer c.Close()
	for {

		var feixiaoData FeixiaoData

		body := getMarket(biAmount)
		json.Unmarshal(body,&feixiaoData)

		ret ,_:=json.Marshal(feixiaoData)
		_, err := c.Write(ret)
		if err != nil {
			break
		}
		time.Sleep(time.Second)
		//3.1 网络数据流读入 buffer
		//cnt, err := c.Read(buf)
		////3.2 数据读尽、读取错误 关闭 socket 连接
		//if cnt == 0 || err != nil {
		//	c.Close()
		//	break
		//}

		////3.3 根据输入流进行逻辑处理
		////buf数据 -> 去两端空格的string
		//inStr := strings.TrimSpace(string(buf[0:cnt]))
		////去除 string 内部空格
		//cInputs := strings.Split(inStr, " ")
		////获取 客户端输入第一条命令
		//fCommand := cInputs[0]

		//fmt.Println("客户端传输->" + fCommand)

		//switch fCommand {
		//case "ping":
		//	c.Write([]byte("服务器端回复-> pong\n"))
		//case "hello":
		//	c.Write([]byte("服务器端回复-> world\n"))
		//default:
		//	c.Write([]byte("服务器端回复" + fCommand + "\n"))
		//}

		//c.Close() //关闭client端的连接，telnet 被强制关闭
	}
	fmt.Printf("来自 %v 的连接关闭\n", c.RemoteAddr())

}

//开启serverSocket
func ServerSocket() {
	//1.监听端口
	server, err := net.Listen("tcp", fmt.Sprintf(":%d",tcpSreverPort))

	if err != nil {
		fmt.Println("开启socket服务失败")
	}

	fmt.Println("正在开启 Server ...")

	for {
		//2.接收来自 client 的连接,会阻塞
		conn, err := server.Accept()

		if err != nil {
			fmt.Println("连接出错")
		}

		//并发模式 接收来自客户端的连接请求，一个连接 建立一个 conn，服务器资源有可能耗尽 BIO模式
		go connHandler(conn)
	}

	wg.Done()

}

func UpdateMarket(){
	for{
		body := getMarket(biAmount)
		UpdateDB(&body)
		time.Sleep(time.Second*updatePreSecond)
	}
	wg.Done()
}