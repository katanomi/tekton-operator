/*
Copyright 2020 The Tekton Authors

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

package common

import (
	mf "github.com/manifestival/manifestival"
	"github.com/tektoncd/operator/pkg/apis/operator/v1alpha1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/client-go/kubernetes/scheme"
)

func isDeploymentAvailable(d *appsv1.Deployment) bool {
	for _, c := range d.Status.Conditions {
		if c.Type == appsv1.DeploymentAvailable && c.Status == corev1.ConditionTrue {
			return true
		}
	}
	return false
}

// DeploymentOverrideTransform configures the resource requests for
// all containers within all deployments in the manifest
func DeploymentOverrideTransform(deploymentOverRides []v1alpha1.DeploymentOverride) mf.Transformer {
	return func(u *unstructured.Unstructured) error {
		if u.GetKind() != "Deployment" {
			return nil
		}

		var deploymentOverRide v1alpha1.DeploymentOverride
		for _, deployment := range deploymentOverRides {
			if deployment.Name == u.GetName() {
				deploymentOverRide = deployment
				break
			}
		}
		if deploymentOverRide.Name == "" {
			return nil
		}

		deployment := &appsv1.Deployment{}
		if err := scheme.Scheme.Convert(u, deployment, nil); err != nil {
			return err
		}
		containers := deployment.Spec.Template.Spec.Containers
		for i := range containers {
			if override := find(deploymentOverRide.Containers, containers[i].Name); override != nil {
				merge(&override.Resource.Limits, &containers[i].Resources.Limits)
				merge(&override.Resource.Requests, &containers[i].Resources.Requests)

				if len(override.Args) > 0 {
					containers[i].Args = append(containers[i].Args, override.Args...)
				}

				if len(override.Env) > 0 {
					containers[i].Env = upsertEnv(containers[i].Env, override.Env)
				}
			}
		}
		if deploymentOverRide.Replicas != nil {
			deployment.Spec.Replicas = deploymentOverRide.Replicas
		}

		if err := scheme.Scheme.Convert(deployment, u, nil); err != nil {
			return err
		}
		// Avoid superfluous updates from converted zero defaults
		u.SetCreationTimestamp(metav1.Time{})

		return nil
	}
}

func merge(src, tgt *corev1.ResourceList) {
	if src == nil || tgt == nil {
		return
	}
	if len(*tgt) > 0 {
		for k, v := range *src {
			(*tgt)[k] = v
		}
	} else {
		*tgt = *src
	}
}

func find(resources []v1alpha1.ContainerOverride, name string) *v1alpha1.ContainerOverride {
	for _, override := range resources {
		if override.Name == name {
			return &override
		}
	}
	return nil
}

func upsertEnv(exists, overrides []corev1.EnvVar) []corev1.EnvVar {
	for _, override := range overrides {
		var found bool
		for i, exist := range exists {
			if override.Name == exist.Name {
				exists[i] = override
				found = true
				break
			}
		}
		if !found {
			exists = append(exists, override)
		}
	}
	return exists
}
