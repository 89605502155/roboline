package main

import (
	"roboline/env"
	"roboline/logger"
	"roboline/read"
	"roboline/write"
)

func main() {
	logger := logger.InitLogger("logger.json")
	envMap := env.InitEnv(logger)
	queries, writes, isQuery, codeMeaning := read.Read(envMap["FILE_WITH_DATA"],
		envMap["JSON_FILE_WITH_OPERATIONS"], logger)

	write.OutPut(queries, writes, isQuery, codeMeaning, logger)
}
