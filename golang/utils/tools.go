package utils

import (
	"baliance.com/gooxml/document"
	"github.com/tealeg/xlsx"
	"log"
)

//读取Excel全部数据
func ReadExcelAll(excelPath string) (string, error) {
	// 打开 Excel 文件
	xlFile, err := xlsx.OpenFile(excelPath)
	if err != nil {
		log.Println(err)
		return "", err
	}
	res := ""
	// 遍历每个 Sheet
	for _, sheet := range xlFile.Sheets {

		// 遍历每行数据
		for _, row := range sheet.Rows {
			// 遍历每个单元格
			for _, cell := range row.Cells {
				// 输出单元格的值
				res += cell.Value
				//fmt.Printf("第%d行，第%d列 ：%s\t", rowKey, colKey, cell.Value)
			}
		}
	}
	return res, nil
}
func ReadDocxAll(fileName string) (string, error) {
	doc, err := document.Open(fileName)
	if err != nil {
		return "", err
	}
	text := ""

	for _, para := range doc.Paragraphs() {
		//run为每个段落相同格式的文字组成的片段
		for _, run := range para.Runs() {
			text += run.Text()
		}
	}
	return text, nil
}
