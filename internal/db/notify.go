// Package db 提供数据库连接管理
package db

import (
	"os/exec"
	"runtime"
	"strings"
)

// ShowError 在桌面环境下弹出错误消息框
// 跨平台实现：Windows 使用 msg, macOS 使用 osascript, Linux 使用 zenity/xmessage
// 返回 error 表示无法弹出对话框（无桌面环境或对应命令不存在）
func ShowError(title, message string) error {
	// 替换双引号避免 shell 错误
	message = strings.ReplaceAll(message, "\"", "'")
	title = strings.ReplaceAll(title, "\"", "'")

	switch runtime.GOOS {
	case "windows":
		// Windows 使用 msg 命令弹窗
		// msg * "标题" "内容"
		cmd := exec.Command("msg", "*", "/TIME:30", title+"\n"+message)
		return cmd.Run()
	case "darwin":
		// macOS 使用 osascript
		cmd := exec.Command("osascript", "-e",
			`display dialog "`+message+`" with title "`+title+`" buttons {"确定"} default button 1`)
		return cmd.Run()
	default:
		// Linux 尝试 zenity，不行用 xmessage
		cmd := exec.Command("zenity", "--error", "--title="+title, "--text="+message, "--width=400")
		if err := cmd.Run(); err != nil {
			cmd = exec.Command("xmessage", "-center", title+"\n"+message)
			return cmd.Run()
		}
		return nil
	}
}
