/*
Copyright 2025.

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

package controller

import (
	"context"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	examplev1 "github.com/helin0815/crd-learn/api/v1"
)

// HlTestReconciler reconciles a HlTest object
type HlTestReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=example.example.com,resources=hltests,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=example.example.com,resources=hltests/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=example.example.com,resources=hltests/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the HlTest object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.20.2/pkg/reconcile
func (r *HlTestReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)
	hltest := &examplev1.HlTest{}
	if err := r.Get(ctx, req.NamespacedName, hltest); err != nil {
		logger.Error(err, "Failed to get HlTest")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}
	pod := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      req.Name,
			Namespace: req.Namespace,
		},
		Spec: corev1.PodSpec{
			InitContainers: []corev1.Container{
				{
					Name:    "init-volume",
					Image:   "busybox",
					Command: []string{"sh", "-c", "mkdir -p /data/hltest"},
					VolumeMounts: []corev1.VolumeMount{
						{Name: "data", MountPath: "/data"},
					},
				},
			},
			Containers: []corev1.Container{{
				Name:    "busybox",
				Image:   "busybox",
				Command: []string{"sh", "-c", "echo 'hello " + hltest.Spec.User + "' > /data/hltest/a.log && sleep 3600"},
				VolumeMounts: []corev1.VolumeMount{
					{Name: "data", MountPath: "/data"},
				},
			}},
			Volumes: []corev1.Volume{
				{Name: "data", VolumeSource: corev1.VolumeSource{EmptyDir: &corev1.EmptyDirVolumeSource{}}},
			},
			RestartPolicy: corev1.RestartPolicyNever,
		},
		Status: corev1.PodStatus{},
	}

	found := &corev1.Pod{}
	err := r.Get(ctx, types.NamespacedName{Name: pod.Name, Namespace: pod.Namespace}, found)
	if err != nil && errors.IsNotFound(err) {
		logger.Info("Creating a new Pod", "Pod.Namespace", pod.Namespace, "Pod.Name", pod.Name)
		err = r.Create(ctx, pod)
		if err != nil {
			return ctrl.Result{}, err
		}
	} else if err != nil {
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *HlTestReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&examplev1.HlTest{}).
		Named("hltest").
		Complete(r)
}
