package mysql

import (
	"fmt"
	"github.com/go-mysql-org/go-mysql/mysql"
	"proxy/modules/config"
	"proxy/modules/log"
	"proxy/modules/redirect"
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
	// find the place where query should go
	target := redirect.FindRedirect(query)
	log.Logger.Tracef("Query: %v will be redirected to: %v (%v)", query, target.Id, target.GetDsn())

	// get connection
	connect, err := Connect(target, *target.GetUser(config.Config.Proxy.DbUsers))
	if err != nil {
		return nil, err
	}

	// execute connection
	// TODO set proper context (proper database)
	execute, err := connect.Execute(fmt.Sprintf("use %v;", h.dbName)) // TODO: absolutely needs to be changed
	execute, err = connect.Execute(query)
	if err != nil {
		return nil, err
	}

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
