package callinfo

import (
	"fmt"
	"html/template"

	"github.com/youtube/vitess/go/rpcwrap/proto"
	"golang.org/x/net/context"
)

// RPCWrapCallInfo takes a context generated by rpcwrap, and
// returns one that has CallInfo filled in.
func RPCWrapCallInfo(ctx context.Context) context.Context {
	remoteAddr, _ := proto.RemoteAddr(ctx)
	username, _ := proto.Username(ctx)
	return NewContext(ctx, &rpcWrapCallInfoImpl{
		remoteAddr: remoteAddr,
		username:   username,
	})
}

type rpcWrapCallInfoImpl struct {
	remoteAddr, username string
}

func (rwci *rpcWrapCallInfoImpl) RemoteAddr() string {
	return rwci.remoteAddr
}

func (rwci *rpcWrapCallInfoImpl) Username() string {
	return rwci.username
}

func (rwci *rpcWrapCallInfoImpl) Text() string {
	return fmt.Sprintf("%s@%s", rwci.username, rwci.remoteAddr)
}

func (rwci *rpcWrapCallInfoImpl) HTML() template.HTML {
	return template.HTML("<b>RemoteAddr:</b> " + rwci.remoteAddr + "</br>\n" + "<b>Username:</b> " + rwci.username + "</br>\n")
}
