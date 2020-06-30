package main

import (
	cfg "TextReplacementTool/utils"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

var (
	configpath      string
	sourcePath      string
	targetPath      string
	fileType        string
	oldString       string
	newString       string
	isProcessSubDir bool
	err             error
)

func main() {
	path := "./config/config.ini"
	config := new(cfg.Config)
	config.InitConfig(path)
	sourcePath = config.Read("replaceconfig", "sourcePath")
	// targetPath = config.Read("replaceconfig", "targetPath")
	fileType = config.Read("replaceconfig", "fileType")
	oldString = config.Read("replaceconfig", "oldString")
	newString = config.Read("replaceconfig", "newString")

	isProcessSubDir, err = strconv.ParseBool(config.Read("replaceconfig", "isProcessSubDir"))
	if err != nil {
		panic(err)
	}
	xfiles, _ := GetAllFiles(sourcePath)

	for _, v := range xfiles {

		content, err := ioutil.ReadFile(v)
		if err != nil {
			fmt.Println("err:", err)
			panic(err)
		}
		newcontent := string(content)
		var isreplace bool
		if strings.Contains(newcontent, oldString) {
			newcontent = strings.Replace(newcontent, oldString, newString, -1)
			isreplace = true
		}

		if isreplace {
			func(file string, content []byte) {
				nf, err := os.OpenFile(file, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
				if err != nil {
					fmt.Println("err:", file, err)
					panic(err)
				}
				defer nf.Close()
				_, err = nf.WriteString(newcontent)
				if err != nil {
					fmt.Println("err:", err)
					panic(err)
				}
			}(v, content)
			fmt.Printf("%v done.\n", v)
		}

	}
}

func GetAllFiles(dirPth string) (files []string, err error) {
	var dirs []string
	dir, err := ioutil.ReadDir(dirPth)
	if err != nil {
		return nil, err
	}

	PthSep := string(os.PathSeparator)
	//suffix = strings.ToUpper(suffix) //忽略后缀匹配的大小写

	for _, fi := range dir {
		if fi.IsDir() && isProcessSubDir { // 目录, 递归遍历
			dirs = append(dirs, dirPth+PthSep+fi.Name())
			GetAllFiles(dirPth + PthSep + fi.Name())
		} else {
			// 过滤指定格式
			ok := strings.HasSuffix(fi.Name(), fileType)
			if ok {
				files = append(files, dirPth+PthSep+fi.Name())
			}
		}
	}

	// 读取子目录下文件
	if isProcessSubDir {
		for _, table := range dirs {
			temp, _ := GetAllFiles(table)
			for _, temp1 := range temp {
				files = append(files, temp1)
			}
		}
	}

	return files, nil
}
