// zyjsxy project main.go
package main

import (
	"zyjsxy-api/database"
)

func main() {
	defer database.Orm.Close()
	router := initRouter()
	router.Run(":8000")
}
