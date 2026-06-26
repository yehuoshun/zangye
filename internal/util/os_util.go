// Package util 提供通用工具函数
package util

import (
	"os/exec"
	"runtime"
)

// OpenWith 使用指定程序打开文件
// 类比 Java 的 Desktop.open() 或 Runtime.exec()
// 参数：
//   - filePath: 要打开的文件路径
//   - exePath: 外部程序路径（为空则使用系统默认打开方式）
//
// 注意：当前仅支持 Windows 平台
func OpenWith(filePath, exePath string) error {
	var cmd *exec.Cmd

	if exePath != "" {
		// 使用指定程序打开
		cmd = exec.Command(exePath, filePath)
	} else {
		// 使用系统默认打开方式
		switch runtime.GOOS {
		case "windows":
			// Windows 下使用 cmd /c start 打开文件
			// 类比 Java 的 Desktop.getDesktop().open(file)
			cmd = exec.Command("cmd", "/c", "start", "", filePath)
		case "darwin":
			cmd = exec.Command("open", filePath)
		default: // linux
			cmd = exec.Command("xdg-open", filePath)
		}
	}

	return cmd.Start()
}
