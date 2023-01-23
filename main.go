package main

import (
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"strings"
	"time"
)

func main() {
	scheduler("./src")
}

func scheduler(dirPath string) {
	for range time.Tick(10 * time.Second) {
		files, err := searchPcapFiles(dirPath)
		if err != nil {
			fmt.Println(err)
			continue
		}
		if len(files) >= 2 {
			target := files[len(files)-2]
			if _, err := compress(dirPath+"/", target); err != nil {
				fmt.Println(err)
				continue
			}
		} else {
			continue
		}
	}
}

func compress(dirPath, pcapFilePath string) (string, error) {
	compFilePath := convertCompressFilePath(pcapFilePath)
	gzipFiles, err := searchGzipFiles(dirPath)
	if err == nil {
		for _, file := range gzipFiles {
			if file == compFilePath {
				return "", nil
			}
		}
	}
	dist, err := os.Create(dirPath + compFilePath)
	if err != nil {
		return "", err
	}
	defer dist.Close()
	compData, err := gzip.NewWriterLevel(dist, gzip.BestCompression)
	if err != nil {
		return "", err
	}
	defer compData.Close()
	src, err := os.Open(dirPath + pcapFilePath)
	if err != nil {
		return "", err
	}
	defer src.Close()
	if _, err := io.Copy(compData, src); err != nil {
		return "", err
	}
	return compFilePath, nil
}

func convertCompressFilePath(pcapFilePath string) string {
	splitPath := strings.Split(pcapFilePath, "/")
	fileName := splitPath[len(splitPath)-1]
	compFilePath := strings.Replace(pcapFilePath, fileName, fileName+".gz", 1)
	return compFilePath
}

func searchPcapFiles(dirPath string) ([]string, error) {
	files, err := os.ReadDir(dirPath)
	if err != nil {
		return []string{}, err
	}
	var res []string
	for _, file := range files {
		fileName := file.Name()
		if strings.Contains(fileName, ".pcap") && !strings.Contains(fileName, ".pcap.gz") {
			res = append(res, fileName)
		}
	}
	return res, nil
}

func searchGzipFiles(dirPath string) ([]string, error) {
	files, err := os.ReadDir(dirPath)
	if err != nil {
		return []string{}, nil
	}
	var res []string
	for _, file := range files {
		fileName := file.Name()
		if strings.Contains(fileName, ".pcap.gz") {
			res = append(res, fileName)
		}
	}
	return res, nil
}
