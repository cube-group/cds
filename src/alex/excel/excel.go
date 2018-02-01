package excel

import (
    "github.com/tealeg/xlsx"
)

//生成表格
//@fileName 文件名
//@data 表格数据
//@sheet sheet名,默认sheet0
//@return error
func Create(fileName string, data [][]string, sheet ...string) error {
    var sheetName string
    file := xlsx.NewFile()
    if len(sheet) <= 0 {
        sheetName = "sheet0"
    } else {
        sheetName = sheet[0]
    }
    sheetObj, err := file.AddSheet(sheetName)
    if err != nil {
        return err
    }
    for _, row := range data {
        rowObj := sheetObj.AddRow()
        for _, cell := range row {
            cellObj := rowObj.AddCell()
            cellObj.Value = cell
        }
    }
    err = file.Save(fileName)
    if err != nil {
        return err
    }
    return nil
}

//读取表格
//@param file 文件路径
//@param sheet 读取的sheet名，默认sheet0
//@return 返回二维slice
func Read(file string, sheet ...string) ([][]string, error) {
    var sheetName string
    if len(sheet) <= 0 {
        sheetName = "sheet0"
    } else {
        sheetName = sheet[0]
    }

    xlFile, err := xlsx.OpenFile(file)
    if err != nil {
        return nil, err
    }

    var res [][]string
    for _, sheetObj := range xlFile.Sheets {
        if sheetObj.Name != sheetName {
            continue
        }
        for _, row := range sheetObj.Rows {
            var temp []string
            for _, cell := range row.Cells {
                temp = append(temp, cell.String())
            }
            res = append(res, temp)
        }
    }
    return res, nil
}
