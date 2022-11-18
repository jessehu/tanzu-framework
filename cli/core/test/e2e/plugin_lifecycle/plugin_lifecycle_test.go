// plugin_e2e_tests provides plugin command specific E2E test cases
package plugin_e2e_tests

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/vmware-tanzu/tanzu-framework/cli/core/test/e2e/framework"
)

var _ = framework.CLICoreDescribe("[Tests:E2E][Feature:Plugin-lifecycle]", func() {
	var (
		tf *framework.Framework
		// pluginOps framework.PluginOps
	)
	BeforeEach(func() {
		tf = framework.NewFramework()
		// pluginOps = framework.NewLocalOCIPluginOps(framework.NewLocalOCIRegistry(framework.DefaultRegistryName, framework.DefaultRegistryPort))
	})
	Context("oci local registry - plugin life cycle basic operations", func() {
		When("Plugin's metadata has given to generate and publish plugins", func() {
			It("should publish plugin's successfully without errors", func() {
				err := tf.ConfigInit()
				Expect(err).To(BeNil())
				// TODO: need install the docker, imgpkg tool in github runner
				/*
					var plugins [1]*framework.PluginMeta
					pm := framework.NewPluginMeta()
					pm.SetName("dummy").SetVersion("1.0").SetDescription("its dummy plugin").SetSHA("345").SetGroup("admin").SetArch("amd64").SetOS("darwin").SetDiscoveryType("oci").SetHidden(false).SetOptional(false)
					plugins[0] = pm

					pluginOps.GenerateScriptBasedPluginBinaries(plugins[:])
					pluginOps.GeneratePluginBundle(plugins[:])
					_, errs := pluginOps.PublishPluginBundle(plugins[:])
					for _, err := range errs {
						Expect(err).To(BeNil())
					}
				*/
			})
		})
	})
})
