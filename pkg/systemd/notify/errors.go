package notify

import "errors"

var (
	ErrNoNotifySocket = errors.New("$NOTIFY_SOCKET was not set")
)
