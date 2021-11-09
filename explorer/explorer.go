package explorer

import (
	"io/ioutil"
	"log"
	"os"
	"path"
	"strings"
)

// Option 遍历选项
type Option struct {
	RootPath   []string `yaml:"rootPath"`   // 目标根目录
	SubFlag    bool     `yaml:"subFlag"`    // 遍历子目录标志 true: 遍历 false: 不遍历
	IgnorePath []string `yaml:"ignorePath"` // 忽略目录
	IgnoreFile []string `yaml:"ignoreFile"` // 忽略文件
}

// Node 树节点
type Node struct {
	Name     string  `json:"name"`     // 目录（或文件）名
	Path     string  `json:"path"`     // 目录（或文件）完整路径
	Children []*Node `json:"children"` // 目录下的文件或子目录
	IsDir    bool    `json:"isDir"`    // 是否为目录 true: 是目录 false: 不是目录
}

// Explorer 遍历多个目录
// 		option : 遍历选项
// 		tree : 遍历结果
func Explorer(option Option) (Node, error) {
	// 根节点
	var root Node

	// 多个目录搜索
	for _, p := range option.RootPath {
		// 空目录跳过
		if strings.TrimSpace(p) == "" {
			continue
		}

		var child Node
		// 目录路径
		child.Path = p
		// 递归
		explorerRecursive(&child, &option)

		root.Children = append(root.Children, &child)
	}

	return root, nil
}

// 递归遍历目录
// 		node : 目录节点
// 		option : 遍历选项
func explorerRecursive(node *Node, option *Option) {
	// 节点的信息
	p, err := os.Stat(node.Path)
	if err != nil {
		log.Println(err)
		return
	}

	// 目录（或文件）名
	node.Name = p.Name()
	// 是否为目录
	node.IsDir = p.IsDir()

	// 非目录，返回
	if !p.IsDir() {
		return
	}

	// 目录中的文件和子目录
	sub, err := ioutil.ReadDir(node.Path)
	if err != nil {
		info := "目录不存在，或打开错误。"
		log.Printf("%v: %v", info, err)
		return
	}

	for _, f := range sub {
		tmp := path.Join(node.Path, f.Name())
		var child Node
		// 完整子目录
		child.Path = tmp
		// 是否为目录
		child.IsDir = f.IsDir()

		// 目录
		if f.IsDir() {
			//查找子目录
			if option.SubFlag {
				// 不在忽略目录中的目录，进行递归查找
				if !IsInSlice(option.IgnorePath, f.Name()) {
					node.Children = append(node.Children, &child)
					explorerRecursive(&child, option)
				}
			}
		} else { // 文件
			// 非忽略文件，添加到结果中
			if !IsInSlice(option.IgnoreFile, f.Name()) {
				child.Name = f.Name()
				node.Children = append(node.Children, &child)
			}
		}
	}
}

// IsInSlice 判断目标字符串是否是在切片中
func IsInSlice(slice []string, s string) bool {
	if len(slice) == 0 {
		return false
	}

	isIn := false
	for _, f := range slice {
		if f == s {
			isIn = true
			break
		}
	}

	return isIn
}
