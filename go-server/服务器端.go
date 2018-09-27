package main

/*服务器*/

import (
	"flag"
	"golib/modules/config"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"io/ioutil"
	"strings"
	"net"
	"time"
	"crypto/tls"
	"runtime"
	"test2/ini"
)

var path = flag.String("http param","./etc/cfg.ini","http cfg")

var PATH string = "./etc/cfg.ini"

func main()  {
	flag.Parse()

	err := config.InitModuleByParams(*path)
	if err != nil {
		fmt.Println("加载配置文件失败",err)
		return
	}

	//方法一
	rout1 := httprouter.New()
	//RoutDel(rout)
	rout1.POST("/test",Deals)
	http.ListenAndServe(":888",rout1)

	//方法二
	http.HandleFunc("/test",Deal)
	svr := http.Server{
		Addr:fmt.Sprintf("%s:%d","192.168.127.181",888),
		ReadTimeout:time.Duration(2)*time.Second,
		WriteTimeout:time.Duration(2)*time.Second,
	}
	svr.ListenAndServe()

	//方法三
	http.HandleFunc("/test",Deal)
	http.ListenAndServe(":888",nil)

	//PostCli()


	//方法四（模板）
	rout2 := httprouter.New()
	ini.DealServicces(rout2)
	//监听判断
	cerF := config.StringDefault("HttpsCertFile","")
	keyF := config.StringDefault("HttpsKeyFile","")
	if cerF != "" && keyF != "" {
		fmt.Println("启动https监听---->")
		go func() {
			fmt.Println(http.ListenAndServeTLS(":8080",cerF,keyF,rout2))
		}()
	}else {
		fmt.Println("启动http监听---->")
		go func() {
			fmt.Println(http.ListenAndServe(":8080",rout2))
		}()
	}

	fmt.Println("程序启动成功！")
	runtime.Goexit()

}

func RoutDel(rout *httprouter.Router)  {
	rout.POST("/test",Deals)
}

func Deals(w http.ResponseWriter,r *http.Request,_ httprouter.Params)  {
	inf,err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Println("读取报文失败",err)
	}
	fmt.Println(string(inf))
	w.Write([]byte("haha"))
}


func Deal(w http.ResponseWriter,r *http.Request)  {
	inf,err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Println("读取报文失败",err)
	}
	fmt.Println(string(inf))
	w.Write([]byte("haha"))
}

func PostCli()  {
	str := "hello world"
	req := strings.NewReader(str)

	tr := &http.Transport{
		Dial:(&net.Dialer{
			Timeout:time.Duration(2)*time.Second,
		}).Dial,
		TLSClientConfig: &tls.Config{InsecureSkipVerify:true},
		DisableKeepAlives:true,
	}

	client := http.Client{Transport:tr,Timeout:time.Duration(2)*time.Second}
	reqMsg,err := http.NewRequest("POST","http:192.168.127.181:888/test",req)
	if err != nil {
		fmt.Println("post err")
	}

	resp,err := client.Do(reqMsg)
	defer resp.Body.Close()

	rsp,_ := ioutil.ReadAll(resp.Body)
	fmt.Println("响应报文",string(rsp))
}
