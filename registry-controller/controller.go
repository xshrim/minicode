package main

import (
	"context"
	"fmt"

	registryv1alpha1 "github.com/rancher/registry-controller/pkg/apis/registry.cattle.io/v1alpha1"
	registryscheme "github.com/rancher/registry-controller/pkg/generated/clientset/versioned/scheme"
	"github.com/rancher/registry-controller/pkg/generated/controllers/registry.cattle.io/v1alpha1"

	v1 "github.com/rancher/wrangler-api/pkg/generated/controllers/apps/v1"
	"github.com/xshrim/gol"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/client-go/kubernetes/scheme"
	typedcorev1 "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/tools/record"
)

//	crdv1 "github.com/rancher/wrangler-api/pkg/generated/controllers/apiextensions.k8s.io/v1beta1"
const controllerAgentName = "registry-controller"

const (
	// ErrResourceExists is used as part of the Event 'reason' when a Foo fails
	// to sync due to a Deployment of the same name already existing.
	ErrResourceExists = "ErrResourceExists"

	// MessageResourceExists is the message used for Events when a resource
	// fails to sync due to a Deployment already existing
	MessageResourceExists = "Resource %q already exists and is not managed by Foo"
)

// Handler is the controller implementation for Foo resources
type Handler struct {
	deployments      v1.DeploymentClient
	deploymentsCache v1.DeploymentCache
	registries       v1alpha1.RegistryController
	registriesCache  v1alpha1.RegistryCache
	recorder         record.EventRecorder
}

// NewController returns a new sample controller
func Register(
	ctx context.Context,
	events typedcorev1.EventInterface,
	deployments v1.DeploymentController,
	registries v1alpha1.RegistryController) {

	controller := &Handler{
		deployments:      deployments,
		deploymentsCache: deployments.Cache(),
		registries:       registries,
		registriesCache:  registries.Cache(),
		recorder:         buildEventRecorder(events),
	}

	// Register handlers
	deployments.OnChange(ctx, "registry-handler", controller.OnDeploymentChanged)
	registries.OnChange(ctx, "registry-handler", controller.OnRegistryChanged)
}

func logf(format string, args ...interface{}) {
	gol.Infof(format, args...)
}

func buildEventRecorder(events typedcorev1.EventInterface) record.EventRecorder {
	// Create event broadcaster
	// Add sample-controller types to the default Kubernetes Scheme so Events can be
	// logged for sample-controller types.
	utilruntime.Must(registryscheme.AddToScheme(scheme.Scheme))
	gol.Info("Creating event broadcaster")
	eventBroadcaster := record.NewBroadcaster()
	eventBroadcaster.StartLogging(logf)
	eventBroadcaster.StartRecordingToSink(&typedcorev1.EventSinkImpl{Interface: events})
	return eventBroadcaster.NewRecorder(scheme.Scheme, corev1.EventSource{Component: controllerAgentName})
}

func (h *Handler) OnRegistryChanged(key string, registry *registryv1alpha1.Registry) (*registryv1alpha1.Registry, error) {
	// foo will be nil if key is deleted from cache
	gol.Prtln(key)
	if registry == nil {
		return nil, nil
	}

	typ := registry.Spec.Type
	if typ == "" {
		typ = "harbor"
	}
	cluster := registry.Spec.Cluster
	if cluster == "" {
		// We choose to absorb the error here as the worker would requeue the
		// resource otherwise. Instead, the next time the resource is updated
		// the resource will be queued again.
		utilruntime.HandleError(fmt.Errorf("%s: cluster id must be specified", key))
		return nil, nil
	}
	owner := registry.Spec.Owner
	if owner == "" {
		utilruntime.HandleError(fmt.Errorf("%s: owner must be specified", key))
		return nil, nil
	}
	url := registry.Spec.Url
	if url == "" {
		utilruntime.HandleError(fmt.Errorf("%s: url must be specified", key))
		return nil, nil
	}
	user := registry.Spec.User
	passwd := registry.Spec.Passwd

	gol.Info(typ, cluster, owner, url, user, passwd)
	// // Get the deployment with the name specified in Foo.spec
	// deployment, err := h.deploymentsCache.Get(registry.Namespace, deploymentName)
	// // If the resource doesn't exist, we'll create it
	// if errors.IsNotFound(err) {
	// 	deployment, err = h.deployments.Create(newDeployment(registry))
	// }

	// // If an error occurs during Get/Create, we'll requeue the item so we can
	// // attempt processing again later. This could have been caused by a
	// // temporary network failure, or any other transient reason.
	// if err != nil {
	// 	return nil, err
	// }

	// // If the Deployment is not controlled by this Foo resource, we should log
	// // a warning to the event recorder and ret
	// if !metav1.IsControlledBy(deployment, registry) {
	// 	msg := fmt.Sprintf(MessageResourceExists, deployment.Name)
	// 	h.recorder.Event(registry, corev1.EventTypeWarning, ErrResourceExists, msg)
	// 	// Notice we don't return an error here, this is intentional because an
	// 	// error means we should retry to reconcile.  In this situation we've done all
	// 	// we could, which was log an error.
	// 	return nil, nil
	// }

	// // If this number of the replicas on the Foo resource is specified, and the
	// // number does not equal the current desired replicas on the Deployment, we
	// // should update the Deployment resource.
	// if registry.Spec.Replicas != nil && *registry.Spec.Replicas != *deployment.Spec.Replicas {
	// 	gol.Infof("Registry %s replicas: %d, deployment replicas: %d", registry.Name, *registry.Spec.Replicas, *deployment.Spec.Replicas)
	// 	deployment, err = h.deployments.Update(newDeployment(registry))
	// }

	// // If an error occurs during Update, we'll requeue the item so we can
	// // attempt processing again later. THis could have been caused by a
	// // temporary network failure, or any other transient reason.
	// if err != nil {
	// 	return nil, err
	// }

	// Finally, we update the status block of the Foo resource to reflect the
	// current state of the world
	// err = h.updateRegistryStatus(registry, deployment)
	// if err != nil {
	// 	return nil, err
	// }

	return nil, nil
}

func (h *Handler) updateRegistryStatus(registry *registryv1alpha1.Registry, deployment *appsv1.Deployment) error {
	// NEVER modify objects from the store. It's a read-only, local cache.
	// You can use DeepCopy() to make a deep copy of original object and modify this copy
	// Or create a copy manually for better performance
	// registryCopy := registry.DeepCopy()
	// registryCopy.Status.AvailableReplicas = deployment.Status.AvailableReplicas
	// // If the CustomResourceSubresources feature gate is not enabled,
	// // we must use Update instead of UpdateStatus to update the Status block of the Foo resource.
	// // UpdateStatus will not allow changes to the Spec of the resource,
	// // which is ideal for ensuring nothing other than resource status has been updated.
	// _, err := h.registries.Update(registryCopy)
	// return err
	return nil
}

func (h *Handler) OnDeploymentChanged(key string, deployment *appsv1.Deployment) (*appsv1.Deployment, error) {
	// When an item is deleted the deployment is nil, just ignore
	// if deployment == nil {
	// 	return nil, nil
	// }

	// if ownerRef := metav1.GetControllerOf(deployment); ownerRef != nil {
	// 	// If this object is not owned by a Foo, we should not do anything more
	// 	// with it.
	// 	if ownerRef.Kind != "Foo" {
	// 		return nil, nil
	// 	}

	// 	registry, err := h.registryCache.Get(deployment.Namespace, ownerRef.Name)
	// 	if err != nil {
	// 		gol.Infof("ignoring orphaned object '%s' of registry '%s'", deployment.GetSelfLink(), ownerRef.Name)
	// 		return nil, nil
	// 	}

	// 	h.registries.Enqueue(registry.Namespace, registry.Name)
	// 	return nil, nil
	// }

	return nil, nil
}

// newDeployment creates a new Deployment for a Foo resource. It also sets
// the appropriate OwnerReferences on the resource so handleObject can discover
// the Foo resource that 'owns' it.
func newDeployment(registry *registryv1alpha1.Registry) *appsv1.Deployment {
	// labels := map[string]string{
	// 	"app":        "nginx",
	// 	"controller": registry.Name,
	// }
	// return &appsv1.Deployment{
	// 	ObjectMeta: metav1.ObjectMeta{
	// 		Name:      registry.Spec.DeploymentName,
	// 		Namespace: registry.Namespace,
	// 		OwnerReferences: []metav1.OwnerReference{
	// 			*metav1.NewControllerRef(registry, schema.GroupVersionKind{
	// 				Group:   registryv1alpha1.SchemeGroupVersion.Group,
	// 				Version: registryv1alpha1.SchemeGroupVersion.Version,
	// 				Kind:    "Registry",
	// 			}),
	// 		},
	// 	},
	// 	Spec: appsv1.DeploymentSpec{
	// 		Replicas: registry.Spec.Replicas,
	// 		Selector: &metav1.LabelSelector{
	// 			MatchLabels: labels,
	// 		},
	// 		Template: corev1.PodTemplateSpec{
	// 			ObjectMeta: metav1.ObjectMeta{
	// 				Labels: labels,
	// 			},
	// 			Spec: corev1.PodSpec{
	// 				Containers: []corev1.Container{
	// 					{
	// 						Name:  "nginx",
	// 						Image: "nginx:latest",
	// 					},
	// 				},
	// 			},
	// 		},
	// 	},
	// }
	return nil
}
