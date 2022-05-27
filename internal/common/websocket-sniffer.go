package common

import (
	"fmt"
	"github.com/gorilla/websocket"
	"io"
)

func websocketsClose(dst, src io.Closer) {
	_ = dst.Close()
	_ = src.Close()
}

func websocketCopy(dst, src *websocket.Conn, sniffer io.Writer) error {
	for {
		messageType, message, readError := src.ReadMessage()
		if readError != nil {
			return readError
		}
		_, sniffError := sniffer.Write(message)
		if sniffError != nil {
			return sniffError
		}
		_, _ = fmt.Fprintf(sniffer, SniffSeparator)
		writeError := dst.WriteMessage(messageType, message)
		if writeError != nil {
			return writeError
		}
	}
}

func ForwardWebsocketsTraffic(dst, src *websocket.Conn, inboundSniff, outboundSniff io.Writer) error {
	defer websocketsClose(dst, src)
	go websocketCopy(dst, src, outboundSniff)
	return websocketCopy(src, dst, inboundSniff)
}
