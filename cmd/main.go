package main

import (
    "fmt"
    "log"
    "net/http"

    "gorm.io/gorm"
    "gorm.io/driver/postgres"
    "github.com/redis/go-redis/v9"

    "github.com/labstack/echo/v4"
    "github.com/labstack/echo/v4/middleware"

    "github.com/totegamma/concurrent/x/association"
    "github.com/totegamma/concurrent/x/character"
    "github.com/totegamma/concurrent/x/message"
    "github.com/totegamma/concurrent/x/socket"
    "github.com/totegamma/concurrent/x/stream"
    "github.com/totegamma/concurrent/x/host"
    "github.com/totegamma/concurrent/x/util"
)

func main() {

    fmt.Print(concurrentBanner)

    e := echo.New()

    config := util.Config{}
    err := config.Load("/etc/concurrent/config.yaml")
    if err != nil {
        e.Logger.Fatal(err)
    }

    log.Print("Config loaded! I am: ", config.CCAddr)

    db, err := gorm.Open(postgres.Open(config.Dsn), &gorm.Config{})
    if err != nil {
        log.Println("failed to connect database");
        panic("failed to connect database")
    }

    // Migrate the schema
    log.Println("start migrate")
    db.AutoMigrate(&message.Message{}, &character.Character{}, &association.Association{},  &stream.Stream{}, &host.Host{})

    rdb := redis.NewClient(&redis.Options{
        Addr:     config.RedisAddr,
        Password: "", // no password set
        DB:       0,  // use default DB
    })

    socketService := socket.NewService();

    socketHandler := SetupSocketHandler(socketService)
    messageHandler := SetupMessageHandler(db, rdb, socketService)
    characterHandler := SetupCharacterHandler(db)
    associationHandler := SetupAssociationHandler(db, rdb, socketService)
    streamHandler := SetupStreamHandler(db, rdb)
    hostHandler := SetupHostHandler(db, config)

    e.HideBanner = true
    e.Use(middleware.CORS())
    e.Use(middleware.Logger())
    e.Use(middleware.Recover())

    apiV1 := e.Group("/api/v1")
    apiV1.GET("/messages/:id", messageHandler.Get)
    apiV1.POST("/messages", messageHandler.Post)
    apiV1.GET("/characters", characterHandler.Get)
    apiV1.PUT("/characters", characterHandler.Put)
    apiV1.GET("/associations/:id", associationHandler.Get)
    apiV1.POST("/associations", associationHandler.Post)
    apiV1.DELETE("/associations", associationHandler.Delete)
    apiV1.GET("/stream", streamHandler.Get)
    apiV1.POST("/stream", streamHandler.Post)
    apiV1.PUT("/stream", streamHandler.Put)
    apiV1.GET("/stream/recent", streamHandler.Recent)
    apiV1.GET("/stream/list", streamHandler.List)
    apiV1.GET("/stream/range", streamHandler.Range)
    apiV1.GET("/socket", socketHandler.Connect)
    apiV1.GET("/host/:id", hostHandler.Get) //TODO deprecated. remove later
    apiV1.PUT("/host", hostHandler.Upsert)
    apiV1.GET("/host", hostHandler.Profile)
    apiV1.GET("/host/list", hostHandler.List)
    apiV1.POST("/host/hello", hostHandler.Hello)
    apiV1.GET("/admin/sayhello/:fqdn", hostHandler.SayHello)

    e.Static("/", "/etc/www/concurrent")
    e.GET("/health", func(c echo.Context) (err error) {
        return c.String(http.StatusOK, "ok")
    })

    e.Logger.Fatal(e.Start(":8000"))
}

