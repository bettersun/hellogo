package step

import (
	"bufio"
	yaml2 "hellogo/yaml"
	"io"
	"log"
	"os"
	"regexp"
	"strings"
)

// DefineFile 注释定义文件
const DefineFile = "./define.yaml"

// Step 统计多个代码文件的代码行数
func Step(file []string) StepSummary {
	var summary StepSummary

	// 读取注释定义
	var commentDef []CommentDefine
	err := yaml2.YamlToStruct(DefineFile, &commentDef)
	if err != nil {
		log.Printf("读取注释定义发生错误：%v \n", err)
		return summary
	}

	// 转换为正则表达式
	commentRegExp := ToCommentRegExp(commentDef)

	// 统计代码行数
	stepInfo := CountAll(file, &commentRegExp)

	// 汇总
	summary = Summary(stepInfo)
	return summary
}

// Summary 代码行数统计结果
func Summary(stepInfo []StepInfo) StepSummary {
	// 代码行数统计信息汇总
	var stepSummary StepSummary
	// 各文件的统计信息
	stepSummary.StepInfo = stepInfo

	// 合计
	for _, step := range stepInfo {
		// 总行数
		stepSummary.TotalStep = stepSummary.TotalStep + step.TotalStep
		// 空行总行数
		stepSummary.EmptyLineStep = stepSummary.EmptyLineStep + step.EmptyLineStep
		// 注释总行数
		stepSummary.CommentStep = stepSummary.CommentStep + step.CommentStep
		// 代码总行数
		stepSummary.SourceStep = stepSummary.SourceStep + step.SourceStep
		// 有效总行数(注释+代码)
		stepSummary.ValidStep = stepSummary.ValidStep + step.ValidStep

		// 无注释定义
		if !step.CommentDefined {
			stepSummary.FlatFile = append(stepSummary.FlatFile, step.File)
		}
		// 未统计文件
		if !step.Counted {
			stepSummary.UnCountedFile = append(stepSummary.UnCountedFile, step.File)
		}

		// 统计文件总数
		if step.Counted {
			stepSummary.CountedFileCount = stepSummary.CountedFileCount + 1
		}
	}

	return stepSummary
}

// CountAll 统计多个文件的代码行数
func CountAll(file []string, mCommentRegExp *map[string]CommentRegExp) []StepInfo {
	var stepInfo []StepInfo
	for _, f := range file {
		Count(f, mCommentRegExp, &stepInfo)
	}
	return stepInfo
}

// Count 统计单个文件的代码行数
func Count(file string, mCommentRegExp *map[string]CommentRegExp, infoAll *[]StepInfo) {
	// 获取代码文件对应的注释定义
	var cmtRegExp CommentRegExp
	isDefined := false
	for k, v := range *mCommentRegExp {
		if strings.HasSuffix(file, k) {
			cmtRegExp = v
			isDefined = true
			break
		}
	}

	var info StepInfo
	if isDefined {
		// 存在注释定义
		info.CommentDefined = true
	} else {
		// 无该类型代码文件注释定义时，使用默认正则表达式，只统计非空行数
		defaultRegEx(&cmtRegExp)
		// 无注释标志定义
		info.CommentDefined = false
		info.ExInfo = "无注释定义, 不统计注释行数。"
	}

	// 文件全路径
	info.File = file
	// 文件信息
	f, err := os.Stat(file)
	if err != nil {
		info.ExInfo = "文件不存在 或 读取文件信息发生错误。"
		*infoAll = append(*infoAll, info)
		return
	}
	// 文件名
	info.FileName = f.Name()

	// 统计
	count(&cmtRegExp, &info, infoAll)
}

// 统计
func count(cmtRegExp *CommentRegExp, info *StepInfo, infoAll *[]StepInfo) {
	var totalStep int     // 总行数
	var emptyLineStep int // 空行数
	var commentStep int   // 注释
	var sourceStep int    // 代码行数

	var isMultiLineComment bool // 多行注释统计中标志
	var matchIndex int          // 多行注释结束 正则表达式 下标

	var isMatch bool // 正则表达式匹配标志

	// 打开文件
	f, err := os.Open(info.File)
	if err != nil {
		info.ExInfo = "打开文件发生错误"
		*infoAll = append(*infoAll, *info)
		return
	}
	defer f.Close()

	reader := bufio.NewReader(f)
	// 按行读取后统计
	for {
		// 读取文件内容行
		line, err := reader.ReadString('\n')

		// 读取文件内容行发生错误,但非文件结尾
		if err != nil && err != io.EOF {
			info.ExInfo = "读取文件内容行发生错误"
			return
		}

		// 总行数
		totalStep++

		if !isMultiLineComment { // 不是多行注释
			if cmtRegExp.RegExEmptyLine.MatchString(line) { //空行
				emptyLineStep++
			} else if _, isMatchSingle := MatchIn(line, cmtRegExp.RegExSingleLine); cmtRegExp.HasSingleLineMark && isMatchSingle { // 单行注释
				commentStep++
			} else if _, isMatchSingleStartEnd := MatchIn(line, cmtRegExp.RegExSingleLineStartEnd); cmtRegExp.HasMultiLineMark && isMatchSingleStartEnd { // 多行注释标志开始结束的单行注释
				commentStep++
			} else if matchIndex, isMatch = MatchIn(line, cmtRegExp.RegExMultiLineStart); cmtRegExp.HasMultiLineMark && isMatch { // 多行注释 开始
				commentStep++
				isMultiLineComment = true
			}
		} else if cmtRegExp.HasMultiLineMark { // 有多行注释
			if isMultiLineComment {
				if cmtRegExp.RegExEmptyLine.MatchString(line) { //多行注释里的空行
					emptyLineStep++
				} else {
					//多行注释
					commentStep++
					//多行注释结束
					if cmtRegExp.RegExMultiLineEnd[matchIndex].MatchString(line) { //多行注释 结束
						isMultiLineComment = false
					}
				}
			}
		}

		// 文件结尾
		if err == io.EOF {
			break
		}
	}

	// 代码行数
	sourceStep = totalStep - commentStep - emptyLineStep
	// 总行数
	info.TotalStep = totalStep
	// 空行数
	info.EmptyLineStep = emptyLineStep
	// 代码行数
	info.SourceStep = sourceStep
	// 注释行数
	info.CommentStep = commentStep
	// 有效行数（总行数 - 空行数 即：注释行数 + 代码行数）
	info.ValidStep = totalStep - emptyLineStep
	// 已统计
	info.Counted = true

	*infoAll = append(*infoAll, *info)
}

// ToCommentRegExp 注释定义转注释正则表达式
func ToCommentRegExp(cmtDef []CommentDefine) map[string]CommentRegExp {
	// 相同注释定义，按照文件类型整理到Map
	var cmtRegex map[string]CommentRegExp
	cmtRegex = make(map[string]CommentRegExp)
	for _, v := range cmtDef {
		for _, ext := range v.FileExtension {
			_, ok := cmtRegex[ext]
			if !ok {
				cmtRegex[ext] = ToRegExp(ext, v)
			}
		}
	}

	return cmtRegex
}

// ToRegExp 转换为正则表达式
func ToRegExp(ext string, def CommentDefine) CommentRegExp {
	var cmtRegExp CommentRegExp

	lenSingle := len(def.SingleLine)
	lenMuilti := len(def.MultiLineStart)

	// 分配切片空间
	cmtRegExp.RegExSingleLine = make([]*regexp.Regexp, lenSingle)
	cmtRegExp.RegExSingleLineStartEnd = make([]*regexp.Regexp, lenMuilti)
	cmtRegExp.RegExMultiLineStart = make([]*regexp.Regexp, lenMuilti)
	cmtRegExp.RegExMultiLineEnd = make([]*regexp.Regexp, lenMuilti)

	// 扩展名
	cmtRegExp.FileExtension = ext
	// 有单行注释
	if !IsAllEmpty(def.SingleLine) {
		cmtRegExp.HasSingleLineMark = true
	}
	// 有多行注释
	if !IsAllEmpty(def.MultiLineStart) {
		cmtRegExp.HasMultiLineMark = true
	}

	// 空行正则表达式
	cmtRegExp.RegExEmptyLine = regexp.MustCompile(`^[\s]*$`)

	// 单行注释正则表达式
	for i, v := range def.SingleLine {

		cmtRegExp.RegExSingleLine[i] =
			regexp.MustCompile(`^[\s]*` + v + `.*`)
	}

	// 单行中使用多行注释的正则表达式 注释开始符和注释结束符在同一行
	for i, _ := range def.MultiLineStart {
		start := Escape(def.MultiLineStart[i])
		end := Escape(def.MultiLineEnd[i])

		cmtRegExp.RegExSingleLineStartEnd[i] =
			regexp.MustCompile(`^[\s]*(` + start + `).*(` + end + `)[\s]*$`)
	}

	// 多行注释开始正则表达式
	for i, v := range def.MultiLineStart {
		v = Escape(v)
		cmtRegExp.RegExMultiLineStart[i] =
			regexp.MustCompile(`^[\s]*(` + v + `).*`)
	}

	// 多行注释结束正则表达式
	for i, v := range def.MultiLineEnd {
		v = Escape(v)
		cmtRegExp.RegExMultiLineEnd[i] =
			regexp.MustCompile(`.*(` + v + `)[\s]*$`)
	}

	return cmtRegExp
}

// 默认正则表达式（不统计注释，只统空行数和非空行数）
func defaultRegEx(cmtRegExp *CommentRegExp) {
	// 有单行注释标志
	cmtRegExp.HasSingleLineMark = true
	// 有多行注释标志
	cmtRegExp.HasMultiLineMark = false
	// 空行正则表达式(字符串使用锐音符(`)不需要转义))
	cmtRegExp.RegExEmptyLine = regexp.MustCompile(`^[\s]*$`)
}

// MatchIn 对字符串判断是否匹配注释正则表达式
func MatchIn(line string, regexList []*regexp.Regexp) (matchIndex int, isMatch bool) {
	isMatch = false
	for i, v := range regexList {
		if v.MatchString(line) {
			isMatch = true
			matchIndex = i
			break
		}
	}

	return matchIndex, isMatch
}

// Escape 转义
func Escape(sRegExp string) string {
	v := sRegExp
	v = strings.ReplaceAll(v, `*`, `\*`)
	v = strings.ReplaceAll(v, `/`, `\/`)
	return v
}
