package services

import "fmt"

/*
	Mock task for go routine to execute and store progress
*/
func TaskMock(id int, start int) {
	ch := ChannelMap[id]
	fmt.Println("job started")
	for i := start + 1; i < 1e9; i++ {
		select {
		case <-ch:
			fmt.Println("job interrupted")
			JobProgress[id] = i // save progress
			return
		default:
			// do something
		}
	}
	fmt.Println("job completed")
	delete(ChannelMap, id)
	delete(JobProgress, id)
	close(ch)
}
