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
	Key        []string `json:"key" `
}

func readOperationCodes(operationCodesFile string, logger *zap.Logger) (map[int]string, map[string][]int, error) {
	codes := jsonFileWithCodes{}
	fileBytes, _ := os.ReadFile(operationCodesFile)
	err := json.Unmarshal(fileBytes, &codes)
	if err != nil {
		logger.Error("Проблемы с файлом с кодами операций", zap.Error(err),
			zap.Namespace("директория"))
		return nil, nil, err
	}
	var meaningMap = make(map[int]string)
	var kindMap = make(map[string][]int)
	for i, operation := range codes.Operations {
		meaningMap[operation] = codes.Meaning[i]
		if _, ok := kindMap[codes.Key[i]]; !ok {
			kindMap[codes.Key[i]] = make([]int, 0)
		}
		kindMap[codes.Key[i]] = append(kindMap[codes.Key[i]], operation)
	}
	return meaningMap, kindMap, nil
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
	codeMeanings, kindMap, err := readOperationCodes(operationCodesFile, logger)
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
			if validCode(kindMap["чтение"], code) {
				logger.Info("Нет операции с таким кодом в списке или вы перепутали коды местами", zap.Int(
					"номер строки", stringNumber), zap.Int("номер операции, введённой вами", code), zap.Any("коды чтения",
					kindMap["чтение"]))
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
			if validCode(kindMap["запись"], code) {
				logger.Info("Нет операции с таким кодом в списке или вы перепутали коды местами", zap.Int(
					"номер строки", stringNumber), zap.Int("номер операции, введённой вами", code), zap.Any("коды чтения",
					kindMap["запись"]))
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
