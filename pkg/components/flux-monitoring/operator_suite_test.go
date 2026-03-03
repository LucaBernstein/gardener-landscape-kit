// SPDX-FileCopyrightText: SAP SE or an SAP affiliate company and Gardener contributors
//
// SPDX-License-Identifier: Apache-2.0

package fluxmonitoring_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestFluxMonitoring(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Components Flux Monitoring Suite")
}
