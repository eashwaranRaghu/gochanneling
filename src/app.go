package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
)

var m = make(map[int]chan string)
var dump = make(map[int]int)
var jid = 0

func task(id int, start int) { // pass channel to this fncn
	ch := m[id]
	fmt.Println("job started")
	for i := start; i < 1e9; i++ {
		select {
		case <-ch:
			fmt.Println("job interrupted")
			dump[id] = i // dump
			return
		default:
			// do nothing
		}
	}
	fmt.Println("job completed")
	delete(m, id)
	delete(dump, id)
	close(ch)
}

func startHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "startHandler, jobid = %d", jid)
	m[jid] = make(chan string)
	go task(jid, 0)
	jid++
}

func terminateHandler(w http.ResponseWriter, r *http.Request) { // rollback
	currentJobid, _ := strconv.Atoi(r.URL.Path[6:])
	delete(dump, currentJobid) // clear the dump
	delete(m, currentJobid)    // delete the channel
	fmt.Fprintf(w, "terminateHandler")
}

func stopHandler(w http.ResponseWriter, r *http.Request) {
	// dump and quit
	currentJobid, _ := strconv.Atoi(r.URL.Path[6:])
	ch := m[currentJobid]
	close(ch)
	delete(m, currentJobid)
	fmt.Fprintf(w, "stopHandler")
}

func resumeHandler(w http.ResponseWriter, r *http.Request) {
	currentJobid, _ := strconv.Atoi(r.URL.Path[6:])
	m[jid] = make(chan string)
	go task(jid, dump[currentJobid]) // start goroutine right where left
	delete(dump, currentJobid)       // clear the dump
	fmt.Fprintf(w, "resumeHandler")
}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hi there, I love %s!", r.URL.Path[1:])
}

func main() {
	http.HandleFunc("/start", stopHandler)
	http.HandleFunc("/terminate/", terminateHandler)
	http.HandleFunc("/stop/", stopHandler)
	http.HandleFunc("/resume/", resumeHandler)
	// http.HandleFunc("/progress", progressHandler)
	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
