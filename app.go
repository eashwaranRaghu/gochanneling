package main

import (
	"log"
	"net/http"

	handler "./handlers"
)

/*
	Exposed 4 REST APIs:
	/start : starts a new job and returns the jobid which can be used to stop the job
	/stop/{id} : stops the running job, dumps the state to map.
	/terminate/{id} : terminates the stopped process by deleting the dump.
	/resume/{id} : Resumes the job by using the state saved in dump.
*/
func main() {
	http.HandleFunc("/start", handler.StartHandler)
	http.HandleFunc("/stop/", handler.StopHandler)
	http.HandleFunc("/terminate/", handler.TerminateHandler)
	http.HandleFunc("/resume/", handler.ResumeHandler)
	http.HandleFunc("/", handler.DefaultHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

/*
	TODO:
	1. Dockerize
	2. Err handling
*/
