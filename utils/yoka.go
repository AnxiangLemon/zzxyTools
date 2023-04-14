package utils

import (
	"fmt"
	"io"
	"math/rand"
	"os"
	"path/filepath"
	"time"
)

var PackageDir = "yokaPackage/"

/*
文件数据
加解密
*/
func yokaxxxTea(filedata []byte, enc bool) []byte {

	//定义一个data
	var data []byte
	var fileLen int64 = int64(len(filedata)) //这个是加密前的长度 也就是明文的长度
	rand.Seed(time.Now().Unix())
	var num2 int64 = rand.Int63n(65535) //好像是个随机数 看是不是

	if enc {
		data = Cocos2d_abc_encrypt(fileLen, num2, filedata)
	} else {
		fileLen = Bin2Int(filedata[0:4])
		num2 = Bin2Int(filedata[4:8])
		data = filedata[8:]
		data, _ = Cocos2d_abc_decrypt(fileLen, num2, data)
	}
	return data
}

// EncryptPackData /* 读取文件 进行加解密
func EncryptPackData(path string, enc bool) {
	fi, err := os.Open(path)
	if err != nil {
		panic(err)
		return
	}

	defer fi.Close()

	filedata, err := io.ReadAll(fi)
	if err != nil {
		fmt.Printf("读取文件 %s 失败, err: %v", path, err)
	}
	//调用加解密
	data := yokaxxxTea(filedata, enc)
	if len(data) == 0 {
		return
	}
	qianzhui := "_解密"
	if enc {
		qianzhui = "_加密"
	}
	file, err := os.OpenFile(path+qianzhui, os.O_CREATE, 0)
	if err != nil {
		panic(err)
		return
	}
	defer file.Close()
	file.Write(data)

}

/*
解资源包 传入一个包
*/
func AssetsUnPackData(path string) {
	fi, _ := os.Open(path)
	defer fi.Close()

	filedata, err := io.ReadAll(fi)
	if err != nil {
		fmt.Printf("读取文件 %s 失败, err: %v", path, err)
	}
	ccUnPack := NewUnPack("little")
	ccUnPack.SetData(filedata)
	//flag
	ccUnPack.GetBin(24)
	/*
		0000000000AA117C | 48:8B9424 98000000       | mov rdx,qword ptr ss:[rsp+98]           | a0
		0000000000AA1184 | 8B72 04                  | mov esi,dword ptr ds:[rdx+4]            | a1
		0000000000AA1187 | 2B72 14                  | sub esi,dword ptr ds:[rdx+14]           | a1-a5 = esi 这就是文件索引
		0000000000AA118A | FFCE                     | dec esi                                 | 减一指令
		0000000000AA118C | 8972 04                  | mov dword ptr ds:[rdx+4],esi            | 赋值给a1
		0000000000AA118F | 8B4A 1C                  | mov ecx,dword ptr ds:[rdx+1C]           | a7给ecx
		0000000000AA1192 | 2B4A 14                  | sub ecx,dword ptr ds:[rdx+14]           | a7-a5 = ecx  0x10  文件名长度
		0000000000AA1195 | 8D79 FF                  | lea edi,qword ptr ds:[rcx-1]            | edi = rcx-1
		0000000000AA1198 | 897A 1C                  | mov dword ptr ds:[rdx+1C],edi           | 0xF赋值给a7
		0000000000AA119B | 44:8B42 18               | mov r8d,dword ptr ds:[rdx+18]           | a6
		0000000000AA119F | 44:2B42 14               | sub r8d,dword ptr ds:[rdx+14]           | a6-a5 = r8d 0x3cd= 973
		0000000000AA11A3 | 45:8D48 FF               | lea r9d,qword ptr ds:[r8-1]             | 3cc 这个减8 就是密文长度
		0000000000AA11A7 | 44:894A 18               | mov dword ptr ds:[rdx+18],r9d           |
		0000000000AA11AB | 44:8B4A 0C               | mov r9d,dword ptr ds:[rdx+C]            | a3
		0000000000AA11AF | 44:2B4A 14               | sub r9d,dword ptr ds:[rdx+14]           | a3-a5 = r9d 3fD  1021
		0000000000AA11B3 | 41:FFC9                  | dec r9d                                 | -1 3fc
		0000000000AA11B6 | 44:894A 0C               | mov dword ptr ds:[rdx+C],r9d            |
		0000000000AA11BA | 4C:8B5424 48             | mov r10,qword ptr ss:[rsp+48]           | 0x18
		0000000000AA11BF | 90                       | nop                                     |
		0000000000AA11C0 | 49:39F2                  | cmp r10,rsi                             | r10:runtime_zerobase
		0000000000AA11C3 | 0F85 DC020000            | jne sgs_pkg.AA14A5                      | 文件索引错误
		0000000000AA11C9 | 81FF 04010000            | cmp edi,104                             | 和0x104比较
		0000000000AA11CF | 0F87 BB020000            | ja sgs_pkg.AA1490                       | 文件名太长
		0000000000AA11D5 | 41:8D3408                | lea esi,qword ptr ds:[r8+rcx]           | 3cD+10 = 3dd
		0000000000AA11D9 | 8D76 1F                  | lea esi,qword ptr ds:[rsi+1F]           | 继续加1f 3fc
		0000000000AA11DC | 0F1F40 00                | nop dword ptr ds:[rax],eax              |
		0000000000AA11E0 | 44:39CE                  | cmp esi,r9d                             |
		0000000000AA11E3 | 0F85 94020000            | jne sgs_pkg.AA147D                      | 总长度错误
	*/

	//尾部32不需要了
	for ccUnPack.Length() > 32 {
		ccUnPack.GetBin(20)
		a5 := ccUnPack.GetInt()
		a6 := ccUnPack.GetInt()
		a7 := ccUnPack.GetInt()

		//fmt.Println("head解析", a0, a1, a2, a3, a4, a5, a6, a7)
		fileNameLen := a7 - a5 - 1 //最后会有个空格\0
		fileName := string(ccUnPack.GetBin(int(fileNameLen)))
		ccUnPack.GetByte()
		resLen := a6 - a5 - 1 //- 8 //这里需要把前面两个整数传过去 所以不能减8
		fmt.Println("正在提取", fileName, "len", resLen-8)

		//不解密直接保存
		luadata := ccUnPack.GetBin(int(resLen))

		//luadata := yokaxxxTea(ccUnPack.GetBin(int(resLen)), false)
		if len(luadata) == 0 {
			fmt.Println("\n数据解密失败", fileName)
			return
		}

		pkgWriteFile(PackageDir+fileName, luadata)

	}

	fmt.Println("正在提取完成！")
}

/*
组资源包 更新文件
*/
func AssetsPackData(path string) {
	//打包的文件名
	asstesname := "new_" + filepath.Base(path)

	// 检查文件是否存在
	if _, err := os.Stat(asstesname); os.IsNotExist(err) {
		fmt.Printf("File %s does not exist\n", asstesname)
	} else {
		// 如果文件存在，删除它
		//文件存在 先删除
		err := os.Remove(asstesname)
		if err != nil {
			fmt.Println(err)
			return
		}
	}

	// 打开压缩文件
	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	// 读取文件头
	flag := make([]byte, 24)
	_, err = io.ReadFull(file, flag)
	if err != nil {
		panic(err)
	}
	ccPack := NewPack("little")
	ccPack.SetBin(flag)
	var fileindex int64 = 24
	// 循环读取每个部分
	for {
		// 读取文件头
		header := make([]byte, 32)
		_, err := io.ReadFull(file, header)
		if err != nil {
			if err == io.EOF {
				break
			}
			panic(err)
		}

		// 解析文件头
		a0 := Bin2Int(header[:4])
		a1 := Bin2Int(header[4:8])
		a2 := Bin2Int(header[8:12])
		a3 := Bin2Int(header[12:16])

		a4 := Bin2Int(header[16:20])
		a5 := Bin2Int(header[20:24])
		a6 := Bin2Int(header[24:28])
		a7 := Bin2Int(header[28:])
		//fmt.Println("head解析", a0, a1, a2, a3, a4, a5, a6, a7)
		if a7 == 0 || a0 > 1 {
			//如果是文件的结尾
			appendToFile(asstesname, header)
			fmt.Println("资源文件更新完成！")
			return
		}
		fileNameLen := a7 - a5
		// 读取文件名
		fileNameBys := make([]byte, fileNameLen) //读取文件名和后四位
		_, err = io.ReadFull(file, fileNameBys)
		if err != nil {
			if err == io.EOF {
				break
			}
			panic(err)
		}
		fileName := string(fileNameBys[:fileNameLen-1])
		fmt.Println("正在更新", fileName)

		/*		0000000000AA1187 | 2B72 14                  | sub esi,dword ptr ds:[rdx+14]           | a1-a5 = esi 这就是文件索引
				0000000000AA118A | FFCE                     | dec esi                                 | 减一指令
				0000000000AA118C | 8972 04                  | mov dword ptr ds:[rdx+4],esi            | 赋值给a1
				0000000000AA118F | 8B4A 1C                  | mov ecx,dword ptr ds:[rdx+1C]           | a7给ecx
				0000000000AA1192 | 2B4A 14                  | sub ecx,dword ptr ds:[rdx+14]           | a7-a5 = ecx  0x10  文件名长度
				0000000000AA1195 | 8D79 FF                  | lea edi,qword ptr ds:[rcx-1]            | edi = rcx-1
				0000000000AA1198 | 897A 1C                  | mov dword ptr ds:[rdx+1C],edi           | 0xF赋值给a7
				0000000000AA119B | 44:8B42 18               | mov r8d,dword ptr ds:[rdx+18]           | a6
				0000000000AA119F | 44:2B42 14               | sub r8d,dword ptr ds:[rdx+14]           | a6-a5 = r8d 0x3cd= 973
				0000000000AA11A3 | 45:8D48 FF               | lea r9d,qword ptr ds:[r8-1]             | 3cc 这个减8 就是密文长度
				0000000000AA11A7 | 44:894A 18               | mov dword ptr ds:[rdx+18],r9d           |
				0000000000AA11AB | 44:8B4A 0C               | mov r9d,dword ptr ds:[rdx+C]            | a3
				0000000000AA11AF | 44:2B4A 14               | sub r9d,dword ptr ds:[rdx+14]           | a3-a5 = r9d 3fD  1021
				0000000000AA11B3 | 41:FFC9                  | dec r9d                                 | -1 3fc

		*/

		//移动文件指针 读取文件
		oldData := make([]byte, a6-a5-1)
		_, err = io.ReadFull(file, oldData)
		if err != nil {
			if err == io.EOF {
				break
			}
			panic(err)
		}

		//读取加密好的的数据  或者是没加密的数据
		luadata, err := pkgReadFile(PackageDir + fileName)
		if err != nil {
			panic(err)
		}
		lualen := len(luadata)
		//进行索引 文件计算 组成header
		a1 = a5 + int64(fileindex) + 1             //索引位置
		a6 = a5 + int64(lualen) + 1                //+ 8                //密文长度
		a3 = a5 + int64(lualen) + 1 + a7 - a5 + 32 // + 8  // 这个包体长度  这个地方要不要加8 取决于打包的文件前面是否是已经加密了
		//		pkhead :=Int2Bin(a0) + Int2Bin(a1) + Int2Bin(a2) + Int2Bin(a3) + Int2Bin(a4) + Int2Bin(a5) + Int2Bin(a6) + Int2Bin(a7)+header2
		ccPack.SetInt(int32(a0))
		ccPack.SetInt(int32(a1))
		ccPack.SetInt(int32(a2))
		ccPack.SetInt(int32(a3))
		ccPack.SetInt(int32(a4))
		ccPack.SetInt(int32(a5))
		ccPack.SetInt(int32(a6))
		ccPack.SetInt(int32(a7))
		ccPack.SetBin(fileNameBys)
		//进行数据替换
		//writeBinaryFile
		ccPack.SetBin(luadata)
		//将数据 写入文件
		//currentTimestamp := time.Now().Unix() strconv.FormatInt(currentTimestamp, 10)+
		fileindex, _ = appendToFile(asstesname, ccPack.GetAll())
		ccPack.Empty()
	}

}
