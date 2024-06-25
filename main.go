package main

import (
	"fmt"

	"roboline/env"
	"roboline/logger"
	"roboline/read"
)

func main() {
	logger := logger.InitLogger("logger.json")
	envMap := env.InitEnv(logger)
	fmt.Println(envMap["FILE_WITH_DATA"])
	read.Read(envMap["FILE_WITH_DATA"], logger)
}
