package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"os"
	"time"

	wwr "github.com/qbeon/webwire-go"
	wwrfasthttp "github.com/qbeon/webwire-go/transport/fasthttp"
	"github.com/valyala/fasthttp"
)

type settings struct {
	HostAddress        string
	Transport          string
	CertFilePath       string
	PrivateKeyFilePath string
	TLSConfig          *tls.Config
	ReadTimeout        time.Duration
}

func newServer(settings settings) (server wwr.Server, err error) {
	// Initialize loggers
	warnLog := log.New(
		os.Stdout,
		"WARN: ",
		log.Ldate|log.Ltime|log.Lshortfile,
	)
	errorLog := log.New(
		os.Stderr,
		"ERR: ",
		log.Ldate|log.Ltime|log.Lshortfile,
	)

	// Initialize server
	if settings.Transport == "https" {
		// Setup a new TLS protected webwire over HTTPS server instance
		server, err = wwr.NewServer(
			&BenchmarkServer{},
			wwr.ServerOptions{
				Host:              settings.HostAddress,
				WarnLog:           warnLog,
				ErrorLog:          errorLog,
				ReadTimeout:       settings.ReadTimeout,
				MessageBufferSize: 1024 * 16,
			},
			&wwrfasthttp.Transport{
				HTTPServer: &fasthttp.Server{
					ReadBufferSize:  1024 * 8,
					WriteBufferSize: 1024 * 8,
				},
				TLS: &wwrfasthttp.TLS{
					CertFilePath:       settings.CertFilePath,
					PrivateKeyFilePath: settings.PrivateKeyFilePath,
					Config:             settings.TLSConfig,
				},
				BeforeUpgrade: func(
					_ *fasthttp.RequestCtx,
				) wwr.ConnectionOptions {
					return wwr.ConnectionOptions{
						ConcurrencyLimit: 10,
					}
				},
			},
		)
		if err != nil {
			return nil, fmt.Errorf("wwr (wss) server setup failure: %s", err)
		}
	} else if settings.Transport == "http" {
		// Setup a new unencrypted webwire over HTTP server instance
		server, err = wwr.NewServer(
			&BenchmarkServer{},
			wwr.ServerOptions{
				Host:              settings.HostAddress,
				WarnLog:           warnLog,
				ErrorLog:          errorLog,
				ReadTimeout:       settings.ReadTimeout,
				MessageBufferSize: 1024 * 16,
			},
			&wwrfasthttp.Transport{
				HTTPServer: &fasthttp.Server{
					ReadBufferSize:  1024 * 8,
					WriteBufferSize: 1024 * 8,
				},
			},
		)
		if err != nil {
			return nil, fmt.Errorf("wwr (ws) server setup failure: %s", err)
		}
	} else {
		return nil, fmt.Errorf(
			"unsupported transport layer: %s",
			settings.Transport,
		)
	}

	return server, nil
}
