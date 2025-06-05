// @title Board Game API
// @version 1.0
// @description API для настольных игр
// @host localhost:8080
// @BasePath /api/v1
package main

import (
	"context"
	"log"

	"github.com/board-box/backend/internal/app"
)

func main() {
	ctx := context.Background()

	a, err := app.NewApp(ctx)
	if err != nil {
		log.Fatalf("could not create app: %v\n", err)
		return
	}

	if err = a.Run(); err != nil {
		log.Fatalf("could not create app: %v\n", err)
	}
}
