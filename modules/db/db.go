package db

import (
	"context"
	"fmt"
	"go-proxy/modules/log"
	"go.uber.org/zap"
)

// Init - should be loaded in logic order, groups -> Servers -> pool
func Init(ctx context.Context) error {
	setup() // clear everything before loading

	if err := LoadGroups(); err != nil {
		return err
	}
	if err := LoadServers(ctx); err != nil {
		return err
	}
	if err := TestServerConfig(); err != nil {
		return err
	}
	if err := TestRequiredServers(ctx); err != nil {
		return err
	}

	return nil
}

func setup() {
	CreateGroups()
	CreatePool()
}

func TestServerConfig() error {
	// test if Servers have the user that they can use
	for _, server := range DbPool.Servers {
		if _, err := server.Config.GetUser(); err != nil {
			return err
		}
	}

	return nil
}

// TestRequiredServers tests if all the Servers that are required are available
func TestRequiredServers(ctx context.Context) error {
	for id, server := range DbPool.Servers {
		if server.Config.Required {
			log.Logger.Debug("server is required, checking connectivity", zap.String("server_id", id), zap.String("dsn", server.Config.GetDsn()))
			if err := server.TestConnection(ctx); err != nil {
				return fmt.Errorf("error encountered while attempting to test connection with the required db - %s", err.Error())
			} else {
				server.Status = OPERATIONAL
			}
		}
	}

	return nil
}
