package main

import (
        "golang.org/x/net/websocket"
        "encoding/base64"
        "flag"
        "io"
        "log"
        "net"
        "net/http"
)

var target = flag.String("t", "localhost:567", "server to proxy to")

var listen = flag.String("l", ":8080", "websocket server bind address")

func main() {
        flag.Parse()

        http.Handle("/telnet", websocket.Handler(func(ws *websocket.Conn) { wsHandler(ws, *target) }))
        if err := http.ListenAndServe(*listen, nil); err != nil {
                log.Fatal(err)
        }
}

func wsHandler(ws *websocket.Conn, addr string) {
        var buf [2048]byte
        var ebuf [4096]byte

        defer ws.Close()
        conn, err := net.Dial("tcp", addr)
        if err != nil {
                return
        }
        defer conn.Close()
        go func() {
                io.Copy(conn, base64.NewDecoder(base64.StdEncoding, ws))
                conn.Close()
                ws.Close()
        }()
        for {
                n, err := conn.Read(buf[:])
                if err != nil {
                        conn.Close()
                        ws.Close()
                        return
                }
                base64.StdEncoding.Encode(ebuf[:], buf[:n])
                ws.Write(ebuf[:base64.StdEncoding.EncodedLen(n)])
        }
}
