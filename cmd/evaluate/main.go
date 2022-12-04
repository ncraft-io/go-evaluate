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

package main

import (
	"flag"
	"log"
	"os"

	evaluate "github.com/ncraft-io/go-evaluate"
)

// update the value for value declaration which has //go:evaluate directive
func main() {
	filename := ""

	flag.Parse()

	args := flag.Args()
	if len(args) == 0 {
		filename = os.Getenv("GOFILE")
	} else if len(args) == 1 {
		filename = flag.Arg(0)
	}

	if len(filename) == 0 {
		log.Fatalln("has no input filename, should using argument or $GOFILE environment to specify the filename")
	}

	if err := evaluate.New().Evaluate(filename); err != nil {
		log.Fatalln(err)
	}
}
