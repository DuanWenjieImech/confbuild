package main

import (
	"flag"
	"strings"

	"github.com/360EntSecGroup-Skylar/excelize/v2"
)

var sheets, excelRoot, pkg, outPath string

func init() {
	flag.StringVar(&sheets, "sheets", "", "-sheets sheets to export, sheet name and file name need to be the same ',' split multiple sheets")
	flag.StringVar(&excelRoot, "excelRoot", "", "-excelRoot  excel file root")
	flag.StringVar(&pkg, "package", "", "-package  struct package name")
	flag.StringVar(&outPath, "outpath", "", "-outpath json data file and go struct file output path")
}

/**
 *@author LanguageY++2013
 *2019/3/10 1:01 AM
 **/
func main() {
	flag.Parse()

	if excelRoot == "" {
		panic("excelRoot can not empty")
	}

	if sheets == "" {
		panic("sheets can not empty")
	}

	if pkg == "" {
		panic("package can not empty")
	}

	sheetSlice := strings.Split(sheets, ",")

	structDescList := make([]*StructDesc, 0)

	for _, sheetName := range sheetSlice {
		xlsx, err := excelize.OpenFile(excelRoot + sheetName + ".xlsx")
		if err != nil {
			panic(err.Error())
		}
		currentSheet :=make([]string, 0)
		currentSheet = append(currentSheet,sheetName)

		//预处理数据，移除空列和注释列
		xlsx_processed := PreProcess(currentSheet, xlsx)

		//数据解析
		Data_Parse(currentSheet, xlsx_processed)

		//结构解析
		tempDesItemList := Struct_Process(currentSheet, xlsx_processed)
		structDescList = append(structDescList,tempDesItemList...)
	}

	Struct_Parse(structDescList)
}
