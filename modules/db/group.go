package db

import (
	"fmt"
	"math/rand"
	"proxy/modules/config"
	"proxy/modules/log"
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

// GetRandomServer TODO - if no servers are available use fallback, and if fallback is dead too then use default server
func (g *Group) GetRandomServer() (*Server, error) {
	log.Logger.Tracef("Looking for random server in group: %v", g)
	if len(g.serverIds) == 0 {
		log.Logger.Tracef("No servers found in group: %v, using default server", g)
		return DbPool.DefaultServer, nil
	}

	// TODO - perf issue
	var activeServerIds []string
	for _, serverID := range g.serverIds {
		if server, found := g.servers[serverID]; found && server.Status == OPERATIONAL {
			activeServerIds = append(activeServerIds, serverID)
		}
	}

	log.Logger.Tracef("Available servers: %v", activeServerIds)
	if len(activeServerIds) == 0 {
		log.Logger.Tracef("There is no operational server in group: %v, using default server", g)
		return DbPool.DefaultServer, nil
	}

	index := rand.Intn(len(activeServerIds))

	s := g.servers[activeServerIds[index]]

	return s, nil
}
