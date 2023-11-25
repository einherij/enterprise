package raftstorage

import (
	"context"
	"errors"
	"fmt"
	"github.com/einherij/enterprise/utils"
	"strings"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/sirupsen/logrus"

	app "github.com/einherij/enterprise"
)

const (
	registrationTimeout = 2 * time.Minute
	updateInterval      = time.Minute
	redisTimeout        = 5 * time.Second
	storagePrefix       = "REGISTER"
)

type Storage interface {
	GetMyAddress() string
	GetReplicas() ([]string, error)
	app.Runner
}

type ReplicaStorage struct {
	myAddress   string
	serviceName string
	redisClient *redis.Client
	logger      *logrus.Logger
	app.Runner
}

func NewReplicaStorage(myAddress string, serviceName string, redisClient *redis.Client) *ReplicaStorage {
	if myAddress == "" {
		panic(errors.New("empty self my_address"))
	}
	if serviceName == "" {
		panic(errors.New("empty service_name"))
	}
	rs := &ReplicaStorage{
		myAddress:   myAddress,
		serviceName: serviceName,
		redisClient: redisClient,
		logger:      logrus.New(),
	}
	return rs
}
func (rs *ReplicaStorage) GetMyAddress() string {
	return rs.myAddress
}

func (rs *ReplicaStorage) selfRegister() error {
	const value = "SET"
	ctx, cancel := context.WithTimeout(context.Background(), redisTimeout)
	defer cancel()
	return rs.redisClient.Set(ctx, rs.makeKey(rs.myAddress), value, registrationTimeout).Err()
}

func (rs *ReplicaStorage) GetReplicas() (replicas []string, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), redisTimeout)
	defer cancel()
	keys, err := rs.redisClient.Keys(ctx, rs.makeKey("*")).Result()
	if err != nil {
		return nil, fmt.Errorf("error getting keys: %w", err)
	}
	for _, key := range keys {
		replicas = append(replicas, strings.TrimPrefix(key, rs.makeKey("")))
	}
	return
}

func (rs *ReplicaStorage) makeKey(address string) string {
	return strings.Join([]string{storagePrefix, rs.serviceName, address}, "_")
}

func (rs *ReplicaStorage) run(ctx context.Context) {
	timer := time.NewTimer(0)
	for {
		select {
		case <-timer.C:
			if err := rs.selfRegister(); err != nil {
				rs.logger.Errorf("error registering replica: %v", err)
			}

			timer.Reset(utils.DurationUntilNextInterval(time.Now(), updateInterval))
		case <-ctx.Done():
			return
		}
	}
}

type dummyStorage struct {
	myAddress string
	addresses []string
	app.Runner
}

func NewDummyStorage(myAddress string, allAddresses ...string) Storage {
	return &dummyStorage{
		addresses: allAddresses,
		myAddress: myAddress,
		Runner:    app.NewRunner("dummy", func(_ context.Context) {}),
	}
}
func (ds *dummyStorage) GetMyAddress() string {
	return ds.myAddress
}
func (ds *dummyStorage) GetReplicas() ([]string, error) {
	return ds.addresses, nil
}
