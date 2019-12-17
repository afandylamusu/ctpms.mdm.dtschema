package main

import (
	"database/sql"
	"fmt"
	"log"
	"net"
	"net/url"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/labstack/echo"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"

	_articleHttpDeliver "github.com/afandylamusu/ctpms.mdm.dtschema/article/delivery/http"
	_articleRepo "github.com/afandylamusu/ctpms.mdm.dtschema/article/repository"
	_articleUcase "github.com/afandylamusu/ctpms.mdm.dtschema/article/usecase"
	_authorRepo "github.com/afandylamusu/ctpms.mdm.dtschema/author/repository"
	_deliverGrpc "github.com/afandylamusu/ctpms.mdm.dtschema/dataset/delivery/delivergrpc"
	"github.com/afandylamusu/ctpms.mdm.dtschema/middleware"
)

func init() {
	viper.SetConfigFile(`config.json`)
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}

	if viper.GetBool(`debug`) {
		fmt.Println("Service RUN on DEBUG mode")
	}
}

// RunGRPCServer to starting GRPC Server
func RunGRPCServer(port string) error {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer(
		grpc.ConnectionTimeout(time.Second),
		grpc.KeepaliveParams(keepalive.ServerParameters{
			MaxConnectionIdle: time.Second * 10,
			Timeout:           time.Second * 20,
		}),
		grpc.KeepaliveEnforcementPolicy(
			keepalive.EnforcementPolicy{
				MinTime:             time.Second,
				PermitWithoutStream: true,
			}),
		grpc.MaxConcurrentStreams(5),
	)

	_deliverGrpc.RegisterAddServiceServer(s, &_deliverGrpc.DataSetServiceHandler{Port: port})
	log.Println("Run GRPC AddServiceServer: " + port)

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

	return err
}

func main() {
	isLocal := viper.GetString("env") == "local"
	var dbHost, dbPort, dbUser, dbPass, dbName string

	if isLocal {
		dbHost = viper.GetString(`database-local.host`)
		dbPort = viper.GetString(`database-local.port`)
		dbUser = viper.GetString(`database-local.user`)
		dbPass = viper.GetString(`database-local.pass`)
		dbName = viper.GetString(`database-local.name`)
	} else {
		dbHost = viper.GetString(`database.host`)
		dbPort = viper.GetString(`database.port`)
		dbUser = viper.GetString(`database.user`)
		dbPass = viper.GetString(`database.pass`)
		dbName = viper.GetString(`database.name`)
	}

	connection := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", dbUser, dbPass, dbHost, dbPort, dbName)
	val := url.Values{}
	val.Add("parseTime", "1")
	val.Add("loc", "Asia/Jakarta")
	dsn := fmt.Sprintf("%s?%s", connection, val.Encode())
	log.Println(dsn)

	dbConn, err := sql.Open(`mysql`, dsn)
	if err != nil && viper.GetBool("debug") {
		fmt.Println(err)
	}
	err = dbConn.Ping()
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	defer func() {
		err := dbConn.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()

	e := echo.New()
	middL := middleware.InitMiddleware()
	e.Use(middL.CORS)
	authorRepo := _authorRepo.NewMysqlAuthorRepository(dbConn)
	ar := _articleRepo.NewMysqlArticleRepository(dbConn)

	timeoutContext := time.Duration(viper.GetInt("context.timeout")) * time.Second
	au := _articleUcase.NewArticleUsecase(ar, authorRepo, timeoutContext)
	_articleHttpDeliver.NewArticleHandler(e, au)

	go RunGRPCServer(viper.GetString("server.grpc-port"))

	log.Fatal(e.Start(viper.GetString("server.http-port")))

	<-make(chan int)
}
