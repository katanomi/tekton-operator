package common

import (
	mf "github.com/manifestival/manifestival"
	"github.com/tektoncd/operator/pkg/apis/operator/v1alpha1"
	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/client-go/kubernetes/scheme"
)

// ResourceRequirementsTransform configures the resource requests for
// all containers within all deployments in the manifest
func ResourceRequirementsTransform(deploymentOverRides []v1alpha1.DeploymentOverride) mf.Transformer {
	return func(u *unstructured.Unstructured) error {
		if u.GetKind() == "Deployment" {

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
				if override := find(deploymentOverRide.Resources, containers[i].Name); override != nil {
					merge(&override.Limits, &containers[i].Resources.Limits)
					merge(&override.Requests, &containers[i].Resources.Requests)
				}
			}
			if err := scheme.Scheme.Convert(deployment, u, nil); err != nil {
				return err
			}
			// Avoid superfluous updates from converted zero defaults
			u.SetCreationTimestamp(metav1.Time{})
		}
		return nil
	}
}

func merge(src, tgt *v1.ResourceList) {
	if len(*tgt) > 0 {
		for k, v := range *src {
			(*tgt)[k] = v
		}
	} else {
		*tgt = *src
	}
}

func find(resources []v1alpha1.ResourceRequirementsOverride, name string) *v1alpha1.ResourceRequirementsOverride {
	for _, override := range resources {
		if override.Container == name {
			return &override
		}
	}
	return nil
}
