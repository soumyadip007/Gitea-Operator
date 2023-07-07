/*
Copyright 2023 Soumyadip Chowdhury.

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

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client/config"

	apiv1alpha1 "github.com/soumyadip007/gitea-operator/api/v1alpha1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
)

// GiteaReconciler reconciles a Gitea object
type GiteaReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=api.gitea.k8s,resources=gitea,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=api.gitea.k8s,resources=gitea/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=api.gitea.k8s,resources=gitea/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Gitea object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.14.1/pkg/reconcile
func (r *GiteaReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = log.FromContext(ctx)
	log := log.FromContext(ctx)

	gitea := &apiv1alpha1.Gitea{}

	log.Info(fmt.Sprintf("%v", gitea.Spec))

	err := r.Get(ctx, req.NamespacedName, gitea)
	if err != nil {
		if apierrors.IsNotFound(err) {
			log.Info("Gitea resource not found. Ignoring since object must be deleted.")
			return ctrl.Result{}, nil
		}
		log.Error(err, "Failed")
		return ctrl.Result{}, err
	}

	containerport := gitea.Spec.ContainerPort
	image := gitea.Spec.Image
	name := gitea.Spec.Name
	nodeport := gitea.Spec.NodePort
	port := gitea.Spec.Port
	replicas := gitea.Spec.Replicas
	targetport := gitea.Spec.TargetPort
	versions := gitea.Spec.Versions

	log.Info(fmt.Sprintf("ContainerPort: %v, Image: %v, Name: %v, NodePort: %v, Port: %v, Replicas: %v, TargetPort: %v, Versions: %v", containerport, image, name, nodeport, port, replicas, targetport, versions))

	/*
		log.Info(fmt.Sprintf("ContainerPort: %v", containerport))
		log.Info(fmt.Sprintf("Image: %v", image))
		log.Info(fmt.Sprintf("Name: %v", name))
		log.Info(fmt.Sprintf("NodePort: %v", nodeport))
		log.Info(fmt.Sprintf("Port: %v", port))
		log.Info(fmt.Sprintf("Replicas: %v", replicas))
		log.Info(fmt.Sprintf("TargetPort: %v", targetport))
		log.Info(fmt.Sprintf("Versions: %v", versions))
	*/

	// Path to your kubeconfig file
	//kubeconfig := "/path/to/kubeconfig"

	// Build the client configuration from the kubeconfig file
	//config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)

	// Get the client configuration
	cfg, err := config.GetConfig()
	if err != nil {
		log.Error(err, "Failed in client configuration")
	}

	// Create a new client
	cl, err := client.New(cfg, client.Options{})
	if err != nil {
		log.Error(err, "Failed in creating a new client")
	}

	// Create a new Namespace object
	namespace := &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: name + "-ns",
		},
	}

	// Create the namespace using the client
	err = cl.Create(context.TODO(), namespace)
	if err != nil {
		log.Error(err, "Failed in creating a namespace")
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *GiteaReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&apiv1alpha1.Gitea{}).
		Complete(r)
}
