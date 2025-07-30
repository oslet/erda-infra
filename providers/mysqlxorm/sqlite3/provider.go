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

package sqlite3

import (
	"fmt"
	"reflect"

	"xorm.io/xorm"
	"xorm.io/xorm/names"

	"github.com/oslet/erda-infra/base/logs"
	"github.com/oslet/erda-infra/base/servicehub"
	"github.com/oslet/erda-infra/providers/mysqlxorm"
)

type config struct {
	DbSourceName string `file:"dbSourceName" env:"SQLITE3_DB_SOURCE_NAME" default:"test.sqlite3"`
	JournalMode  string `file:"journalMode" env:"SQLITE3_JOURNAL_MODE" default:""`
	RandomName   bool   `file:"randomName" env:"SQLITE3_RANDOM_NAME" default:"false"`
}

type provider struct {
	Cfg *config
	Log logs.Logger
	*Sqlite3
}

var _ servicehub.ProviderInitializer = (*provider)(nil)

// Init .
func (p *provider) Init(ctx servicehub.Context) error {
	var dbSourceName = p.Cfg.DbSourceName
	var err error
	if p.Cfg.RandomName {
		dbSourceName, err = randomName(p.Cfg.DbSourceName)
		if err != nil {
			return err
		}
	}

	engine, err := xorm.NewEngine("sqlite3", dbSourceName)
	if err != nil {
		return fmt.Errorf("failed to connect to sqlite3 engine,err : %s", err)
	}

	engine.SetMapper(names.GonicMapper{})

	if p.Cfg.JournalMode != "" {
		_, err = engine.Exec(fmt.Sprintf("PRAGMA journal_mode = %s", p.Cfg.JournalMode))
		if err != nil {
			return err
		}
	}

	p.Sqlite3 = &Sqlite3{db: engine}

	return nil
}

func init() {
	servicehub.Register("sqlite3-xorm", &servicehub.Spec{
		Services: []string{"sqlite3-xorm"},
		Types: []reflect.Type{
			reflect.TypeOf((*mysqlxorm.Interface)(nil)).Elem(),
		},
		Description: "sqlite3-xorm",
		ConfigFunc:  func() interface{} { return &config{} },
		Creator: func() servicehub.Provider {
			return &provider{}
		},
	})
}
