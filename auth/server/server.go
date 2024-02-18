package server

import (
	"context"
	"errors"
	"log"
	"net"
	"time"

	"github.com/04Akaps/go-util/auth/paseto"
	auth "github.com/04Akaps/go-util/auth/proto"
	"github.com/04Akaps/go-util/types"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type AuthGrpcServer struct {
	auth.AuthServiceServer
	pasetoMaker    paseto.PasetoInterface
	tokenVerifyMap map[string]*auth.AuthData
}

func NewGrpcServer(url, pasetoKey string) error {
	if lis, err := net.Listen("tcp", url); err != nil {
		return err
	} else {
		s := grpc.NewServer([]grpc.ServerOption{}...)

		auth.RegisterAuthServiceServer(s, &AuthGrpcServer{
			pasetoMaker:    paseto.NewPasetoMaker(pasetoKey),
			tokenVerifyMap: make(map[string]*auth.AuthData),
		})
		reflection.Register(s)

		go func() {
			log.Printf("Success To Create GRPC Server URL : %s", url)
			if err = s.Serve(lis); err != nil {
				panic(err)
			}
		}()

		return nil
	}
}

func (a *AuthGrpcServer) CreateAuth(_ context.Context, req *auth.CreateNewPasetoTokenRequest) (*auth.CreateNewPasetoTokenResponse, error) {
	data := req.Auth
	token := data.PasetoToken

	a.tokenVerifyMap[token] = data
	return &auth.CreateNewPasetoTokenResponse{
		Auth: data,
	}, nil
}

func (a *AuthGrpcServer) VerifyAuth(_ context.Context, req *auth.VerifyPasetoTokenRequest) (*auth.VerifyPasetoTokenResponse, error) {
	token := req.PasetoToken
	res := &auth.VerifyPasetoTokenResponse{Status: &auth.VerifyResponse{Auth: nil}}

	if authData, ok := a.tokenVerifyMap[token]; !ok {
		// 데이터가 없는 경우
		res.Status.Status = []auth.ResponseType{auth.ResponseType_FAILED}
		return res, errors.New(types.VerifyTokenErrMap[types.FailedVerify])
	} else if authData.ExpireDate < time.Now().Unix() {
		// 만약 만료가 되었다면, Map을 비워주고, response를 수정해서 내려 준다.
		delete(a.tokenVerifyMap, token)
		res.Status.Status = []auth.ResponseType{auth.ResponseType_EXPIRED_DATE}
		return res, errors.New(types.VerifyTokenErrMap[types.Expired])
	} else {
		res.Status.Status = []auth.ResponseType{auth.ResponseType_SUCCESS}
		return res, nil
	}
}
