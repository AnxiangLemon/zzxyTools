package main

import (
	"flag"
	"fmt"
	"game/utils"
)

/*
*
2023年3月3日15:52:19
*/
func main() {

	fmt.Println("yoka游戏资源解析工具-zzxy 版本:1.0.1 ")

	updateFileDecrypt := flag.String("d", "", "解密文件")
	updateFileEncrypt := flag.String("e", "", "加密文件")
	packageFileCreate := flag.String("p", "", "资源打包")
	packageFileExtract := flag.String("u", "", "资源解包")
	flag.Parse()
	if *updateFileDecrypt != "" {
		fmt.Printf("解密文件: %s\n", *updateFileDecrypt)
		utils.EncryptPackData(*updateFileDecrypt, false)
	} else if *updateFileEncrypt != "" {
		fmt.Printf("加密文件: %s\n", *updateFileEncrypt)
		utils.EncryptPackData(*updateFileEncrypt, true)
	} else if *packageFileCreate != "" {
		fmt.Printf("资源打包: %s\n", *packageFileCreate)
		utils.AssetsPackData(*packageFileCreate)
	} else if *packageFileExtract != "" {
		fmt.Printf("资源解包: %s\n", *packageFileExtract)
		utils.AssetsUnPackData(*packageFileExtract)
	} else {
		fmt.Println("Usage:")
		flag.PrintDefaults()

	}
}
