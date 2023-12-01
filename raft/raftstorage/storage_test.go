package raftstorage

import (
	"github.com/stretchr/testify/suite"
	"testing"
)

type StorageSuite struct {
	suite.Suite
}

func TestStorage(t *testing.T) {
	suite.Run(t, new(StorageSuite))
}

func (s *StorageSuite) TestMakeKey() {
	rs := ReplicaStorage{serviceName: "ssp_proxy"}

	prefix := rs.makeKey("")
	s.Equal(storagePrefix+"_ssp_proxy_", prefix)

	s.Equal(storagePrefix+"_ssp_proxy_192.168.1.1:50051", rs.makeKey("192.168.1.1:50051"))

	s.Equal(storagePrefix+"_ssp_proxy_*", rs.makeKey("*"))
}
