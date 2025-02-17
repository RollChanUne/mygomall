//user客户端的rpc

package rpc

import (
	"context"
	"fmt"
	"time"

	"github.com/Group-lifelong-youth-training/mygomall/pkg/HTTPviper"
	"github.com/Group-lifelong-youth-training/mygomall/pkg/middleware"
	etcd "github.com/Group-lifelong-youth-training/mygomall/pkg/registry-etcd"
	"github.com/Group-lifelong-youth-training/mygomall/rpc_gen/kitex_gen/user"
	"github.com/Group-lifelong-youth-training/mygomall/rpc_gen/kitex_gen/user/userservice"
	"github.com/cloudwego/kitex/client"
	"github.com/cloudwego/kitex/pkg/retry"
	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/kitex-contrib/obs-opentelemetry/provider"
	"github.com/kitex-contrib/obs-opentelemetry/tracing"
)

var userClient userservice.Client

// User RPC 客户端初始化
func initUserRpc(Config *HTTPviper.Config) {
	EtcdAddress := fmt.Sprintf("%s:%d", Config.Viper.GetString("Etcd.Address"), Config.Viper.GetInt("Etcd.Port"))
	// 服务发现
	r, err := etcd.NewEtcdResolver([]string{EtcdAddress})
	if err != nil {
		panic(err)
	}
	ServiceName := Config.Viper.GetString("Server.Name")

	p := provider.NewOpenTelemetryProvider(
		provider.WithServiceName(ServiceName),
		provider.WithExportEndpoint("localhost:4317"),
		provider.WithInsecure(),
	)
	defer p.Shutdown(context.Background())

	c, err := userservice.NewClient(
		ServiceName,
		client.WithMiddleware(middleware.CommonMiddleware),
		client.WithInstanceMW(middleware.ClientMiddleware),
		client.WithMuxConnection(1),                       // mux
		client.WithRPCTimeout(30*time.Second),             // rpc timeout
		client.WithConnectTimeout(30000*time.Millisecond), // conn timeout
		client.WithFailureRetry(retry.NewFailurePolicy()), // retry
		client.WithSuite(tracing.NewClientSuite()),        // tracer
		client.WithResolver(r),                            // resolver
		// Please keep the same as provider.WithServiceName
		client.WithClientBasicInfo(&rpcinfo.EndpointBasicInfo{ServiceName: ServiceName}),
	)
	if err != nil {
		panic(err)
	}
	userClient = c
}

// 传递 注册操作 的上下文, 并获取 RPC Server 端的响应.
func Register(ctx context.Context, req *user.RegisterReq) (resp *user.RegisterResp, err error) {
	resp, err = userClient.Register(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// 传递 登录操作 的上下文, 并获取 RPC Server 端的响应.
func Login(ctx context.Context, req *user.LoginReq) (resp *user.LoginResp, err error) {
	resp, err = userClient.Login(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
