package main

func main() {
    //创建服务端 设置服务端的ip 和端口
    server := NewServer("127.0.0.1", 8888)

    //开启服务端
    server.Start()
}

