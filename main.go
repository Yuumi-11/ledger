package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"ledger/database"
	"ledger/routes"

	"golang.org/x/term"
)

/* func main(){
	fmt.Println("begin to work > <")
	fmt.Println("end to work > <")
} */

func main() {
	fmt.Print("请输入 MySQL root 密码：")
	password, err := term.ReadPassword(int(os.Stdin.Fd()))
	fmt.Println()
	if err != nil {
		log.Fatal("读取密码失败，请在终端运行程序：", err)
	}
	dsn := fmt.Sprintf("root:%s@tcp(127.0.0.1:3306)/ledger?charset=utf8mb4&parseTime=true&loc=Local", string(password))

	db, err := database.Connect(dsn)
	if err != nil {
		log.Fatal("failed to connect to database:", err)
	}

	fmt.Println("successfully connected to database")

	router := gin.Default()
	router.GET("/ping", ping)
	routes.Register(router, db)

	err = router.Run("127.0.0.1:8080")
	if err != nil {
		log.Fatal("启动服务失败：", err)
	}
} // main 在这里结束

func ping(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "记账服务已启动",
	})
}
