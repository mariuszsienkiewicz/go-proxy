package db

import (
	"fmt"
	"proxy/modules/log"
)

// Init - should be loaded in logic order, groups -> servers -> pool
func Init() error {
	if err := LoadGroups(); err != nil {
		return err
	}
	if err := LoadServers(); err != nil {
		return err
	}
	if err := TestServerConfig(); err != nil {
		return err
	}
	if err := TestRequiredServers(); err != nil {
		return err
	}

	return nil
}

func TestServerConfig() error {
	// test if servers have the user that they can use
	for _, server := range DbPool.Servers {
		if _, err := server.Config.GetUser(); err != nil {
			return err
		}
	}

	return nil
}

// TestRequiredServers tests if all the servers that are required are available
func TestRequiredServers() error {
	for id, server := range DbPool.Servers {
		if server.Config.Required {
			log.Logger.Tracef("Server %v is required, checking connectivity via %v", id, server.Config.GetDsn())
			if err := server.TestConnection(); err != nil {
				return fmt.Errorf("error encountered while attempting to test connection with the required db - %s", err.Error())
			} else {
				server.Status = OPERATIONAL
			}

			log.Logger.Tracef("Server %v is recheable", id)
		}
	}

	return nil
}
