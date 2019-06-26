/*
 * Copyright (c) 2019.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package main

import (
	"github.com/Flaque/filet"
	_ "github.com/Flaque/filet"
	"os"
	"testing"
)

func TestParseCommandLineDatabase(t *testing.T) {
	os.Args = []string{"command", "--host=hostname", "--port=1234", "--user=username", "--password=passwd", "--name=icmdb"}

	m := &Config{}
	m.ParseCommandLine()

	if m.timeout != 200 {
		t.Errorf("Timeout is not configured with default value. It is %d.", m.timeout)
	}
}

func TestParseCommandLineLockFile(t *testing.T) {
	os.Args = []string{"command", "--host=hostname", "--port=1234", "--user=username", "--password=passwd", "--name=icmdb", "--lockfile=test"}

	m := &Config{}
	m.ParseCommandLine()

	if m.lockfile != "test" {
		t.Errorf("Lockfile is not configured. It is %s.", m.lockfile)
	}
}

func TestFileExist(t *testing.T) {
	defer filet.CleanUp(t)

	// Creates a temporary file with string "some content"
	file := filet.TmpFile(t, "", "some content")

	m := &Config{}
	m.lockfile = file.Name()

	if !m.LockFileExists() {
		t.Errorf("Lockfile is there but not correct identified. File is %s", file.Name())
	}
}
