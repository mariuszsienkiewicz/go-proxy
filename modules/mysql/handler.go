package mysql

import (
	"errors"
	"fmt"
	"github.com/go-mysql-org/go-mysql/mysql"
	"proxy/modules/db"
	"proxy/modules/log"
	iquery "proxy/modules/query"
	"proxy/modules/redirect"
	"proxy/modules/stats"
	"time"
)

type ProxyHandler struct {
	dbName string
}

func (h *ProxyHandler) UseDB(dbName string) error {
	log.Logger.Tracef("DB: use: %v", dbName)
	h.dbName = dbName
	return nil
}

func (h *ProxyHandler) HandleQuery(query string) (*mysql.Result, error) {
	// normalize and hash query
	normalizedQuery, hash := iquery.NormalizeAndHashQuery(query)

	// find the group which should handle the query
	targetGroup := redirect.FindRedirect(normalizedQuery, hash)
	serverGroup, groupFound := db.Groups[targetGroup]
	if groupFound == false {
		log.Logger.Tracef("Target group: %v", targetGroup)
		log.Logger.Tracef("Groups: %v", db.Groups)
		return nil, errors.New("group to use not found")
	}

	// get server that should be used for this query
	target, errRandomServer := serverGroup.GetRandomServer()
	if errRandomServer != nil {
		return nil, errRandomServer
	}
	log.Logger.Tracef("Query \"%v\" will be redirected to: %v (Host: %v, Hash: %s)", normalizedQuery, target.Config.Id, target.Config.GetDsn(), hash)

	// get connection
	conn, err := target.Connect()
	if err != nil {
		return nil, err
	}

	// TODO set proper proper database
	execute, err := conn.Execute(fmt.Sprintf("use %v;", h.dbName)) // TODO: absolutely needs to be changed

	// start timer
	start := time.Now()

	// execute query
	execute, err = conn.Execute(query)
	if err != nil {
		return nil, err
	}
	execTime := time.Since(start)

	// save query to statistics
	stats.SaveQuery(normalizedQuery, hash, execTime)

	return execute, nil
}

func (h *ProxyHandler) HandleFieldList(table string, fieldWildcard string) ([]*mysql.Field, error) {
	return nil, fmt.Errorf("not supported nowwww")
}

func (h *ProxyHandler) HandleStmtPrepare(query string) (int, int, interface{}, error) {
	return 0, 0, nil, fmt.Errorf("not supported nowwww")
}

func (h *ProxyHandler) HandleStmtExecute(context interface{}, query string, args []interface{}) (*mysql.Result, error) {
	return nil, fmt.Errorf("not supported nowwww")
}

func (h *ProxyHandler) HandleStmtClose(context interface{}) error {
	return nil
}

func (h *ProxyHandler) HandleOtherCommand(cmd byte, data []byte) error {
	return mysql.NewError(
		mysql.ER_UNKNOWN_ERROR,
		fmt.Sprintf("command %d is not supported now", cmd),
	)
}
