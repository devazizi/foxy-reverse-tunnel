package internal

import (
	"context"
	"foxy-tunnel/config"
	"foxy-tunnel/pkg/log"
	"foxy-tunnel/pkg/net_helper"
	"io"
	"net"
	"sync"
)

func NewServer(ctx context.Context, serverCfg config.ServerConfig) {
	// TODO we have to part client side server and server side server
	connChan := make(chan net.Conn)

	go func() {
		err := ClientSideServer(ctx, serverCfg.ListenOn, connChan)
		if err != nil {
			log.Error("client server ", err.Error(), nil)
		}
	}()

	err := ServerSideServer(ctx, serverCfg.ServerOn, connChan)
	if err != nil {
		log.Error("server server ", err.Error(), nil)
	}

}

func ServerSideServer(ctx context.Context, serverListenOn string, connChan chan net.Conn) error {
	listen, err := net.Listen("tcp", serverListenOn)

	log.Debug("server side", "server side strarted on port "+serverListenOn, nil)

	if err != nil {
		log.Error("client server ", err.Error(), nil)
		return err
	}

	defer listen.Close()

	serverConn, err := listen.Accept()

	if err != nil {
		log.Error("client server ", err.Error(), nil)
	}

	for {
		clientConn := <-connChan
		go HandleClientConnection(clientConn, serverConn)
	}
	//for {
	//	conn, err := listen.Accept()
	//	if err != nil {
	//		log.Error("client server ", err.Error(), nil)
	//	}
	//
	//	// send connection to channel
	//	connChan <- conn
	//}

}

func ClientSideServer(ctx context.Context, serverListenOn string, connChan chan net.Conn) error {
	listen, err := net.Listen("tcp", serverListenOn)
	if err != nil {
		log.Error("client server ", err.Error(), nil)
		return err
	}

	log.Debug("server side", "client side started on port "+serverListenOn, nil)

	defer listen.Close()

	for {
		conn, err := listen.Accept()
		if err != nil {
			log.Error("client server ", err.Error(), nil)
		}

		log.Debug("client side", "accept connection", map[string]interface{}{
			"client_ip": net_helper.GetIp(conn.RemoteAddr()),
		})

		// send connection to channel
		connChan <- conn
	}

	//<-ctx.Done()
	//return nil
}

func HandleClientConnection(clientConn net.Conn, serverConn net.Conn) {
	defer clientConn.Close()

	wg := new(sync.WaitGroup)
	wg.Add(2)

	go func() {
		io.Copy(clientConn, serverConn)
		clientConn.Close()
		wg.Done()
	}()

	go func() {
		io.Copy(serverConn, clientConn)
		wg.Done()
	}()

	wg.Wait()
}

//func HandleClientConnection(conn net.Conn) {
//
//}
