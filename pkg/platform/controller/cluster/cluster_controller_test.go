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

package cluster

import (
	"testing"
	"time"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	platformv2 "tkestack.io/tke/api/platform/v2"
)

func newClusterForTest(resourcesVersion string, spec *platformv2.ClusterSpec, phase platformv2.ClusterPhase, conditions []platformv2.ClusterCondition) *platformv2.Cluster {
	mc := &platformv2.Cluster{
		ObjectMeta: v1.ObjectMeta{ResourceVersion: resourcesVersion},
		Spec: platformv2.ClusterSpec{
			TenantID:    "default",
			DisplayName: "global",
			Type:        "Baremetal",
			Version:     "1.21.1-tke.1",
		},
		Status: platformv2.ClusterStatus{
			Phase: platformv2.ClusterRunning,
			Conditions: []platformv2.ClusterCondition{
				{
					Type:          conditionTypeHealthCheck,
					Status:        platformv2.ConditionTrue,
					LastProbeTime: v1.Now(),
				},
			},
		},
	}
	if spec != nil {
		mc.Spec = *spec
	}
	if len(phase) != 0 {
		mc.Status.Phase = phase
	}
	if conditions != nil {
		mc.Status.Conditions = conditions
	}
	return mc
}

func TestController_needsUpdate(t *testing.T) {
	resyncInternal := time.Minute
	// type fields struct {
	// 	queue             workqueue.RateLimitingInterface
	// 	lister            platformv1lister.ClusterLister
	// 	listerSynced      cache.InformerSynced
	// 	log               log.Logger
	// 	platformClient    platformversionedclient.PlatformV1Interface
	// 	deleter           deletion.ClusterDeleterInterface
	// 	healthCheckPeriod time.Duration
	// }
	type args struct {
		old *platformv2.Cluster
		new *platformv2.Cluster
	}
	tests := []struct {
		name string
		// fields fields
		args args
		want bool
	}{
		{
			name: "change spec",
			args: args{
				old: newClusterForTest("old", &platformv2.ClusterSpec{Version: "old"}, platformv2.ClusterPhase(""), nil),
				new: newClusterForTest("new", &platformv2.ClusterSpec{Version: "new"}, platformv2.ClusterPhase(""), nil),
			},
			want: true,
		},
		{
			name: "Initializing to Running",
			args: args{
				old: newClusterForTest("old", nil, platformv2.ClusterInitializing, nil),
				new: newClusterForTest("new", nil, platformv2.ClusterRunning, nil),
			},
			want: true,
		},
		{
			name: "Initializing to Failed",
			args: args{
				old: newClusterForTest("old", nil, platformv2.ClusterInitializing, nil),
				new: newClusterForTest("new", nil, platformv2.ClusterFailed, nil),
			},
			want: true,
		},
		{
			name: "Running to Failed",
			args: args{
				old: newClusterForTest("old", nil, platformv2.ClusterRunning, nil),
				new: newClusterForTest("new", nil, platformv2.ClusterFailed, nil),
			},
			want: true,
		},
		{
			name: "Running to Terminating",
			args: args{
				old: newClusterForTest("old", nil, platformv2.ClusterRunning, nil),
				new: newClusterForTest("new", nil, platformv2.ClusterTerminating, nil),
			},
			want: true,
		},
		{
			name: "Failed to Terminating",
			args: args{
				old: newClusterForTest("old", nil, platformv2.ClusterFailed, nil),
				new: newClusterForTest("new", nil, platformv2.ClusterTerminating, nil),
			},
			want: true,
		},
		{
			name: "Failed to Running",
			args: args{
				old: newClusterForTest("old", nil, platformv2.ClusterFailed, nil),
				new: newClusterForTest("new", nil, platformv2.ClusterRunning, nil),
			},
			want: true,
		},
		{
			name: "Failed to Initializing",
			args: args{
				old: newClusterForTest("old", nil, platformv2.ClusterFailed, nil),
				new: newClusterForTest("new", nil, platformv2.ClusterInitializing, nil),
			},
			want: true,
		},
		{
			name: "Initializing last conditon unkonwn to false",
			args: args{
				old: newClusterForTest("old", nil, platformv2.ClusterInitializing, []platformv2.ClusterCondition{{Status: platformv2.ConditionUnknown}}),
				new: newClusterForTest("new", nil, platformv2.ClusterInitializing, []platformv2.ClusterCondition{{Status: platformv2.ConditionFalse}}),
			},
			want: false,
		},
		{
			name: "last conditon unkonwn to true",
			args: args{
				old: newClusterForTest("old", nil, platformv2.ClusterPhase(""), []platformv2.ClusterCondition{{Status: platformv2.ConditionUnknown}}),
				new: newClusterForTest("new", nil, platformv2.ClusterPhase(""), []platformv2.ClusterCondition{{Status: platformv2.ConditionFalse}}),
			},
			want: true,
		},
		{
			name: "Initializing last conditon false retrun true if resync",
			args: func() args {
				// resource version equal
				new := newClusterForTest("new", nil, platformv2.ClusterPhase(""), []platformv2.ClusterCondition{{Status: platformv2.ConditionFalse}})
				return args{new, new}
			}(),
			want: true,
		},
		{
			name: "last conditon true to unknown",
			args: args{
				old: newClusterForTest("old", nil, platformv2.ClusterPhase(""), []platformv2.ClusterCondition{{Status: platformv2.ConditionTrue}}),
				new: newClusterForTest("new", nil, platformv2.ClusterPhase(""), []platformv2.ClusterCondition{{Status: platformv2.ConditionUnknown}}),
			},
			want: true,
		},
		{
			name: "last conditon false to unknown",
			args: args{
				old: newClusterForTest("old", nil, platformv2.ClusterPhase(""), []platformv2.ClusterCondition{{Status: platformv2.ConditionFalse}}),
				new: newClusterForTest("new", nil, platformv2.ClusterPhase(""), []platformv2.ClusterCondition{{Status: platformv2.ConditionUnknown}}),
			},
			want: true,
		},
		{
			name: "health check is not long enough",
			args: func() args {
				new := newClusterForTest("old", nil, platformv2.ClusterPhase(""), []platformv2.ClusterCondition{{
					Type:          conditionTypeHealthCheck,
					Status:        platformv2.ConditionTrue,
					LastProbeTime: v1.NewTime(time.Now().Add(-resyncInternal / 2))}})
				return args{new, new}
			}(),
			want: false,
		},
		{
			name: "health check is long enough",
			args: func() args {
				new := newClusterForTest("old", nil, platformv2.ClusterPhase(""), []platformv2.ClusterCondition{{
					Type:          conditionTypeHealthCheck,
					Status:        platformv2.ConditionTrue,
					LastProbeTime: v1.NewTime(time.Now().Add(-resyncInternal - 1))}})
				return args{new, new}
			}(),
			want: true,
		},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Controller{
				// queue:             tt.fields.queue,
				// lister:            tt.fields.lister,
				// listerSynced:      tt.fields.listerSynced,
				// log:               tt.fields.log,
				// platformClient:    tt.fields.platformClient,
				// deleter:           tt.fields.deleter,
				healthCheckPeriod: resyncInternal,
			}
			if got := c.needsUpdate(tt.args.old, tt.args.new); got != tt.want {
				t.Errorf("Controller.needsUpdate() = %v, want %v", got, tt.want)
			}
		})
	}
}
