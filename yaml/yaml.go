package yaml

import (
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"log"
)

// YamlToStruct   Yaml文件转struct
//  file: yaml文件
//  s   : 要转换的 struct 变量的地址(调用处需要加&)
func YamlToStruct(file string, s interface{}) error {
	// 读取文件
	b, err := ioutil.ReadFile(file)
	if err != nil {
		log.Print(err)
		return err
	}

	// 转换成Struct
	err = yaml.Unmarshal(b, s)
	if err != nil {
		log.Printf("%v\n", err.Error())
	}

	return nil
}
