// Copyright 2025 Baptiste Gouhoury <baptiste.gouhoury@scalend.fr>
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

package errors

import (
	"os"

	log "github.com/sirupsen/logrus"
)

// Check logs a fatal message and exits with ErrorGeneric if err is not nil
func Check(err error) {
	if err != nil {
		exitfunc := func() {
			os.Exit(1)
		}
		log.RegisterExitHandler(exitfunc)
		log.Fatal(err)
	}
}
