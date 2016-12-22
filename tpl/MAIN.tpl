package main

import (
	"database/sql"
	"fmt"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

func main() {
	var err error
	db, err = sql.Open("mysql", "$USER$:$PASSWORD$@tcp($HOST$:$PORT$)/$DB$")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer db.Close()
	err = db.Ping()
	if err != nil {
		fmt.Println(err)
		return
	}
	router := gin.Default()
	//=======ADD ROUTER
	$ROUTERS$

	//======END  ADD ROUTER
	router.Run(":3000")

}
