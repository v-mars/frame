package utils

import (
	"compress/gzip"
	"io"
	"os"
	"time"
)

func FileExists(name string) bool {
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

func EnsureDirExist(name string) error {
	if !FileExists(name) {
		return os.MkdirAll(name, os.ModePerm)
	}
	return nil
}

func GzipCompressFile(srcPath, dstPath string) error {
	sf, err := os.Open(srcPath)
	if err != nil {
		return err
	}
	defer func(sf *os.File) {
		err := sf.Close()
		if err != nil {

		}
	}(sf)
	df, err := os.Create(dstPath)
	if err != nil {
		return err
	}
	defer func(df *os.File) {
		err := df.Close()
		if err != nil {

		}
	}(df)
	writer := gzip.NewWriter(df)
	writer.Name = dstPath
	writer.ModTime = time.Now().UTC()
	_, err = io.Copy(writer, sf)
	if err != nil {
		return err
	}
	if err := writer.Close(); err != nil {
		return err
	}
	return nil
}

func Sum(i []int) int {
	sum := 0
	for _, v := range i {
		sum += v
	}
	return sum
}

func Abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func CurrentUTCTime() string {
	return time.Now().UTC().Format("2006-01-02 15:04:05 +0000")
}
