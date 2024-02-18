package auth

import (
	"os/exec"

	"github.com/04Akaps/go-util/auth/client"
	"github.com/04Akaps/go-util/auth/server"
)

type Auth struct {
	*client.AuthGrpcClient
}

func NewAuth(serverUrl, clientUrl, pasetoKey string, init bool) (*Auth, error) {
	a := new(Auth)
	var err error

	if init {
		cmd := exec.Command("protoc",
			"--go_out=.",
			"--go_opt=paths=source_relative",
			"--go-grpc_out=.",
			"--go-grpc_opt=paths=source_relative",
			"proto/auth.proto")

		if _, err := cmd.CombinedOutput(); err != nil {
			panic(err)
		}
	}

	if a.AuthGrpcClient, err = client.NewGrpcClient(clientUrl, pasetoKey); err != nil {
		panic(err)
	} else {
		server.NewGrpcServer(serverUrl, pasetoKey)
		return a, nil
	}
}
