/*
 * Tencent is pleased to support the open source community by making TKEStack
 * available.
 *
 * Copyright (C) 2012-2019 Tencent. All Rights Reserved.
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

package storage

import (
	"context"

	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"tkestack.io/tke/api/application"
	v1 "tkestack.io/tke/api/application/v1"
	applicationinternalclient "tkestack.io/tke/api/client/clientset/internalversion/typed/application/internalversion"
	platformversionedclient "tkestack.io/tke/api/client/clientset/versioned/typed/platform/v2"
	helmaction "tkestack.io/tke/pkg/application/helm/action"
	"tkestack.io/tke/pkg/application/util"
)

// HistoryREST adapts a service registry into apiserver's RESTStorage model.
type HistoryREST struct {
	application       ApplicationStorage
	applicationClient *applicationinternalclient.ApplicationClient
	platformClient    platformversionedclient.PlatformV2Interface
}

// NewHistoryREST returns a wrapper around the underlying generic storage and performs
// allocations and deallocations of various helm releases related histories.
// TODO: all transactional behavior should be supported from within generic storage
//
//	or the strategy.
func NewHistoryREST(
	application ApplicationStorage,
	applicationClient *applicationinternalclient.ApplicationClient,
	platformClient platformversionedclient.PlatformV2Interface,
) *HistoryREST {
	rest := &HistoryREST{
		application:       application,
		applicationClient: applicationClient,
		platformClient:    platformClient,
	}
	return rest
}

// New returns an empty object that can be used with Create and Update after request data has been put into it.
func (rs *HistoryREST) New() runtime.Object {
	return rs.application.New()
}

// Get retrieves the object from the storage. It is required to support Patch.
func (rs *HistoryREST) Get(ctx context.Context, name string, options *metav1.GetOptions) (runtime.Object, error) {
	obj, err := rs.application.Get(ctx, name, &metav1.GetOptions{})
	if err != nil {
		return nil, errors.NewInternalError(err)
	}
	app := obj.(*application.App)
	appv1 := &v1.App{}
	if err := v1.Convert_application_App_To_v1_App(app, appv1, nil); err != nil {
		return nil, err
	}
	client, err := util.NewHelmClientWithProvider(ctx, rs.platformClient, appv1)
	if err != nil {
		return nil, errors.NewInternalError(err)
	}
	history, err := client.History(&helmaction.HistoryOptions{
		Namespace:   app.Spec.TargetNamespace,
		ReleaseName: app.Spec.Name,
	})
	if err != nil {
		return nil, errors.NewInternalError(err)
	}
	appHistory := &application.AppHistory{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: app.Namespace,
			Name:      app.Name,
		},
		Spec: application.AppHistorySpec{
			Type:            app.Spec.Type,
			TenantID:        app.Spec.TenantID,
			Name:            app.Spec.Name,
			TargetCluster:   app.Spec.TargetCluster,
			TargetNamespace: app.Spec.TargetNamespace,
			Histories:       make([]application.History, len(history)),
		},
	}
	appHistory.Spec.Histories = make([]application.History, len(history))
	for k, h := range history {
		appHistory.Spec.Histories[k] = application.History{
			Revision:    int64(h.Revision),
			Updated:     metav1.NewTime(h.Updated.Time),
			Status:      h.Status,
			Chart:       h.Chart,
			AppVersion:  h.AppVersion,
			Description: h.Description,
			Manifest:    h.Manifest,
		}
	}
	return appHistory, nil
}
