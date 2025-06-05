package app

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/board-box/backend/docs"
	"github.com/board-box/backend/internal/auth"
	"github.com/board-box/backend/internal/config"
	chatHandler "github.com/board-box/backend/internal/handler/chat"
	collectionHandler "github.com/board-box/backend/internal/handler/collection"
	gameHandler "github.com/board-box/backend/internal/handler/game"
	userHandler "github.com/board-box/backend/internal/handler/user"
	"github.com/board-box/backend/internal/service/chat"
	"github.com/board-box/backend/internal/service/collection"
	"github.com/board-box/backend/internal/service/game"
	"github.com/board-box/backend/internal/service/user"
	"github.com/gin-gonic/gin"
	pgx "github.com/jackc/pgx/v5"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type App struct {
	cfg *config.Config
	jwt *auth.JWTManager

	r  *gin.Engine
	db *pgx.Conn

	authMW func(c *gin.Context)

	chatSvc       *chat.Service
	gameSvc       *game.Service
	userSvc       *user.Service
	collectionSvc *collection.Service
}

func NewApp(ctx context.Context) (*App, error) {
	a := &App{}

	if err := a.initDeps(ctx); err != nil {
		return nil, err
	}

	return a, nil
}

func (a *App) Run() error {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := a.r.Run(a.cfg.Addr()); err != nil {
			panic("server error: " + err.Error())
		}
	}()

	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := a.db.Close(ctx); err != nil {
		return fmt.Errorf("error closing db: %w", err)
	}

	return nil
}

func (a *App) initDeps(ctx context.Context) error {
	inits := []func(context.Context) error{
		a.initConfigs,
		a.initMiddleware,
		a.initDB,
		a.initService,
		a.initRouter,
	}

	for _, fn := range inits {
		if err := fn(ctx); err != nil {
			return err
		}
	}

	return nil
}

func (a *App) initConfigs(_ context.Context) error {
	var err error
	a.cfg, err = config.New()
	if err != nil {
		return err
	}

	a.jwt = auth.NewJWTManager(a.cfg.JWT.SecretKey, a.cfg.JWT.TokenDuration)
	return nil
}

func (a *App) initMiddleware(_ context.Context) error {
	a.authMW = auth.Middleware(a.jwt.SecretKey)
	return nil
}

func (a *App) initDB(_ context.Context) error {
	var err error
	ctx := context.Background()

	a.db, err = pgx.Connect(ctx, a.cfg.DSN())
	if err != nil {
		return fmt.Errorf("unable to connect to database: %v\n", err)
	}

	return nil
}

func (a *App) initService(_ context.Context) error {
	a.chatSvc = chat.NewService(a.cfg.ChatApiKey)
	a.gameSvc = game.NewService(a.db)
	a.userSvc = user.NewService(a.db, a.jwt)
	a.collectionSvc = collection.NewService(a.db, a.gameSvc)
	return nil
}

func (a *App) initRouter(_ context.Context) error {
	a.r = gin.Default()

	api := a.r.Group("/api/v1")

	gameRouter := gameHandler.New(a.gameSvc)
	gameRouter.RegisterRoutes(api)

	collectionRouter := collectionHandler.New(a.collectionSvc, a.authMW)
	collectionRouter.RegisterRoutes(api)

	userRouter := userHandler.New(a.userSvc, a.authMW)
	userRouter.RegisterRoutes(api)

	chatRouter := chatHandler.New(a.chatSvc, a.authMW)
	chatRouter.RegisterRoutes(api)

	a.r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return nil
}
