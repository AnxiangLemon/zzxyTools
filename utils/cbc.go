package utils

import (
	"bytes"
	"compress/zlib"
	"crypto/md5"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/xxtea/xxtea-go/xxtea"
	"io"
)

// 加密文件
func Cocos2d_abc_encrypt(num1, num2 int64, encryptedData []byte) []byte {

	//进行zlib压缩.
	// 创建一个字节数组缓冲区
	var buffer bytes.Buffer

	// 创建一个zlib压缩器
	writer := zlib.NewWriter(&buffer)
	// 将原始数据写入压缩器中
	// 将要压缩的数据写入压缩器中
	_, err := writer.Write(encryptedData)
	if err != nil {
		panic(err)
		return nil
	}
	// 关闭压缩器，确保所有数据都被刷新到缓冲区中  必须先关闭 不然数据没刷新进去！！！
	err = writer.Close()
	if err != nil {
		panic(err)
		return nil
	}
	content := buffer.Bytes()
	key := Cocos2d_cbc_Key(num1, num2)
	data := xxtea.Encrypt(content, key)
	if len(data) == 0 {
		fmt.Printf("数据加密失败 data =0")
		return nil
	}
	head := append(Int2Bin(num1), Int2Bin(num2)...)
	data = append(head, data...)
	return data
}

// 解密文件
func Cocos2d_abc_decrypt(num1, num2 int64, encryptedData []byte) ([]byte, error) {
	key := Cocos2d_cbc_Key(num1, num2)
	data := xxtea.Decrypt(encryptedData, key)
	//解压数据
	buf := new(bytes.Buffer)
	buf.Write(data)
	reader, err := zlib.NewReader(buf)
	if err != nil {
		//WriteFile("test111.lua", encryptedData)
		fmt.Printf("数据解压NewReader失败,len:%d , err: %v", len(data), err)
		return nil, err
	}
	defer reader.Close()
	content, err := io.ReadAll(reader)
	if err != nil && !errors.Is(err, io.ErrUnexpectedEOF) {
		fmt.Printf("数据解压失败, err: %v", err)
		return nil, err
	}

	return content, err
}

// 设置key
func Cocos2d_cbc_Key(num1, num2 int64) []byte {
	// 两个num顺序还不一样
	//s := fmt.Sprintf("k=%ds=%d%s", num2, num1, "WxZwQEmMOYNLXDInuA1PoxsKGPEVFY9d")
	// sprintf(s, "k=%ds=%d%s", a2, a1, "YFdjRYXNEYSIcOGhNcHimxjuKwpvrsRe");
	s := fmt.Sprintf("k=%ds=%d%s", num1, num2, "YFdjRYXNEYSIcOGhNcHimxjuKwpvrsRe")
	key := md5.Sum([]byte(s))
	return []byte(hex.EncodeToString(key[:]))
}

func Bin2Int(b []byte) int64 {
	nb := []byte{b[3], b[2], b[1], b[0]}
	bytesBuffer := bytes.NewBuffer(nb)
	var x int32
	binary.Read(bytesBuffer, binary.BigEndian, &x)
	return int64(x)
}

func Int2Bin(n int64) []byte {
	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b, uint64(n))
	//fmt.Println(b, "数子", n)
	return b[0:4]
}
