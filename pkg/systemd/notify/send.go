package notify

import (
	"fmt"
	"net"
	"os"
)

func Send(message string) error {
	socketAddr := &net.UnixAddr{
		Name: os.Getenv("NOTIFY_SOCKET"),
		Net:  "unixgram",
	}

	if socketAddr.Name == "" {
		return ErrNoNotifySocket
	}

	conn, err := net.DialUnix(socketAddr.Net, nil, socketAddr)
	if err != nil {
		return fmt.Errorf("connect to unix socket: %w", err)
	}
	defer conn.Close()

	_, err = conn.Write([]byte(message))
	if err != nil {
		return fmt.Errorf("write to unix socket: %w", err)
	}

	return nil
}
