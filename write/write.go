package write

import (
	"fmt"

	"go.uber.org/zap"

	"roboline/read"
)

func OutPut(queries []read.Query, writes []read.WriteQuery, isQueryType []bool,
	codeMeanings map[int]string, logger *zap.Logger) {
	length := len(queries) + len(writes)
	var q, w int
	for i := 0; i < length; i++ {
		if isQueryType[i] {
			outPutQueries(queries[q], codeMeanings)
			q++
		} else {
			outPutWrites(writes[w], codeMeanings)
			w++
		}
	}
}

func outPutQueries(query read.Query, codeMeanings map[int]string) {
	fmt.Println("Адрес устройства:", query.Adress)
	fmt.Printf("Код функции: %d (%s)\n", query.Code, codeMeanings[query.Code])
	fmt.Printf("Адрес ячейки памяти, откуда идёт чтение: %d\n", query.MemoryAdress)
	fmt.Println()
	fmt.Println()
}
func outPutWrites(write read.WriteQuery, codeMeanings map[int]string) {
	fmt.Println("Адрес устройства:", write.Adress)
	fmt.Printf("Код функции: %d (%s)\n", write.Code, codeMeanings[write.Code])
	fmt.Printf("Адрес ячейки памяти, куда идёт запись: %d\n", write.MemoryAdress)
	fmt.Println("Записываемое значение: ", write.Value)
	fmt.Println()
	fmt.Println()
}
