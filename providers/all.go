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

package providers

import (
	_ "github.com/oslet/erda-infra/providers/cassandra"             //
	_ "github.com/oslet/erda-infra/providers/clickhouse"            //
	_ "github.com/oslet/erda-infra/providers/elasticsearch"         //
	_ "github.com/oslet/erda-infra/providers/etcd"                  //
	_ "github.com/oslet/erda-infra/providers/etcd-election"         //
	_ "github.com/oslet/erda-infra/providers/etcd-mutex"            //
	_ "github.com/oslet/erda-infra/providers/expvar"                //
	_ "github.com/oslet/erda-infra/providers/grpcclient"            //
	_ "github.com/oslet/erda-infra/providers/grpcserver"            //
	_ "github.com/oslet/erda-infra/providers/health"                //
	_ "github.com/oslet/erda-infra/providers/httpserver"            //
	_ "github.com/oslet/erda-infra/providers/i18n"                  //
	_ "github.com/oslet/erda-infra/providers/kafka"                 //
	_ "github.com/oslet/erda-infra/providers/kafkav2"               //
	_ "github.com/oslet/erda-infra/providers/kubernetes"            //
	_ "github.com/oslet/erda-infra/providers/mysql"                 //
	_ "github.com/oslet/erda-infra/providers/pprof"                 //
	_ "github.com/oslet/erda-infra/providers/prometheus"            //
	_ "github.com/oslet/erda-infra/providers/redis"                 //
	_ "github.com/oslet/erda-infra/providers/remote-forward/client" //
	_ "github.com/oslet/erda-infra/providers/remote-forward/server" //
	_ "github.com/oslet/erda-infra/providers/serviceregister"       //
	// _ "github.com/oslet/erda-infra/providers/zk-master-election"   //
	// _ "github.com/oslet/erda-infra/providers/zookeeper"            //
	// _ "github.com/oslet/erda-infra/providers/legacy/httpendpoints" //
)
