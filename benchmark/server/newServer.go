package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"os"
	"time"

	wwr "github.com/qbeon/webwire-go"
)

type settings struct {
	HostAddress        string
	HTTPSEnabled       bool
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
	if settings.HTTPSEnabled {
		// Setup a new TLS protected webwire server instance
		server, err = wwr.NewServerSecure(
			&BenchmarkServer{},
			wwr.ServerOptions{
				Host:              settings.HostAddress,
				WarnLog:           warnLog,
				ErrorLog:          errorLog,
				ReadTimeout:       settings.ReadTimeout,
				ReadBufferSize:    1024 * 8,
				WriteBufferSize:   1024 * 8,
				MessageBufferSize: 1024 * 16,
			},
			settings.CertFilePath,
			settings.PrivateKeyFilePath,
			settings.TLSConfig,
		)
		if err != nil {
			return nil, fmt.Errorf("wwr secure server setup failure: %s", err)
		}
	} else {
		// Setup a new unencrypted webwire server instance
		server, err = wwr.NewServer(
			&BenchmarkServer{},
			wwr.ServerOptions{
				Host:              settings.HostAddress,
				WarnLog:           warnLog,
				ErrorLog:          errorLog,
				ReadTimeout:       settings.ReadTimeout,
				MessageBufferSize: 1024 * 16,
			},
		)
		if err != nil {
			return nil, fmt.Errorf("wwr server setup failure: %s", err)
		}
	}
	return server, nil
}
