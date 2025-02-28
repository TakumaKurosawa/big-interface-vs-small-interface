package main

import (
	"fmt"

	"github.com/yourusername/simple-go-project/pkg/greeting"
)

func main() {
	fmt.Println("シンプルなGoアプリケーションを開始します")
	message := greeting.GetGreeting("ユーザー")
	fmt.Println(message)
}
