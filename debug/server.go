package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"time"
)

func main() {
	http.HandleFunc("/choose", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
			return
		}

		// 解析表单数据
		err := r.ParseForm()
		if err != nil {
			http.Error(w, "Error parsing form data", http.StatusInternalServerError)
			return
		}

		podname := r.FormValue("podname")
		fmt.Println("接收到Pod: ", podname)

		options := []string{"volcano", "volcano-node01"}
		rand.Seed(time.Now().UnixNano()) // 初始化随机数生成器
		time.Sleep(5 * time.Second)
		choice := options[rand.Intn(len(options))]
		fmt.Println("调度到: ", choice)
		// 向客户端发送响应
		fmt.Fprintf(w, choice)
	})

	fmt.Println("Server is running on http://localhost:1234")
	err := http.ListenAndServe(":1234", nil)
	if err != nil {
		fmt.Printf("Server failed to start: %v\n", err)
	}
}
