// Package proxy provides functionalities to handle MySQL database connections and queries through a proxy.
package proxy

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-mysql-org/go-mysql/client"
	"github.com/go-mysql-org/go-mysql/mysql"
	"go-proxy/modules/db"
	"go-proxy/modules/db/util"
	"go-proxy/modules/log"
	"go-proxy/modules/redirect"
	"go.uber.org/zap"
)

// ProxyHandler represents a handler for MySQL proxy queries.
type ProxyHandler struct {
	Id                string             // UUID of the handler
	ctx               context.Context    // Context of the app
	ConnectionManager *ConnectionManager // Manages the connections used by ProxyHandler
	dbName            string             // Name of the currently selected database
	charsetClient     string             // Charset set by the client
	transaction       bool               // Indicates whether a transaction is ongoing
	sendInTransaction bool               // Indicates if a query should still be sent in transaction even if transaction is false (for example COMMIT)
}

// StmtContext represents the context of a statement, containing the connection and statement itself.
type StmtContext struct {
	connection *DbConnection // Connection used to prepare the statement
	statement  *client.Stmt  // Prepared statement
}

// NewProxyHandler creates a new ProxyHandler instance.
func NewProxyHandler(ctx context.Context, uuid string) *ProxyHandler {
	return &ProxyHandler{
		Id:                uuid,
		ctx:               ctx,
		ConnectionManager: NewConnectionManager(ctx),
	}
}

// UseDB selects the specified database for subsequent queries.
func (h *ProxyHandler) UseDB(dbName string) error {
	log.Logger.Debug("Use DB", zap.String("handler", h.Id), zap.String("name", dbName))
	h.dbName = dbName
	return nil
}

// HandleQuery processes a given query.
func (h *ProxyHandler) HandleQuery(query string) (*mysql.Result, error) {
	log.Logger.Debug("Query", zap.String("handler", h.Id), zap.String("query", query))

	// Check if the context is done
	select {
	case <-h.ctx.Done():
		return nil, h.ctx.Err()
	default:
	}

	// Analyze query content
	h.analyzeQuery(query)

	// Find the connection that should be used
	var dbConnection *DbConnection
	if !h.transaction && !h.sendInTransaction {
		normalizedQuery, hash := util.NormalizeAndHashQuery(query)

		var err error
		dbConnection, err = h.getTargetConnection(normalizedQuery, hash)
		if err != nil {
			return nil, err
		}
	} else {
		log.Logger.Debug("Query is in the transaction", zap.String("query", query))
		var err error
		dbConnection, err = h.ConnectionManager.getDefaultConnection()
		if err != nil {
			return nil, err
		}
	}

	// Setup connection
	err := h.setupConnection(dbConnection)
	if err != nil {
		log.Logger.Warn("Error setting up connection", zap.Error(err))
		return nil, err
	}

	// Execute query
	execute, err := dbConnection.connection.Execute(query)
	if err != nil {
		log.Logger.Warn("Error executing query: %v, reason: %v", zap.String("query", query), zap.Error(err))
		return nil, err
	}

	// Reset the ProxyHandler sendInTransaction flag
	h.sendInTransaction = false

	log.Logger.Debug("Successfully executed query", zap.String("handler", h.Id), zap.String("query", query))

	return execute, nil
}

// HandleFieldList returns the list of fields for the specified table.
func (h *ProxyHandler) HandleFieldList(table string, fieldWildcard string) ([]*mysql.Field, error) {
	log.Logger.Debug("List fields", zap.String("handler", h.Id), zap.String("table", table), zap.String("fieldWildcard", fieldWildcard))

	// Check if the context is done
	select {
	case <-h.ctx.Done():
		return nil, h.ctx.Err()
	default:
	}

	dbConnection, err := h.ConnectionManager.getRandomConnection()
	if err != nil {
		return nil, err
	}

	fields, err := dbConnection.connection.FieldList(table, fieldWildcard)
	if err != nil {
		log.Logger.Warn("Error getting fields", zap.Error(err))
	}
	return fields, nil
}

// HandleStmtPrepare prepares a statement for execution.
func (h *ProxyHandler) HandleStmtPrepare(query string) (int, int, interface{}, error) {
	log.Logger.Debug("Stmt prepare", zap.String("query", query))

	// Check if the context is done
	select {
	case <-h.ctx.Done():
		return 0, 0, nil, h.ctx.Err()
	default:
	}

	// Find the target for the statement
	var dbConnection *DbConnection
	if !h.transaction && !h.sendInTransaction {
		normalizedQuery, hash := util.NormalizeAndHashQuery(query)

		var err error
		dbConnection, err = h.getTargetConnection(normalizedQuery, hash)
		if err != nil {
			return 0, 0, nil, err
		}
	} else {
		log.Logger.Debug("Query is in the transaction", zap.String("query", query))
		var err error
		dbConnection, err = h.ConnectionManager.getDefaultConnection()
		if err != nil {
			return 0, 0, nil, err
		}
	}

	stmt, err := dbConnection.connection.Prepare(query)
	if err != nil {
		log.Logger.Warn("Error preparing statement", zap.Error(err))
	}

	return stmt.ParamNum(), stmt.ColumnNum(), StmtContext{
		connection: dbConnection,
		statement:  stmt,
	}, nil
}

// HandleStmtExecute executes a prepared statement.
func (h *ProxyHandler) HandleStmtExecute(context interface{}, query string, args []interface{}) (*mysql.Result, error) {
	log.Logger.Debug("Stmt execute")

	// Check if the context is done
	select {
	case <-h.ctx.Done():
		return nil, h.ctx.Err()
	default:
	}

	stmtContext, ok := context.(StmtContext)
	if !ok {
		log.Logger.Warn("Error getting statement context")
		return nil, errors.New("go-proxy error, while getting the statement context")
	}

	execute, err := stmtContext.statement.Execute(args...)
	if err != nil {
		log.Logger.Warn("Error while executing the statement", zap.String("query", query), zap.Error(err))
	}

	return execute, nil
}

// HandleStmtClose closes a prepared statement.
func (h *ProxyHandler) HandleStmtClose(context interface{}) error {
	log.Logger.Debug("Stmt close")

	stmtContext, ok := context.(StmtContext)
	if !ok {
		log.Logger.Error("Error getting statement context")
		return errors.New("go-proxy error, while getting the statement context")
	}
	return stmtContext.connection.connection.Close()
}

// HandleOtherCommand handles unsupported MySQL commands.
func (h *ProxyHandler) HandleOtherCommand(cmd byte, data []byte) error {
	log.Logger.Error("Command executed but not supported", zap.String("cmd", fmt.Sprintf("%c", cmd)))
	return mysql.NewError(
		mysql.ER_UNKNOWN_ERROR,
		fmt.Sprintf("command %d is not supported", cmd),
	)
}

// analyzeQuery analyzes the query, checks special queries, and sets the handler state.
func (h *ProxyHandler) analyzeQuery(query string) {
	command := Analyze(query)
	switch command.Type {
	case SetNames:
		log.Logger.Debug("Set names", zap.String("value", command.Value))
		h.charsetClient = command.Value
		return
	case UseDatabase:
		log.Logger.Debug("Use database", zap.String("value", command.Value))
		h.dbName = command.Value
		return
	default:
	}

	// Check transaction
	AnalyzeTransaction(h, query)
}

// setupConnection sets up the connection to be compatible with the current context.
func (h *ProxyHandler) setupConnection(connection *DbConnection) error {
	// Change charset if it's wrong
	if h.charsetClient != "" {
		if h.charsetClient != connection.charset {
			log.Logger.Debug(
				"Connection uses the wrong charset",
				zap.String("handler", h.Id),
				zap.String("server", connection.server.Config.Id),
				zap.String("client charset", h.charsetClient),
				zap.String("connection charset", connection.charset),
			)
			_, err := connection.connection.Execute(fmt.Sprintf("SET NAMES %s;", h.charsetClient))
			if err != nil {
				log.Logger.Warn("Couldn't set proper charset", zap.String("handler", h.Id), zap.Error(err))
				return err
			}
			connection.charset = h.charsetClient
		}
	}

	if h.dbName != "" {
		if h.dbName != connection.dbName {
			log.Logger.Debug(
				"Connection uses the wrong database",
				zap.String("handler", h.Id),
				zap.String("server", connection.server.Config.Id),
				zap.String("client database", h.dbName),
				zap.String("connection database", connection.dbName),
			)
			_, err := connection.connection.Execute(fmt.Sprintf("USE %v;", h.dbName))
			if err != nil {
				log.Logger.Warn("Couldn't set proper database", zap.String("handler", h.Id), zap.Error(err))
				return err
			}
			connection.dbName = h.dbName
		}
	}

	return nil
}

// getTargetGroup gets the database which should be used for the query.
func (h *ProxyHandler) getTargetConnection(query string, hash string) (*DbConnection, error) {
	// Find the group which should handle the query
	targetGroup := redirect.FindRedirect(query, hash)
	serverGroup, groupFound := db.Groups[targetGroup]
	if !groupFound {
		log.Logger.Debug("Target group not found", zap.String("group", targetGroup))
		return nil, errors.New("proxy error")
	}

	log.Logger.Debug(
		"Query redirection",
		zap.String("query", query),
		zap.String("group", serverGroup.Id),
		zap.String("hash", hash),
	)

	// get connection
	connection, err := h.ConnectionManager.getConnection(serverGroup)
	if err != nil {
		log.Logger.Warn("Couldn't get needed connection", zap.String("handler", h.Id), zap.Error(err))
		return nil, err
	}

	return connection, nil
}
