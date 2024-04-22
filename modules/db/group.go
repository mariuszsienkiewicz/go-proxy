package db

import (
	"fmt"
	"math/rand"
	"proxy/modules/config"
	"proxy/modules/log"
	"time"
)

// TODO add ServerGroupId

type Group struct {
	servers   map[string]*Server
	serverIds []string // used for randomized getter
}

var (
	Groups map[string]*Group
)

func init() {
	Groups = make(map[string]*Group)
}

func LoadGroups() error {
	for _, group := range config.Config.Proxy.ServerGroups {
		_, groupFound := Groups[group.Id]
		if groupFound {
			return fmt.Errorf("group %s already exists", group.Id)
		}
		Groups[group.Id] = NewGroup() // what about the type?
	}

	return nil
}

func NewGroup() *Group {
	return &Group{
		servers:   make(map[string]*Server),
		serverIds: make([]string, 0),
	}
}

func (g *Group) GetServer(serverId string) (*Server, bool) {
	s, found := g.servers[serverId]
	return s, found
}

func (g *Group) AddServer(server *Server) {
	g.servers[server.Config.Id] = server
	g.serverIds = append(g.serverIds, server.Config.Id)
	log.Logger.Tracef("New server added, group: %v", g)
}

func (g *Group) GetRandomServer() (*Server, error) {
	log.Logger.Tracef("Looking for random server, available servers: %v, group: %v", g.serverIds, g)
	if len(g.serverIds) == 0 {
		return nil, fmt.Errorf("no servers found")
	}

	rand.Seed(time.Now().UnixNano())
	index := rand.Intn(len(g.serverIds))

	s, found := g.servers[g.serverIds[index]]
	if found == false {
		return nil, fmt.Errorf("server %s not found", g.serverIds[index])
	}

	return s, nil
}
