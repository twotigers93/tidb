// Copyright 2022 PingCAP, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package ddlhelper

import (
	"github.com/twotigers93/tidb/ddl"
	"github.com/twotigers93/tidb/parser/ast"
	"github.com/twotigers93/tidb/parser/model"
)

// BuildTableInfoFromAST builds model.TableInfo from a SQL statement.
// Note: TableID and PartitionID are left as uninitialized value.
func BuildTableInfoFromAST(s *ast.CreateTableStmt) (*model.TableInfo, error) {
	return ddl.BuildTableInfoFromAST(s)
}
