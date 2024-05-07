package db

import (
	"context"
	"proxy/modules/log"
	"time"
)

func MonitorServers(ctx context.Context) {
	go func() {
		ticker := time.NewTicker(1 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				log.Logger.Tracef("Context canceled, shutting down the server monitoring")
				return
			default:
				for range ticker.C {
					for _, server := range DbPool.Servers {
						log.Logger.Tracef("Monitoring server: %s", server.Config.Id)
						err := server.TestConnection()
						if err != nil {
							server.Status = SHUNNED
							log.Logger.Tracef("No connection with the server %s, server is shunned", err)
						} else {
							server.Status = OPERATIONAL
						}
					}
					log.Logger.Tracef("Server: %v", DbPool.Servers)
				}
			}
		}
	}()
}
