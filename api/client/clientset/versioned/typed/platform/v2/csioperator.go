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
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	rest "k8s.io/client-go/rest"
	scheme "tkestack.io/tke/api/client/clientset/versioned/scheme"
	v2 "tkestack.io/tke/api/platform/v2"
)

// CSIOperatorsGetter has a method to return a CSIOperatorInterface.
// A group's client should implement this interface.
type CSIOperatorsGetter interface {
	CSIOperators() CSIOperatorInterface
}

// CSIOperatorInterface has methods to work with CSIOperator resources.
type CSIOperatorInterface interface {
	Create(ctx context.Context, cSIOperator *v2.CSIOperator, opts v1.CreateOptions) (*v2.CSIOperator, error)
	Update(ctx context.Context, cSIOperator *v2.CSIOperator, opts v1.UpdateOptions) (*v2.CSIOperator, error)
	UpdateStatus(ctx context.Context, cSIOperator *v2.CSIOperator, opts v1.UpdateOptions) (*v2.CSIOperator, error)
	Delete(ctx context.Context, name string, opts v1.DeleteOptions) error
	Get(ctx context.Context, name string, opts v1.GetOptions) (*v2.CSIOperator, error)
	List(ctx context.Context, opts v1.ListOptions) (*v2.CSIOperatorList, error)
	Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error)
	Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v2.CSIOperator, err error)
	CSIOperatorExpansion
}

// cSIOperators implements CSIOperatorInterface
type cSIOperators struct {
	client rest.Interface
}

// newCSIOperators returns a CSIOperators
func newCSIOperators(c *PlatformV2Client) *cSIOperators {
	return &cSIOperators{
		client: c.RESTClient(),
	}
}

// Get takes name of the cSIOperator, and returns the corresponding cSIOperator object, and an error if there is any.
func (c *cSIOperators) Get(ctx context.Context, name string, options v1.GetOptions) (result *v2.CSIOperator, err error) {
	result = &v2.CSIOperator{}
	err = c.client.Get().
		Resource("csioperators").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do(ctx).
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of CSIOperators that match those selectors.
func (c *cSIOperators) List(ctx context.Context, opts v1.ListOptions) (result *v2.CSIOperatorList, err error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	result = &v2.CSIOperatorList{}
	err = c.client.Get().
		Resource("csioperators").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Do(ctx).
		Into(result)
	return
}

// Watch returns a watch.Interface that watches the requested cSIOperators.
func (c *cSIOperators) Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	opts.Watch = true
	return c.client.Get().
		Resource("csioperators").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Watch(ctx)
}

// Create takes the representation of a cSIOperator and creates it.  Returns the server's representation of the cSIOperator, and an error, if there is any.
func (c *cSIOperators) Create(ctx context.Context, cSIOperator *v2.CSIOperator, opts v1.CreateOptions) (result *v2.CSIOperator, err error) {
	result = &v2.CSIOperator{}
	err = c.client.Post().
		Resource("csioperators").
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(cSIOperator).
		Do(ctx).
		Into(result)
	return
}

// Update takes the representation of a cSIOperator and updates it. Returns the server's representation of the cSIOperator, and an error, if there is any.
func (c *cSIOperators) Update(ctx context.Context, cSIOperator *v2.CSIOperator, opts v1.UpdateOptions) (result *v2.CSIOperator, err error) {
	result = &v2.CSIOperator{}
	err = c.client.Put().
		Resource("csioperators").
		Name(cSIOperator.Name).
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(cSIOperator).
		Do(ctx).
		Into(result)
	return
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().
func (c *cSIOperators) UpdateStatus(ctx context.Context, cSIOperator *v2.CSIOperator, opts v1.UpdateOptions) (result *v2.CSIOperator, err error) {
	result = &v2.CSIOperator{}
	err = c.client.Put().
		Resource("csioperators").
		Name(cSIOperator.Name).
		SubResource("status").
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(cSIOperator).
		Do(ctx).
		Into(result)
	return
}

// Delete takes name of the cSIOperator and deletes it. Returns an error if one occurs.
func (c *cSIOperators) Delete(ctx context.Context, name string, opts v1.DeleteOptions) error {
	return c.client.Delete().
		Resource("csioperators").
		Name(name).
		Body(&opts).
		Do(ctx).
		Error()
}

// Patch applies the patch and returns the patched cSIOperator.
func (c *cSIOperators) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v2.CSIOperator, err error) {
	result = &v2.CSIOperator{}
	err = c.client.Patch(pt).
		Resource("csioperators").
		Name(name).
		SubResource(subresources...).
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(data).
		Do(ctx).
		Into(result)
	return
}
