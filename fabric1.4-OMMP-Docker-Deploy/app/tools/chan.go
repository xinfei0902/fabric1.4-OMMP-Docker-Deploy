package tools

//SafeCloseQuit 线程启动函数 安全退出
func SafeCloseQuit(quit chan struct{}) {
	select {
	case _, ok := <-quit:
		if false == ok {
			return
		}
	default:
	}

	close(quit)
}
