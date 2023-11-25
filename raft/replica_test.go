package raft

//
//import (
//	"context"
//	"fmt"
//	"strconv"
//	"testing"
//	"time"
//
//	"github.com/sirupsen/logrus"
//	"github.com/stretchr/testify/assert"
//
//)
//
//func Test_Replica(t *testing.T) {
//	a := assert.New(t)
//
//	redisClient, err := db.NewRedisClient(db.RedisConfig{
//		Host: "localhost",
//		Port: 6379,
//	})
//	a.NoError(err)
//
//	for i := 1; i <= 5; i++ {
//		servicePort := "414" + strconv.Itoa(i)
//
//		replicasStorage := raft_storage.NewReplicaStorage(raft_storage.ReplicaStorageConfig{
//			MyAddress:   "localhost:" + servicePort,
//			ServiceName: "raft",
//		}, redisClient)
//		replicasStorage.Start()
//
//		srv := raft_grpc.NewReplicaServer(replicasStorage)
//
//		grpcConfig := grpcserver.GRPCServerConfig{Port: servicePort}
//		grpcSrv, err := grpcserver.New(grpcConfig)
//		a.NoError(err)
//		grpcSrv.RegisterService(&protocol.Follower_ServiceDesc, protocol.FollowerServer(srv))
//		go grpcSrv.Run()
//
//		log := logrus.New()
//		log.Level = logrus.DebugLevel
//		r := NewReplica(replicasStorage, srv, log)
//
//		r.RegisterCommand("test_command", func(ctx context.Context, myAddress string, myState raft_grpc.State, replicaCount int) error {
//			fmt.Println(myState, myAddress, "task started", replicaCount)
//			<-time.After(electionTimeout())
//			fmt.Println(myState, myAddress, "task executed", replicaCount)
//			return nil
//		})
//
//		r.Start()
//
//		go func() {
//			<-time.After(500 * time.Millisecond)
//			if r.server.GetState() == raft_grpc.Leader {
//				r.Stop()
//			}
//			a.NoError(r.ExecuteCommand("test_command"))
//		}()
//
//		go func() {
//			<-time.After(5 * time.Second)
//			r.Stop()
//			grpcSrv.Stop()
//			replicasStorage.Stop()
//		}()
//	}
//
//	<-time.After(6 * time.Second)
//}
