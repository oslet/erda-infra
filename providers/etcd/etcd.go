// Copyright (c) 2021 Terminus, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package etcd

import (
	"crypto/tls"
	"crypto/x509"
	"os"
	"reflect"
	"strings"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
	"google.golang.org/grpc"

	"github.com/oslet/erda-infra/base/logs"
	"github.com/oslet/erda-infra/base/servicehub"
)

// Interface .
type Interface interface {
	Connect() (*clientv3.Client, error)
	Client() *clientv3.Client
	Timeout() time.Duration
}

type config struct {
	Endpoints             string        `file:"endpoints" env:"ETCD_ENDPOINTS"`
	Timeout               time.Duration `file:"timeout" default:"5s"`
	AutoSyncInterval      time.Duration `file:"auto_sync_interval" env:"ETCD_AUTO_SYNC_INTERVAL"`
	MaxCallSendMsgSize    int           `file:"max_call_send_msg_size" env:"ETCD_MAX_CALL_SEND_MSG_SIZE"`
	MaxCallRecvMsgSize    int           `file:"max_call_recv_msg_size" env:"ETCD_MAX_CALL_RECV_MSG_SIZE"`
	RejectOldCluster      bool          `file:"reject_old_cluster" env:"ETCD_REJECT_OLD_CLUSTER"`
	PermitWithoutStream   bool          `file:"permit_without_stream" env:"ETCD_PERMIT_WITHOUT_STREAM"`
	BackoffWaitBetween    time.Duration `file:"backoff_wait_between" env:"ETCD_BACKOFF_WAIT_BETWEEN"`
	BackoffJitterFraction float64       `file:"backoff_jitter_fraction" env:"ETCD_BACKOFF_JITTER_FRACTION"`
	TLS                   struct {
		CertFile    string `file:"cert_file"`
		CertKeyFile string `file:"cert_key_file"`
		CaFile      string `file:"ca_file"`
	} `file:"tls"`
	SyncConnect bool   `file:"sync_connect" default:"true"`
	Username    string `file:"username"`
	Password    string `file:"password"`
}

var clientType = reflect.TypeOf((*clientv3.Client)(nil))

type provider struct {
	Cfg       *config
	Log       logs.Logger
	client    *clientv3.Client
	tlsConfig *tls.Config
}

func (p *provider) Init(ctx servicehub.Context) error {
	err := p.initTLSConfig()
	if err != nil {
		return err
	}
	client, err := p.Connect()
	if err != nil {
		return err
	}
	p.client = client
	return nil
}

func (p *provider) Connect() (*clientv3.Client, error) {
	config := clientv3.Config{
		Endpoints:   strings.Split(p.Cfg.Endpoints, ","),
		DialTimeout: p.Cfg.Timeout,
		TLS:         p.tlsConfig,
		Username:    p.Cfg.Username,
		Password:    p.Cfg.Password,
	}
	if p.Cfg.SyncConnect {
		config.DialOptions = append(config.DialOptions, grpc.WithBlock())
	}
	if p.Cfg.AutoSyncInterval > 0 {
		config.AutoSyncInterval = p.Cfg.AutoSyncInterval
	}
	if p.Cfg.MaxCallSendMsgSize > 0 {
		config.MaxCallSendMsgSize = p.Cfg.MaxCallSendMsgSize
	}
	if p.Cfg.MaxCallRecvMsgSize > 0 {
		config.MaxCallRecvMsgSize = p.Cfg.MaxCallRecvMsgSize
	}
	if p.Cfg.RejectOldCluster {
		config.RejectOldCluster = p.Cfg.RejectOldCluster
	}
	if p.Cfg.PermitWithoutStream {
		config.PermitWithoutStream = p.Cfg.PermitWithoutStream
	}
	if p.Cfg.BackoffWaitBetween > 0 {
		config.BackoffWaitBetween = p.Cfg.BackoffWaitBetween
	}
	if p.Cfg.BackoffJitterFraction > 0 {
		config.BackoffJitterFraction = p.Cfg.BackoffJitterFraction
	}
	return clientv3.New(config)
}

func (p *provider) Client() *clientv3.Client { return p.client }

func (p *provider) Timeout() time.Duration { return p.Cfg.Timeout }

func (p *provider) initTLSConfig() error {
	if len(p.Cfg.TLS.CertFile) > 0 || len(p.Cfg.TLS.CertKeyFile) > 0 {
		cfg, err := readTLSConfig(p.Cfg.TLS.CertFile, p.Cfg.TLS.CertKeyFile, p.Cfg.TLS.CaFile)
		if err != nil {
			if os.IsNotExist(err) {
				p.Log.Warnf("fail to load tls files: %s", err)
				return nil
			}
			return err
		}
		p.tlsConfig = cfg
	}
	return nil
}

func readTLSConfig(certFile, certKeyFile, caFile string) (*tls.Config, error) {
	cert, err := tls.LoadX509KeyPair(certFile, certKeyFile)
	if err != nil {
		return nil, err
	}
	caData, err := os.ReadFile(caFile)
	if err != nil {
		return nil, err
	}
	pool := x509.NewCertPool()
	pool.AppendCertsFromPEM(caData)
	return &tls.Config{
		Certificates: []tls.Certificate{cert},
		RootCAs:      pool,
	}, nil
}

// Provide .
func (p *provider) Provide(ctx servicehub.DependencyContext, args ...interface{}) interface{} {
	if ctx.Type() == clientType || ctx.Service() == "etcd-client" {
		return p.client
	}
	return p
}

func init() {
	servicehub.Register("etcd", &servicehub.Spec{
		Services: []string{"etcd", "etcd-client"},
		Types: []reflect.Type{
			reflect.TypeOf((*Interface)(nil)).Elem(),
			clientType,
		},
		Description: "etcd",
		ConfigFunc:  func() interface{} { return &config{} },
		Creator: func() servicehub.Provider {
			return &provider{}
		},
	})
}
