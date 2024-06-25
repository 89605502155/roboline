package env

import (
	"flag"
	"os"

	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

func InitEnv(logger *zap.Logger) map[string]string {
	var envFileDirectory string
	flag.StringVar(&envFileDirectory, "env", "./../", "use .env file")
	flag.Parse()
	if _, err := os.Stat(envFileDirectory); os.IsNotExist(err) {
		logger.Error("", zap.String("File %s does not exist", envFileDirectory))
	}
	godotenv.Load(envFileDirectory)
	myEnv, err := godotenv.Read(".env")
	if err != nil {
		logger.Error("Проблемы с конфигом", zap.Error(err), zap.String("env ", envFileDirectory),
			zap.Namespace("директория env"), zap.Any("myEnv", myEnv))
	}
	return myEnv
}
