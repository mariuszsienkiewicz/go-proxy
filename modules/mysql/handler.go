// Package mysql provides functionalities to handle MySQL database connections and queries through a proxy.
package mysql

import (
	"errors"
	"fmt"
	"github.com/go-mysql-org/go-mysql/client"
	"github.com/go-mysql-org/go-mysql/mysql"
	"proxy/modules/db"
	"proxy/modules/db/util"
	"proxy/modules/log"
	"proxy/modules/redirect"
)

// DbConnection represents a connection to a MySQL server, later used to return connections to pool
type DbConnection struct {
	Connection *client.Conn // Connection is the client connection to the MySQL server.
	Server     *db.Server   // Server is the MySQL server associated with this connection.
}

// ProxyHandler represents a handler for MySQL proxy queries.
type ProxyHandler struct {
	dbName        string                   // dbName is the name of the currently selected database.
	transaction   bool                     // transaction indicates whether a transaction is ongoing.
	dbConnections map[string]*DbConnection // dbConnections maps server IDs to their respective DbConnection instances.
}

// NewProxyHandler creates a new ProxyHandler instance.
func NewProxyHandler() *ProxyHandler {
	return &ProxyHandler{
		dbConnections: make(map[string]*DbConnection),
	}
}

// ReturnConnectionsToPool returns all database connections to the connection pool.
func (h *ProxyHandler) ReturnConnectionsToPool() {
	for _, dbConn := range h.dbConnections {
		dbConn.Server.Pool.PutConn(dbConn.Connection)
	}
}

// UseDB selects the specified database for subsequent queries.
func (h *ProxyHandler) UseDB(dbName string) error {
	log.Logger.Tracef("DB: use: %v", dbName)
	h.dbName = dbName
	return nil
}

// HandleQuery handles a MySQL util, redirecting it to the appropriate server.
func (h *ProxyHandler) HandleQuery(query string) (*mysql.Result, error) {
	// TODO check if query is in transaction / check if transaction is started/ended

	var target *db.Server
	if !h.transaction {
		// Normalize and hash util
		normalizedQuery, hash := util.NormalizeAndHashQuery(query)

		// Find the group which should handle the util
		targetGroup := redirect.FindRedirect(normalizedQuery, hash)
		serverGroup, groupFound := db.Groups[targetGroup]
		if !groupFound {
			log.Logger.Tracef("Target group: %v", targetGroup)
			log.Logger.Tracef("Groups: %v", db.Groups)
			return nil, errors.New("proxy error")
		}

		// Get server that should be used for this util
		var errRandomServer error
		target, errRandomServer = serverGroup.GetRandomServer()
		if errRandomServer != nil {
			log.Logger.Errorf("Error while getting random server: %v", errRandomServer)
			return nil, errors.New("proxy error")
		}
		log.Logger.Tracef("Query \"%v\" will be redirected to: %v (Host: %v, Hash: %s)", normalizedQuery, target.Config.Id, target.Config.GetDsn(), hash)
	} else {
		target = db.DbPool.DefaultServer
	}

	dbConnection, found := h.dbConnections[target.Config.Id]
	if !found {
		// Get connection
		conn, err := target.Connect()
		if err != nil {
			return nil, err
		}

		// Save connection to this context
		dbConnection = &DbConnection{
			Connection: conn,
			Server:     target,
		}

		h.dbConnections[target.Config.Id] = dbConnection
		_, err = dbConnection.Connection.Execute(fmt.Sprintf("use %v;", h.dbName))
		if err != nil {
			return nil, err
		}
	}

	// Execute util
	execute, err := dbConnection.Connection.Execute(query)
	if err != nil {
		return nil, err
	}

	return execute, nil
}

// HandleFieldList is not supported currently.
func (h *ProxyHandler) HandleFieldList(table string, fieldWildcard string) ([]*mysql.Field, error) {
	return nil, fmt.Errorf("not supported now")
}

// HandleStmtPrepare is not supported currently.
func (h *ProxyHandler) HandleStmtPrepare(query string) (int, int, interface{}, error) {
	return 0, 0, nil, fmt.Errorf("not supported now")
}

// HandleStmtExecute is not supported currently.
func (h *ProxyHandler) HandleStmtExecute(context interface{}, query string, args []interface{}) (*mysql.Result, error) {
	return nil, fmt.Errorf("not supported now")
}

// HandleStmtClose is not supported currently.
func (h *ProxyHandler) HandleStmtClose(context interface{}) error {
	return nil
}

// HandleOtherCommand handles unsupported MySQL commands.
func (h *ProxyHandler) HandleOtherCommand(cmd byte, data []byte) error {
	return mysql.NewError(
		mysql.ER_UNKNOWN_ERROR,
		fmt.Sprintf("command %d is not supported now", cmd),
	)
}
