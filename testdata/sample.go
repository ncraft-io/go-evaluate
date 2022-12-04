// Copyright 2022 ncraft.io
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

package testdata

//go:generate go run github.com/ncraft-io/go-evaluate/cmd/evaluate

//go:evaluate date "+%Y-%m-%d %H:%M:%S %Z"
const BuildTime = "2022-12-03 17:22:06 CST"

//go:evaluate git rev-list -1 HEAD
const GitHash = ""

//go:evaluate git branch --show-current
const GitBranch = ""
