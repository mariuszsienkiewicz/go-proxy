package stats

import (
	"time"
)

type QueryStat struct {
	Query    string
	ExecTime time.Duration
}

var (
	// Stats temporary placeholder for statistics
	Stats map[string][]QueryStat
)

func init() {
	Stats = make(map[string][]QueryStat)
}

func SaveQuery(query string, hash string, execTime time.Duration) {
	qStat, ok := Stats[hash]
	if ok == false {
		qStat = []QueryStat{}
	}

	qStat = append(qStat, QueryStat{Query: query, ExecTime: execTime})
	Stats[hash] = qStat
}
