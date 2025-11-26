package csv

import (
	"encoding/csv"
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"

	"cloud.google.com/go/bigquery"
)

// 自动将 data 写入 CSV，自动生成表头
func WriteStructToCSV(path string, data interface{}) error {
	v := reflect.ValueOf(data)     // 获取 data 的反射对象，方便动态读取类型和字段
	if v.Kind() != reflect.Slice { // 检查 data 是否为切片
		return fmt.Errorf("data must be a slice") // 如果不是切片，返回错误
	}

	if v.Len() == 0 { // 检查切片是否为空
		return fmt.Errorf("data slice is empty") // 空切片无法生成表头或数据
	}

	// 创建 CSV 文件
	f, err := os.Create(path) // 在指定路径创建文件
	if err != nil {
		return err // 如果创建失败，返回错误
	}
	defer f.Close() // 函数返回时关闭文件

	w := csv.NewWriter(f) // 创建 CSV 写入器
	defer w.Flush()       // 写入完成时刷新缓冲区，确保内容写入文件

	// 获取元素类型
	elemType := v.Index(0).Type()    // 获取切片第一个元素的类型
	numFields := elemType.NumField() // 获取结构体字段数量

	// 自动生成表头
	header := make([]string, numFields) // 创建表头字符串切片
	for i := 0; i < numFields; i++ {
		header[i] = elemType.Field(i).Name // 使用字段名作为表头
	}

	if err := w.Write(header); err != nil { // 写入表头到 CSV
		return err
	}

	// 写每行数据
	for i := 0; i < v.Len(); i++ { // 遍历切片每个元素
		row := v.Index(i)                   // 获取当前元素
		record := make([]string, numFields) // 创建当前行字符串切片

		for j := 0; j < numFields; j++ { // 遍历每个字段
			field := row.Field(j)          // 获取字段值
			record[j] = formatField(field) // 转换为字符串（支持各种类型）
		}

		if err := w.Write(record); err != nil { // 写入 CSV
			return err
		}
	}

	return nil // 成功写入
}

// 处理各种类型并返回 Excel-friendly 字符串
func formatField(v reflect.Value) string {
	// 先处理 bigquery.NullXXX 类型
	switch v.Type().String() {
	case "bigquery.NullString":
		ns := v.Interface().(bigquery.NullString) // 转换为 bigquery.NullString
		if ns.Valid {                             // 有效值
			return ns.String() // 返回字符串
		}
		return "" // 无效返回空字符串
	case "bigquery.NullFloat64":
		nf := v.Interface().(bigquery.NullFloat64)
		if !nf.Valid {
			return "" // 无效返回空
		}
		// float64 转字符串，保留 Excel-friendly 格式
		s := strconv.FormatFloat(nf.Float64, 'f', 10, 64) // 保留10位小数
		s = strings.TrimRight(s, "0")                     // 去掉末尾多余 0
		s = strings.TrimRight(s, ".")                     // 去掉末尾小数点
		if s == "" {
			s = "0" // 防止空字符串
		}
		return s
	case "bigquery.NullDate":
		nd := v.Interface().(bigquery.NullDate)
		if nd.Valid {
			return nd.Date.String() // 返回日期字符串
		}
		return ""
	}

	// 基础类型处理
	switch v.Kind() {
	case reflect.String:
		return v.String()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return strconv.FormatInt(v.Int(), 10)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return strconv.FormatUint(v.Uint(), 10)
	case reflect.Float32, reflect.Float64:
		s := strconv.FormatFloat(v.Float(), 'f', 6, 64) // 保留6位小数
		s = strings.TrimRight(s, "0")
		s = strings.TrimRight(s, ".")
		if s == "" {
			s = "0"
		}
		return s
	case reflect.Bool:
		return strconv.FormatBool(v.Bool())
	case reflect.Struct:
		if t, ok := v.Interface().(time.Time); ok { // 判断是否为 time.Time
			return t.Format("2006-01-02 15:04:05") // 格式化时间
		}
	}

	return "" // 其他类型返回空字符串
}
