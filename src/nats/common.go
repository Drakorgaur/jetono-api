package nats

import (
	"crypto/tls"
	"fmt"
	"github.com/Drakorgaur/jetono-api/src/storage"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nkeys"
	"os"
	"strconv"
	"strings"
	"time"
)

var (
	DefaultConfigInt = 0
)

func boolize(v string) bool {
	return v == "true" || v == "1"
}

func intize(v string) int {
	if i, err := strconv.Atoi(v); err != nil {
		return DefaultConfigInt
	} else {
		return i
	}
}

func timeize(v string, def int) time.Duration {
	if duration := time.Duration(intize(v)); duration != 0 {
		return duration
	} else {
		return time.Duration(def)
	}
}

type Configurator interface {
	TLSConfig() *tls.Config

	CustomReconnectDelayCB() nats.ReconnectDelayHandler

	ClosedCB() nats.ConnHandler

	DisconnectedErrCB() nats.ConnErrHandler

	ReconnectedCB() nats.ConnHandler

	DiscoveredServersCB() nats.ConnHandler

	AsyncErrorCB() nats.ErrHandler

	UserJWT() nats.UserJWTHandler

	SignatureCB() nats.SignatureHandler

	TokenHandler() nats.AuthTokenHandler

	CustomDialer() nats.CustomDialer

	LameDuckModeHandler() nats.ConnHandler
}

type DefaultConfigurator struct{}

func (c *DefaultConfigurator) TLSConfig() *tls.Config {
	return nil
}

func (c *DefaultConfigurator) CustomReconnectDelayCB() nats.ReconnectDelayHandler {
	return nil
}

func (c *DefaultConfigurator) ClosedCB() nats.ConnHandler {
	return nil
}

func (c *DefaultConfigurator) DisconnectedErrCB() nats.ConnErrHandler {
	return nil
}

func (c *DefaultConfigurator) ReconnectedCB() nats.ConnHandler {
	return nil
}

func (c *DefaultConfigurator) DiscoveredServersCB() nats.ConnHandler {
	return nil
}

func (c *DefaultConfigurator) AsyncErrorCB() nats.ErrHandler {
	return nil
}

func (c *DefaultConfigurator) UserJWT() nats.UserJWTHandler {
	return nil
}

func (c *DefaultConfigurator) SignatureCB() nats.SignatureHandler {
	return nil
}

func (c *DefaultConfigurator) TokenHandler() nats.AuthTokenHandler {
	return nil
}

func (c *DefaultConfigurator) CustomDialer() nats.CustomDialer {
	return nil
}

func (c *DefaultConfigurator) LameDuckModeHandler() nats.ConnHandler {
	return nil
}

func configurationTemplate(c Configurator) func(*nats.Options) error {
	return func(opts *nats.Options) error {
		opts.Url = os.Getenv("NATS_URL")
		opts.Servers = strings.Split(os.Getenv("NATS_URL"), ",")
		opts.NoRandomize = boolize(os.Getenv("NATS_NO_RANDOMIZE"))
		opts.NoEcho = boolize(os.Getenv("NATS_NO_ECHO"))
		opts.Name = os.Getenv("NATS_NAME")
		opts.Verbose = boolize(os.Getenv("NATS_VERBOSE"))
		opts.Pedantic = boolize(os.Getenv("NATS_PEDANTIC"))
		opts.Secure = boolize(os.Getenv("NATS_SECURE"))
		opts.TLSConfig = c.TLSConfig()
		opts.AllowReconnect = boolize(os.Getenv("NATS_ALLOW_RECONNECT"))
		opts.MaxReconnect = intize(os.Getenv("NATS_MAX_RECONNECT"))
		opts.ReconnectWait = timeize(os.Getenv("NATS_RECONNECT_WAIT"), 0)
		opts.CustomReconnectDelayCB = c.CustomReconnectDelayCB()
		opts.ReconnectJitter = timeize(os.Getenv("NATS_RECONNECT_JITTER"), 0)
		opts.ReconnectJitterTLS = timeize(os.Getenv("NATS_RECONNECT_JITTER_TLS"), 0)
		opts.Timeout = timeize(os.Getenv("NATS_TIMEOUT"), 0)
		opts.DrainTimeout = timeize(os.Getenv("NATS_DRAIN_TIMEOUT"), 0)
		opts.FlusherTimeout = timeize(os.Getenv("NATS_FLUSHER_TIMEOUT"), 0)
		opts.PingInterval = timeize(os.Getenv("NATS_PING_INTERVAL"), 0)
		opts.MaxPingsOut = intize(os.Getenv("NATS_MAX_PINGS_OUT"))
		opts.ClosedCB = c.ClosedCB()
		opts.DisconnectedErrCB = c.DisconnectedErrCB()
		opts.ReconnectedCB = c.ReconnectedCB()
		opts.DiscoveredServersCB = c.DiscoveredServersCB()
		opts.AsyncErrorCB = c.AsyncErrorCB()
		opts.ReconnectBufSize = intize(os.Getenv("NATS_RECONNECT_BUF_SIZE"))
		opts.SubChanLen = intize(os.Getenv("NATS_SUB_CHAN_LEN"))
		opts.UserJWT = c.UserJWT()
		opts.Nkey = os.Getenv("NATS_NKEY")
		opts.SignatureCB = c.SignatureCB()
		opts.User = os.Getenv("NATS_USER")
		opts.Password = os.Getenv("NATS_PASSWORD")
		opts.Token = os.Getenv("NATS_TOKEN")
		opts.TokenHandler = c.TokenHandler()
		opts.CustomDialer = c.CustomDialer()
		opts.UseOldRequestStyle = boolize(os.Getenv("NATS_USE_OLD_REQUEST_STYLE"))
		opts.NoCallbacksAfterClientClose = boolize(os.Getenv("NATS_NO_CALLBACKS_AFTER_CLIENT_CLOSE"))
		opts.LameDuckModeHandler = c.LameDuckModeHandler()
		opts.RetryOnFailedConnect = boolize(os.Getenv("NATS_RETRY_ON_FAILED_CONNECT"))
		opts.Compression = boolize(os.Getenv("NATS_COMPRESSION"))
		opts.ProxyPath = os.Getenv("NATS_PROXY_PATH")
		opts.InboxPrefix = os.Getenv("NATS_INBOX_PREFIX")
		return nil
	}
}

type UserNatsConn struct {
	*storage.AccountServerMap
	*nats.Conn
	conf  Configurator
	creds []byte
}

func (u *UserNatsConn) SetCreds(creds []byte) {
	u.creds = creds
}

func (u *UserNatsConn) UserCredentials() nats.Option {
	userCB := func() (string, error) {
		return nkeys.ParseDecoratedJWT(u.creds)
	}
	sigCB := func(nonce []byte) ([]byte, error) {
		kp, err := nkeys.ParseDecoratedNKey(u.creds)
		if err != nil {
			return nil, fmt.Errorf("unable to extract key pair from file %v", err)
		}
		defer kp.Wipe()

		sig, _ := kp.Sign(nonce)
		return sig, nil
	}
	return nats.UserJWT(userCB, sigCB)
}

func (u *UserNatsConn) authTemplate(opts []nats.Option) {
	opts = append(opts, u.UserCredentials())
}

func (u *UserNatsConn) GetNats() (*nats.Conn, error) {
	if u.conf == nil {
		u.conf = &DefaultConfigurator{}
	}
	options := []nats.Option{
		configurationTemplate(u.conf),
	}
	u.authTemplate(options)
	if nc, err := nats.Connect(u.ServersList, options...); err != nil {
		return nil, err
	} else {
		return nc, nil
	}
}
