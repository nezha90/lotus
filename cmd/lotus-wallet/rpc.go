package main

import (
	"encoding/json"
	"fmt"
	"github.com/filecoin-project/go-address"
	"github.com/filecoin-project/lotus/api"
	"io"
	"net/http"
)

const (
	MethodWalletBalance = "Filecoin.WalletBalance"
	MethodWalletSign    = "Filecoin.WalletSign"
	MethodWalletList    = "Filecoin.WalletList"
)

// 定义用于解析 JSON-RPC 请求的结构体
type rpcRequest struct {
	JSONRPC string            `json:"jsonrpc"`
	Method  string            `json:"method"`
	Params  []json.RawMessage `json:"params"`
	ID      int               `json:"id"`
	Meta    json.RawMessage   `json:"meta"`
}

func methodFilterMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 读取请求体
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "1 Failed to read request body", http.StatusBadRequest)
			return
		}
		defer r.Body.Close()

		// 解析 JSON-RPC 请求
		var req rpcRequest
		fmt.Println(body)
		if err := json.Unmarshal(body, &req); err != nil {
			http.Error(w, "2 Invalid JSON-RPC request", http.StatusBadRequest)
			return
		}

		// 检查 method 字段是否为允许的 RPC 方法
		switch req.Method {
		case MethodWalletBalance, MethodWalletList:
		case MethodWalletSign:
			// Filecoin.WalletSign 继续处理 params
			handleWalletSign(w, req.Params)
		default:
			// 不支持的 method 返回错误
			http.Error(w, fmt.Sprintf("3 Unsupported method: %s", req.Method), http.StatusBadRequest)
			return
		}

		http.Error(w, "stop", http.StatusBadRequest)

		//next.ServeHTTP(w, r)
	})
}

// 处理 WalletSign 请求，解析参数
func handleWalletSign(w http.ResponseWriter, params []json.RawMessage) {
	// 校验参数长度
	if len(params) != 2 {
		http.Error(w, "5 Invalid number of params", http.StatusBadRequest)
		return
	}
	var msgMeta api.MsgMeta
	msgMetaParamsBytes, err := json.Marshal(params[2]) // 参数列表中第三个参数为消息对象
	if err != nil {
		http.Error(w, "6 Failed to parse sign params", http.StatusBadRequest)
		return
	}

	if err := json.Unmarshal(msgMetaParamsBytes, &msgMeta); err != nil {
		http.Error(w, "7 Failed to parse sign message", http.StatusBadRequest)
		return
	}

	if msgMeta.Type == api.MTBlock {
		return
	} else if msgMeta.Type != api.MTChainMsg {
		http.Error(w, "8 Failed to sign this type", http.StatusBadRequest)
	}

	var addr address.Address
	var msgByte []byte

	addrParamsBytes, err := json.Marshal(params[0]) // 参数列表中第三个参数为消息对象
	if err != nil {
		http.Error(w, "9 Failed to parse sign params", http.StatusBadRequest)
		return
	}

	if err := json.Unmarshal(addrParamsBytes, &addr); err != nil {
		http.Error(w, "10 Failed to parse sign message", http.StatusBadRequest)
		return
	}

	msgByte, err = json.Marshal(params[0]) // 参数列表中第三个参数为消息对象
	if err != nil {
		http.Error(w, "11 Failed to parse sign params", http.StatusBadRequest)
		return
	}
	fmt.Println(msgByte)
	fmt.Println(addr)

	//// 记录 from, to 和消息类型
	//fmt.Printf("Signing from: %s, to: %s, message type: %s\n", signParams.From, signParams.To, signParams.Msg.MessageType)
	//
	//// 模拟返回签名结果
	//result := map[string]string{
	//	"signature": "mocked_signature",
	//}
	//response, _ := json.Marshal(result)
	//w.Write(response)
}
