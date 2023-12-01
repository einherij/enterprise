package raft

import (
	"context"
	"crypto/tls"
	"fmt"
	"math/rand"
	"sync"
	"sync/atomic"
	"time"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/einherij/enterprise/raft/raftgrpc"
	"github.com/einherij/enterprise/raft/raftgrpc/protocol"
	"github.com/einherij/enterprise/raft/raftstorage"
)

var _ = ReplicaInterface(&Replica{})

type ReplicaInterface interface {
	RegisterCommand(commandName string, command raftgrpc.Command)
	ExecuteCommand(commandName string, sharedData []byte) error
}

const (
	grpcCommandTimeout = time.Minute     // rise this, if it's not enough time to execute command
	grpcConnectTimeout = 5 * time.Second // rise this, if not enough time to communicate between replicas

	// don't change options below if you're not sure

	grpcMaxExecutionDuration = grpcConnectTimeout * 15 / 10 // +50% to finish grpc command execution

	appendEntries = grpcMaxExecutionDuration * 2 // to finish grpc command execution and don't do it to often

	electionFrom = grpcMaxExecutionDuration * 3 // so it won't become candidate and will wait for append entries
	electionTo   = grpcMaxExecutionDuration * 6 // enough time difference between candidates' elections
)

// time until follower becomes a candidate
func electionTimeout() time.Duration {
	return timeoutInRange(electionFrom, electionTo)
}

// time between heartbeats
func appendEntriesTimeout() time.Duration {
	return appendEntries
}

func timeoutInRange(min, max time.Duration) time.Duration {
	var (
		difference           = max - min
		randomTimeDifference = time.Duration(rand.Int()) % difference
	)
	return min + randomTimeDifference
}

// Replica participates in elections, if it becomes a leader sends commands to follower servers
type Replica struct {
	log *logrus.Logger

	server  *raftgrpc.ReplicaServer
	storage raftstorage.Storage
}

func NewReplica(storage raftstorage.Storage, server *raftgrpc.ReplicaServer, logger *logrus.Logger) *Replica {
	r := &Replica{
		log:     logger,
		server:  server,
		storage: storage,
	}
	return r
}

func (r *Replica) Run(ctx context.Context) {
	electionTimer := time.NewTimer(electionTimeout())
	appendEntriesTimer := time.NewTimer(0)
	<-appendEntriesTimer.C               // need empty timer
	r.server.SetState(raftgrpc.Follower) // default state

mainLoop:
	for {
		switch r.server.GetState() {
		case raftgrpc.Follower:
			// listen for coordinator command or election requests
			// if there is no command change state to Candidate
			// if there is leader command, then do it, and send done message to leader
			// if there is election request change state to candidate and send vote to server, that requested
			select {
			case votedFor := <-r.server.IncomingElectionRequests():
				r.log.Debugf("%s accepted election request from: %s, term: %d", r.storage.GetMyAddress(), votedFor, r.server.GetTerm())

				electionTimer.Reset(electionTimeout())
			case leaderAddress := <-r.server.IncomingHeartbeats():
				r.log.Debugf("%s received heartbeat from %s, term: %d", r.storage.GetMyAddress(), leaderAddress, r.server.GetTerm())

				electionTimer.Reset(electionTimeout())
			case <-electionTimer.C:
				r.server.SetState(raftgrpc.Candidate)
				r.log.Debugf("%s state is changed to %v", r.storage.GetMyAddress(), r.server.GetState())
			case <-ctx.Done():
				return
			}
		case raftgrpc.Candidate:
			// start elections timeout
			// vote for myself, if I've started elections
			// send election requests
			var (
				term      = r.server.NewTerm()
				myAddress = r.storage.GetMyAddress()

				started     = time.Now()
				wg          sync.WaitGroup
				votesAmount atomic.Int32
			)
			replicas, err := r.storage.GetReplicas()
			if err != nil {
				r.log.Errorf("error getting replicas: %v", err)
				r.server.SetState(raftgrpc.Follower)
				r.log.Debugf("%s state is changed to %v", r.storage.GetMyAddress(), r.server.GetState())
				electionTimer.Reset(electionTimeout())
				continue mainLoop
			}

			votesAmount.Store(1) // self vote
			for _, replica := range replicas {
				replica := replica
				if replica == myAddress {
					continue
				}
				wg.Add(1)
				go func() {
					defer wg.Done()
					if voted := r.sendElectionRequest(replica, myAddress, term); voted {
						votesAmount.Add(1) // replica vote
					}
				}()
			}
			wg.Wait()
			r.log.Debugf("%s sending election requests duration: %v", r.storage.GetMyAddress(), time.Now().Sub(started))
			r.log.Debugf("%s voted %d followers from %d", r.storage.GetMyAddress(), votesAmount.Load(), len(replicas))
			if float32(votesAmount.Load()) > float32(len(replicas))/2. {
				r.server.SetState(raftgrpc.Leader)
				appendEntriesTimer.Reset(0)
			} else {
				r.server.SetState(raftgrpc.Follower)
				electionTimer.Reset(electionTimeout())
			}
			r.log.Debugf("%s state is changed to %v", r.storage.GetMyAddress(), r.server.GetState())
		case raftgrpc.Leader:
			// send keep alive or send commands to followers, wait for respond
			// perform command on current node if more than 50% of followers did command
			select {
			case <-appendEntriesTimer.C:
				var (
					started             = time.Now()
					myAddress           = r.storage.GetMyAddress()
					wg                  sync.WaitGroup
					heartbeatsResponded atomic.Int32
				)
				replicas, err := r.storage.GetReplicas()
				if err != nil {
					r.server.SetState(raftgrpc.Follower)
					electionTimer.Reset(electionTimeout())
					continue mainLoop
				}

				heartbeatsResponded.Store(1) // self heartbeat
				for _, follower := range replicas {
					follower := follower
					if follower == myAddress {
						continue
					}
					wg.Add(1)
					go func() {
						defer wg.Done()
						if ok := r.sendHeartBeat(follower, r.storage.GetMyAddress(), r.server.GetTerm()); ok {
							heartbeatsResponded.Add(1) // replica vote
						}
					}()
				}
				wg.Wait()
				r.log.Debugf("%s sending heartbeats duration: %v", r.storage.GetMyAddress(), time.Now().Sub(started))
				r.log.Debugf("%s responded %d heartbeats", r.storage.GetMyAddress(), heartbeatsResponded.Load())
				if float32(heartbeatsResponded.Load()) <= float32(len(replicas))/2. {
					r.server.SetState(raftgrpc.Follower)
					r.log.Debugf(r.storage.GetMyAddress(), "state is changed to", r.server.GetState())
					electionTimer.Reset(electionTimeout())
				}

				appendEntriesTimer.Reset(appendEntriesTimeout())
			case <-ctx.Done():
				return
			}
		}
	}
}

func (r *Replica) RegisterCommand(commandName string, command raftgrpc.Command) {
	r.server.AddCommand(commandName, command)
}

func (r *Replica) ExecuteCommand(commandName string, sharedData []byte) error {
	if r.server.GetState() != raftgrpc.Leader {
		return nil
	}
	replicas, err := r.storage.GetReplicas()
	if err != nil {
		return fmt.Errorf("error getting replilcas: %w", err)
	}
	for _, replica := range replicas {
		if ok := r.sendExecuteCommand(replica, commandName, sharedData); !ok {
			return fmt.Errorf("%s replica %s didn't execute the command %s", r.storage.GetMyAddress(), replica, commandName)
		}
	}
	return nil
}

func (r *Replica) sendExecuteCommand(toReplica, commandName string, sharedData []byte) (done bool) {
	var err = grpcSingleCall(toReplica, func(ctx context.Context, client protocol.FollowerClient) error {
		_, err := client.SendExecuteCommand(ctx, &protocol.Command{Name: commandName, SharedData: sharedData})
		return err
	})

	return err == nil
}

func (r *Replica) sendElectionRequest(toReplica, myAddress string, term uint64) (voted bool) {
	err := grpcSingleCall(toReplica, func(ctx context.Context, client protocol.FollowerClient) error {
		response, err := client.SendElectionRequest(ctx, &protocol.ElectionRequest{
			Address: myAddress,
			Term:    term,
		})
		voted = response.GetVote()
		return err
	})
	if err != nil {
		r.log.Errorf("error sending grpc single call: %v", err)
	}
	return voted
}

func (r *Replica) sendHeartBeat(toReplica, myAddress string, term uint64) (ok bool) {
	var err = grpcSingleCall(toReplica, func(ctx context.Context, client protocol.FollowerClient) error {
		response, err := client.SendHeartBeat(ctx, &protocol.HeartbeatRequest{
			LeaderAddress: myAddress,
			Term:          term,
		})
		ok = response.GetOk()
		return err
	})
	if err != nil {
		r.log.Errorf("error sending grpc single call: %v", err)
	}
	return ok
}

type singleCallFunc func(ctx context.Context, client protocol.FollowerClient) error

func grpcSingleCall(toReplica string, call singleCallFunc) error {
	client, conn, err := establishConnection(toReplica)
	if err != nil {
		return err
	}
	defer func() {
		_ = conn.Close()
	}()
	ctx, cancel := contextWithCommandTimeout()
	defer cancel()
	err = call(ctx, client)
	if err != nil {
		return err
	}
	return nil
}

func establishConnection(serverAddress string) (protocol.FollowerClient, *grpc.ClientConn, error) {
	const externalServer = false
	// use tls credentials for external grpc server
	var transportCreds credentials.TransportCredentials
	if externalServer {
		transportCreds = credentials.NewTLS(&tls.Config{})
	} else {
		transportCreds = insecure.NewCredentials()
	}

	ctx, cancel := context.WithTimeout(context.Background(), grpcConnectTimeout)
	defer cancel()
	conn, err := grpc.DialContext(ctx, serverAddress, grpc.WithTransportCredentials(transportCreds))
	if err != nil {
		return nil, nil, fmt.Errorf("cannot connect to replica at %q: %+v", serverAddress, err)
	}
	return protocol.NewFollowerClient(conn), conn, nil
}

func contextWithCommandTimeout() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), grpcCommandTimeout)
}
