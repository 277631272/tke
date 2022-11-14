/*
 * Tencent is pleased to support the open source community by making TKEStack
 * available.
 *
 * Copyright (C) 2012-2021 Tencent. All Rights Reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the “License”); you may not use
 * this file except in compliance with the License. You may obtain a copy of the
 * License at
 *
 * https://opensource.org/licenses/Apache-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an “AS IS” BASIS, WITHOUT
 * WARRANTIES OF ANY KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations under the License.
 */

package vendor

import (
	"strings"

	platformv2 "tkestack.io/tke/api/platform/v2"
)

// GetKubeVendor get k8s vendor from k8s version
// ref https://github.com/open-cluster-management/multicloud-operators-foundation/blob/e94b719de6d5f3541e948dd70ad8f1ff748aa452/pkg/klusterlet/clusterinfo/clusterinfo_controller.go#L326
func GetKubeVendor(version string) (kubeVendor platformv2.KubeVendorType) {
	version = strings.ToUpper(version)
	switch {
	case strings.Contains(version, string(platformv2.KubeVendorTKE)):
		kubeVendor = platformv2.KubeVendorTKE
		return
	case strings.Contains(version, string(platformv2.KubeVendorIKS)):
		kubeVendor = platformv2.KubeVendorIKS
		return
	case strings.Contains(version, string(platformv2.KubeVendorEKS)):
		kubeVendor = platformv2.KubeVendorEKS
		return
	case strings.Contains(version, string(platformv2.KubeVendorGKE)):
		kubeVendor = platformv2.KubeVendorGKE
		return
	case strings.Contains(version, string(platformv2.KubeVendorICP)):
		kubeVendor = platformv2.KubeVendorICP
	default:
		kubeVendor = platformv2.KubeVendorOther
	}
	return
}
