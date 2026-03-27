/*
 * Licensed to the Apache Software Foundation (ASF) under one or more
 * contributor license agreements.  See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The ASF licenses this file to You under the Apache License, Version 2.0
 * (the "License"); you may not use this file except in compliance with
 * the License.  You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package hessian

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
)

var (
	javaCommandOnce sync.Once
	javaCommandPath string
	javaCommandArgs []string
	javaCommandErr  error
)

func newJavaCommand(args ...string) (*exec.Cmd, error) {
	javaCommandOnce.Do(func() {
		javaCommandPath, javaCommandErr = findJavaBinary()
		if javaCommandErr != nil {
			return
		}
		javaCommandArgs, javaCommandErr = javaCompatibilityArgs(javaCommandPath)
	})
	if javaCommandErr != nil {
		return nil, javaCommandErr
	}

	cmdArgs := append(append([]string{}, javaCommandArgs...), args...)
	return exec.Command(javaCommandPath, cmdArgs...), nil
}

func findJavaBinary() (string, error) {
	if path, err := exec.LookPath("java"); err == nil {
		if _, err := javaMajorVersion(path); err == nil {
			return path, nil
		}
	}

	candidates := []string{
		filepath.Join(os.Getenv("JAVA_HOME"), "bin", "java"),
		"/opt/homebrew/opt/openjdk/bin/java",
		"/usr/local/opt/openjdk/bin/java",
	}
	for _, candidate := range candidates {
		if candidate == "" {
			continue
		}
		info, err := os.Stat(candidate)
		if err == nil && !info.IsDir() {
			if _, err := javaMajorVersion(candidate); err == nil {
				return candidate, nil
			}
		}
	}

	return "", fmt.Errorf("java executable not found in PATH or common JDK locations")
}

func javaCompatibilityArgs(javaPath string) ([]string, error) {
	major, err := javaMajorVersion(javaPath)
	if err != nil {
		return nil, err
	}
	if major < 9 {
		return nil, nil
	}

	return []string{
		"--add-opens=java.base/java.lang=ALL-UNNAMED",
		"--add-opens=java.base/java.lang.annotation=ALL-UNNAMED",
		"--add-opens=java.base/java.lang.instrument=ALL-UNNAMED",
		"--add-opens=java.base/java.lang.invoke=ALL-UNNAMED",
		"--add-opens=java.base/java.lang.reflect=ALL-UNNAMED",
		"--add-opens=java.base/java.io=ALL-UNNAMED",
		"--add-opens=java.base/java.math=ALL-UNNAMED",
		"--add-opens=java.base/java.time=ALL-UNNAMED",
		"--add-opens=java.base/java.time.format=ALL-UNNAMED",
		"--add-opens=java.base/java.time.temporal=ALL-UNNAMED",
		"--add-opens=java.base/java.util=ALL-UNNAMED",
		"--add-opens=java.base/java.util.concurrent=ALL-UNNAMED",
		"--add-opens=java.base/java.util.jar=ALL-UNNAMED",
		"--add-opens=java.prefs/java.util.prefs=ALL-UNNAMED",
		"--add-opens=java.base/java.util.zip=ALL-UNNAMED",
		"--add-opens=java.base/java.time.zone=ALL-UNNAMED",
	}, nil
}

func javaMajorVersion(javaPath string) (int, error) {
	output, err := exec.Command(javaPath, "-version").CombinedOutput()
	if err != nil {
		return 0, fmt.Errorf("inspect java version: %w", err)
	}

	version, err := parseJavaVersion(string(output))
	if err != nil {
		return 0, err
	}

	if strings.HasPrefix(version, "1.") {
		version = strings.TrimPrefix(version, "1.")
	}
	majorPart := version
	if idx := strings.IndexByte(version, '.'); idx >= 0 {
		majorPart = version[:idx]
	}

	major, err := strconv.Atoi(majorPart)
	if err != nil {
		return 0, fmt.Errorf("parse java major version %q: %w", version, err)
	}
	return major, nil
}

func parseJavaVersion(output string) (string, error) {
	firstQuote := strings.IndexByte(output, '"')
	if firstQuote < 0 {
		return "", fmt.Errorf("parse java version from output: %q", output)
	}
	rest := output[firstQuote+1:]
	secondQuote := strings.IndexByte(rest, '"')
	if secondQuote < 0 {
		return "", fmt.Errorf("parse java version from output: %q", output)
	}

	return rest[:secondQuote], nil
}
