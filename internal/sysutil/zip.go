package sysutil

import (
	"fmt"
	"os/exec"
)

// ZipDir 使用 PowerShell Compress-Archive 压缩目录。
func ZipDir(srcPath, destPath string) error {
	out, err := exec.Command("powershell", "Compress-Archive",
		"-Path", srcPath, "-DestinationPath", destPath, "-Force").CombinedOutput()
	if err != nil {
		return fmt.Errorf("压缩文件夹失败: %s", DecodeGBK(out))
	}
	return nil
}

// Unzip 使用 PowerShell Expand-Archive 解压文件。
func Unzip(zipPath, destPath string) error {
	out, err := exec.Command("powershell", "Expand-Archive",
		"-Path", zipPath, "-DestinationPath", destPath, "-Force").CombinedOutput()
	if err != nil {
		return fmt.Errorf("解压文件失败: %s", DecodeGBK(out))
	}
	return nil
}
