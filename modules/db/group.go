package db

import (
	"fmt"
	"go-proxy/modules/config"
	"go-proxy/modules/log"
	"go.uber.org/zap"
	"math/rand"
)

type Group struct {
	Id        string
	servers   map[string]*Server
	serverIds []string // used for randomized getter
}

var (
	Groups map[string]*Group
)

func init() {
	CreateGroups()
}

func CreateGroups() {
	Groups = make(map[string]*Group)
}

func LoadGroups() error {
	for _, group := range config.Config.Proxy.ServerGroups {
		_, groupFound := Groups[group.Id]
		if groupFound {
			return fmt.Errorf("group %s already exists", group.Id)
		}
		Groups[group.Id] = NewGroup(group.Id) // what about the type?
	}

	return nil
}

func NewGroup(id string) *Group {
	return &Group{
		Id:        id,
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
	log.Logger.Debug("New server added, group", zap.String("group", server.Config.Id))
}

func (g *Group) GetRandomServer() (*Server, error) {
	log.Logger.Debug("Looking for random server")
	if len(g.serverIds) == 0 {
		log.Logger.Debug("No servers found in group, using default server")
		return DbPool.DefaultServer, nil
	}

	var activeServerIds []string
	for _, serverID := range g.serverIds {
		if server, found := g.servers[serverID]; found && server.Status == OPERATIONAL {
			activeServerIds = append(activeServerIds, serverID)
		}
	}

	if len(activeServerIds) == 0 {
		log.Logger.Debug("There is no operational server in group, using default server")
		return DbPool.DefaultServer, nil
	}

	index := rand.Intn(len(activeServerIds))

	s := g.servers[activeServerIds[index]]

	return s, nil
}
