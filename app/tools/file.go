package tools

import (
	"errors"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	yaml "gopkg.in/yaml.v2"
)

//STDPath 检测是否绝对地址 不是前缀加当前目录
func STDPath(name string) string {
	if len(name) == 0 {
		return name
	}
	if false == filepath.IsAbs(name) {
		base, _ := filepath.Split(os.Args[0])
		if true == filepath.IsAbs(base) {
			name = filepath.Join(base, name)
		}
	}
	return name
}

//CopyFileOne 拷贝文件
func CopyFileOne(dst, src string) (err error) {
	info, err := os.Stat(dst)
	if err == nil || false == os.IsNotExist(err) {
		if err == nil {
			err = os.ErrExist
		}
		return
	}

	info, err = os.Stat(src)
	if err != nil {
		return
	}
	if info.IsDir() {
		err = errors.New("not support yet")
		return
	}

	f, err := os.Open(src)
	if err != nil {
		return
	}
	defer f.Close()

	fo, err := os.OpenFile(dst, os.O_RDWR|os.O_CREATE|os.O_TRUNC, os.ModePerm&info.Mode())
	if err != nil {
		return
	}
	defer fo.Close()

	_, err = io.Copy(fo, f)
	if err != nil {
		return
	}

	fo.Sync()

	return nil
}

//CopyFolder 拷贝目录
func CopyFolder(dst, src string) (err error) {
	n := len(src)
	err = filepath.Walk(src, func(path string, f os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		child := filepath.Join(dst, path[n:])

		if f.IsDir() {
			err = os.MkdirAll(child, 0777)
			return err
		}

		err = CopyFileOne(child, path)
		return err
	})
	return
}

//WriteYamlFile 用来往 .yaml格式文件写入的
func WriteYamlFile(filename string, obj interface{}) (err error) {
	buff, err := yaml.Marshal(obj)
	if err != nil {
		return
	}
	err = ioutil.WriteFile(filename, buff, 0644)
	return err
}
