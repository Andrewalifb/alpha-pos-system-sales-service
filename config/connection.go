package config

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/Andrewalifb/alpha-pos-system-sales-service/entity"
	"github.com/go-redis/redis/v8"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/streadway/amqp"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Config struct {
	SQLDB   *gorm.DB
	RedisDB *redis.Client
}

type GrpcClientConfig struct {
	ProductServiceConn *grpc.ClientConn
	CompanyServiceConn *grpc.ClientConn
}

type RabbitMqConfig struct {
	RabbitMQConn *amqp.Connection
}

func connectPostgres() *gorm.DB {
	host := os.Getenv("SQL_HOST")
	port, _ := strconv.Atoi(os.Getenv("SQL_PORT"))
	user := os.Getenv("SQL_USER")
	dbname := os.Getenv("SQL_DB_NAME")
	pass := os.Getenv("SQL_PASSWORD")

	psqlSetup := fmt.Sprintf("host=%s port=%d user=%s dbname=%s password=%s sslmode=disable",
		host, port, user, dbname, pass)

	sqlDB, err := gorm.Open("postgres", psqlSetup)

	if err != nil {
		fmt.Println("Failed to connect to PostgreSQL:", err)
		return nil
	} else {
		fmt.Println("Successfully connected to PostgreSQL")
		sqlDB.AutoMigrate(entity.PosCashDrawer{}, entity.PosInvoice{}, entity.PosOnlinePayment{}, entity.PosPaymentMethod{}, entity.PosReturn{}, entity.PosSale{}, entity.PosCustomer{})
		return sqlDB
	}
}

func connectRedis() *redis.Client {
	redisDB := redis.NewClient(&redis.Options{
		Addr: os.Getenv("REDIS_ADDR"),
		// Password: os.Getenv("REDIS_PASSWORD"),
		DB: 0,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_, err := redisDB.Ping(ctx).Result()
	if err != nil {
		fmt.Println("Failed to connect to Redis:", err)
		return nil
	} else {
		fmt.Println("Successfully connected to Redis")
		return redisDB
	}
}

func connectRabbitMQ() *amqp.Connection {
	conn, err := amqp.Dial(os.Getenv("RABBITMQ_URL"))
	if err != nil {
		fmt.Println("Failed to connect to RabbitMQ:", err)
		return nil
	} else {
		fmt.Println("Successfully connected to RabbitMQ")
		return conn
	}
}

func connectProductServiceGRPC() *grpc.ClientConn {
	productGrpcServicePort := os.Getenv("PRODUCT_GRPC")
	addr := fmt.Sprintf("localhost:%s", productGrpcServicePort)
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		fmt.Println("Failed to connect to Product Service gRPC:", err)
		return nil
	} else {
		fmt.Println("Successfully connected to Product Service gRPC")
		return conn
	}
}

func connectCompanyServiceGRPC() *grpc.ClientConn {
	companyGrpcServicePort := os.Getenv("COMPANY_GRPC")
	addr := fmt.Sprintf("localhost:%s", companyGrpcServicePort)
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		fmt.Println("Failed to connect to Company Service gRPC:", err)
		return nil
	} else {
		fmt.Println("Successfully connected to Company Service gRPC")
		return conn
	}
}

func NewConfig() *Config {
	return &Config{
		SQLDB:   connectPostgres(),
		RedisDB: connectRedis(),
	}
}

func NewRabbitMqCofig() *RabbitMqConfig {
	return &RabbitMqConfig{
		RabbitMQConn: connectRabbitMQ(),
	}
}

func NewGRPCConfig() *GrpcClientConfig {
	return &GrpcClientConfig{
		ProductServiceConn: connectProductServiceGRPC(),
		CompanyServiceConn: connectCompanyServiceGRPC(),
	}
}
