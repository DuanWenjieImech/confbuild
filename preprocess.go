package main

import (
	"log"

	"github.com/360EntSecGroup-Skylar/excelize/v2"
)

func PreProcess(sheetSlice []string, xlsx *excelize.File) *excelize.File {
	//mapData := make(map[string]interface{})
	for _, sheet := range sheetSlice {
		rows, err := xlsx.GetRows(sheet)

		if err != nil {
			panic(err)
		}

		if len(rows) == 0 {
			break
		}

		columnNum := len(rows[0])
		indexToDel := make([]int, 0)

		//skip first column, always keep it
		for j := 1; j < columnNum; j++ {
			cell := rows[0][j]

			if CheckIsValidRule(cell) == false {
				indexToDel = append(indexToDel, j)
				log.Printf("Remove Column Index: %d", j)
			}
		}
		var additionMove int = 0
		for _, index := range indexToDel {
			columnString, _ := excelize.ColumnNumberToName(index + 1 - additionMove)
			log.Printf("Remove Column Index String: %s", columnString)
			err = xlsx.RemoveCol(sheet, columnString)
			if err != nil {
				panic(err)
			}
			additionMove++
		}
	}

	return xlsx
}

func CheckIsValidRule(rule string) bool {
	switch rule {
	case "required", "optional", "repeated", "optional_struct":
		return true
	case "comment", "":
		return false
	default:
		return false
	}
}
