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

// Code generated by lister-gen. DO NOT EDIT.

package v2

import (
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/tools/cache"
	v2 "tkestack.io/tke/api/platform/v2"
)

// PersistentEventLister helps list PersistentEvents.
// All objects returned here must be treated as read-only.
type PersistentEventLister interface {
	// List lists all PersistentEvents in the indexer.
	// Objects returned here must be treated as read-only.
	List(selector labels.Selector) (ret []*v2.PersistentEvent, err error)
	// Get retrieves the PersistentEvent from the index for a given name.
	// Objects returned here must be treated as read-only.
	Get(name string) (*v2.PersistentEvent, error)
	PersistentEventListerExpansion
}

// persistentEventLister implements the PersistentEventLister interface.
type persistentEventLister struct {
	indexer cache.Indexer
}

// NewPersistentEventLister returns a new PersistentEventLister.
func NewPersistentEventLister(indexer cache.Indexer) PersistentEventLister {
	return &persistentEventLister{indexer: indexer}
}

// List lists all PersistentEvents in the indexer.
func (s *persistentEventLister) List(selector labels.Selector) (ret []*v2.PersistentEvent, err error) {
	err = cache.ListAll(s.indexer, selector, func(m interface{}) {
		ret = append(ret, m.(*v2.PersistentEvent))
	})
	return ret, err
}

// Get retrieves the PersistentEvent from the index for a given name.
func (s *persistentEventLister) Get(name string) (*v2.PersistentEvent, error) {
	obj, exists, err := s.indexer.GetByKey(name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.NewNotFound(v2.Resource("persistentevent"), name)
	}
	return obj.(*v2.PersistentEvent), nil
}
