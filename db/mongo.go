package db

import (
	"context"
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

const (
	defaultConnectTimeout = time.Minute
)

type MongoDBConfig struct {
	AppName   string `mapstructure:"app_name"`
	Host      string `mapstructure:"host"`
	Port      string `mapstructure:"port"`
	Username  string `mapstructure:"username"`
	Password  string `mapstructure:"password"`
	AuthDB    string `mapstructure:"auth_db"`
	UseDBName string `mapstructure:"database"`
}

type MongoClient struct {
	Client   *mongo.Client
	Database *mongo.Database
}

func NewMongoClient(cfg MongoDBConfig) (*MongoClient, error) {
	mongoDBAddress := fmt.Sprintf("mongodb://%s:%s/%s", cfg.Host, cfg.Port, cfg.UseDBName)
	opt := options.Client()
	opt.ApplyURI(mongoDBAddress)
	opt.SetAppName(cfg.AppName)
	opt.Auth = &options.Credential{
		Username:   cfg.Username,
		Password:   cfg.Password,
		AuthSource: cfg.AuthDB,
	}

	ctx, cancel := context.WithTimeout(context.Background(), defaultConnectTimeout)
	defer cancel()
	client, err := mongo.Connect(ctx, opt)
	if err != nil {
		return nil, fmt.Errorf("[NewMongoClient] cannot create mongodb client: %v", err)
	}
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		return nil, fmt.Errorf("[NewMongoClient] mongodb ping at %s failed: %v", mongoDBAddress, err)
	}

	db := client.Database(cfg.UseDBName)
	return &MongoClient{
		Client:   client,
		Database: db,
	}, nil
}

// Disconnect from mongodb
// If error occur on disconnection nothing happen
func (c *MongoClient) Disconnect() {
	ctx, cancel := context.WithTimeout(context.Background(), defaultConnectTimeout)
	defer cancel()

	err := c.Client.Disconnect(ctx)
	if err != nil {
		logrus.Errorf("[Disconnect] Error while try to close connection to mongodb: %v", err)
		return
	}
	logrus.Info("[Disconnect] success disconnect mongodb")
}
