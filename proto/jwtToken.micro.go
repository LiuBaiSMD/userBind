// Code generated by protoc-gen-micro. DO NOT EDIT.
// source: jwtToken.proto

package heartbeat

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	math "math"
)

import (
	context "context"
	client "github.com/micro/go-micro/client"
	server "github.com/micro/go-micro/server"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ client.Option
var _ server.Option

// Client API for TokenCreator service

type TokenCreatorService interface {
	GetToken(ctx context.Context, in *TokenRequest, opts ...client.CallOption) (*TokenResponse, error)
}

type tokenCreatorService struct {
	c    client.Client
	name string
}

func NewTokenCreatorService(name string, c client.Client) TokenCreatorService {
	if c == nil {
		c = client.NewClient()
	}
	if len(name) == 0 {
		name = "tokencreator"
	}
	return &tokenCreatorService{
		c:    c,
		name: name,
	}
}

func (c *tokenCreatorService) GetToken(ctx context.Context, in *TokenRequest, opts ...client.CallOption) (*TokenResponse, error) {
	req := c.c.NewRequest(c.name, "TokenCreator.GetToken", in)
	out := new(TokenResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for TokenCreator service

type TokenCreatorHandler interface {
	GetToken(context.Context, *TokenRequest, *TokenResponse) error
}

func RegisterTokenCreatorHandler(s server.Server, hdlr TokenCreatorHandler, opts ...server.HandlerOption) error {
	type tokenCreator interface {
		GetToken(ctx context.Context, in *TokenRequest, out *TokenResponse) error
	}
	type TokenCreator struct {
		tokenCreator
	}
	h := &tokenCreatorHandler{hdlr}
	return s.Handle(s.NewHandler(&TokenCreator{h}, opts...))
}

type tokenCreatorHandler struct {
	TokenCreatorHandler
}

func (h *tokenCreatorHandler) GetToken(ctx context.Context, in *TokenRequest, out *TokenResponse) error {
	return h.TokenCreatorHandler.GetToken(ctx, in, out)
}
