// Copyright 2019 Google LLC
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
//

// +build go1.12

package runner29

import (
	"fmt"
	"log"
	"runtime/debug"
)

func printModuleVersion() {
	if bi, exists := debug.ReadBuildInfo(); exists {
		fmt.Println(bi.Main.Version)
	} else {
		log.Printf("No version information found. Make sure to use " +
			"GO111MODULE=on when running 'go get' in order to use specific " +
			"version of the binary.")
	}

}
