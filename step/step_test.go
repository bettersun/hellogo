package step

import (
	"log"
	"testing"
)

func TestStep(t *testing.T) {
	file := []string{
		"./test/test.c",
		"./test/test.go",
		"./test/test.xxx",
		"./test/xxss.go",
		"./test/abcd.efg",
		"./test/test.sql"}

	summary := Step(file)
	err := OutJson("step.json", summary)
	if err != nil {
		log.Println("结果输出到文件时发生错误")
	}
}
