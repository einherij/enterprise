package raftgrpc

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"github.com/einherij/enterprise/raft/raftgrpc/protocol"
	"github.com/einherij/enterprise/raft/raftstorage"
)

const (
	waitFollowerStateTimeout = 20 * time.Millisecond
)

type State uint

const (
	Follower State = iota
	Candidate
	Leader
)

func (s State) String() string {
	switch s {
	case Follower:
		return "Follower"
	case Candidate:
		return "Candidate"
	case Leader:
		return "Leader"
	}
	return "UnknownState"
}

type ReplicaServer struct {
	storage raftstorage.Storage

	commandsMux sync.RWMutex
	commands    map[string]Command
	state       atomic.Uint32
	term        atomic.Uint64

	heartbeats       chan string
	electionRequests chan string
	protocol.UnimplementedFollowerServer
}

type Command func(ctx context.Context, myAddress string, myState State, replicaCount int, sharedData []byte) error

func NewReplicaServer(storage raftstorage.Storage) *ReplicaServer {
	return &ReplicaServer{
		storage:          storage,
		commands:         make(map[string]Command),
		heartbeats:       make(chan string),
		electionRequests: make(chan string),
	}
}

func (rs *ReplicaServer) AddCommand(commandName string, command Command) {
	rs.commandsMux.Lock()
	defer rs.commandsMux.Unlock()
	rs.commands[commandName] = command
}

func (rs *ReplicaServer) SendExecuteCommand(ctx context.Context, command *protocol.Command) (*protocol.Nothing, error) {
	replicas, err := rs.storage.GetReplicas()
	if err != nil {
		return nil, fmt.Errorf("error getting replicas: %w", err)
	}
	rs.commandsMux.RLock()
	defer rs.commandsMux.RUnlock()
	err = rs.commands[command.GetName()](ctx, rs.storage.GetMyAddress(), State(rs.state.Load()), len(replicas), command.GetSharedData())
	if err != nil {
		return nil, fmt.Errorf("error executing command: %w", err)
	}
	return &protocol.Nothing{}, nil
}

func (rs *ReplicaServer) SendHeartBeat(ctx context.Context, heartbeatRequest *protocol.HeartbeatRequest) (*protocol.HeartbeatResponse, error) {

	if requestTerm := heartbeatRequest.GetTerm(); requestTerm >= rs.term.Load() {
		rs.term.Store(requestTerm)
		select {
		case rs.heartbeats <- heartbeatRequest.GetLeaderAddress():
		case <-time.After(waitFollowerStateTimeout):
			return nil, errors.New("replica is not in follower state")
		case <-ctx.Done():
			return nil, ctx.Err()
		}
		return &protocol.HeartbeatResponse{Ok: true}, nil
	}
	return &protocol.HeartbeatResponse{Ok: false}, nil
}

func (rs *ReplicaServer) SendElectionRequest(ctx context.Context, request *protocol.ElectionRequest) (*protocol.ElectionResponse, error) {

	if requestTerm := request.GetTerm(); requestTerm > rs.term.Load() {
		rs.term.Store(requestTerm)
		select {
		case rs.electionRequests <- request.GetAddress():
			return &protocol.ElectionResponse{Vote: true}, nil
		case <-time.After(waitFollowerStateTimeout):
			return nil, errors.New("replica is not in follower state")
		case <-ctx.Done():
			return nil, ctx.Err()
		}
	}
	return &protocol.ElectionResponse{Vote: false}, nil
}

func (rs *ReplicaServer) SetState(state State) {
	rs.state.Store(uint32(state))
}

func (rs *ReplicaServer) GetState() State {
	return State(rs.state.Load())
}

func (rs *ReplicaServer) NewTerm() uint64 {
	rs.term.Add(1)
	return rs.term.Load()
}

func (rs *ReplicaServer) GetTerm() uint64 {
	return rs.term.Load()
}

func (rs *ReplicaServer) IncomingElectionRequests() chan string {
	return rs.electionRequests
}

func (rs *ReplicaServer) IncomingHeartbeats() chan string {
	return rs.heartbeats
}
