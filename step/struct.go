package step

import (
	"regexp"
)

// Config 程序的配置文件
type Config struct {
	// 统计目标文件扩展名(.+小写，多个用逗号连接)
	CodeFileExt string `yaml:"codeFileExt"`
	// 忽略目录(多个用逗号连接)
	IgnorePath string `yaml:"ignorePath"`
	// 忽略文件(多个用逗号连接)
	IgnoreFile string `yaml:"ignoreFile"`
	// 结果文件类型
	// 0: txt 1:json
	ResultFileType string `yaml:"resultFileType"`
	// 结果文件名(不含扩展名)
	ResultFileName string `yaml:"resultFileName"`
}

// CommentDefine 代码注释定义
type CommentDefine struct {
	// 扩展名(多个)
	FileExtension []string `yaml:"fileExtension"`
	// 单行注释标志
	SingleLine []string `yaml:"singleLine"`
	// 多行注释开始
	MultiLineStart []string `yaml:"multiLineStart"`
	// 多行注释结束
	MultiLineEnd []string `yaml:"multiLineEnd"`
}

// CommentRegExp 代码注释统计用正则表达式
type CommentRegExp struct {
	// 扩展名
	FileExtension string
	// 有单行注释
	HasSingleLineMark bool
	// 有多行注释
	HasMultiLineMark bool

	// 空行正则表达式
	RegExEmptyLine *regexp.Regexp
	// 单行注释正则表达式
	RegExSingleLine []*regexp.Regexp
	// 写在单行的多行注释正则表达式 即使用多行注释开始和多行注释结束来写的一行注释
	RegExSingleLineStartEnd []*regexp.Regexp
	// 多行注释开始正则表达式
	RegExMultiLineStart []*regexp.Regexp
	// 多行注释结束正则表达式
	RegExMultiLineEnd []*regexp.Regexp
}

// StepInfo 代码行数信息
type StepInfo struct {
	CommentDefined bool   `json:"commentDefined"` // 存在注释配置 true:存在 false:不存在
	Counted        bool   `json:"counted"`        // 已统计标志 true:已统计 false:未统计
	ExInfo         string `json:"exInfo"`         // 扩展信息

	File          string `json:"file"`          // 文件全路径
	FileName      string `json:"fileName"`      // 文件名
	TotalStep     int    `json:"totalStep"`     // 总行数
	EmptyLineStep int    `json:"emptyLineStep"` // 空行数
	CommentStep   int    `json:"commentStep"`   // 注释行数
	SourceStep    int    `json:"sourceStep"`    // 代码行数
	ValidStep     int    `json:"validStep"`     // 有效行数(注释+代码)
}

// StepSummary 代码行数信息汇总
type StepSummary struct {
	FileCount int `json:"fileCount"` // 文件总数

	StepInfo      []StepInfo `json:"stepInfo"`      // 代码行数统计结果
	FlatFile      []string   `json:"flatFile"`      // 无注释定义文件
	UnCountedFile []string   `json:"unCountedFile"` // 未统计文件

	TotalStep     int `json:"totalStep"`     // 总行数
	EmptyLineStep int `json:"emptyLineStep"` // 空行总行数
	CommentStep   int `json:"commentStep"`   // 注释总行数
	SourceStep    int `json:"sourceStep"`    // 代码总行数
	ValidStep     int `json:"validStep"`     // 有效总行数(注释+代码)
}
