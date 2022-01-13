package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/RoyXiang/putcallback/rclone"
)

func handleCallback(w http.ResponseWriter, r *http.Request) {
	fileId := r.FormValue("file_id")
	if fileId != "" {
		id, err := strconv.ParseInt(fileId, 10, 64)
		if err != nil {
			return
		}
		log.Printf("Callback received (file_id: %d, name: %s)", id, r.FormValue("name"))
		go rclone.SendFileIdToWorker(id)
	}

	_, _ = fmt.Fprint(w, "OK")
}

func main() {
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		http.HandleFunc("/putio/callback", handleCallback)
		log.Fatal(http.ListenAndServe(":1880", nil))
	}()
	log.Print("Server started on :1880")

	<-done
	rclone.Stop()
}
