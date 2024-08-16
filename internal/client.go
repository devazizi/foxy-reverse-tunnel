package internal

import (
	"context"
	"foxy-tunnel/config"
	"foxy-tunnel/pkg/log"
	"io"
	"net"
	"sync"
	"time"
)

func NewClient(ctx context.Context, clientCfg config.ClientConfig) {
	time.Sleep(time.Second * 4)
	packetChan := make(chan net.Conn)

	go func() {
		foreignServerDial, err := net.Dial("tcp", clientCfg.ForeignServer)
		if err != nil {
			log.Debug("client side", "server side strarted on port "+clientCfg.ForeignServer, nil)
		}

		packetChan <- foreignServerDial
	}()

	func() {
		localServer, err := net.Dial("tcp", clientCfg.LocalServer)
		if err != nil {
			log.Debug("client side local", "server side strarted on port "+clientCfg.LocalServer, nil)
		}

		foreignServer := <-packetChan

		RelayBetweenClients(foreignServer, localServer)
	}()

}
func RelayBetweenClients(localServer net.Conn, foreignServer net.Conn) {

	defer localServer.Close()
	defer foreignServer.Close()

	wg := new(sync.WaitGroup)
	wg.Add(2)

	go func() {
		io.Copy(localServer, foreignServer)
		localServer.Close()
		wg.Done()
	}()

	go func() {
		io.Copy(foreignServer, localServer)
		foreignServer.Close()
		wg.Done()
	}()

	wg.Wait()

}
