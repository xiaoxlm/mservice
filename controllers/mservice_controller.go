/*
Copyright 2021.

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

package controllers

import (
	"context"
	"fmt"
	"github.com/go-logr/logr"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	testv1 "r.kubebuilder.io/api/v1"
	"r.kubebuilder.io/pkg/apply"
	"r.kubebuilder.io/pkg/components"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"time"
)

const (
	APP_LABEL_KEY = "app"
)

// MServiceReconciler reconciles a MService object
type MServiceReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=test.r.kubebuilder.io,resources=mservices,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=test.r.kubebuilder.io,resources=mservices/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=test.r.kubebuilder.io,resources=mservices/finalizers,verbs=update
//+kubebuilder:rbac:groups=test.r.kubebuilder.io,resources=deployment,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=test.r.kubebuilder.io,resources=deployment/status,verbs=get

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the MService object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.8.3/pkg/reconcile
func (r *MServiceReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	r.Log = r.Log.WithValues("mservice", req.NamespacedName)
	// your logic here
	var msvc testv1.MService
	if err := r.Get(ctx, req.NamespacedName, &msvc); err != nil {
		r.Log.Error(err, "unable to fetch MService")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	// apply
	if err := r.apply(ctx, &msvc); err != nil {
		return ctrl.Result{}, err
	}

	// status update
	if err := r.updateStatus(ctx, &msvc); err != nil {
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *MServiceReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&testv1.MService{}).
		Owns(&corev1.Pod{}).
		Complete(r)
}

func (r *MServiceReconciler) apply(ctx context.Context, msvc *testv1.MService) error {
	objectWithGVKs := convert(msvc)

	if len(objectWithGVKs) < 1 {
		return fmt.Errorf("objects is empty")
	}

	var errs []error

	for _, o := range objectWithGVKs {
		// owner
		if err := controllerutil.SetControllerReference(msvc, o.Object, r.Scheme); err != nil {
			errs = append(errs, err)
			continue
		}

		if err := apply.Action(ctx, r.Client, o.Object.GetNamespace(), o.Object); err != nil {
			errs = append(errs, err)
		}
	}

	if len(errs) > 0 {
		errString := ""
		for _, e := range errs {
			errString += e.Error() + ";"
		}
		return fmt.Errorf(errString)
	}

	return nil
}

func (r *MServiceReconciler) updateStatus(ctx context.Context, msvc *testv1.MService) error {
	// wait 5 second
	time.Sleep(5 * time.Second)

	objectWithGVKs := convert(msvc)

	// check resource
	for _, o := range objectWithGVKs {
		currentObject, err := components.GetCurrentObject(o.Object)
		if err != nil {
			return err
		}

		if err := r.Client.Get(ctx, client.ObjectKeyFromObject(o.Object), currentObject); err != nil {
			if apierrors.IsNotFound(err) {
				return fmt.Errorf("%v not found", o.GVK)
			}
		}
	}

	// update deployment status
	//deployment := &appsv1.Deployment{}
	//if err := r.Client.Get(ctx, client.ObjectKeyFromObject(msvc), deployment); err != nil {
	//	return err
	//}
	//
	//podList := new(corev1.PodList)
	//if err := r.Client.List(
	//	ctx, podList,
	//	client.InNamespace(msvc.GetNamespace()),
	//	client.MatchingLabels(map[string]string{
	//		APP_LABEL_KEY: msvc.GetName(),
	//	}),
	//); err != nil {
	//	return err
	//}
	//
	msvc.Status.Stage = "PROCESSING"
	//
	if err := r.Client.Status().Update(ctx, msvc); err != nil {
		return err
	}

	// check pod status for apply result

	return nil
}

type DeploymentStage = string

const (
	DeploymentStageDone       DeploymentStage = "DONE"
	DeploymentStageFail       DeploymentStage = "FAIL"
	DeploymentStageProcessing DeploymentStage = "PROCESSING"
)

type ObjectWithGVK struct {
	Object client.Object
	GVK    schema.GroupVersionKind
}

func convert(msvc *testv1.MService) []ObjectWithGVK {
	//var ret []ObjectWithGVK

	formatObjectWithGVK := func(resources ...components.ISelfResource) []ObjectWithGVK {
		var ret []ObjectWithGVK
		for _, r := range resources {
			o, gvk := r.Convert(&msvc.ObjectMeta, msvc.GetLabels(), msvc.GetAnnotations())

			ret = append(ret, ObjectWithGVK{
				Object: o,
				GVK:    gvk,
			})
		}

		return ret
	}

	//return formatObjectWithGVK(msvc.Spec.Ingress, msvc.Spec.Ports, &msvc.Spec.MDeployment, msvc.Spec.Secret)
	return formatObjectWithGVK(msvc.Spec.Ingress, msvc.Spec.Ports, msvc.Spec.Secret)
}
