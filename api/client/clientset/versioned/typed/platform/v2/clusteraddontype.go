/*
 * Tencent is pleased to support the open source community by making TKEStack
 * available.
 *
 * Copyright (C) 2012-2020 Tencent. All Rights Reserved.
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

// Code generated by client-gen. DO NOT EDIT.

package v2

import (
	"context"
	"time"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	rest "k8s.io/client-go/rest"
	scheme "tkestack.io/tke/api/client/clientset/versioned/scheme"
	v2 "tkestack.io/tke/api/platform/v2"
)

// ClusterAddonTypesGetter has a method to return a ClusterAddonTypeInterface.
// A group's client should implement this interface.
type ClusterAddonTypesGetter interface {
	ClusterAddonTypes() ClusterAddonTypeInterface
}

// ClusterAddonTypeInterface has methods to work with ClusterAddonType resources.
type ClusterAddonTypeInterface interface {
	List(ctx context.Context, opts v1.ListOptions) (*v2.ClusterAddonTypeList, error)
	ClusterAddonTypeExpansion
}

// clusterAddonTypes implements ClusterAddonTypeInterface
type clusterAddonTypes struct {
	client rest.Interface
}

// newClusterAddonTypes returns a ClusterAddonTypes
func newClusterAddonTypes(c *PlatformV2Client) *clusterAddonTypes {
	return &clusterAddonTypes{
		client: c.RESTClient(),
	}
}

// List takes label and field selectors, and returns the list of ClusterAddonTypes that match those selectors.
func (c *clusterAddonTypes) List(ctx context.Context, opts v1.ListOptions) (result *v2.ClusterAddonTypeList, err error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	result = &v2.ClusterAddonTypeList{}
	err = c.client.Get().
		Resource("clusteraddontypes").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Do(ctx).
		Into(result)
	return
}
