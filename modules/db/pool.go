package db

import "fmt"

var (
	DbPool Pool
)

// Pool - populated by LoadServers
type Pool struct {
	Servers map[string]*Server
}

func init() {
	DbPool = Pool{
		Servers: make(map[string]*Server),
	}
}

func (p Pool) String() string {
	result := "{"
	for key, value := range p.Servers {
		result += fmt.Sprintf("%s:%v ", key, value)
	}
	result += "}"
	return result
}
