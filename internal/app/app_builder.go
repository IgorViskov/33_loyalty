package app

import (
	"flag"
	"github.com/IgorViskov/33_loyalty/internal/api"
	"github.com/IgorViskov/33_loyalty/internal/apperrors"
	"github.com/IgorViskov/33_loyalty/internal/config"
	"github.com/IgorViskov/33_loyalty/internal/data"
	"github.com/IgorViskov/33_loyalty/internal/data/migrator"
	"github.com/IgorViskov/33_loyalty/internal/services"
	"github.com/caarlos0/env/v11"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	"runtime"
)

const defaultHost = "localhost:8090"

type App struct {
	conf      *config.AppConfig
	router    *echo.Echo
	connector data.Connector
}

func Create() *App {
	return &App{
		router: echo.New(),
	}
}

func (app *App) ApplyMigrations() *App {
	err := migrator.AutoMigrate(app.connector)
	if err != nil {
		log.Fatal(err)
	}

	return app
}

func (app *App) Configure() *App {
	conf := config.AppConfig{
		PeriodRequests:      5,
		MaxParallelRequests: runtime.NumCPU(),
	}
	readFlags(&conf)
	err := readEnvironments(&conf)
	if err != nil {
		log.Error(err)
	}

	err = checkConfig(&conf)
	if err != nil {
		panic(err)
	}

	app.conf = &conf
	app.connector = data.NewConnector(&conf)
	return app
}

func (app *App) Build() *App {
	r := app.router
	useGzip(r)
	bindUserContext(r, app.connector)
	r.Use(AuthMiddleware)

	bindUserAPI(app.router)
	bindOrderAPI(app.router, app.connector, app.conf)

	return app
}

func (app *App) Start() {
	err := app.router.Start(app.conf.RunHost)
	if err != nil {
		panic(err)
	}
}

func readFlags(conf *config.AppConfig) {
	flag.Func("a", "Адрес запуска HTTP-сервера", config.HostNameParser(conf))
	flag.Func("d", "Адрес подключения к базе данных", config.DBURIParser(conf))
	flag.Func("r", "Адрес системы расчёта начислений", config.AccrualHostParser(conf))

	// запускаем парсинг
	flag.Parse()
}

func readEnvironments(conf *config.AppConfig) error {
	return env.Parse(conf)
}

func checkConfig(conf *config.AppConfig) error {
	if conf.RunHost == "" {
		log.Info(apperrors.InfoEmptyRunHost)
		conf.RunHost = defaultHost
	}
	if conf.AccrualHost == nil || conf.AccrualHost.Host == "" {
		return apperrors.ErrNotValidAccrualHost
	}
	if conf.DBURI == "" {
		return apperrors.ErrDBURIIsEmpty
	}

	return nil
}

func useGzip(r *echo.Echo) {
	r.Use(middleware.Gzip())
	r.Use(middleware.Decompress())
}

func bindUserContext(r *echo.Echo, con data.Connector) {
	service := services.NewUserService(data.NewUsersRepository(con), nil)
	r.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			cc := &api.UserContext{Context: c, UserService: service}
			return next(cc)
		}
	})
}

func bindUserAPI(r *echo.Echo) {
	r.POST("/api/user/register", api.Register)
	r.POST("/api/user/login", api.Login)
}

func bindOrderAPI(r *echo.Echo, con data.Connector, conf *config.AppConfig) {
	accruals := data.NewAccrualRepository(con)
	tasks := data.NewAccrualTasksRepository(con)
	account := services.NewAccountService(data.NewWithdrawalsRepository(con))
	withdrawals := data.NewWithdrawalsRepository(con)
	tasksPool := services.NewAccrualTasksPool(conf, tasks, accruals, services.NewExternalAccrualService(conf), account)
	withdrawalService := services.NewWithdrawService(withdrawals)
	orderService := services.NewOrdersService(accruals, tasksPool, withdrawals)
	controller := api.NewController(orderService, withdrawalService)

	r.POST("/api/user/orders", controller.RegisterOrder)
	r.GET("/api/user/orders", controller.GetAllRegisteredOrders)
	r.GET("/api/user/balance", controller.Balance)
	r.POST("/api/user/balance/withdraw", controller.Withdraw)
	r.GET("/api/user/withdrawals", controller.AllWithdraw)
}
