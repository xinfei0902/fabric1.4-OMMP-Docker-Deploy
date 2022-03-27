package main

import (
	"deploy-server/app/dcmd"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	dcmd.Execute()
}
