package main

import (
	"encoding/base64"
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/openimsdk/openim-sdk-core/v3/internal/interaction"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/constant"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/sdkerrs"
	"github.com/openimsdk/protocol/sdkws"
	"github.com/sqzxcv/glog"
)

var (
	target = "H4sIAAAAAAAA/0TKTYvTQBgH8DzbuoW4RVhwBfEQAqIoyLxmMnvUXbCHRelB9BQmM09iRJN2Jr4UEQsi+Dl68nt49OLdTyD4GVIJCHv6v/ArPk8ApksMK9hvIziE+RLXC4dt31QNephGMLsI9aK1Ho4iuPpkhd70TdcuzsY9O/f+UedwdIfn3l+EerynZ6Y3EEdRNMDt/RaGa7uDGx2Tm9evNm1YUyWY0DrLCBN0MgDEA8TzUBdUEJERIrk6/gvxH4hvccooKk1LK5gxKDOTuYoJZ60kN+NLnyZZ5aTlVFlqFHE81yXnNpeC5YrltuT3EkXQaZtJ6QSTFaVWGZs7wrUqmcyFJVdO58+a0PSdp0ToXD88amxh3pne+ILQx5On7jmWJx9T27U9tn16miLHgMFV6afV9Q8/vv78/ptt4deXMb8B7OBOfPKie5u8NCExSYvvkzcYgqnxwfHsf0sP7tO7EP0LAAD//0RWAWKIAQAA"
	//target = "H4sIAAAAAAAA/3ySz0ocMRzHJzoojbW2C/1z8CA5Fg87M5mZjD2ulm5BK4rQw4Dkz290YDfZJllxEekehEKfoPToqW/S0lfoC/Qp1jKzFnooewrfD5/fNz9ITj8uIxQegRuhu2mAVtD6EXzoK9C+rmqwKAzQ6r4762tp0cMArb0bgeW+Nrq/2+TVPWt7RkHjrexZu+/OGhzucs8RDoJgtqTvpmi2cbv03NTp5dDzSeGinMZxzoqc5rRYni297XwP8St9mqQiVlmaxowyyiCOklwpoRiniThNWcx5LkRGqYICiqQq8oLKNI+k6HwK8U2INxcVdDYXNZAtkVZpl3PKWVYwUVAmaV5lWSZkFhdVCi+3KiYyJVgXoqhikrNUdWMJuZBJl4kEijfo8Ad6/wWLz+iKKPC8HpAdclWSyprhiQPb3y3JTkkWbVmS7ZJ484+9aOvWdp77sSvJTrRdktGA+8rYYTucXpPry5vfP7/9iqdofn5FG3htrC1w1TNj7R8HLTFVNag1HI7deUMeYTwC62rnYa48wx1p9AVY1z7+yUhxD/emA63AHk+0bMAmfjEHvf/66/hB7Y5Bq3131uQ1vHpeO2/spElP8ZPaHRh/YJrvJ9vZe26Byya+tmbY4/L87+3a+ENbX7T1tygI/gQAAP//7deSmtECAAA="
)

func main() {
	// 1. Base64 解码
	decodedData, err := decodeBase64(target)
	if err != nil {
		fmt.Println("Base64 解码错误:", err)
		return
	}

	// 2. Gzip 解压
	handleMessage(decodedData)
}

func handleMessage(message []byte) error {
	var decompressErr error
	com := interaction.NewGzipCompressor()
	message, decompressErr = com.DeCompress(message)
	if decompressErr != nil {
		fmt.Println("DeCompress failed", decompressErr, message)
		return sdkerrs.ErrMsgDeCompression
	}
	var wsResp interaction.GeneralWsResp
	err := interaction.NewGobEncoder().Decode(message, &wsResp)
	if err != nil {
		fmt.Println("decodeBinaryWs err", err, "message", message)
		return sdkerrs.ErrMsgDecodeBinaryWs
	}
	//ctx := context.WithValue(c.ctx, "operationID", wsResp.OperationID)
	fmt.Println("recv msg", "errCode", wsResp.ErrCode, "errMsg", wsResp.ErrMsg,
		"reqIdentifier", wsResp.ReqIdentifier)
	switch wsResp.ReqIdentifier {
	case constant.PushMsg:
		//if err = c.doPushMsg(ctx, wsResp); err != nil {
		//	log.ZError(ctx, "doWSPushMsg failed", err, "wsResp", wsResp)
		//}
		doPushMsg(wsResp)
	case constant.LogoutMsg:
		fmt.Println("client logout")
	case constant.KickOnlineMsg:
		fmt.Println("kick online")
	case constant.GetNewestSeq:
		fallthrough
	case constant.PullMsgBySeqList:
		fallthrough
	case constant.SendMsg:
		fallthrough
	case constant.SendSignalMsg:
		fallthrough
	case constant.WsHeartbeat:
		fallthrough
	case constant.SetBackgroundStatus:
		fmt.Println("msg type:", wsResp.ReqIdentifier)
	default:
		// log.Error(wsResp.OperationID, "type failed, ", wsResp.ReqIdentifier)
		return sdkerrs.ErrMsgBinaryTypeNotSupport
	}
	return nil
}

// decodeBase64 解码 Base64 字符串
func decodeBase64(encoded string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(encoded)
}

func doPushMsg(wsResp interaction.GeneralWsResp) error {
	var msg sdkws.PushMessages
	err := proto.Unmarshal(wsResp.Data, &msg)
	if err != nil {
		return err
	}
	fmt.Println("push msg", msg.String())
	for _, v := range msg.Msgs {
		for _, m := range v.Msgs {
			fmt.Println("Msgs: content:", string(m.Content))
		}
	}

	for _, v := range msg.NotificationMsgs {
		for _, m := range v.Msgs {
			glog.Info("NotificationMsgs: content:", string(m.Content))
		}
	}

	return nil
}
