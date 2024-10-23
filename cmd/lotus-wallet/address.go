package main

import (
	"fmt"
	"github.com/xuri/excelize/v2"
	"sync"
	"time"
)

var FromToMap = make(map[string]map[string]struct{})
var lk = sync.Mutex{}

func loadAddressExcel() {
	ticker := time.NewTicker(time.Minute * 10)
	defer ticker.Stop() // 程序结束时停止 ticker

	for range ticker.C {
		readAddressExcel()
	}
}
func readAddressExcel() {
	f, err := excelize.OpenFile("地址.xlsx")
	if err != nil {
		panic(fmt.Sprintf("无法打开文件: %v", err))
	}
	defer f.Close()

	sheetName := f.GetSheetName(0)

	rows, err := f.GetRows(sheetName)
	if err != nil {
		panic(fmt.Sprintf("无法读取行: %v", err))
	}

	var tmpMap = make(map[string]map[string]struct{})
	for i, row := range rows {
		if i == 0 {
			continue
		}
		if len(row) < 2 {
			panic(fmt.Sprintf("文件格式错误, row: %d", i))
		}
		from := row[0]
		to := row[1]

		// 如果 `from` 不在 map 中，初始化
		if _, exists := tmpMap[from]; !exists {
			tmpMap[from] = make(map[string]struct{})
		}

		// 将 `to` 添加到 from 对应的 map 中
		tmpMap[from][to] = struct{}{}
	}

	log.Infof("reload address: %v", tmpMap)
	lk.Lock()
	defer lk.Unlock()

	FromToMap = tmpMap
}

func checkAddress(from, to string) bool {
	lk.Lock()
	defer lk.Unlock()
	if toMap, ok := FromToMap[from]; !ok {
		return false
	} else if _, ok2 := toMap[to]; !ok2 {
		return false
	} else {
		return true
	}
}
