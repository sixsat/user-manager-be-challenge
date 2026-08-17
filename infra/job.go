package infra

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/sixsat/user-manager-be-challenge/port"
)

func StartBackgroundJob(ctx context.Context, userSvc port.UserService) {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			count, _ := userSvc.Count(ctx)
			slog.Info(fmt.Sprintf("[job] total number of users: %d", count))
		case <-ctx.Done():
			return
		}
	}

}
