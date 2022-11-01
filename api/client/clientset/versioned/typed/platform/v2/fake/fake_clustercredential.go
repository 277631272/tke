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

package fake

import (
	"context"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	labels "k8s.io/apimachinery/pkg/labels"
	schema "k8s.io/apimachinery/pkg/runtime/schema"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	testing "k8s.io/client-go/testing"
	v2 "tkestack.io/tke/api/platform/v2"
)

// FakeClusterCredentials implements ClusterCredentialInterface
type FakeClusterCredentials struct {
	Fake *FakePlatformV2
}

var clustercredentialsResource = schema.GroupVersionResource{Group: "platform.tkestack.io", Version: "v2", Resource: "clustercredentials"}

var clustercredentialsKind = schema.GroupVersionKind{Group: "platform.tkestack.io", Version: "v2", Kind: "ClusterCredential"}

// Get takes name of the clusterCredential, and returns the corresponding clusterCredential object, and an error if there is any.
func (c *FakeClusterCredentials) Get(ctx context.Context, name string, options v1.GetOptions) (result *v2.ClusterCredential, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootGetAction(clustercredentialsResource, name), &v2.ClusterCredential{})
	if obj == nil {
		return nil, err
	}
	return obj.(*v2.ClusterCredential), err
}

// List takes label and field selectors, and returns the list of ClusterCredentials that match those selectors.
func (c *FakeClusterCredentials) List(ctx context.Context, opts v1.ListOptions) (result *v2.ClusterCredentialList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootListAction(clustercredentialsResource, clustercredentialsKind, opts), &v2.ClusterCredentialList{})
	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &v2.ClusterCredentialList{ListMeta: obj.(*v2.ClusterCredentialList).ListMeta}
	for _, item := range obj.(*v2.ClusterCredentialList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested clusterCredentials.
func (c *FakeClusterCredentials) Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewRootWatchAction(clustercredentialsResource, opts))
}

// Create takes the representation of a clusterCredential and creates it.  Returns the server's representation of the clusterCredential, and an error, if there is any.
func (c *FakeClusterCredentials) Create(ctx context.Context, clusterCredential *v2.ClusterCredential, opts v1.CreateOptions) (result *v2.ClusterCredential, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootCreateAction(clustercredentialsResource, clusterCredential), &v2.ClusterCredential{})
	if obj == nil {
		return nil, err
	}
	return obj.(*v2.ClusterCredential), err
}

// Update takes the representation of a clusterCredential and updates it. Returns the server's representation of the clusterCredential, and an error, if there is any.
func (c *FakeClusterCredentials) Update(ctx context.Context, clusterCredential *v2.ClusterCredential, opts v1.UpdateOptions) (result *v2.ClusterCredential, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootUpdateAction(clustercredentialsResource, clusterCredential), &v2.ClusterCredential{})
	if obj == nil {
		return nil, err
	}
	return obj.(*v2.ClusterCredential), err
}

// Delete takes name of the clusterCredential and deletes it. Returns an error if one occurs.
func (c *FakeClusterCredentials) Delete(ctx context.Context, name string, opts v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewRootDeleteAction(clustercredentialsResource, name), &v2.ClusterCredential{})
	return err
}

// Patch applies the patch and returns the patched clusterCredential.
func (c *FakeClusterCredentials) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v2.ClusterCredential, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootPatchSubresourceAction(clustercredentialsResource, name, pt, data, subresources...), &v2.ClusterCredential{})
	if obj == nil {
		return nil, err
	}
	return obj.(*v2.ClusterCredential), err
}
