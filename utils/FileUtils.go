package utils

import (
	"fmt"
	"os"
	"path/filepath"
)

// 读入数据
func pkgReadFile(tmppath string) ([]byte, error) {
	return os.ReadFile(tmppath)
	//data,err :=
}

// 写入数据
func pkgWriteFile(tmppath string, data []byte) error {
	// 检查目录是否存在，不存在则创建
	if err := os.MkdirAll(filepath.Dir(tmppath), os.ModePerm); err != nil {
		fmt.Println("创建目录失败", err.Error())
		return err
	}
	err := os.WriteFile(tmppath, data, os.ModePerm)
	if err != nil {
		fmt.Println("写入文件失败", err.Error())
		return err
	}

	return nil
}

// 向文件追加数据 并且返回文件已知长度
func appendToFile(filename string, data []byte) (int64, error) {
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return 0, err
	}
	defer file.Close()

	_, err = file.Write(data)
	if err != nil {
		return 0, err
	}

	fileInfo, err := file.Stat()
	if err != nil {
		return 0, err
	}
	fileSize := fileInfo.Size()

	return fileSize, nil
}

// 分片读取文件数据
func readBinaryFileChunk(filePath string, chunkSize int64, offset int64) ([]byte, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Move file pointer to the specified offset
	_, err = file.Seek(offset, 0)
	if err != nil {
		return nil, err
	}

	// Create a buffer with the specified chunk size
	buf := make([]byte, chunkSize)

	// Read a chunk of binary data from the file
	n, err := file.Read(buf)
	if err != nil {
		return nil, err
	}

	// Return the chunk of binary data
	return buf[:n], nil
}

// 覆写数据
func writeBinaryFile(filePath string, data []byte, offset int64) error {
	// Open the file in read-write mode
	file, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE, os.ModePerm)
	if err != nil {
		return err
	}
	defer file.Close()

	// Move file pointer to the specified offset
	_, err = file.Seek(offset, 0)
	if err != nil {
		return err
	}

	// Write the data to the file
	_, err = file.Write(data)
	if err != nil {
		return err
	}

	return nil
}

func main2() {
	// Write data to a file starting from offset 1024
	data := []byte("Hello, World!")
	err := writeBinaryFile("example.bin", data, 1024)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
}

func main3() {
	// Read a chunk of binary data from a file starting from offset 1024 with a chunk size of 1024 bytes
	chunk, err := readBinaryFileChunk("example.bin", 1024, 1024)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Println(chunk)
}
