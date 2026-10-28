// SPDX-FileCopyrightText: Contributors to the Gardener project
//
// SPDX-License-Identifier: Apache-2.0

package registry

import (
	"fmt"
	"slices"
	"strings"

	"github.com/elliotchance/orderedmap/v3"
	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/util/sets"

	"github.com/gardener/gardener-landscape-kit/pkg/apis/config/v1alpha1"
	"github.com/gardener/gardener-landscape-kit/pkg/components"
	"github.com/gardener/gardener-landscape-kit/pkg/components/flux"
	calico "github.com/gardener/gardener-landscape-kit/pkg/components/gardener-extensions/networking-calico"
	cilium "github.com/gardener/gardener-landscape-kit/pkg/components/gardener-extensions/networking-cilium"
	gardenlinux "github.com/gardener/gardener-landscape-kit/pkg/components/gardener-extensions/os-gardenlinux"
	suse "github.com/gardener/gardener-landscape-kit/pkg/components/gardener-extensions/os-suse-chost"
	alicloud "github.com/gardener/gardener-landscape-kit/pkg/components/gardener-extensions/provider-alicloud"
	aws "github.com/gardener/gardener-landscape-kit/pkg/components/gardener-extensions/provider-aws"
	azure "github.com/gardener/gardener-landscape-kit/pkg/components/gardener-extensions/provider-azure"
	gcp "github.com/gardener/gardener-landscape-kit/pkg/components/gardener-extensions/provider-gcp"
	openstack "github.com/gardener/gardener-landscape-kit/pkg/components/gardener-extensions/provider-openstack"
	gvisor "github.com/gardener/gardener-landscape-kit/pkg/components/gardener-extensions/runtime-gvisor"
	certservice "github.com/gardener/gardener-landscape-kit/pkg/components/gardener-extensions/shoot-cert-service"
	dnsservice "github.com/gardener/gardener-landscape-kit/pkg/components/gardener-extensions/shoot-dns-service"
	networkingproblemdetector "github.com/gardener/gardener-landscape-kit/pkg/components/gardener-extensions/shoot-networking-problemdetector"
	oidcservice "github.com/gardener/gardener-landscape-kit/pkg/components/gardener-extensions/shoot-oidc-service"
	traefik "github.com/gardener/gardener-landscape-kit/pkg/components/gardener-extensions/shoot-traefik"
	"github.com/gardener/gardener-landscape-kit/pkg/components/gardener/garden"
	"github.com/gardener/gardener-landscape-kit/pkg/components/gardener/operator"
	virtualgardenaccess "github.com/gardener/gardener-landscape-kit/pkg/components/gardener/virtual-garden-access"
	githubcomponent "github.com/gardener/gardener-landscape-kit/pkg/components/github"
	gardenconfig "github.com/gardener/gardener-landscape-kit/pkg/components/virtual-garden/garden-config"
)

// ComponentList contains all available components.
var ComponentList = []func() (components.Interface, error){
	flux.NewComponent,
	githubcomponent.NewComponent,
	operator.NewComponent,
	garden.NewComponent,
	calico.NewComponent,
	cilium.NewComponent,
	alicloud.NewComponent,
	aws.NewComponent,
	azure.NewComponent,
	gcp.NewComponent,
	openstack.NewComponent,
	gardenlinux.NewComponent,
	suse.NewComponent,
	certservice.NewComponent,
	dnsservice.NewComponent,
	oidcservice.NewComponent,
	traefik.NewComponent,
	networkingproblemdetector.NewComponent,
	gvisor.NewComponent,
	virtualgardenaccess.NewComponent,
	gardenconfig.NewComponent,
}

// RegisterAllComponents registers all available components.
func RegisterAllComponents(log logr.Logger, registry Interface, config *v1alpha1.LandscapeKitConfiguration) error {
	orderedComponents := orderedmap.NewOrderedMap[string, components.Interface]()
	for _, newComponent := range ComponentList {
		component, err := newComponent()
		if err != nil {
			return fmt.Errorf("failed to create component: %w", err)
		}
		orderedComponents.Set(component.GetComponentMetadata().Name, component)
	}

	if err := selectComponents(config, orderedComponents); err != nil {
		return err
	}

	for _, component := range orderedComponents.AllFromFront() {
		if err := registry.RegisterComponent(log, component); err != nil {
			return fmt.Errorf("failed to register component: %w", err)
		}
	}

	return nil
}

func selectComponents(config *v1alpha1.LandscapeKitConfiguration, orderedComponents *orderedmap.OrderedMap[string, components.Interface]) error {
	// An omitted or empty `components` list enables no components: only explicitly listed components are registered.
	if config == nil || len(config.Components) == 0 {
		for _, name := range slices.Collect(orderedComponents.Keys()) {
			orderedComponents.Delete(name)
		}
		return nil
	}

	// Non-empty -> positive allow-list: validate names, then delete anything not listed.
	listedComponents := sets.New[string]()
	for _, comp := range config.Components {
		listedComponents.Insert(comp.Name)
	}

	availableComponents := slices.Collect(orderedComponents.Keys())
	invalidComponentNames := listedComponents.Difference(sets.New(availableComponents...))
	if len(invalidComponentNames) > 0 {
		return fmt.Errorf(
			"configuration contains invalid component names: %s - available component names are: %s",
			strings.Join(invalidComponentNames.UnsortedList(), ", "),
			strings.Join(availableComponents, ", "),
		)
	}

	for _, componentName := range availableComponents {
		if !listedComponents.Has(componentName) {
			orderedComponents.Delete(componentName)
		}
	}

	return nil
}
