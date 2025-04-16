//go:build js && wasm
// +build js,wasm

package utils

import (
	"github.com/openimsdk/tools/utils"
	"syscall/js"
	"time"
)

func ReloadProgress() {
	js.Global().Get("location").Call("reload")
}

func ConfirmWithTimeout(msg string) {
	done := make(chan bool, 1)

	// 启动 10 秒超时定时器
	go func() {
		time.Sleep(10 * time.Second)
		select {
		case done <- false:
			// 自动触发 reload
			ReloadProgress()
		default:
			// 已处理，无需处理超时
		}
	}()

	// 在主线程中调用 JS 的 confirm（是同步的）
	confirmed := js.Global().Call("confirm", msg).Bool()

	// 尝试写入 channel，防止超时协程也写
	select {
	case done <- true:
		if confirmed {
			ReloadProgress()
		}
	default:
		// 超时已经触发了 reload，无需再执行
	}
}

type EventData struct {
	Event       string      `json:"event"`
	ErrCode     int32       `json:"errCode"`
	ErrMsg      string      `json:"errMsg"`
	Data        interface{} `json:"data,omitempty"`
	OperationID string      `json:"operationID"`
}

func SendGlobalNotification(event EventData) {
	js.Global().Call("sendGlobalNotification", utils.StructToJsonString(event))
}

func ReloadLib() {
	event := EventData{
		Event: "OnNeedReloadLib",
	}
	SendGlobalNotification(event)
}
