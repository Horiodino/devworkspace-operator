// Copyright (c) 2019-2021 Red Hat, Inc.
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package sync

import (
	"context"
	"fmt"
	"reflect"

	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/runtime"
	crclient "sigs.k8s.io/controller-runtime/pkg/client"
)

type ClusterAPI struct {
	Client crclient.Client
	Scheme *runtime.Scheme
	Logger logr.Logger
	Ctx    context.Context
}

type NotInSyncError struct {
	Reason NotInSyncReason
	Object crclient.Object
}

type NotInSyncReason string

const (
	UpdatedObjectReason NotInSyncReason = "Updated object"
	CreatedObjectReason NotInSyncReason = "Created object"
	NeedRetryReason     NotInSyncReason = "Need to retry"
)

func (e *NotInSyncError) Error() string {
	return fmt.Sprintf("%s %s is not ready: %s", reflect.TypeOf(e.Object).Elem().String(), e.Object.GetName(), e.Reason)
}

func NewNotInSync(obj crclient.Object, reason NotInSyncReason) *NotInSyncError {
	return &NotInSyncError{
		Reason: reason,
		Object: obj,
	}
}

type UnrecoverableSyncError struct {
	Cause error
}

func (e *UnrecoverableSyncError) Error() string {
	return e.Cause.Error()
}
