package grpc

import (
	"context"
	_ "gitee.com/geekbang/basic-go/webook/pkg/grpcx/balancer/wrr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/naming/endpoints"
	"go.etcd.io/etcd/client/v3/naming/resolver"
	"google.golang.org/grpc"
	_ "google.golang.org/grpc/balancer/weightedroundrobin"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"net"
	"testing"
	"time"
)

type LoadBalancerTestSuite struct {
	suite.Suite
	// 借助 etcd 来做服务发现
	cli *clientv3.Client
}

func (s *LoadBalancerTestSuite) SetupSuite() {
	cli, err := clientv3.NewFromURL("http://localhost:12379")
	assert.NoError(s.T(), err)
	s.cli = cli
}

// TestServerFailover 启动测试 failover 的服务器
func (s *LoadBalancerTestSuite) TestServerFailover() {
	go func() {
		s.startFailoverServer(":8090")
	}()
	s.startWeightedServer(":8091", 20)
}

// TestServer 会启动两个服务器，一个监听 8090，一个监听 8091
func (s *LoadBalancerTestSuite) TestServer() {
	go func() {
		s.startWeightedServer(":8090", 10)
	}()
	s.startWeightedServer(":8091", 20)
}

func (s *LoadBalancerTestSuite) TestClientWeightedRoundRobin() {
	t := s.T()
	etcdResolver, err := resolver.NewBuilder(s.cli)
	assert.NoError(t, err)
	cc, err := grpc.Dial("etcd:///service/user",
		grpc.WithResolvers(etcdResolver),
		grpc.WithDefaultServiceConfig(`{"loadBalancingConfig": [{"weighted_round_robin":{}}]}`),
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	assert.NoError(t, err)
	userClient := NewUserServiceClient(cc)
	for i := 0; i < 10; i++ {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		resp, err := userClient.GetById(ctx, &GetByIdReq{
			Id: 123,
		})
		cancel()
		require.NoError(t, err)
		t.Log(resp.User)
	}
}

func (s *LoadBalancerTestSuite) TestRoundRobinFailover() {
	cfg := `
{
  "loadBalancingConfig": [{"round_robin":{}}],
  "methodConfig": [{
    "name": [{"service": "UserService"}],
    "retryPolicy": {
      "maxAttempts": 4,
      "initialBackoff": "0.01s",
      "maxBackoff": "0.1s",
      "backoffMultiplier": 2.0,
      "retryableStatusCodes": [ "UNAVAILABLE" ]
    }
  }]
}
`
	t := s.T()
	etcdResolver, err := resolver.NewBuilder(s.cli)
	assert.NoError(t, err)
	cc, err := grpc.Dial("etcd:///service/user",
		grpc.WithResolvers(etcdResolver),
		grpc.WithDefaultServiceConfig(cfg),
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	assert.NoError(t, err)
	userClient := NewUserServiceClient(cc)
	for i := 0; i < 10; i++ {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		resp, err := userClient.GetById(ctx, &GetByIdReq{
			Id: 123,
		})
		cancel()
		require.NoError(t, err)
		t.Log(resp.User)
	}
}

func (s *LoadBalancerTestSuite) TestClientRoundRobin() {
	t := s.T()
	etcdResolver, err := resolver.NewBuilder(s.cli)
	assert.NoError(t, err)
	cc, err := grpc.Dial("etcd:///service/user",
		grpc.WithResolvers(etcdResolver),
		grpc.WithDefaultServiceConfig(`{"loadBalancingConfig": [{"round_robin":{}}]}`),
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	assert.NoError(t, err)
	userClient := NewUserServiceClient(cc)
	for i := 0; i < 10; i++ {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		resp, err := userClient.GetById(ctx, &GetByIdReq{
			Id: 123,
		})
		cancel()
		require.NoError(t, err)
		t.Log(resp.User)
	}
}

// TestClient 默认情况下是啥负载均衡策略都没有的
func (s *LoadBalancerTestSuite) TestClientPickFirst() {
	t := s.T()
	etcdResolver, err := resolver.NewBuilder(s.cli)
	assert.NoError(t, err)
	cc, err := grpc.Dial("etcd:///service/user",
		grpc.WithResolvers(etcdResolver),
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	assert.NoError(t, err)
	userClient := NewUserServiceClient(cc)
	for i := 0; i < 10; i++ {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		resp, err := userClient.GetById(ctx, &GetByIdReq{
			Id: 123,
		})
		cancel()
		require.NoError(t, err)
		t.Log(resp.User)
	}
}

// TestClientCustomWeightedRondRobin 测试自定义的基于权重的负载均衡算法
func (s *LoadBalancerTestSuite) TestClientCustomWeightedRondRobin() {
	t := s.T()
	etcdResolver, err := resolver.NewBuilder(s.cli)
	assert.NoError(t, err)
	cc, err := grpc.Dial("etcd:///service/user",
		grpc.WithResolvers(etcdResolver),
		grpc.WithDefaultServiceConfig(`{"loadBalancingConfig": [{"custom_weighted_round_robin":{}}]}`),
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	assert.NoError(t, err)
	userClient := NewUserServiceClient(cc)
	for i := 0; i < 10; i++ {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		resp, err := userClient.GetById(ctx, &GetByIdReq{
			Id: 123,
		})
		cancel()
		require.NoError(t, err)
		t.Log(resp.User)
	}
}

func (s *LoadBalancerTestSuite) startFailoverServer(addr string) {
	t := s.T()
	em, err := endpoints.NewManager(s.cli,
		"service/user")
	assert.NoError(t, err)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	// 要以 /service/user 为前缀
	addr = "127.0.0.1" + addr
	key := "service/user/" + addr
	// 5s
	var ttl int64 = 5
	leaseResp, err := s.cli.Grant(ctx, ttl)
	// metadata 一般用在客户端
	err = em.AddEndpoint(ctx, key,
		endpoints.Endpoint{
			Addr: addr,
		}, clientv3.WithLease(leaseResp.ID))
	assert.NoError(t, err)

	// 忽略掉 ctx，因为在测试环境下，我们不需要手动控制退出续约
	kaCtx, _ := context.WithCancel(context.Background())
	go func() {
		_, err1 := s.cli.KeepAlive(kaCtx, leaseResp.ID)
		require.NoError(t, err1)
		//for resp := range ch {
		//	t.Log(resp.String())
		//}
	}()

	server := grpc.NewServer()
	RegisterUserServiceServer(server,
		&FailoverServer{code: codes.Unavailable})
	l, err := net.Listen("tcp", addr)
	require.NoError(t, err)
	server.Serve(l)
}

func (s *LoadBalancerTestSuite) startWeightedServer(addr string, weight int) {
	t := s.T()
	em, err := endpoints.NewManager(s.cli,
		"service/user")
	assert.NoError(t, err)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	// 要以 /service/user 为前缀
	addr = "127.0.0.1" + addr
	key := "service/user/" + addr
	// 5s
	var ttl int64 = 5
	leaseResp, err := s.cli.Grant(ctx, ttl)
	// metadata 一般用在客户端
	err = em.AddEndpoint(ctx, key,
		endpoints.Endpoint{
			Addr: addr,
			Metadata: map[string]any{
				"weight": weight,
			},
		}, clientv3.WithLease(leaseResp.ID))
	assert.NoError(t, err)

	// 忽略掉 ctx，因为在测试环境下，我们不需要手动控制退出续约
	kaCtx, _ := context.WithCancel(context.Background())
	go func() {
		_, err1 := s.cli.KeepAlive(kaCtx, leaseResp.ID)
		require.NoError(t, err1)
		//for resp := range ch {
		//	t.Log(resp.String())
		//}
	}()

	server := grpc.NewServer()
	RegisterUserServiceServer(server, &Server{name: addr})
	l, err := net.Listen("tcp", addr)
	require.NoError(t, err)
	server.Serve(l)
}

func TestLoadBalancer(t *testing.T) {
	suite.Run(t, new(LoadBalancerTestSuite))
}
