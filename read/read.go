package read

import (
	"bufio"
	"os"
	"strconv"
	"strings"

	"go.uber.org/zap"
)

type Query struct {
	Id           int
	Adress       int
	Code         int
	MemoryAdress int
	Control      int
}
type WriteQuery struct {
	Id           int
	Adress       int
	Code         int
	MemoryAdress int
	Control      int
}

func formatingToHex(hexaString string) string {
	numberStr := strings.Replace(hexaString, "0x", "", -1)
	numberStr = strings.Replace(numberStr, "0X", "", -1)
	return numberStr
}
func middle(str string) int {
	n, _ := strconv.ParseInt(formatingToHex(str), 16, 64)
	return int(n)
}

func Read(fileName string, logger *zap.Logger) {
	file, err := os.Open(fileName)

	if err != nil {
		logger.Error("Проблемы с файлом с данными", zap.Error(err),
			zap.Namespace("директория бд"))
	}
	defer file.Close()
	scan := bufio.NewScanner(file)
	// var counter int
	for scan.Scan() {
		str := strings.Split(scan.Text(), " ")
		if len(str) == 4 {
			logger.Info("Данные из файла", zap.Int("имя", middle(str[0])), zap.String("s", str[0]))
		} else if len(str) == 5 {
			logger.Info("Данные из файла", zap.Int("фамилия", middle(str[1])), zap.String("s", str[1]))
		}
	}
}
