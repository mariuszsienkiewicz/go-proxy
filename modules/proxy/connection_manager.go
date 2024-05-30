package proxy

import (
	"context"
	"fmt"
	"github.com/go-mysql-org/go-mysql/client"
	"go-proxy/modules/db"
	"go-proxy/modules/log"
	"go.uber.org/zap"
	"math/rand"
)

// DbConnection represents a connection to a MySQL server, later used to return connections to the pool.
type DbConnection struct {
	connection *client.Conn // connection is the client connection to the MySQL server.
	server     *db.Server   // server is the MySQL server associated with this connection.
	charset    string       // charset used in this connection.
	dbName     string       // dbName is the database name used in this connection.
}

// ConnectionManager manages multiple database connections.
type ConnectionManager struct {
	ctx             context.Context          // ctx context of the app
	dbConnections   map[string]*DbConnection // dbConnections maps group IDs to their respective DbConnection instances.
	dbConnectionIds []string                 // dbConnectionIds is a list of group IDs used for random selection.
}

// NewConnectionManager creates and returns a new ConnectionManager instance.
func NewConnectionManager(ctx context.Context) *ConnectionManager {
	return &ConnectionManager{
		ctx:             ctx,
		dbConnections:   make(map[string]*DbConnection),
		dbConnectionIds: make([]string, 0),
	}
}

// ReturnConnectionsToPool returns all connections in the manager back to their respective connection pools.
func (m *ConnectionManager) ReturnConnectionsToPool() {
	for _, dbConn := range m.dbConnections {
		log.Logger.Debug("Returning connection to pool", zap.String("server", dbConn.server.Config.Id))
		dbConn.server.Pool.PutConn(dbConn.connection)
	}
}

func (m *ConnectionManager) ReturnConnectionById(id string) {
	// get server from dbConn
	dbConn, ok := m.dbConnections[id]
	if ok {
		dbConn.server.Pool.PutConn(dbConn.connection)
		delete(m.dbConnections, id)

		index := -1
		for i, k := range m.dbConnectionIds {
			if k == id {
				index = i
				break
			}
		}

		if index != -1 {
			m.dbConnectionIds = append(m.dbConnectionIds[:index], m.dbConnectionIds[index+1:]...)
		}
	}
}

// getConnection retrieves a connection from the manager or creates a new one if it does not exist.
func (m *ConnectionManager) getConnection(group *db.Group) (*DbConnection, error) {
	// Check if the context is done
	select {
	case <-m.ctx.Done():
		return nil, m.ctx.Err()
	default:
	}

	dbConnection, found := m.dbConnections[group.Id]
	if !found {
		var err error
		dbConnection, err = m.createRandomConnection(group)
		if err != nil {
			return nil, err
		}
	}

	return dbConnection, nil
}

func (m *ConnectionManager) getDefaultConnection() (*DbConnection, error) {
	// Check if the context is done
	select {
	case <-m.ctx.Done():
		return nil, m.ctx.Err()
	default:
	}

	dbConnection, found := m.dbConnections[db.DbPool.DefaultServer.Config.ServerGroup]
	if !found {
		var err error
		dbConnection, err = m.createConnection(db.DbPool.DefaultServer.Config.ServerGroup, db.DbPool.DefaultServer)
		if err != nil {
			return nil, err
		}
	}

	return dbConnection, nil
}

// createConnection establishes a new connection to the target server and adds it to the manager.
func (m *ConnectionManager) createRandomConnection(group *db.Group) (*DbConnection, error) {
	// Check if the context is done
	select {
	case <-m.ctx.Done():
		return nil, m.ctx.Err()
	default:
	}

	target, err := group.GetRandomServer()
	if err != nil {
		return nil, err
	}

	return m.createConnection(group.Id, target)
}

func (m *ConnectionManager) createConnection(id string, target *db.Server) (*DbConnection, error) {
	// Check if the context is done
	select {
	case <-m.ctx.Done():
		return nil, m.ctx.Err()
	default:
	}

	// Establish a connection to the server
	conn, err := target.Connect(m.ctx)
	if err != nil {
		return nil, err
	}

	// Create a new DbConnection instance
	dbConnection := &DbConnection{
		connection: conn,
		server:     target,
	}

	// Store the new connection in the manager
	m.dbConnections[id] = dbConnection
	m.dbConnectionIds = append(m.dbConnectionIds, id)

	return dbConnection, nil
}

// getRandomConnection retrieves a random connection from the manager.
func (m *ConnectionManager) getRandomConnection() (*DbConnection, error) {
	// Check if the context is done
	select {
	case <-m.ctx.Done():
		return nil, m.ctx.Err()
	default:
	}

	if len(m.dbConnectionIds) == 0 {
		return m.getDefaultConnection()
	}

	// Select a random ID from the list of connection IDs
	index := rand.Intn(len(m.dbConnectionIds))
	randomId := m.dbConnectionIds[index]
	conn, found := m.dbConnections[randomId]
	if !found {
		return nil, fmt.Errorf("connection not found")
	}

	return conn, nil
}
