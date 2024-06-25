package main

import (
	"fmt"


	"roboline/env"
	"roboline/logger"
)

func main() {
	logger := logger.InitLogger("logger.json")
	envMap:= env.InitEnv(logger)
	fmt.Println(envMap["FILE_WITH_DATA"])
}
