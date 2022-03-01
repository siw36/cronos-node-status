package router

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"github.com/siw36/cronos-node-status/ldat"
	"github.com/siw36/cronos-node-status/rdat"
)

type blockdata struct {
	Rblock int `json:"remote_block"`
	Lblock int `json:"local_block"`
}

var source = "https://cronos.org/explorer/chain-blocks"

func Serve() {
	r := mux.NewRouter()
	// API
	r.HandleFunc("/api/v1/data", data)

	srv := &http.Server{
		Addr: fmt.Sprintf("0.0.0.0:%d", 8081),
		// Good practice to set timeouts to avoid Slowloris attacks.
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      r, // Pass our instance of gorilla/mux in.
	}

	// Run our server in a goroutine so that it doesn't block.
	go func() {
		log.Info("Starting server at ", srv.Addr)
		//if err := srv.ListenAndServeTLS(Config.SSLCertLocation, Config.SSLKeyLocation); err != nil {
		if err := srv.ListenAndServe(); err != nil {
			log.Fatal(err)
		}
	}()

	c := make(chan os.Signal, 1)
	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	// SIGKILL, SIGQUIT or SIGTERM (Ctrl+/) will not be caught.
	signal.Notify(c, os.Interrupt)

	// Block until we receive our signal.
	<-c

	// Create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()
	// Doesn't block if no connections, but will otherwise wait
	// until the timeout deadline.
	srv.Shutdown(ctx)
	// Optionally, you could run srv.Shutdown in a goroutine and block on
	// <-ctx.Done() if your application should wait for other services
	// to finalize based on context cancellation.
	log.Info("Shutting down server gracefully")
	os.Exit(0)

}

func data(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	var data blockdata
	log.Info("API getting data")
	data.Lblock = ldat.Exec()
	data.Rblock = rdat.Fetch(source)
	log.Info("API writing data")
	json.NewEncoder(w).Encode(data)
}
