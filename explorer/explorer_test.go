package explorer

import (
	"bytes"
	"encoding/json"
	"log"
	"os"
	"runtime"
	"testing"
)

func TestExplorer(t *testing.T) {

	var option Option
	if runtime.GOOS == "darwin" {
		option.RootPath = []string{
			`/Users/xx/Documents/Develop/github.com/bettersun/hellogo`,
			`/Users/xx/Documents/Develop/github.com/bettersun/helloflutter/hellokiwi/`,
		}
	}
	if runtime.GOOS == "windows" {
		option.RootPath = []string{`E:\BS\Mac`}
	}

	option.SubFlag = true
	option.IgnorePath = []string{`.git`, `.svn`}
	option.IgnoreFile = []string{`.DS_Store`, `.gitignore`}

	result, err := Explorer(option)
	if err != nil {
		panic(err)
	}

	OutTerminal(result)

}

func OutTerminal(obj interface{}) {
	b, err := json.Marshal(obj)
	if err != nil {
		log.Println(err)
	}

	var out bytes.Buffer
	err = json.Indent(&out, b, "", "  ")

	out.WriteTo(os.Stdout)
}
