// Copyright 2020 PingCAP, Inc. Licensed under Apache-2.0.

package gluetikv

import (
	"context"

	pd "github.com/tikv/pd/client"
	"github.com/twotigers93/tidb/br/pkg/glue"
	"github.com/twotigers93/tidb/br/pkg/summary"
	"github.com/twotigers93/tidb/br/pkg/utils"
	"github.com/twotigers93/tidb/br/pkg/version/build"
	"github.com/twotigers93/tidb/config"
	"github.com/twotigers93/tidb/domain"
	"github.com/twotigers93/tidb/kv"
	"github.com/twotigers93/tidb/store/driver"
)

// Asserting Glue implements glue.ConsoleGlue and glue.Glue at compile time.
var (
	_ glue.ConsoleGlue = Glue{}
	_ glue.Glue        = Glue{}
)

// Glue is an implementation of glue.Glue that accesses only TiKV without TiDB.
type Glue struct {
	glue.StdIOGlue
}

// GetDomain implements glue.Glue.
func (Glue) GetDomain(store kv.Storage) (*domain.Domain, error) {
	return nil, nil
}

// CreateSession implements glue.Glue.
func (Glue) CreateSession(store kv.Storage) (glue.Session, error) {
	return nil, nil
}

// Open implements glue.Glue.
func (Glue) Open(path string, option pd.SecurityOption) (kv.Storage, error) {
	if option.CAPath != "" {
		conf := config.GetGlobalConfig()
		conf.Security.ClusterSSLCA = option.CAPath
		conf.Security.ClusterSSLCert = option.CertPath
		conf.Security.ClusterSSLKey = option.KeyPath
		config.StoreGlobalConfig(conf)
	}
	return driver.TiKVDriver{}.Open(path)
}

// OwnsStorage implements glue.Glue.
func (Glue) OwnsStorage() bool {
	return true
}

// StartProgress implements glue.Glue.
func (Glue) StartProgress(ctx context.Context, cmdName string, total int64, redirectLog bool) glue.Progress {
	return utils.StartProgress(ctx, cmdName, total, redirectLog, nil)
}

// Record implements glue.Glue.
func (Glue) Record(name string, val uint64) {
	summary.CollectSuccessUnit(name, 1, val)
}

// GetVersion implements glue.Glue.
func (Glue) GetVersion() string {
	return "BR\n" + build.Info()
}

// UseOneShotSession implements glue.Glue.
func (g Glue) UseOneShotSession(store kv.Storage, closeDomain bool, fn func(glue.Session) error) error {
	return nil
}
