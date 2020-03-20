package main

import (
	"flag"
	"strings"

	"github.com/360EntSecGroup-Skylar/excelize/v2"
)

var sheets, excel, pkg, outPath string

func init() {
	flag.StringVar(&sheets, "sheets", "", "-sheets sheets to export, ',' split multiple sheets")
	flag.StringVar(&excel, "excel", "", "-excel  excel filename to parse")
	flag.StringVar(&pkg, "package", "", "-package  struct package name")
	flag.StringVar(&outPath, "outpath", "", "-outpath json data file and go struct file output path")
}

/**
 *@author LanguageY++2013
 *2019/3/10 1:01 AM
 **/
func main() {
	flag.Parse()

	if excel == "" {
		panic("excel can not empty")
	}

	if sheets == "" {
		panic("sheets can not empty")
	}

	if pkg == "" {
		panic("package can not empty")
	}

	sheetSlice := strings.Split(sheets, ",")

	xlsx, err := excelize.OpenFile(excel)
	if err != nil {
		panic(err.Error())
	}

	//预处理数据，移除空列和注释列
	xlsx_processed := PreProcess(sheetSlice, xlsx)

	//数据解析
	Data_Parse(sheetSlice, xlsx_processed)

	//结构解析
	Struct_Parse(sheetSlice, xlsx_processed)
}
