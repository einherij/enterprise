package raft

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
//	"github.com/einherij/enterprise/db"
//	"github.com/einherij/enterprise/raft/raftgrpc"
//	"github.com/einherij/enterprise/raft/raftgrpc/protocol"
//	"github.com/einherij/enterprise/raft/raftstorage"
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
//		replicasStorage := raftstorage.NewReplicaStorage("localhost:"+servicePort,
//			"raft",
//			redisClient)
//		ctx, cancel := context.WithCancel(context.Background())
//		go replicasStorage.Run(ctx)
//
//		srv := raftgrpc.NewReplicaServer(replicasStorage)
//
//		grpcConfig := grpcserver.GRPCServerConfig{Port: servicePort}
//		grpcSrv, err := grpcserver.New(grpcConfig)
//		a.NoError(err)
//		grpcSrv.RegisterService(&protocol.Follower_ServiceDesc, protocol.FollowerServer(srv))
//		go grpcSrv.Run(ctx)
//
//		log := logrus.New()
//		log.Level = logrus.DebugLevel
//		r := NewReplica(replicasStorage, srv, log)
//
//		r.RegisterCommand("test_command", func(ctx context.Context, myAddress string, myState raftgrpc.State, replicaCount int) error {
//			fmt.Println(myState, myAddress, "task started", replicaCount)
//			<-time.After(electionTimeout())
//			fmt.Println(myState, myAddress, "task executed", replicaCount)
//			return nil
//		})
//
//		r.Run(ctx)
//
//		go func() {
//			<-time.After(20 * time.Second)
//			a.NoError(r.ExecuteCommand("test_command"))
//		}()
//
//		go func() {
//			<-time.After(120 * time.Second)
//			cancel()
//		}()
//	}
//
//	<-time.After(70 * time.Second)
//}
