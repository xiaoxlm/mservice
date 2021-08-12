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
	"k8s.io/apimachinery/pkg/runtime"
	"r.kubebuilder.io/pkg/apply"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	testv1 "r.kubebuilder.io/api/v1"
)

// MServiceReconciler reconciles a MService object
type MServiceReconciler struct {
	client.Client
	Log logr.Logger
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

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *MServiceReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&testv1.MService{}).
		Complete(r)
}

func (r *MServiceReconciler) apply(ctx context.Context, msvc *testv1.MService) error {
	objects := convert(msvc)

	if len(objects) < 1 {
		return fmt.Errorf("objects is empty")
	}

	var errs []error

	for _, o := range objects {
		if err := apply.Action(ctx, r.Client, o.GetNamespace(), o); err != nil {
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
	//objects := convert(msvc)
	return nil
}

func convert(msvc *testv1.MService) []client.Object {
	var ret []client.Object

	// apply ingress
	ingress := msvc.Spec.Ingress.Convert(&msvc.ObjectMeta, msvc.GetLabels(), msvc.GetAnnotations())
	// apply service
	service := msvc.Spec.Ports.Convert(&msvc.ObjectMeta, msvc.GetLabels(), msvc.GetAnnotations())
	// apply deployment
	deployment := msvc.Spec.MDeployment.Convert(&msvc.ObjectMeta, msvc.GetLabels(), msvc.GetAnnotations())
	// apply secret
	secret := msvc.Spec.Secret.Convert(&msvc.ObjectMeta, msvc.GetLabels(), msvc.GetAnnotations())

	ret = append(ret, ingress, service, deployment, secret)

	return ret
}



