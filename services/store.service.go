package services

var ChannelMap = make(map[int]chan struct{}) // Global int to struct map to store all job id to channel mappings
var JobProgress = make(map[int]int)          // Global int to int map to store all job ids to job progress mappings
var JobId = 0                                // Global job id tracker
