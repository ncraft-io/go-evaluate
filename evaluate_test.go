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

package evaluate

import (
	"io/ioutil"
	"regexp"
	"testing"

	"github.com/alecthomas/assert"
)

func TestValuer_ValueFile(t *testing.T) {
	err := New().Evaluate("./testdata/sample.go")
	assert.NoError(t, err)

	content, err := ioutil.ReadFile("./testdata/sample.go")
	assert.NoError(t, err)

	assert.True(t, regexp.MustCompile(`const BuildTime = ".+"`).Match(content))
}
