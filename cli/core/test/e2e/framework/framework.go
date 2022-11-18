// Copyright 2022 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

// Package test defines the integration and end-to-end test case for cli core
package framework

import (
	"os"
	"path/filepath"

	"github.com/onsi/ginkgo"
)

const (
	CliCore = "[CLI-Core]"

	TanzuInit    = "tanzu init"
	TanzuVersion = "tanzu version"

	ConfigSet          = "tanzu config set "
	ConfigGet          = "tanzu config get "
	ConfigUnset        = "tanzu config unset "
	ConfigInit         = "tanzu config init"
	ConfigServerList   = "tanzu config server list"
	ConfigServerDelete = "tanzu config server delete "

	KindCreateCluster = "kind create cluster --name "
	DockerInfo        = "docker info"
	StartDockerUbuntu = "sudo systemctl start docker"
	StopDockerUbuntu  = "sudo systemctl stop docker"

	TestDir        = ".tanzu-cli-e2e"
	TestPluginsDir = ".e2e-test-plugins"
)

var TestDirPath string
var TestPluginsDirPath string

// CLICoreDescribe annotates the test with the CLICore label.
func CLICoreDescribe(text string, body func()) bool {
	return ginkgo.Describe(CliCore+text, body)
}

// Framework has all helper functions to write CLI e2e test cases
type Framework struct {
	CliOps
	ConfigOps
	ClusterOps
}

func NewFramework() *Framework {
	return &Framework{
		CliOps:     NewCliOps(),
		ConfigOps:  NewConfOps(),
		ClusterOps: NewKindCluster(NewDocker()),
	}
}

func init() {
	homeDir, _ := os.UserHomeDir()
	TestDirPath = filepath.Join(homeDir, TestDir)
	os.Setenv("HOME", TestDirPath)
	TestPluginsDirPath = filepath.Join(TestDirPath, TestPluginsDir)
}

const ScriptBasedPluginTemplate = `#!/bin/bash

# Minimally Viable Dummy Tanzu CLI 'Plugin'

info() {
   cat << EOF
{
  "name": "%s",
  "description": "%s",
  "version": "%s",
  "buildSHA": "%s",
  "group": "%s",
  "hidden": %s,
  "aliases": %s
}
EOF
  exit 0
}
version() { 
	echo "%s"
}

case "$1" in
    info)  $1;;
    version)   $1;;
    *) cat << EOF
%s functionality

Usage:
  tanzu %s [command]

Available Commands:
  info     plugin details
  version   plugin version
EOF
       exit 0
       ;;
esac
`

const ImagesTemplate = `---
apiVersion: imgpkg.carvel.dev/v1alpha1
images:
kind: ImagesLock`

const GeneratedValuesTemplate = `#@data/values
#@overlay/match-child-defaults missing_ok=True

---`