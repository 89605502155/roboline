package read

import (
	"bufio"
	"encoding/json"
	"os"
	"strconv"
	"strings"

	"go.uber.org/zap"
)

type Query struct {
	Adress       int
	Code         int
	MemoryAdress int
	Control      int
}
type WriteQuery struct {
	Adress       int
	Code         int
	MemoryAdress int
	Value        int
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

type jsonFileWithCodes struct {
	Operations []int    `json:"operations"`
	Meaning    []string `json:"meanings"`
}

func readOperationCodes(operationCodesFile string, logger *zap.Logger) ([]int, map[int]string, error) {
	codes := jsonFileWithCodes{}
	fileBytes, _ := os.ReadFile(operationCodesFile)
	err := json.Unmarshal(fileBytes, &codes)
	if err != nil {
		logger.Error("Проблемы с файлом с кодами операций", zap.Error(err),
			zap.Namespace("директория"))
		return nil, nil, err
	}
	var meaningMap = make(map[int]string)
	for i, operation := range codes.Operations {
		meaningMap[operation] = codes.Meaning[i]
	}
	return codes.Operations, meaningMap, nil
}

func validCode(codes []int, operation int) bool {
	resoult := true
	for _, code := range codes {
		if code == operation {
			resoult = false
			break
		}
	}
	return resoult
}

func Read(fileName string, operationCodesFile string, logger *zap.Logger) ([]Query, []WriteQuery, []bool, map[int]string) {
	file, err := os.Open(fileName)
	if err != nil {
		logger.Error("Проблемы с файлом с данными", zap.Error(err),
			zap.Namespace("директория бд"))
		return nil, nil, nil, nil
	}
	defer file.Close()
	codes, codeMeanings, err := readOperationCodes(operationCodesFile, logger)
	if err != nil {
		return nil, nil, nil, nil
	}

	scan := bufio.NewScanner(file)
	querySlice := make([]Query, 0)
	writeQuerySlice := make([]WriteQuery, 0)
	isQueryTypeSlice := make([]bool, 0)
	var code, stringNumber int
	for scan.Scan() {
		stringNumber++
		str := strings.Split(scan.Text(), " ")
		if len(str) == 4 {
			code = middle(str[1])
			if validCode(codes, code) {
				logger.Info("Нет операции с таким кодом в списке", zap.Int(
					"номер строки", stringNumber), zap.Int("номер операции", code))
				continue
			}
			object := Query{
				Adress: middle(str[0]), Code: middle(str[1]),
				MemoryAdress: middle(str[2]), Control: middle(str[3]),
			}
			querySlice = append(querySlice, object)
			isQueryTypeSlice = append(isQueryTypeSlice, true)
		} else if len(str) == 5 {
			code = middle(str[1])
			if validCode(codes, code) {
				logger.Info("Нет операции с таким кодом в списке", zap.Int(
					"номер строки", stringNumber), zap.Int("номер операции", code))
				continue
			}
			object := WriteQuery{
				Adress: middle(str[0]), Code: middle(str[1]),
				MemoryAdress: middle(str[2]), Value: middle(str[3]), Control: middle(str[4]),
			}
			writeQuerySlice = append(writeQuerySlice, object)
			isQueryTypeSlice = append(isQueryTypeSlice, false)
		} else {
			logger.Info("Неверный формат строки с данными", zap.Int(
				"номер строки", stringNumber), zap.Int("количество переданных вами аргументов",
				len(str)))
		}
	}
	return querySlice, writeQuerySlice, isQueryTypeSlice, codeMeanings
}
