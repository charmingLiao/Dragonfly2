/*
 *     Copyright 2020 The Dragonfly Authors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package client

import (
	"context"
	"fmt"

	grpc_retry "github.com/grpc-ecosystem/go-grpc-middleware/retry"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"

	"d7y.io/dragonfly/v2/internal/dfnet"
	"d7y.io/dragonfly/v2/pkg/rpc"
	"d7y.io/dragonfly/v2/pkg/rpc/base"
	"d7y.io/dragonfly/v2/pkg/rpc/cdnsystem"
)

func GetClientByAddrs(addrs []dfnet.NetAddr, opts ...grpc.DialOption) (CDNClient, error) {
	if len(addrs) == 0 {
		return nil, errors.New("address list of cdn is empty")
	}

	r := rpc.NewD7yResolver("cdn", addrs)

	dialOpts := append(append(
		rpc.DefaultClientOpts,
		grpc.WithDefaultServiceConfig(fmt.Sprintf(`{"loadBalancingPolicy": "%s"}`, rpc.D7yBalancerPolicy)),
		grpc.WithResolvers(r)),
		opts...)

	// target is "cdnsystem.Seeder" is the cdnsystem._Seeder_serviceDesc.ServiceName
	clientConn, err := grpc.Dial(
		fmt.Sprintf("%s:///%s", "cdn", "cdnsystem.Seeder"),
		dialOpts...)

	if err != nil {
		return nil, err
	}

	return &cdnClient{
		cc:           clientConn,
		seederClient: cdnsystem.NewSeederClient(clientConn),
		resolver:     r,
	}, nil
}

type CDNClient interface {
	ObtainSeeds(ctx context.Context, req *cdnsystem.SeedRequest, opts ...grpc.CallOption) (cdnsystem.Seeder_ObtainSeedsClient, error)

	GetPieceTasks(ctx context.Context, addr dfnet.NetAddr, req *base.PieceTaskRequest, opts ...grpc.CallOption) (*base.PiecePacket, error)

	UpdateAddresses(addrs []dfnet.NetAddr)

	Close() error
}

type cdnClient struct {
	cc           *grpc.ClientConn
	seederClient cdnsystem.SeederClient
	resolver     *rpc.D7yResolver
}

var _ CDNClient = (*cdnClient)(nil)

func (cc *cdnClient) ObtainSeeds(ctx context.Context, req *cdnsystem.SeedRequest, opts ...grpc.CallOption) (cdnsystem.Seeder_ObtainSeedsClient, error) {
	opts = append([]grpc.CallOption{
		grpc_retry.WithCodes(codes.ResourceExhausted, codes.Aborted, codes.Unavailable, codes.Unknown, codes.Internal),
	}, opts...)
	ctx = rpc.NewContext(ctx, &rpc.PickRequest{
		HashKey: req.TaskId,
	})
	return cc.seederClient.ObtainSeeds(ctx, req, opts...)
}

func (cc *cdnClient) GetPieceTasks(ctx context.Context, addr dfnet.NetAddr, req *base.PieceTaskRequest, opts ...grpc.CallOption) (*base.PiecePacket, error) {
	ctx = rpc.NewContext(ctx, &rpc.PickRequest{
		TargetAddr: addr.String(),
	})
	opts = append([]grpc.CallOption{
		grpc_retry.WithCodes(codes.ResourceExhausted, codes.Aborted, codes.Unavailable, codes.Unknown, codes.Internal),
	}, opts...)
	return cc.seederClient.GetPieceTasks(ctx, req, opts...)
}

func (cc *cdnClient) UpdateAddresses(addrs []dfnet.NetAddr) {
	cc.resolver.UpdateAddresses(addrs)
}

func (cc *cdnClient) Close() error {
	return cc.cc.Close()
}

func getClientByAddr(ctx context.Context, addr dfnet.NetAddr, opts ...grpc.DialOption) (*grpc.ClientConn, error) {
	dialOpts := append(append(rpc.DefaultClientOpts, grpc.WithDisableServiceConfig()), opts...)
	return grpc.DialContext(ctx, addr.GetEndpoint(), dialOpts...)
}

func (cc *cdnClient) getCdnClientByAddr(addr dfnet.NetAddr) (cdnsystem.SeederClient, error) {
	conn, err := getClientByAddr(context.Background(), addr)
	if err != nil {
		return nil, err
	}
	return cdnsystem.NewSeederClient(conn), nil
}
