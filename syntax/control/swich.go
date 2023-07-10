package main

func Switch(status int) string {
	switch status {
	case 0:
		return "初始化"
	case 1:
		return "运行中"
	default:
		return "未知状态"
	}
}
