package api

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/OzyKleyton/studio-api/config"
	"github.com/OzyKleyton/studio-api/config/db"
	"github.com/OzyKleyton/studio-api/internal/api/handler"
	"github.com/OzyKleyton/studio-api/internal/api/router"
	"github.com/OzyKleyton/studio-api/internal/model/user"
	"github.com/OzyKleyton/studio-api/internal/repository"
	"github.com/OzyKleyton/studio-api/internal/service"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
)

func Run(host, port string) error {
	address := fmt.Sprintf("%s:%s", host, port)
	log.Println("Listen app in port ", address)

	app := fiber.New(fiber.Config{
		JSONEncoder: json.Marshal,
		JSONDecoder: json.Unmarshal,
		Prefork:     config.GetConfig().Prefork,
		ProxyHeader: fiber.HeaderXForwardedFor,
		ReadTimeout: 10 * time.Second,
	})

	// ========== MIDDLEWARES GLOBAIS ==========

	// Middleware de Recovery (captura panics)
	app.Use(recover.New(recover.Config{
		EnableStackTrace: true,
	}))

	// Middleware de CORS
	// app.Use(cors.New(cors.Config{
	// 	AllowOrigins:     "*", // Ou defina origens específicas
	// 	AllowMethods:     "GET,POST,PUT,DELETE,OPTIONS",
	// 	AllowHeaders:     "Origin,Content-Type,Accept,Authorization",
	// 	ExposeHeaders:    "Content-Length",
	// 	AllowCredentials: true,
	// 	MaxAge:           86400,
	// }))

	// Middleware de Logger
	app.Use(logger.New(logger.Config{
		Format:     "${time} ${ip} ${status} - ${method} ${path}\n",
		TimeFormat: "2006-01-02 15:04:05",
		TimeZone:   "America/Sao_Paulo",
	}))

	// Middleware de Request ID
	app.Use(requestid.New())

	// Middleware de health check (rota pública)
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":    "ok",
			"timestamp": time.Now(),
			"service":   "studio-api",
		})
	})

	// ========== CONEXÃO COM BANCO E MIGRAÇÕES ==========
	db, err := db.ConnectDB(config.GetConfig().DBURL)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	db = db.WithContext(ctx)

	if err := db.AutoMigrate(
		&user.User{},
	); err != nil {
		return err
	}

	// ========== INJEÇÃO DE DEPENDÊNCIAS ==========
	userRepo := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepo)
	userHandler := handler.NewUserHandler(userService)

	// ========== ROTAS ==========
	router.SetupRouter(app, userHandler.Routes())

	// ========== MIDDLEWARE DE 404 ==========
	app.Use(func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error":   "Endpoint não encontrado",
			"path":    c.Path(),
			"method":  c.Method(),
			"message": "Verifique a documentação da API",
		})
	})

	// ========== GRACEFUL SHUTDOWN ==========
	c := make(chan os.Signal, 1)
	errc := make(chan error, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)

	go func() {
		<-c
		fmt.Println("Gracefully shutting down...")
		cancel()
		errc <- app.Shutdown()
	}()

	if err := app.Listen(address); err != nil {
		return err
	}

	err = <-errc
	return err
}
