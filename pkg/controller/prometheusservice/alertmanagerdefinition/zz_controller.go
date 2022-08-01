/*
Copyright 2021 The Crossplane Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Code generated by ack-generate. DO NOT EDIT.

package alertmanagerdefinition

import (
	"context"

	svcapi "github.com/aws/aws-sdk-go/service/prometheusservice"
	svcsdk "github.com/aws/aws-sdk-go/service/prometheusservice"
	svcsdkapi "github.com/aws/aws-sdk-go/service/prometheusservice/prometheusserviceiface"
	"github.com/google/go-cmp/cmp"
	"github.com/pkg/errors"
	"sigs.k8s.io/controller-runtime/pkg/client"

	xpv1 "github.com/crossplane/crossplane-runtime/apis/common/v1"
	"github.com/crossplane/crossplane-runtime/pkg/meta"
	"github.com/crossplane/crossplane-runtime/pkg/reconciler/managed"
	cpresource "github.com/crossplane/crossplane-runtime/pkg/resource"

	svcapitypes "github.com/crossplane-contrib/provider-aws/apis/prometheusservice/v1alpha1"
	awsclient "github.com/crossplane-contrib/provider-aws/pkg/clients"
)

const (
	errUnexpectedObject = "managed resource is not an AlertManagerDefinition resource"

	errCreateSession = "cannot create a new session"
	errCreate        = "cannot create AlertManagerDefinition in AWS"
	errUpdate        = "cannot update AlertManagerDefinition in AWS"
	errDescribe      = "failed to describe AlertManagerDefinition"
	errDelete        = "failed to delete AlertManagerDefinition"
)

type connector struct {
	kube client.Client
	opts []option
}

func (c *connector) Connect(ctx context.Context, mg cpresource.Managed) (managed.ExternalClient, error) {
	cr, ok := mg.(*svcapitypes.AlertManagerDefinition)
	if !ok {
		return nil, errors.New(errUnexpectedObject)
	}
	sess, err := awsclient.GetConfigV1(ctx, c.kube, mg, cr.Spec.ForProvider.Region)
	if err != nil {
		return nil, errors.Wrap(err, errCreateSession)
	}
	return newExternal(c.kube, svcapi.New(sess), c.opts), nil
}

func (e *external) Observe(ctx context.Context, mg cpresource.Managed) (managed.ExternalObservation, error) {
	cr, ok := mg.(*svcapitypes.AlertManagerDefinition)
	if !ok {
		return managed.ExternalObservation{}, errors.New(errUnexpectedObject)
	}
	if meta.GetExternalName(cr) == "" {
		return managed.ExternalObservation{
			ResourceExists: false,
		}, nil
	}
	input := GenerateDescribeAlertManagerDefinitionInput(cr)
	if err := e.preObserve(ctx, cr, input); err != nil {
		return managed.ExternalObservation{}, errors.Wrap(err, "pre-observe failed")
	}
	resp, err := e.client.DescribeAlertManagerDefinitionWithContext(ctx, input)
	if err != nil {
		return managed.ExternalObservation{ResourceExists: false}, awsclient.Wrap(cpresource.Ignore(IsNotFound, err), errDescribe)
	}
	currentSpec := cr.Spec.ForProvider.DeepCopy()
	if err := e.lateInitialize(&cr.Spec.ForProvider, resp); err != nil {
		return managed.ExternalObservation{}, errors.Wrap(err, "late-init failed")
	}
	GenerateAlertManagerDefinition(resp).Status.AtProvider.DeepCopyInto(&cr.Status.AtProvider)

	upToDate, err := e.isUpToDate(cr, resp)
	if err != nil {
		return managed.ExternalObservation{}, errors.Wrap(err, "isUpToDate check failed")
	}
	return e.postObserve(ctx, cr, resp, managed.ExternalObservation{
		ResourceExists:          true,
		ResourceUpToDate:        upToDate,
		ResourceLateInitialized: !cmp.Equal(&cr.Spec.ForProvider, currentSpec),
	}, nil)
}

func (e *external) Create(ctx context.Context, mg cpresource.Managed) (managed.ExternalCreation, error) {
	cr, ok := mg.(*svcapitypes.AlertManagerDefinition)
	if !ok {
		return managed.ExternalCreation{}, errors.New(errUnexpectedObject)
	}
	cr.Status.SetConditions(xpv1.Creating())
	input := GenerateCreateAlertManagerDefinitionInput(cr)
	if err := e.preCreate(ctx, cr, input); err != nil {
		return managed.ExternalCreation{}, errors.Wrap(err, "pre-create failed")
	}
	resp, err := e.client.CreateAlertManagerDefinitionWithContext(ctx, input)
	if err != nil {
		return managed.ExternalCreation{}, awsclient.Wrap(err, errCreate)
	}

	if resp.Status.StatusCode != nil {
		cr.Status.AtProvider.StatusCode = resp.Status.StatusCode
	} else {
		cr.Status.AtProvider.StatusCode = nil
	}
	if resp.Status.StatusReason != nil {
		cr.Status.AtProvider.StatusReason = resp.Status.StatusReason
	} else {
		cr.Status.AtProvider.StatusReason = nil
	}

	return e.postCreate(ctx, cr, resp, managed.ExternalCreation{}, err)
}

func (e *external) Update(ctx context.Context, mg cpresource.Managed) (managed.ExternalUpdate, error) {
	return e.update(ctx, mg)

}

func (e *external) Delete(ctx context.Context, mg cpresource.Managed) error {
	cr, ok := mg.(*svcapitypes.AlertManagerDefinition)
	if !ok {
		return errors.New(errUnexpectedObject)
	}
	cr.Status.SetConditions(xpv1.Deleting())
	input := GenerateDeleteAlertManagerDefinitionInput(cr)
	ignore, err := e.preDelete(ctx, cr, input)
	if err != nil {
		return errors.Wrap(err, "pre-delete failed")
	}
	if ignore {
		return nil
	}
	resp, err := e.client.DeleteAlertManagerDefinitionWithContext(ctx, input)
	return e.postDelete(ctx, cr, resp, awsclient.Wrap(cpresource.Ignore(IsNotFound, err), errDelete))
}

type option func(*external)

func newExternal(kube client.Client, client svcsdkapi.PrometheusServiceAPI, opts []option) *external {
	e := &external{
		kube:           kube,
		client:         client,
		preObserve:     nopPreObserve,
		postObserve:    nopPostObserve,
		lateInitialize: nopLateInitialize,
		isUpToDate:     alwaysUpToDate,
		preCreate:      nopPreCreate,
		postCreate:     nopPostCreate,
		preDelete:      nopPreDelete,
		postDelete:     nopPostDelete,
		update:         nopUpdate,
	}
	for _, f := range opts {
		f(e)
	}
	return e
}

type external struct {
	kube           client.Client
	client         svcsdkapi.PrometheusServiceAPI
	preObserve     func(context.Context, *svcapitypes.AlertManagerDefinition, *svcsdk.DescribeAlertManagerDefinitionInput) error
	postObserve    func(context.Context, *svcapitypes.AlertManagerDefinition, *svcsdk.DescribeAlertManagerDefinitionOutput, managed.ExternalObservation, error) (managed.ExternalObservation, error)
	lateInitialize func(*svcapitypes.AlertManagerDefinitionParameters, *svcsdk.DescribeAlertManagerDefinitionOutput) error
	isUpToDate     func(*svcapitypes.AlertManagerDefinition, *svcsdk.DescribeAlertManagerDefinitionOutput) (bool, error)
	preCreate      func(context.Context, *svcapitypes.AlertManagerDefinition, *svcsdk.CreateAlertManagerDefinitionInput) error
	postCreate     func(context.Context, *svcapitypes.AlertManagerDefinition, *svcsdk.CreateAlertManagerDefinitionOutput, managed.ExternalCreation, error) (managed.ExternalCreation, error)
	preDelete      func(context.Context, *svcapitypes.AlertManagerDefinition, *svcsdk.DeleteAlertManagerDefinitionInput) (bool, error)
	postDelete     func(context.Context, *svcapitypes.AlertManagerDefinition, *svcsdk.DeleteAlertManagerDefinitionOutput, error) error
	update         func(context.Context, cpresource.Managed) (managed.ExternalUpdate, error)
}

func nopPreObserve(context.Context, *svcapitypes.AlertManagerDefinition, *svcsdk.DescribeAlertManagerDefinitionInput) error {
	return nil
}

func nopPostObserve(_ context.Context, _ *svcapitypes.AlertManagerDefinition, _ *svcsdk.DescribeAlertManagerDefinitionOutput, obs managed.ExternalObservation, err error) (managed.ExternalObservation, error) {
	return obs, err
}
func nopLateInitialize(*svcapitypes.AlertManagerDefinitionParameters, *svcsdk.DescribeAlertManagerDefinitionOutput) error {
	return nil
}
func alwaysUpToDate(*svcapitypes.AlertManagerDefinition, *svcsdk.DescribeAlertManagerDefinitionOutput) (bool, error) {
	return true, nil
}

func nopPreCreate(context.Context, *svcapitypes.AlertManagerDefinition, *svcsdk.CreateAlertManagerDefinitionInput) error {
	return nil
}
func nopPostCreate(_ context.Context, _ *svcapitypes.AlertManagerDefinition, _ *svcsdk.CreateAlertManagerDefinitionOutput, cre managed.ExternalCreation, err error) (managed.ExternalCreation, error) {
	return cre, err
}
func nopPreDelete(context.Context, *svcapitypes.AlertManagerDefinition, *svcsdk.DeleteAlertManagerDefinitionInput) (bool, error) {
	return false, nil
}
func nopPostDelete(_ context.Context, _ *svcapitypes.AlertManagerDefinition, _ *svcsdk.DeleteAlertManagerDefinitionOutput, err error) error {
	return err
}
func nopUpdate(context.Context, cpresource.Managed) (managed.ExternalUpdate, error) {
	return managed.ExternalUpdate{}, nil
}
