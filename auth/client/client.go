package client

import (
	"context"
	"github.com/04Akaps/go-util/auth/paseto"
	auth "github.com/04Akaps/go-util/auth/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"time"
)

type AuthGrpcClient struct {
	gRPCClient auth.AuthServiceClient
	client     *grpc.ClientConn

	pasetoMaker paseto.PasetoInterface
}

func NewGrpcClient(url, key string) (*AuthGrpcClient, error) {
	c := new(AuthGrpcClient)

	if client, err := grpc.Dial(url, grpc.WithTransportCredentials(insecure.NewCredentials())); err != nil {
		return nil, err
	} else {

		c.client = client
		c.gRPCClient = auth.NewAuthServiceClient(client)
		c.pasetoMaker = paseto.NewPasetoMaker(key)

		log.Printf("Success To Create GRPC Client URL : %s", url)

		return c, nil
	}
}

func (a *AuthGrpcClient) CreateAuth(address string) (*auth.AuthData, error) {
	now := time.Now()
	expiredTime := now.Add(30 * time.Minute)

	d := &auth.AuthData{
		Address:     address,
		CreatedDate: now.Unix(),
		ExpireDate:  expiredTime.Unix(),
	}

	if token, err := a.pasetoMaker.CreateToken(d); err != nil {
		return nil, err
	} else {
		d.PasetoToken = token

		if res, err := a.gRPCClient.CreateAuth(context.Background(), &auth.CreateNewPasetoTokenRequest{Auth: d}); err != nil {
			return nil, err
		} else {
			return res.Auth, nil
		}
	}
}

func (a *AuthGrpcClient) VerifyAuth(pasetoToken string) (*auth.VerifyResponse, error) {
	if res, err := a.gRPCClient.VerifyAuth(context.Background(), &auth.VerifyPasetoTokenRequest{PasetoToken: pasetoToken}); err != nil {
		return nil, err
	} else {
		return res.Status, nil
	}
}
