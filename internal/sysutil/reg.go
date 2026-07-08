package sysutil

import (
	"fmt"
	"os/exec"

	"golang.org/x/sys/windows/registry"
)

// ExportReg 导出注册表分支到文件。
func ExportReg(regBranch, destPath string) error {
	out, err := exec.Command("REG", "EXPORT", regBranch, destPath, "/y").CombinedOutput()
	if err != nil {
		return fmt.Errorf("导出注册表失败: %s", DecodeGBK(out))
	}
	return nil
}

// ImportReg 从文件导入注册表。
func ImportReg(regPath string) error {
	out, err := exec.Command("REG", "IMPORT", regPath).CombinedOutput()
	if err != nil {
		return fmt.Errorf("导入注册表失败: %s", DecodeGBK(out))
	}
	return nil
}

// DeleteReg 删除注册表分支。
func DeleteReg(regBranch string) error {
	out, err := exec.Command("reg", "delete", regBranch, "/f").CombinedOutput()
	if err != nil {
		return fmt.Errorf("删除注册表失败: %s", DecodeGBK(out))
	}
	return nil
}

// ReadHsrUID 读取崩坏：星穹铁道当前账号 UID。
func ReadHsrUID() (int64, error) {
	key, err := registry.OpenKey(registry.CURRENT_USER, `Software\miHoYo\崩坏：星穹铁道`, registry.QUERY_VALUE)
	if err != nil {
		return 0, fmt.Errorf("打开注册表项失败: %w", err)
	}
	defer key.Close()

	const valueName = "App_LastUserID_h2841727341"
	val, _, err := key.GetIntegerValue(valueName)
	if err != nil {
		return 0, fmt.Errorf("读取注册表值失败: %w", err)
	}
	return int64(val), nil
}
