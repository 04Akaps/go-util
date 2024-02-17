package auth

import (
	"os/exec"

	"github.com/04Akaps/go-util/auth/client"
)

type Auth struct {
	*client.AuthGrpcClient
}

func NewAuth(serverUrl, clientUrl, pasetoKey string) (*Auth, error) {
	a := new(Auth)

	cmd := exec.Command("protoc",
		"--go_out=.",
		"--go_opt=paths=source_relative",
		"--go-grpc_out=.",
		"--go-grpc_opt=paths=source_relative",
		"proto/auth.proto")

	if _, err := cmd.CombinedOutput(); err != nil {
		panic(err)
	} else if a.AuthGrpcClient, err = client.NewGrpcClient(clientUrl, pasetoKey); err != nil {
		panic(err)
	} else {
		return a, nil
	}
}
