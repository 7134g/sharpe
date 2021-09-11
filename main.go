package main

import (
	"sharpe/model"
	"sharpe/serve"
)

func main() {
	model.DBInit()
	serve.Run()
}
