package main

import (
	"hellogo/step"
	yaml2 "hellogo/yaml"
	"log"
	"os"
	"path/filepath"
	"strings"
)

// 全局变量
var (
	codeFile    []string // 统计目标文件切片
	codeFileExt []string // 统计目标文件扩展名(.+小写，多个用逗号连接)为了在WalkFunc中使用，定义成全局变量
	ignorePath  []string // 忽略目录(多个用逗号连接) 为了在WalkFunc中使用，定义成全局变量
	ignoreFile  []string // 忽略文件(多个用逗号连接) 为了在WalkFunc中使用，定义成全局变量
)

// 常量
const (
	ConfigFile      = "./config.yaml" // 配置文件
	FileTypeText    = "0"             // 输出结果文件类型：文本
	FileTypeJson    = "1"             // 输出结果文件类型：json
	DefaultFileName = "step"          // 默认输出结果文件名（不含扩展名）
)

func main() {
	if len(os.Args) == 0 {
		log.Println("目标目录或文件未指定")
	}

	// 程序配置
	var config step.Config
	config.ResultFileType = FileTypeText
	config.ResultFileName = DefaultFileName

	// 读取配置
	err := yaml2.YamlToStruct(ConfigFile, &config)
	if err != nil {
		log.Println("读取程序配置发生错误，程序已停止运行。")
		return
	}

	// 统计目标文件扩展名
	ext := strings.Split(config.CodeFileExt, ",")
	// 去除空格或空字符串后赋给全局变量
	codeFileExt = step.RemoveSpaceEmpty(ext)
	// 忽略目录
	ignoreP := strings.Split(config.IgnorePath, ",")
	// 去除空格或空字符串后赋给全局变量
	ignorePath = step.RemoveSpaceEmpty(ignoreP)
	// 忽略文件
	ignoreF := strings.Split(config.IgnoreFile, ",")
	// 去除空格或空字符串后赋给全局变量
	ignoreFile = step.RemoveSpaceEmpty(ignoreF)

	// 程序运行参数指定的目标目录或文件
	var path []string
	for k, v := range os.Args {
		// 忽略程序文件本身
		if k == 0 {
			continue
		}

		p, err := os.Stat(v)
		if err != nil {
			// 错误文件也添加到统计目标文件切片
			codeFile = append(codeFile, v)
		} else {
			if p.IsDir() {
				// 目录添加到目录切片
				path = append(path, v)
			} else {
				// 文件添加到统计目标文件切片
				if len(codeFileExt) == 0 {
					codeFile = append(codeFile, v)
				} else if step.IsInSuffix(v, codeFileExt) {
					codeFile = append(codeFile, v)
				}
			}
		}
	}

	// 遍历目录下的文件
	for _, p := range path {
		filepath.Walk(p, filter)
	}

	// 对统计目标文件切片里的所有文件进行统计
	summary := step.Step(codeFile)

	// 输出结果类型
	if config.ResultFileType != FileTypeText && config.ResultFileType != FileTypeJson {
		log.Println("目前只支持输出统计结果到文本文件或 json 文件，默认输出到文本文件")
		config.ResultFileType = FileTypeText
	}

	stepFile := config.ResultFileName
	// 输出为文本文件
	if config.ResultFileType == FileTypeText {
		stepFile = stepFile + ".txt"
		s := step.SummaryToText(summary)
		err = step.WriteFile(stepFile, s)
	}
	// 输出为 json 文件
	if config.ResultFileType == FileTypeJson {
		stepFile = stepFile + ".json"
		err = step.OutJson(stepFile, summary)
	}

	if err != nil {
		log.Println("输出结果到文件时发生错误")
	}
}

// 过滤
func filter(path string, info os.FileInfo, err error) error {
	// 忽略目录
	if step.IsInSlice(ignorePath, info.Name()) {
		return filepath.SkipDir
	}

	// 文件且不在忽略文件中，添加到全局变量
	if !info.IsDir() && !step.IsInSlice(ignoreFile, info.Name()) {
		if len(codeFileExt) == 0 {
			codeFile = append(codeFile, path)
		} else if step.IsInSuffix(path, codeFileExt) {
			codeFile = append(codeFile, path)
		}
	}

	return err
}
