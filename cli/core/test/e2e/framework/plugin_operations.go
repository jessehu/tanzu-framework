// Copyright 2021 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package framework

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	plugintypes "github.com/vmware-tanzu/tanzu-framework/apis/cli/v1alpha1"
)

// PluginOps helps to perform plugin operations like create plugin bundle, and publish
type PluginOps interface {

	// GenerateScriptBasedPluginBinaries generates script based plugin binaries for given plugin metadata
	GenerateScriptBasedPluginBinaries(pluginsMeta []*PluginMeta) ([]string, []error)

	// GeneratePluginBundle generates plugin bundle in local file system for given plugin metadata
	GeneratePluginBundle(pluginsMeta []*PluginMeta) ([]string, []error)

	// PublishPluginBundle publishes the plugin bundle to given registry bucket
	PublishPluginBundle(pluginsInfo []*PluginMeta) (discoveryUrls []string, errs []error)
}

// CLIPlugin for plugin overlay info
type CLIPlugin struct {
	metav1.TypeMeta `yaml:",inline"`
	Metadata        metav1.ObjectMeta         `json:"metadata"`
	Spec            plugintypes.CLIPluginSpec `json:"spec"`
}

// localOCIPluginOps is the implementation of PluginOps interface
type localOCIPluginOps struct {
	cmdExe    CmdOps
	registry  PluginRegistry
	imgpkgOps ImgpkgOps
}

func NewLocalOCIPluginOps(registry PluginRegistry) PluginOps {
	return &localOCIPluginOps{
		cmdExe:    NewCmdOps(),
		registry:  registry,
		imgpkgOps: NewImgpkgOps(),
	}
}

func (po *localOCIPluginOps) GenerateScriptBasedPluginBinaries(pluginsMeta []*PluginMeta) ([]string, []error) {
	pluginsProcessed := make(map[string]bool)
	size := len(pluginsMeta)
	pluginBinaryFilePaths := make([]string, size)
	errs := make([]error, size)
	for i, pm := range pluginsMeta {
		if _, exists := pluginsProcessed[pm.name]; exists {
			errs[i] = errors.New("plugin name already exists, currently multiple versions of same plugin not supported")
			continue
		}
		pluginsProcessed[pm.name] = true

		// Set plugin local dir path if not, to generate binary image and bundle to publish to registry
		if pm.pluginLocalPath == "" {
			pm.pluginLocalPath = filepath.Join(TestPluginsDirPath, pm.name)
		}

		pluginBinaryFilePath, err := po.generatePluginBinary(pm)
		if err != nil {
			errs[i] = err
			continue
		}
		pm.pluginBinaryFilePath = pluginBinaryFilePath
		pluginBinaryFilePaths[i] = pm.pluginBinaryFilePath
	}
	return pluginBinaryFilePaths, errs
}

func (po *localOCIPluginOps) GeneratePluginBundle(pluginsMeta []*PluginMeta) ([]string, []error) {
	pluginsProcessed := make(map[string]bool)
	size := len(pluginsMeta)
	pluginBundlePath := make([]string, size)
	errs := make([]error, size)
	for i, pm := range pluginsMeta {
		if _, exists := pluginsProcessed[pm.name]; exists {
			errs[i] = errors.New("plugin name already exists, currently multiple versions of same plugin not supported")
			continue
		}
		pluginsProcessed[pm.name] = true

		// Set registry discovery url if not set already
		if pm.registryDiscoveryURL == "" {
			pm.registryDiscoveryURL = filepath.Join(po.registry.GetRegistryURLWithDefaultCLIPluginsBucket(), ("/" + pm.name + "/"))
		}

		imageRegistryURL, err := po.imgpkgOps.PushBinary(pm.pluginBinaryFilePath, pm.registryDiscoveryURL)
		if err != nil {
			errs[i] = err
			continue
		}
		pm.binaryDistributionURL = imageRegistryURL

		pluginOverlayObj, err := po.generatePluginDiscoveryOverlay(pm)
		if err != nil {
			errs[i] = err
			continue
		}

		pluginBundlePathLocal, err := po.createLocalPluginBundle(pm, pluginOverlayObj)
		if err != nil {
			errs[i] = err
			continue
		}

		pluginBundlePath[i] = pluginBundlePathLocal
	}
	return pluginBundlePath, errs
}

func (po *localOCIPluginOps) PublishPluginBundle(pluginsMeta []*PluginMeta) ([]string, []error) {
	size := len(pluginsMeta)
	discoveryUrls := make([]string, size)
	errs := make([]error, size)
	for i, pm := range pluginsMeta {
		// Set registry discovery url if not set already
		if pm.registryDiscoveryURL == "" {
			pm.registryDiscoveryURL = filepath.Join(po.registry.GetRegistryURLWithDefaultCLIPluginsBucket(), ("/" + pm.name + "/"))
		}

		_, err := po.imgpkgOps.PushBundle(pm.pluginLocalPath, pm.registryDiscoveryURL)
		if err != nil {
			errs[i] = err
			continue
		}
		discoveryUrls[i] = pm.registryDiscoveryURL
	}

	return discoveryUrls, errs
}

// createLocalPluginBundle creates plugin bundle in local file system for given plugin metadata and plugin overlay object
func (po *localOCIPluginOps) createLocalPluginBundle(pluginInfo *PluginMeta, pluginOverlayObj *CLIPlugin) (string, error) {
	dirPath := filepath.Join(pluginInfo.pluginLocalPath, ".imgpkg")
	if err := os.MkdirAll(dirPath, os.ModePerm); err != nil {
		return pluginInfo.pluginLocalPath, err
	}

	imagesFile := filepath.Join(dirPath, "images.yml")
	f, err := os.OpenFile(imagesFile, os.O_CREATE|os.O_WRONLY, 0755)
	if err != nil {
		return pluginInfo.pluginLocalPath, err
	}
	defer f.Close()
	fmt.Fprintf(f, ImagesTemplate)

	configDirPath := filepath.Join(pluginInfo.pluginLocalPath, "config")
	if err := os.MkdirAll(configDirPath, os.ModePerm); err != nil {
		return pluginInfo.pluginLocalPath, err
	}

	generatedValuesFile := filepath.Join(configDirPath, "zz_generated_values.yaml")
	gf, err := os.OpenFile(generatedValuesFile, os.O_CREATE|os.O_WRONLY, 0755)
	if err != nil {
		return pluginInfo.pluginLocalPath, err
	}
	defer gf.Close()
	fmt.Fprintf(gf, GeneratedValuesTemplate)

	overlayDirPath := filepath.Join(configDirPath, "overlay")
	if err := os.MkdirAll(overlayDirPath, os.ModePerm); err != nil {
		return pluginInfo.pluginLocalPath, err
	}
	overlayFile := filepath.Join(configDirPath, (pluginInfo.name + ".yaml"))
	of, err := os.OpenFile(overlayFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0755)
	if err != nil {
		return pluginInfo.pluginLocalPath, err
	}
	defer of.Close()
	yamlData, err := yaml.Marshal(&pluginOverlayObj)
	if err != nil {
		return pluginInfo.pluginLocalPath, err
	}
	err = os.WriteFile(overlayFile, yamlData, 0644)
	if err != nil {
		return pluginInfo.pluginLocalPath, err
	}
	return pluginInfo.pluginLocalPath, nil
}

// generatePluginBinary creates the script based plugin binary file for given plugin metadata, saves in local file system
func (po *localOCIPluginOps) generatePluginBinary(pluginInfo *PluginMeta) (string, error) {
	pluginInfo.pluginBinaryFileName = pluginInfo.name + "-" + pluginInfo.os + "-" + pluginInfo.version
	if err := CreateDir(pluginInfo.pluginLocalPath); err != nil {
		return "", err
	}
	filePath := filepath.Join(pluginInfo.pluginLocalPath, (pluginInfo.pluginBinaryFileName))
	f, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY, 0755)
	if err != nil {
		return "", err
	}
	defer f.Close()

	fmt.Fprintf(f, ScriptBasedPluginTemplate, pluginInfo.name, pluginInfo.description, pluginInfo.version,
		pluginInfo.sha, pluginInfo.group, strconv.FormatBool(pluginInfo.hidden), pluginInfo.aliases, pluginInfo.version, pluginInfo.name, pluginInfo.name)
	return filePath, nil
}

// CreateDir creates directory if not exists
func CreateDir(dir string) error {
	err := os.MkdirAll(dir, 0750)
	if err != nil && !os.IsExist(err) {
		return err
	}
	return nil
}

// generatePluginDiscoveryOverlay creates plugin overly object for given plugin metadata
func (po *localOCIPluginOps) generatePluginDiscoveryOverlay(pluginInfo *PluginMeta) (plugin *CLIPlugin, err error) {
	plugin = &CLIPlugin{}
	plugin.TypeMeta.Kind = "CLIPlugin"
	plugin.TypeMeta.APIVersion = "cli.tanzu.vmware.com/v1alpha1"
	plugin.Metadata.Name = pluginInfo.name

	var artifactsMap = make(map[string]plugintypes.ArtifactList)
	artifacts := make([]plugintypes.Artifact, 1)
	artifacts[0].OS = pluginInfo.os
	artifacts[0].Image = pluginInfo.binaryDistributionURL
	artifacts[0].Arch = pluginInfo.arch
	artifacts[0].Type = pluginInfo.discoveryType

	artifactsMap[pluginInfo.version] = artifacts
	plugin.Spec.Artifacts = artifactsMap
	plugin.Spec.Description = pluginInfo.description
	plugin.Spec.Optional = pluginInfo.optional
	plugin.Spec.RecommendedVersion = pluginInfo.version

	return plugin, nil
}
