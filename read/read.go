package read

import (
	"bufio"
	"os"
)

func Read(fileName string) {
	file, err := os.Open(fileName)
	defer file.Close()
	scan := bufio.NewScanner(file)
	for scan.Scan() {
		scan.Text()
	}
}
