package jalapeno

import (
	"time"

	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"
)

func newConnection(url string) *grpc.ClientConn {
	kacp := keepalive.ClientParameters{
		Time:    10 * time.Second, // send ping every 10 seconds if there is no activity
		Timeout: time.Second,      // wait 1 second for ping back
	}

	log.Debugln("Dialing grpc")
	connection, err := grpc.Dial(url,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithKeepaliveParams(kacp))

	if err != nil {
		log.Fatalf("Failed to setup connection: %s", err)
	}

	return connection
}
