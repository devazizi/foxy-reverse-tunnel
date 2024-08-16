package internal

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"foxy-tunnel/config"
	"foxy-tunnel/pkg/log"
	"io"
	"net"
	"os"
	"sync"
	"time"
)

func NewClient(ctx context.Context, clientCfg config.ClientConfig) {
	time.Sleep(time.Second * 4)
	packetChan := make(chan net.Conn)

	go func() {
		caCert, err := os.ReadFile(clientCfg.ClientCertificate)

		if err != nil {
			log.Debug("client side", "fail to load client cert", map[string]interface{}{
				"error": err.Error(),
			})
		}

		caCertPool := x509.NewCertPool()
		caCertPool.AppendCertsFromPEM(caCert)

		tlsCfg := &tls.Config{
			InsecureSkipVerify: true,
			RootCAs:            caCertPool,
			ServerName:         clientCfg.Sni,
		}

		foreignServerDial, foreignRadialErr := tls.Dial("tcp", clientCfg.ForeignServer, tlsCfg)
		if foreignRadialErr != nil {
			log.Error("client side", "server side started on port "+clientCfg.ForeignServer, nil)
		}

		handShakeErr := foreignServerDial.Handshake()
		if handShakeErr != nil {

			log.Error("client side local srv to foreign srv", handShakeErr.Error(), nil)
			return
		}

		packetChan <- foreignServerDial
	}()

	func() {
		localServer, err := net.Dial("tcp", clientCfg.LocalServer)
		if err != nil {
			log.Debug("client side local", "server side started on port "+clientCfg.LocalServer, nil)
		}

		foreignServer := <-packetChan

		RelayBetweenClients(foreignServer, localServer)
	}()

}
func RelayBetweenClients(localServer net.Conn, foreignServer net.Conn) {

	//defer localServer.Close()
	//defer foreignServer.Close()

	wg := new(sync.WaitGroup)
	wg.Add(2)

	go func() {
		_, err := io.Copy(localServer, foreignServer)
		if err != nil {
			log.Error("client side local srv to foreign srv", err.Error(), nil)
		}
		//localServer.Close()
		wg.Done()
	}()

	go func() {
		_, err := io.Copy(foreignServer, localServer)
		if err != nil {
			log.Error("lient side foreign srv to local srv", err.Error(), nil)
		}
		//foreignServer.Close()
		wg.Done()
	}()

	wg.Wait()

}
