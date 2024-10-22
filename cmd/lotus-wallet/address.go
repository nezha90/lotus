package main

import (
	"fmt"
	"github.com/xuri/excelize/v2"
)

var FromToMap = make(map[string]map[string]struct{})

func loadExcel() {
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
		if _, exists := FromToMap[from]; !exists {
			FromToMap[from] = make(map[string]struct{})
		}

		// 将 `to` 添加到 from 对应的 map 中
		FromToMap[from][to] = struct{}{}
	}
}

func checkAddress(from, to string) bool {
	if toMap, ok := FromToMap[from]; !ok {
		return false
	} else if _, ok2 := toMap[to]; !ok2 {
		return false
	} else {
		return true
	}
}
