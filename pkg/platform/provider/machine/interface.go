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

package machine

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"runtime"
	"strings"
	"time"

	"tkestack.io/tke/pkg/util/apiclient"
	"tkestack.io/tke/pkg/util/log"

	"github.com/thoas/go-funk"
	"k8s.io/apimachinery/pkg/util/validation/field"
	"tkestack.io/tke/api/platform"

	platformv2 "tkestack.io/tke/api/platform/v2"
	typesv2 "tkestack.io/tke/pkg/platform/types/v2"
)

const (
	ReasonWaiting      = "Waiting"
	ReasonSkip         = "Skip"
	ReasonFailedInit   = "FailedInit"
	ReasonFailedUpdate = "FailedUpdate"
	ReasonFailedDelete = "FailedDelete"

	ConditionTypeDone = "EnsureDone"

	ConditionTypeHealthCheck = "HealthCheck"
	FailedHealthCheckReason  = "FailedHealthCheck"
)

type APIProvider interface {
	Validate(machine *platform.Machine) field.ErrorList
	ValidateUpdate(machine *platform.Machine, oldMachine *platform.Machine) field.ErrorList
}

// ControllerProvider ControllerProvider
type ControllerProvider interface {
	// NeedUpdate could be implemented by user to judge whether machine need update or not.
	NeedUpdate(old, new *platformv2.Machine) bool

	PreCreate(machine *platform.Machine) error
	AfterCreate(machine *platform.Machine) error

	OnCreate(ctx context.Context, machine *platformv2.Machine, cluster *typesv2.Cluster) error
	OnUpdate(ctx context.Context, machine *platformv2.Machine, cluster *typesv2.Cluster) error
	OnDelete(ctx context.Context, machine *platformv2.Machine, cluster *typesv2.Cluster) error
	// OnHealthCheck could be implemented by user, and default implementation is checking
	// tenant cluster node status by machine IP
	OnHealthCheck(ctx context.Context, machine *platformv2.Machine, cluster *typesv2.Cluster) *platformv2.Machine
}

// Provider defines a set of response interfaces for specific machine
// types in machine management.
type Provider interface {
	Name() string

	APIProvider
	ControllerProvider
}

var _ Provider = &DelegateProvider{}

type Handler func(context.Context, *platformv2.Machine, *typesv2.Cluster) error

func (h Handler) Name() string {
	name := runtime.FuncForPC(reflect.ValueOf(h).Pointer()).Name()
	i := strings.LastIndex(name, ".")
	if i == -1 {
		return "Unknown"
	}
	return strings.TrimSuffix(name[i+1:], "-fm")
}

type DelegateProvider struct {
	ProviderName string

	ValidateFunc       func(machine *platform.Machine) field.ErrorList
	ValidateUpdateFunc func(machine *platform.Machine, oldMachine *platform.Machine) field.ErrorList
	PreCreateFunc      func(machine *platform.Machine) error
	AfterCreateFunc    func(machine *platform.Machine) error

	CreateHandlers []Handler
	DeleteHandlers []Handler
	UpdateHandlers []Handler
}

func (p *DelegateProvider) Name() string {
	if p.ProviderName == "" {
		return "unknown"
	}
	return p.ProviderName
}

func (p *DelegateProvider) Validate(machine *platform.Machine) field.ErrorList {
	if p.ValidateFunc != nil {
		return p.ValidateFunc(machine)
	}

	return nil
}

func (p *DelegateProvider) ValidateUpdate(machine *platform.Machine, oldMachine *platform.Machine) field.ErrorList {
	if p.ValidateUpdateFunc != nil {
		return p.ValidateUpdateFunc(machine, oldMachine)
	}

	return nil
}

func (p *DelegateProvider) PreCreate(machine *platform.Machine) error {
	if p.PreCreateFunc != nil {
		return p.PreCreateFunc(machine)
	}

	return nil
}

func (p *DelegateProvider) AfterCreate(machine *platform.Machine) error {
	if p.AfterCreateFunc != nil {
		return p.AfterCreateFunc(machine)
	}

	return nil
}

func (p *DelegateProvider) OnCreate(ctx context.Context, machine *platformv2.Machine, cluster *typesv2.Cluster) error {
	condition, err := p.getCreateCurrentCondition(machine)
	if err != nil {
		return err
	}

	if cluster.Spec.Features.SkipConditions != nil &&
		funk.ContainsString(cluster.Spec.Features.SkipConditions, condition.Type) {
		machine.SetCondition(platformv2.MachineCondition{
			Type:    condition.Type,
			Status:  platformv2.ConditionTrue,
			Reason:  ReasonSkip,
			Message: "Skip current condition",
		})
	} else {
		handler := p.getCreateHandler(condition.Type)
		if handler == nil {
			return fmt.Errorf("can't get handler by %s", condition.Type)
		}
		ctx := log.FromContext(ctx).WithName("MachineProvider.OnCreate").WithName(handler.Name()).WithContext(ctx)
		log.FromContext(ctx).Info("Doing")
		startTime := time.Now()
		err = handler(ctx, machine, cluster)
		log.FromContext(ctx).Info("Done", "error", err, "cost", time.Since(startTime).String())
		if err != nil {
			machine.SetCondition(platformv2.MachineCondition{
				Type:    condition.Type,
				Status:  platformv2.ConditionFalse,
				Message: err.Error(),
				Reason:  ReasonFailedInit,
			})
			return err
		}

		machine.SetCondition(platformv2.MachineCondition{
			Type:   condition.Type,
			Status: platformv2.ConditionTrue,
		})
	}

	nextConditionType := p.getNextConditionType(condition.Type)
	if nextConditionType == ConditionTypeDone {
		machine.Status.Phase = platformv2.MachineRunning
	} else {
		machine.SetCondition(platformv2.MachineCondition{
			Type:    nextConditionType,
			Status:  platformv2.ConditionUnknown,
			Message: "waiting execute",
			Reason:  ReasonWaiting,
		})
	}

	return nil
}

func (p *DelegateProvider) OnUpdate(ctx context.Context, machine *platformv2.Machine, cluster *typesv2.Cluster) error {
	if machine.Status.Phase != platformv2.MachineUpgrading {
		return nil
	}
	for _, handler := range p.UpdateHandlers {
		ctx := log.FromContext(ctx).WithName("MachineProvider.OnUpdate").WithName(handler.Name()).WithContext(ctx)
		log.FromContext(ctx).Info("Doing")
		startTime := time.Now()
		err := handler(ctx, machine, cluster)
		log.FromContext(ctx).Info("Done", "error", err, "cost", time.Since(startTime).String())
		if err != nil {
			machine.Status.Reason = ReasonFailedUpdate
			machine.Status.Message = fmt.Sprintf("%s error: %v", handler.Name(), err)
			return err
		}
	}
	machine.Status.Reason = ""
	machine.Status.Message = ""

	return nil
}

func (p *DelegateProvider) OnDelete(ctx context.Context, machine *platformv2.Machine, cluster *typesv2.Cluster) error {
	for _, handler := range p.DeleteHandlers {
		ctx := log.FromContext(ctx).WithName("MachineProvider.OnDelete").WithName(handler.Name()).WithContext(ctx)
		log.FromContext(ctx).Info("Doing")
		startTime := time.Now()
		err := handler(ctx, machine, cluster)
		log.FromContext(ctx).Info("Done", "error", err, "cost", time.Since(startTime).String())
		if err != nil {
			cluster.Status.Reason = ReasonFailedDelete
			cluster.Status.Message = fmt.Sprintf("%s error: %v", handler.Name(), err)
			return err
		}
	}
	cluster.Status.Reason = ""
	cluster.Status.Message = ""

	return nil
}

func (p *DelegateProvider) OnHealthCheck(ctx context.Context, machine *platformv2.Machine, cluster *typesv2.Cluster) *platformv2.Machine {
	if !(machine.Status.Phase == platformv2.MachineRunning ||
		machine.Status.Phase == platformv2.MachineFailed) {
		return machine
	}

	healthCheckCondition := platformv2.MachineCondition{
		Type:   ConditionTypeHealthCheck,
		Status: platformv2.ConditionFalse,
	}

	clientset, err := cluster.Clientset()
	if err != nil {
		machine.Status.Phase = platformv2.MachineFailed

		healthCheckCondition.Reason = FailedHealthCheckReason
		healthCheckCondition.Message = err.Error()
	} else {
		_, err = apiclient.GetNodeByMachineIP(ctx, clientset, machine.Spec.IP)
		if err != nil {
			machine.Status.Phase = platformv2.MachineFailed

			healthCheckCondition.Reason = FailedHealthCheckReason
			healthCheckCondition.Message = err.Error()
		} else {
			machine.Status.Phase = platformv2.MachineRunning

			healthCheckCondition.Status = platformv2.ConditionTrue
		}
	}

	machine.SetCondition(healthCheckCondition)

	log.FromContext(ctx).Info("Update machine health status", "phase", machine.Status.Phase)

	return machine
}

func (p *DelegateProvider) NeedUpdate(old, new *platformv2.Machine) bool {
	return false
}

func (p *DelegateProvider) getNextConditionType(conditionType string) string {
	var (
		i       int
		handler Handler
	)
	for i, handler = range p.CreateHandlers {
		name := handler.Name()
		if name == conditionType {
			break
		}
	}
	if i == len(p.CreateHandlers)-1 {
		return ConditionTypeDone
	}
	next := p.CreateHandlers[i+1]

	return next.Name()
}

func (p *DelegateProvider) getCreateHandler(conditionType string) Handler {
	for _, f := range p.CreateHandlers {
		if conditionType == f.Name() {
			return f
		}
	}

	return nil
}

func (p *DelegateProvider) getCreateCurrentCondition(c *platformv2.Machine) (*platformv2.MachineCondition, error) {
	if c.Status.Phase == platformv2.MachineRunning {
		return nil, errors.New("machine phase is running now")
	}
	if len(p.CreateHandlers) == 0 {
		return nil, errors.New("no create handlers")
	}

	if len(c.Status.Conditions) == 0 {
		return &platformv2.MachineCondition{
			Type:    p.CreateHandlers[0].Name(),
			Status:  platformv2.ConditionUnknown,
			Message: "waiting process",
			Reason:  ReasonWaiting,
		}, nil
	}

	for _, condition := range c.Status.Conditions {
		if condition.Status == platformv2.ConditionFalse || condition.Status == platformv2.ConditionUnknown {
			return &condition, nil
		}
	}

	return nil, errors.New("no condition need process")
}
