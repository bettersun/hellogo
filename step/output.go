package step

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

// OutJson 接口内容输出到json文件
func OutJson(file string, obj interface{}) error {
	// HTML原文输出
	bf := bytes.NewBuffer([]byte{})
	jsonEncoder := json.NewEncoder(bf)
	jsonEncoder.SetEscapeHTML(false)
	// 缩进
	jsonEncoder.SetIndent("", "  ")
	jsonEncoder.Encode(obj)

	err := ioutil.WriteFile(file, bf.Bytes(), os.ModePerm) // 覆盖所有Unix权限位（用于通过&获取类型位）
	if err != nil {
		log.Println(err)
	}

	return err
}

// SummaryToText 代码汇总转为字符串切片
func SummaryToText(summary StepSummary) []string {
	var s []string

	// 汇总
	headerAll := "统计文件总数\t总行数\t总空行数\t总注释行数\t总代码行数\t总有效行数"
	s = append(s, headerAll)
	stepAll := fmt.Sprintf("%v\t%v\t%v\t%v\t%v\t%v",
		summary.CountedFileCount, summary.TotalStep, summary.EmptyLineStep, summary.CommentStep, summary.SourceStep, summary.ValidStep)
	s = append(s, stepAll)
	s = append(s, "")

	// 注释定义不存在文件
	flatFile := "注释定义不存在文件："
	s = append(s, flatFile)
	for _, v := range summary.FlatFile {
		s = append(s, v)
	}
	s = append(s, "")

	// 未统计文件
	uncountedFile := "未统计文件："
	s = append(s, uncountedFile)
	errHeader := "文件全路径\t备考"
	s = append(s, errHeader)
	for _, v := range summary.StepInfo {
		if !v.Counted {
			stepFile := fmt.Sprintf("%v\t%v",
				v.File, v.ExInfo)
			s = append(s, stepFile)
		}
	}
	s = append(s, "")

	// 各文件代码行数
	s = append(s, "各文件代码行数：")
	s = append(s, "=====")
	stepHeader := "文件全路径\t文件名\t总行数\t空行数\t注释行数\t代码行数\t有效行数\t备考"
	s = append(s, stepHeader)
	for _, v := range summary.StepInfo {
		if v.Counted {
			stepFile := fmt.Sprintf("%v\t%v\t%v\t%v\t%v\t%v\t%v\t%v",
				v.File, v.FileName, v.TotalStep, v.EmptyLineStep, v.CommentStep, v.SourceStep, v.ValidStep, v.ExInfo)
			s = append(s, stepFile)
		}
	}

	return s
}

// WriteFile 将字符串切片写入文件，各字符串间添加换行符
func WriteFile(filename string, s []string) error {
	var buffer bytes.Buffer

	l := len(s)
	for i, v := range s {
		buffer.WriteString(v)
		if i < l-1 {
			buffer.WriteString("\n")
		}
	}

	err := ioutil.WriteFile(filename, buffer.Bytes(), 0644)
	if err != nil {
		log.Printf("写入文件发生错误 ： %v\n", err)
		return err
	}

	return nil
}
