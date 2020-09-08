package handler

import (
	"fmt"
	"net/http"
	"strconv"

	services "../services"
)

/*
	Starts a job by genrating a channel and passing it to a go routine meant to execute mock task
*/
func StartHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "startHandler, jobid = %d", services.JobId)
	services.ChannelMap[services.JobId] = make(chan struct{})
	go services.TaskMock(services.JobId, 0)
	services.JobId++
}

/*
	Terminate a job by removings it corresponding channel and progress from map
*/
func TerminateHandler(w http.ResponseWriter, r *http.Request) { // rollback
	currentJobid, _ := strconv.Atoi(r.URL.Path[6:])
	delete(services.JobProgress, currentJobid) // clear the services.JobProgress
	delete(services.ChannelMap, currentJobid)  // delete the channel
	fmt.Println("Job terminated")
	fmt.Fprintf(w, "terminateHandler")
}

/*
	Signals goroutine to stop by closing the channel and dump the current state to map
*/
func StopHandler(w http.ResponseWriter, r *http.Request) {
	// services.JobProgress and quit
	currentJobid, _ := strconv.Atoi(r.URL.Path[6:])
	ch := services.ChannelMap[currentJobid]
	close(ch)
	delete(services.ChannelMap, currentJobid)
	fmt.Fprintf(w, "stopHandler")
}

/*
	Spawns a goroutine and resumes the job by reading saved state from map
*/
func ResumeHandler(w http.ResponseWriter, r *http.Request) {
	currentJobid, _ := strconv.Atoi(r.URL.Path[6:])
	services.ChannelMap[services.JobId] = make(chan struct{})
	go services.TaskMock(services.JobId, services.JobProgress[currentJobid]) // start goroutine right where left
	delete(services.JobProgress, currentJobid)                               // clear the services.JobProgress
	fmt.Fprintf(w, "resumeHandler")
}

func DefaultHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hi, Exposed 4 REST APIs:\n/start : starts a new job and returns the jobid which can be used to stop the job\n/stop/{id} : stops the running job, dumps the state to map.\n/terminate/{id} : terminates the stopped process by deleting the dump.\n/resume/{id} : Resumes the job by using the state saved in dump.")
}
