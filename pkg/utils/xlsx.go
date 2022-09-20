package utils

import (
	"encoding/json"
	"fmt"
	"github.com/xuri/excelize/v2"
	"reflect"
)

func StructWriteXlsx(sheet string, records interface{}) *excelize.File {
	xlsx := excelize.NewFile()    // new file
	index := xlsx.NewSheet(sheet) // new sheet
	xlsx.SetActiveSheet(index)    // set active (default) sheet
	firstCharacter := 65          // start from 'A' line
	s := reflect.ValueOf(records)
	t := reflect.TypeOf(records)
	if s.Kind() == reflect.Ptr {
		//s := reflect.ValueOf(data)
		//fmt.Println(reflect.Value.Elem(s).Interface())
		records = reflect.Value.Elem(s).Interface()
		s = reflect.ValueOf(records)
		t = reflect.TypeOf(records)
	}
	// log.Info(t)
	if t.Kind() != reflect.Slice {
		return xlsx
	}

	// log.Info(s)
	for i := 0; i < s.Len(); i++ {
		elem := s.Index(i).Interface()
		elemType := reflect.TypeOf(elem)
		elemValue := reflect.ValueOf(elem)
		// log.Info(elemValue)
		for j := 0; j < elemType.NumField(); j++ {
			field := elemType.Field(j)
			//fieldVal := elemValue.Field(j)
			//fmt.Println("field:", field.Name)

			name := field.Name
			// 获取struct的tag
			tag := field.Tag.Get("json")
			if len(tag) > 0 {
				name = tag
			}
			//fmt.Println("xlsx name:", name)
			column := string(rune(firstCharacter + j))
			// 把读到的tag写入表头
			if i == 0 {
				if err := xlsx.SetCellValue(sheet, fmt.Sprintf("%s%d", column, i+1), name); err != nil {
					return nil
				}
			}

			// 写入struct值
			if err := xlsx.SetCellValue(sheet, fmt.Sprintf("%s%d", column, i+2), elemValue.Field(j).Interface()); err != nil {
				return nil
			}
		}
	}
	return xlsx
}

func MapWriteXlsx(sheet string, records interface{}) (*excelize.File, error) {
	xlsx := excelize.NewFile()    // new file
	index := xlsx.NewSheet(sheet) // new sheet
	xlsx.SetActiveSheet(index)    // set active (default) sheet

	s := reflect.ValueOf(records)
	t := reflect.TypeOf(records)
	if s.Kind() == reflect.Ptr {
		records = reflect.Value.Elem(s).Interface()
		s = reflect.ValueOf(records)
		t = reflect.TypeOf(records)
	}

	// log.Info(t)
	if t.Kind() != reflect.Slice {
		return xlsx, fmt.Errorf("records is not slice")
	}
	//var array = make([]map[string]interface{}, 0)
	var array = make([]map[string]json.RawMessage, 0)
	//var array = make([]map[string][]byte, 0)
	if err := AnyToAnyV2(records, &array); err != nil {
		return xlsx, fmt.Errorf("records convert to slice map, %s", err)
	}
	//fmt.Println("array:", array)
	var headers = make(map[string]string, 0)
	for i, m := range array {
		i := i
		m := m
		j := 0
		for k, _ := range m {
			name := k
			column, err := excelize.ColumnNumberToName(j + 1)
			//cellName, err := excelize.CoordinatesToCellName(j+1, i+1)
			if err != nil {
				return nil, err
			}
			// 写入表头
			if i == 0 {
				if err = xlsx.SetCellValue(sheet, fmt.Sprintf("%s%d", column, i+1), name); err != nil {
					return nil, err
				}
				headers[column] = name
			}

			// 写入map值
			var axis = fmt.Sprintf("%s%d", column, i+2)
			var val = string(m[headers[column]])
			fmt.Println("val:", val)
			if err = xlsx.SetCellValue(sheet, axis, val); err != nil {
				return nil, err
			}
			j = j + 1
		}
	}
	return xlsx, nil
}

// String 将 `[]byte` 转换为 `string`
func String(b []byte) string {
	for idx, c := range b {
		if c == 0 {
			return string(b[:idx])
		}
	}
	return string(b)
}

// StringWithoutZero 将 `[]byte` 转换为 `string`
func StringWithoutZero(b []byte) string {
	s := make([]rune, len(b))
	offset := 0
	for i, c := range b {
		if c == 0 {
			offset++
		} else {
			s[i-offset] = rune(c)
		}
	}
	return string(s[:len(b)-offset-1])
}

type Fruit struct {
	Id    uint
	Name  string
	Price float64
}

type Style struct {
	Border        []Border    `json:"border"`
	Fill          Fill        `json:"fill"`
	Font          *Font       `json:"font"`
	Alignment     *Alignment  `json:"alignment"`
	Protection    *Protection `json:"protection"`
	NumFmt        int         `json:"number_format"`
	DecimalPlaces int         `json:"decimal_places"`
	CustomNumFmt  *string     `json:"custom_number_format"`
	Lang          string      `json:"lang"`
	NegRed        bool        `json:"negred"`
}

// Border 边框
type Border struct {
	Type  string `json:"type"`
	Color string `json:"color"`
	Style int    `json:"style"`
}

// Fill 填充
type Fill struct {
	Type    string   `json:"type"`
	Pattern int      `json:"pattern"`
	Color   []string `json:"color"`
	Shading int      `json:"shading"`
}

// Font 字体
type Font struct {
	Bold      bool    `json:"bold"`      // 是否加粗
	Italic    bool    `json:"italic"`    // 是否倾斜
	Underline string  `json:"underline"` // single    double
	Family    string  `json:"family"`    // 字体样式
	Size      float64 `json:"size"`      // 字体大小
	Strike    bool    `json:"strike"`    // 删除线
	Color     string  `json:"color"`     // 字体颜色
}

// Protection 保护
type Protection struct {
	Hidden bool `json:"hidden"`
	Locked bool `json:"locked"`
}

// Alignment 对齐
type Alignment struct {
	Horizontal      string `json:"horizontal"`        // 水平对齐方式
	Indent          int    `json:"indent"`            // 缩进  只要设置了值，就变成了左对齐
	JustifyLastLine bool   `json:"justify_last_line"` // 两端分散对齐，只有在水平对齐选择 distributed 时起作用
	ReadingOrder    uint64 `json:"reading_order"`     // 文字方向 不知道值范围和具体的含义
	RelativeIndent  int    `json:"relative_indent"`   // 不知道具体的含义
	ShrinkToFit     bool   `json:"shrink_to_fit"`     // 缩小字体填充
	TextRotation    int    `json:"text_rotation"`     // 文本旋转
	Vertical        string `json:"vertical"`          // 垂直对齐
	WrapText        bool   `json:"wrap_text"`         // 自动换行
}

/*
https://blog.csdn.net/weixin_41546513/article/details/121272718

常用的工具函数
打开文件
func OpenFile(filename string, opt ...Options) (*File, error)

新建文件
func excelize.NewFile() *excelize.File

拆分单元格坐标 单元格坐标字符串拆分成 列名称 行索引
func excelize.SplitCellName(cell string) (string, int, error)

通过行列名称组合单元格坐标
func excelize.JoinCellName(col string, row int) (string, error)

列名转索引
func excelize.ColumnNameToNumber(name string) (int, error)

列索引转列名
func excelize.ColumnNumberToName(num int) (string, error)

坐标字符串转索引 行，列的数字索引
func excelize.CellNameToCoordinates(cell string) (int, int, error)

行列数字索引转坐标字符串 最后一个是否绝对坐标 例如:$A$1
func excelize.CoordinatesToCellName(col int, row int, abs ...bool) (string, error)

其他一些知识点
设置行样式
func (*excelize.File).SetColStyle(sheet string, columns string, styleID int) error

设置列宽
func (*excelize.File).SetColWidth(sheet string, startcol string, endcol string, width float64) error

设置行高
func (*excelize.File).SetRowHeight(sheet string, row int, height float64) error

创建表格
func (*excelize.File).AddTable(sheet string, hcell string, vcell string, format string) error

*/
