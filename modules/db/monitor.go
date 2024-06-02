package db

import (
	"context"
	"go-proxy/modules/log"
	"go.uber.org/zap"
	"time"
)

func MonitorServers(ctx context.Context) {
	go func() {
		ticker := time.NewTicker(1 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				log.Logger.Info("Context canceled, shutting down the server monitoring")
				return
			default:
				for range ticker.C {
					for _, server := range DbPool.Servers {
						err := server.TestConnection(ctx)
						if err != nil {
							server.Status = SHUNNED
							log.Logger.Warn("No connection with the server, server is shunned", zap.NamedError("reason", err))
						} else {
							server.Status = OPERATIONAL
						}
					}
				}
			}
		}
	}()
}
