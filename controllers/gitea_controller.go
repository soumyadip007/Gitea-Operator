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

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"sigs.k8s.io/controller-runtime/pkg/client/config"

	apiv1alpha1 "github.com/soumyadip007/gitea-operator/api/v1alpha1"
	//"k8s.io/apimachinery/pkg/api/errors"
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
	namespace := gitea.Spec.NameSpace
	servicename := gitea.Spec.Service
	deploymentname := gitea.Spec.DeploymentName
	nodeport := gitea.Spec.NodePort
	port := gitea.Spec.Port
	replicas := gitea.Spec.Replicas
	targetport := gitea.Spec.TargetPort
	versions := gitea.Spec.Versions

	log.Info(fmt.Sprintf("ContainerPort: %v, Image: %v, Name: %v, NodePort: %v, Port: %v, Replicas: %v, TargetPort: %v, Versions: %v, Serivce: %v, Namespace: %v, DeploymentName : %v", containerport, image, name, nodeport, port, replicas, targetport, versions, servicename, namespace, deploymentname))

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

	if err := createNamespace(cl, namespace); err != nil {
		log.Error(err, "Failed in creating a namespace")
	} else {
		log.Info("Namespace created successfully")
	}

	if err := createService(cl, name, namespace, servicename, port, targetport, nodeport); err != nil {
		log.Error(err, "Failed in creating a service")
	} else {
		log.Info("Service created successfully")
	}

	if err := createDeployment(cl, name, deploymentname, namespace, image, containerport); err != nil {
		log.Error(err, "Failed in creating a deployment")
	} else {
		log.Info("Deployment created successfully")
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *GiteaReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&apiv1alpha1.Gitea{}).
		Complete(r)
}

func createNamespace(cl client.Client, name string) error {
	// Check if the namespace already exists
	existingNamespace := &corev1.Namespace{}
	err := cl.Get(context.TODO(), client.ObjectKey{Name: name}, existingNamespace)
	if err == nil {
		// Namespace already exists, no need to create it
		return nil
	} else if !apierrors.IsNotFound(err) {
		// An error occurred other than "not found", return the error
		return err
	}

	// Create a new Namespace object
	namespace := &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
	}

	// Create the namespace using the client
	err = cl.Create(context.TODO(), namespace)
	if err != nil {
		return err
	}

	return nil
}

func createService(cl client.Client, name, namespace, servicename string, port int32, targetport int, nodeport int32) error {
	// Check if the service already exists
	existingService := &corev1.Service{}
	err := cl.Get(context.TODO(), client.ObjectKey{Name: servicename, Namespace: namespace}, existingService)
	if err == nil {
		// Service already exists, no need to create it
		return nil
	} else if !apierrors.IsNotFound(err) {
		// An error occurred other than "not found", return the error
		return err
	}

	// Create a new Service object
	service := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      servicename,
			Namespace: namespace,
		},
		Spec: corev1.ServiceSpec{
			Type: corev1.ServiceTypeNodePort,
			Ports: []corev1.ServicePort{
				{
					Port:       port,
					TargetPort: intstr.FromInt(int(targetport)),
					NodePort:   nodeport,
				},
			},
			Selector: map[string]string{
				"server": name,
			},
		},
	}

	// Create the service using the client
	err = cl.Create(context.TODO(), service)
	if err != nil {
		return nil
	}

	return nil
}

// create deployment from user input
func createDeployment(cl client.Client, name string, deploymentname string, namespace string, image string, containerport int32) error {
	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      deploymentname,
			Namespace: namespace,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: int32Ptr(3),
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"server": name,
				},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"server": name,
					},
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  name,
							Image: image,
							Ports: []corev1.ContainerPort{
								{
									ContainerPort: containerport,
								},
							},
							Env: []corev1.EnvVar{
								{
									Name:  "NODE_ENV",
									Value: "production",
								},
							},
						},
					},
				},
			},
		},
	}
	err := cl.Create(context.TODO(), deployment)
	if err != nil {
		// It won't throw any error
		return nil
	}
	return nil
}

func int32Ptr(i int32) *int32 {
	return &i
}
