package main

import (
	"context"
	"crypto/tls"
	"flag"
	"fmt"
	"net/http"
	"os"
	"runtime"
	"strconv"
	"time"

	"entgo.io/ent/dialect/sql/schema"
	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
	"github.com/msrevive/nexus2/internal/controller"
	"github.com/msrevive/nexus2/internal/ent"
	"github.com/msrevive/nexus2/internal/log"
	"github.com/msrevive/nexus2/internal/middleware"
	"github.com/msrevive/nexus2/internal/system"
	"golang.org/x/crypto/acme"
	"golang.org/x/crypto/acme/autocert"
)

// TODO: Set this value to via CI pipeline (dynamically using the git tag)
var version = "v1.0.5"

func initPrint() {
	fmt.Printf(`
    _   __                    ___
   / | / /__  _  ____  Nexus2|__ \
  /  |/ / _ \| |/_/ / / / ___/_/ /
 / /|  /  __/>  </ /_/ (__  ) __/
/_/ |_/\___/_/|_|\__,_/____/____/

Copyright © %d, Team MSRebirth

Version: %s
Website: https://msrebirth.net/
License: GPL-3.0 https://github.com/MSRevive/nexus2/blob/main/LICENSE %s`, time.Now().Year(), version, "\n\n")
}

func main() {
	// initial print
	initPrint()

	if err := run(); err != nil {
		fmt.Errorf("critical error detected: %v", err)
		os.Exit(1)
	}
}

func run() error {
	var cfile string
	var dbg bool
	flag.StringVar(&cfile, "cfile", "./runtime/config.toml", "Where to load the config file.")
	flag.BoolVar(&dbg, "dbg", false, "Run with debug mode.")
	flag.Parse()

	// Load configuration file
	config, err := system.LoadApiConfig(cfile, dbg)
	if err != nil {
		return err
	}

	// Initiate logging
	log.InitLogging("server.log", config.Log.Dir, config.Log.Level, config.Log.ExpireTime)

	if dbg {
		log.Log.Warnln("Running in Debug mode, do not use in production!")
	}

	// Max threads allowed
	if config.Core.MaxThreads != 0 {
		runtime.GOMAXPROCS(config.Core.MaxThreads)
	}

	// Load json files
	if config.ApiAuth.EnforceIP {
		log.Log.Printf("Loading IP list from %s", config.ApiAuth.IPListFile)
		if err := config.LoadIPList(); err != nil {
			log.Log.Warnln("Failed to load IP list.")
		}
	}

	if config.Verify.EnforceMap {
		log.Log.Printf("Loading Map list from %s", config.Verify.MapListFile)
		if err := config.LoadMapList(); err != nil {
			log.Log.Warnln("Failed to load Map list.")
		}
	}

	if config.Verify.EnforceBan {
		log.Log.Printf("Loading Ban list from %s", config.Verify.BanListFile)
		if err := config.LoadBanList(); err != nil {
			log.Log.Warnln("Failed to load Ban list.")
		}
	}

	log.Log.Printf("Loading Admin list from %s", config.Verify.AdminListFile)
	if err := config.LoadAdminList(); err != nil {
		log.Log.Warnln("Failed to load Admin list.")
	}

	// Connect database
	log.Log.Println("Connecting to database")
	client, err := ent.Open("sqlite3", config.Core.DBString)
	if err != nil {
		log.Log.Errorf("failed to open connection to sqlite3: %v", err)
		return err
	}
	defer client.Close()

	if err := client.Schema.Create(context.Background(), schema.WithAtlas(true)); err != nil {
		log.Log.Errorf("failed to create schema resources: %v", err)
		return err
	}

	// Variables for web server
	var srv *http.Server
	router := mux.NewRouter()
	srv = &http.Server{
		Handler:      api.NewRouter(config),
		Addr:         config.Core.Address + ":" + strconv.Itoa(config.Core.Port),
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
		// DefaultTLSConfig sets sane defaults to use when configuring the internal
		// webserver to listen for public connections.
		//
		// @see https://blog.cloudflare.com/exposing-go-on-the-internet
		// credit to https://github.com/pterodactyl/wings/blob/develop/config/config.go
		TLSConfig: &tls.Config{
			NextProtos: []string{"h2", "http/1.1"},
			CipherSuites: []uint16{
				tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
				tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
				tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
				tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
				tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305_SHA256,
				tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305_SHA256,
			},
			PreferServerCipherSuites: true,
			MinVersion:               tls.VersionTLS12,
			MaxVersion:               tls.VersionTLS13,
			CurvePreferences:         []tls.CurveID{tls.X25519, tls.CurveP256},
		},
	}

	if config.Cert.Enable {
		cm := autocert.Manager{
			Prompt:     autocert.AcceptTOS,
			HostPolicy: autocert.HostWhitelist(config.Cert.Domain),
			Cache:      autocert.DirCache("./runtime/certs"),
		}

		srv.TLSConfig = &tls.Config{
			GetCertificate: cm.GetCertificate,
			NextProtos:     append(srv.TLSConfig.NextProtos, acme.ALPNProto), // enable tls-alpn ACME challenges
		}

		go func() {
			if err := http.ListenAndServe(":http", cm.HTTPHandler(nil)); err != nil {
				log.Log.Fatalf("failed to serve autocert server: %v", err)
			}
		}()

		log.Log.Printf("Listening on: %s TLS", srv.Addr)
		if err := srv.ListenAndServeTLS("", ""); err != nil {
			log.Log.Fatalf("failed to serve over HTTPS: %v", err)
		}
	} else {
		log.Log.Printf("Listening on: %s", srv.Addr)
		if err := srv.ListenAndServe(); err != nil {
			log.Log.Fatalf("failed to serve over HTTP: %v", err)
		}
	}
}
