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

package install

import (
	"fmt"
	"k8s.io/apimachinery/pkg/runtime"
	runtimeutil "k8s.io/apimachinery/pkg/util/runtime"
	"tkestack.io/tke/api/platform"
	"tkestack.io/tke/api/platform/v1"
	"tkestack.io/tke/api/platform/v2"
)

func init() {
	fmt.Println("init, tkestack register platform scheme")
	Install(platform.Scheme)
}

// Install registers the API group and adds types to a scheme
func Install(scheme *runtime.Scheme) {
	fmt.Println("install, tkestack register platform scheme")
	runtimeutil.Must(platform.AddToScheme(scheme))
	runtimeutil.Must(v1.AddToScheme(scheme))
	runtimeutil.Must(v2.AddToScheme(scheme))
	runtimeutil.Must(scheme.SetVersionPriority(v2.SchemeGroupVersion, v1.SchemeGroupVersion))
}
