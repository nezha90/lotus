package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/filecoin-project/go-address"
	"github.com/filecoin-project/lotus/api"
	"github.com/filecoin-project/lotus/chain/types"
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
		var bodyBuffer bytes.Buffer
		// 使用 io.TeeReader 读取 body，同时将数据复制到 bodyBuffer 中
		r.Body = io.NopCloser(io.TeeReader(r.Body, &bodyBuffer))

		body := bodyBuffer.Bytes()

		fmt.Println("unmarshal json body")
		// 解析 JSON-RPC 请求
		var req rpcRequest
		if err := json.Unmarshal(body, &req); err != nil {
			http.Error(w, "2 Invalid JSON-RPC request", http.StatusBadRequest)
			return
		}

		fmt.Println("switch message type")
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

		next.ServeHTTP(w, r)
	})
}

// 处理 WalletSign 请求，解析参数
func handleWalletSign(w http.ResponseWriter, params []json.RawMessage) {
	// 校验参数长度
	fmt.Println("check len")
	if len(params) != 3 {
		http.Error(w, "5 Invalid number of params", http.StatusBadRequest)
		return
	}

	var msgMeta api.MsgMeta

	fmt.Println("unmarshal msg meta")
	if err := json.Unmarshal(params[2], &msgMeta); err != nil {
		http.Error(w, "7 Failed to parse sign message", http.StatusBadRequest)
		return
	}

	fmt.Println("ok")
	if msgMeta.Type == api.MTBlock {
		return
	} else if msgMeta.Type != api.MTChainMsg {
		http.Error(w, "8 Failed to sign this type", http.StatusBadRequest)
		return
	}

	var addr address.Address

	if err := json.Unmarshal(params[0], &addr); err != nil {
		http.Error(w, "10 Failed to parse sign message", http.StatusBadRequest)
		return
	}
	fmt.Println(addr)

	var msgByte []byte
	if err := json.Unmarshal(params[1], &msgByte); err != nil {
		http.Error(w, "12 Failed to parse sign message", http.StatusBadRequest)
		return
	}

	msg, err := types.DecodeMessage(msgByte)
	if err != nil {
		http.Error(w, "11 Failed to parse sign message", http.StatusBadRequest)
		return
	}

	fmt.Println(msg)
}
