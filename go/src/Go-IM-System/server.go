package main

import (
    "fmt"
    "net"
    "sync"
    "io"
    "time"
)

type Server struct {
    Ip   string
    Port int

    OnlineMap map[string]*User
    mapLock  sync.RWMutex

    Message chan string
}

//创建一个server接口
func NewServer(ip string, port int) *Server {
    server := &Server{
        Ip: ip,
        Port: port,
        OnlineMap: make(map[string]*User),
        Message: make(chan string),
    }
    return server
}


//消息广播的方法
func (this *Server) BroadCast(user *User, msg string){
    sendMsg := "[" + user.Addr + "]" + user.Name + ":" + msg
    this.Message <- sendMsg

}

//监听Message广播消息channel的goroutine，一旦有消息就发送给全部的在线User
func (this *Server) ListenMessager(){
    for{
        msg := <-this.Message

        this.mapLock.Lock()
        for _, cli := range this.OnlineMap{
            cli.C <- msg
        }
        this.mapLock.Unlock()
    }
}


//处理业务逻辑
func (this *Server) Handler(conn net.Conn) {
    //fmt.Println("链接建立成功")
    //用户上线，将用户加入到OnlineMap中
    user := NewUser(conn, this)

    user.Online()
    //监听用户是否活跃的channel
    isLive := make(chan bool)

    //接受客户端发送的消息
    go func(){
        buf := make([]byte, 4096)
        for{
            n, err := conn.Read(buf)
            if n == 0 {
                user.Offline()
                return
            }

            if err != nil && err != io.EOF {
                fmt.Println("Conn Read err:", err)
                return
            }
            msg := string(buf[:n-1])

            //fmt.Println(buf[:n-1])
            user.DoMessage(msg)

            //用户的任意消息，代表当前用户是一个活跃的
            isLive <- true
        }
    }()
    //当前handler阻塞
    for {
        select {
        case <-isLive:

        case <-time.After(time.Second * 100):

        user.SendMsg("你被踢了")

        //销毁资源
        close(user.C)

        //关闭连接
        conn.Close()

        //退出当前Handler
        return //runtime.Goexit()

        }

    }

}

//启动服务器的接口
func (this *Server) Start() {
    listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", this.Ip, this.Port))
    if err != nil {
        fmt.Println("net.Listen err:", err)
        return
    }
    defer listener.Close()

    go this.ListenMessager()

    for {
        conn, err := listener.Accept()
        if err != nil {
            fmt.Println("listener accept err:", err)
            continue
        }

        go this.Handler(conn)
    }
}

