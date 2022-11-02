/*
 * Tencent is pleased to support the open source community by making TKEStack
 * available.
 *
 * Copyright (C) 2012-2021 Tencent. All Rights Reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the "License"); you may not use
 * this file except in compliance with the License. You may obtain a copy of the
 * License at
 *
 * https://opensource.org/licenses/Apache-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
 * WARRANTIES OF ANY KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations under the License.
 */

package credential

import (
	"context"
	platformv2 "tkestack.io/tke/api/platform/v2"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	platforminternalclient "tkestack.io/tke/api/client/clientset/internalversion/typed/platform/internalversion"
	platformversionedclient "tkestack.io/tke/api/client/clientset/versioned/typed/platform/v2"
	"tkestack.io/tke/api/platform"
)

// GetClusterCredential returns the cluster's credential
func GetClusterCredential(ctx context.Context, client platforminternalclient.PlatformInterface, cluster *platform.Cluster, username string) (*platform.ClusterCredential, error) {
	var (
		credential *platform.ClusterCredential
		err        error
	)

	if cluster.Spec.ClusterCredentialRef != nil {
		credential, err = client.ClusterCredentials().Get(ctx, cluster.Spec.ClusterCredentialRef.Name, metav1.GetOptions{})
		if err != nil {
			return nil, err
		}
	} else if client != nil {
		return nil, apierrors.NewNotFound(platform.Resource("ClusterCredential"), cluster.Name)
	}

	return credential, nil
}

// GetClusterCredentialV2 returns the versioned cluster's credential
func GetClusterCredentialV2(ctx context.Context, client platformversionedclient.PlatformV2Interface, cluster *platformv2.Cluster, username string) (*platformv2.ClusterCredential, error) {
	var (
		credential *platformv2.ClusterCredential
		err        error
	)

	if cluster.Spec.ClusterCredentialRef != nil {
		credential, err = client.ClusterCredentials().Get(ctx, cluster.Spec.ClusterCredentialRef.Name, metav1.GetOptions{})
		if err != nil && !apierrors.IsNotFound(err) {
			return nil, err
		}
	} else if client != nil {
		return nil, apierrors.NewNotFound(platform.Resource("ClusterCredential"), cluster.Name)
	}

	return credential, nil
}
