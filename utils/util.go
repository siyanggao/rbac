package utils

import (
	"strings"
)

func ParsePathFile(pathFile string) (path, fileName, subfix string) {
	pathIndex := strings.LastIndex(pathFile, "/")
	if pathIndex != -1 {
		path = pathFile[0 : pathIndex+1]
	}
	subfixIndex := strings.LastIndex(pathFile, ".")
	if subfixIndex == -1 {
		subfixIndex = len(pathFile)
	} else {
		subfix = pathFile[subfixIndex:]
	}
	fileName = pathFile[pathIndex+1 : subfixIndex]
	return
}

func IntToInt32Slice(source []int) (dst []int32) {
	for _, item := range source {
		dst = append(dst, int32(item))
	}
	return
}
