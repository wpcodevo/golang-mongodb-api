package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/wpcodevo/golang-mongodb/config"
	"github.com/wpcodevo/golang-mongodb/controllers"
	"github.com/wpcodevo/golang-mongodb/gapi"
	"github.com/wpcodevo/golang-mongodb/pb"
	"github.com/wpcodevo/golang-mongodb/routes"
	"github.com/wpcodevo/golang-mongodb/services"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var (
	server      *gin.Engine
	ctx         context.Context
	mongoclient *mongo.Client
	redisclient *redis.Client

	userService         services.UserService
	UserController      controllers.UserController
	UserRouteController routes.UserRouteController

	authCollection      *mongo.Collection
	authService         services.AuthService
	AuthController      controllers.AuthController
	AuthRouteController routes.AuthRouteController

	// 👇 Create the Post Variables
	postService         services.PostService
	PostController      controllers.PostController
	postCollection      *mongo.Collection
	PostRouteController routes.PostRouteController
)

func init() {
	config, err := config.LoadConfig(".")
	if err != nil {
		log.Fatal("Could not load environment variables", err)
	}

	ctx = context.TODO()

	// Connect to MongoDB
	mongoconn := options.Client().ApplyURI(config.DBUri)
	mongoclient, err := mongo.Connect(ctx, mongoconn)

	if err != nil {
		panic(err)
	}

	if err := mongoclient.Ping(ctx, readpref.Primary()); err != nil {
		panic(err)
	}

	fmt.Println("MongoDB successfully connected...")

	// Connect to Redis
	redisclient = redis.NewClient(&redis.Options{
		Addr: config.RedisUri,
	})

	if _, err := redisclient.Ping(ctx).Result(); err != nil {
		panic(err)
	}

	err = redisclient.Set(ctx, "test", "Welcome to Golang with Redis and MongoDB", 0).Err()
	if err != nil {
		panic(err)
	}

	fmt.Println("Redis client connected successfully...")

	// Collections
	authCollection = mongoclient.Database("golang_mongodb").Collection("users")
	userService = services.NewUserServiceImpl(authCollection, ctx)
	authService = services.NewAuthService(authCollection, ctx)
	AuthController = controllers.NewAuthController(authService, userService, ctx, authCollection)
	AuthRouteController = routes.NewAuthRouteController(AuthController)

	UserController = controllers.NewUserController(userService)
	UserRouteController = routes.NewRouteUserController(UserController)

	// 👇 Instantiate the Constructors
	postCollection = mongoclient.Database("golang_mongodb").Collection("posts")
	postService = services.NewPostService(postCollection, ctx)
	PostController = controllers.NewPostController(postService)
	PostRouteController = routes.NewPostControllerRoute(PostController)

	server = gin.Default()
}

func main() {
	config, err := config.LoadConfig(".")

	if err != nil {
		log.Fatal("Could not load config", err)
	}

	defer mongoclient.Disconnect(ctx)

	// Start gRPC server in a separate goroutine
	go startGrpcServer(config)

	// Start Gin server in the main goroutine
	startGinServer(config)

}

func startGrpcServer(config config.Config) {
	authServer, err := gapi.NewGrpcAuthServer(config, authService, userService, authCollection)
	if err != nil {
		log.Fatal("cannot create grpc authServer: ", err)
	}

	userServer, err := gapi.NewGrpcUserServer(config, userService, authCollection)
	if err != nil {
		log.Fatal("cannot create grpc userServer: ", err)
	}

	postServer, err := gapi.NewGrpcPostServer(postCollection, postService)
	if err != nil {
		log.Fatal("cannot create grpc postServer: ", err)
	}

	grpcServer := grpc.NewServer()

	pb.RegisterAuthServiceServer(grpcServer, authServer)
	pb.RegisterUserServiceServer(grpcServer, userServer)
	// 👇 Register the Post gRPC service
	pb.RegisterPostServiceServer(grpcServer, postServer)
	reflection.Register(grpcServer)

	listener, err := net.Listen("tcp", config.GrpcServerAddress)
	if err != nil {
		log.Fatal("cannot create grpc server: ", err)
	}

	log.Printf("start gRPC server on %s", listener.Addr().String())
	err = grpcServer.Serve(listener)
	if err != nil {
		log.Fatal("cannot create grpc server: ", err)
	}
}

func startGinServer(config config.Config) {
	value, err := redisclient.Get(ctx, "test").Result()

	if err == redis.Nil {
		fmt.Println("key: test does not exist")
	} else if err != nil {
		panic(err)
	}

	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{config.Origin}
	corsConfig.AllowCredentials = true

	server.Use(cors.New(corsConfig))

	router := server.Group("/api")
	router.GET("/healthchecker", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{"status": "success", "message": value})
	})

	AuthRouteController.AuthRoute(router, userService)
	UserRouteController.UserRoute(router, userService)
	// 👇 Post Route
	PostRouteController.PostRoute(router)
	log.Fatal(server.Run(":" + config.Port))
}
