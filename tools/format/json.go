package format

import (
	"bytes"
	"encoding/json"
	// "encoding/json/jsontext"
	// jsonv2 "encoding/json/v2"
	"fmt"
)

const jsonIndent = "  "

// ToJSON marshals the input into a json string.
//
// If marshal fails, it falls back to fmt.Sprintf("%+v").
func ToJSON(e interface{}) string {
	ret, err := json.MarshalIndent(e, "", jsonIndent)
	if err != nil {
		return fmt.Sprintf("%+v", e)
	}
	return string(ret)
}

// FormatJSON formats the input json bytes with indentation.
//
// If Indent fails, it returns the unchanged input as string.
func FormatJSON(b []byte) string {
	var out bytes.Buffer
	err := json.Indent(&out, b, "", jsonIndent)
	if err != nil {
		return string(b)
	}
	return out.String()
}

// // ToJSON marshals the input into a json string.
// //
// // If marshal fails, it falls back to fmt.Sprintf("%+v").
// func ToJSONV2(e interface{}) string {
// 	opts := jsonv2.JoinOptions(
// 		//控制输出是否确定性排序（字典序字段顺序）
// 		jsonv2.Deterministic(true),
// 		// 控制字段名匹配是否忽略大小写（默认区分）
// 		jsonv2.MatchCaseInsensitiveNames(false),
// 		// 数字以字符串形式序列化（避免精度丢失）
// 		// jsonv2.StringifyNumbers(true),
// 		// 解码时遇到未知字段报错
// 		jsonv2.RejectUnknownMembers(true),
// 		// 解码时丢弃未知字段，不报错
// 		jsonv2.DiscardUnknownMembers(true),
// 		// 启用零值字段省略
// 		jsonv2.OmitZeroStructFields(false),
// 		jsontext.WithIndent("  "),     // <- 指定缩进（等价于 v1 的 indent）
// 		jsontext.WithIndentPrefix(""), // <- 可选：行前缀（前缀通常为空）
// 	)
// 	ret, err := jsonv2.Marshal(e, opts)
// 	if err != nil {
// 		return fmt.Sprintf("%+v", e)
// 	}
// 	return string(ret)
// }

// // FormatJSON formats the input json bytes with indentation.
// //
// // If Indent fails, it returns the unchanged input as string.
// func FormatJSONV2(b []byte) string {
// 	var out bytes.Buffer
// 	// err := jsonv2.MarshalIndent(obj, "", "  ", opts)
// 	// if err != nil {
// 	// 	return string(b)
// 	// }
// 	return out.String()
// }
