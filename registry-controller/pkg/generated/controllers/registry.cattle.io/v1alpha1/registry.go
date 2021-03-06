/*
Copyright 2019 Wrangler Sample Controller Authors

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

// Code generated by main. DO NOT EDIT.

package v1alpha1

import (
	"context"
	"time"

	v1alpha1 "github.com/rancher/registry-controller/pkg/apis/registry.cattle.io/v1alpha1"
	clientset "github.com/rancher/registry-controller/pkg/generated/clientset/versioned/typed/registry.cattle.io/v1alpha1"
	informers "github.com/rancher/registry-controller/pkg/generated/informers/externalversions/registry.cattle.io/v1alpha1"
	listers "github.com/rancher/registry-controller/pkg/generated/listers/registry.cattle.io/v1alpha1"
	"github.com/rancher/wrangler/pkg/apply"
	"github.com/rancher/wrangler/pkg/condition"
	"github.com/rancher/wrangler/pkg/generic"
	"github.com/rancher/wrangler/pkg/kv"
	"k8s.io/apimachinery/pkg/api/equality"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/tools/cache"
)

type RegistryHandler func(string, *v1alpha1.Registry) (*v1alpha1.Registry, error)

type RegistryController interface {
	generic.ControllerMeta
	RegistryClient

	OnChange(ctx context.Context, name string, sync RegistryHandler)
	OnRemove(ctx context.Context, name string, sync RegistryHandler)
	Enqueue(name string)
	EnqueueAfter(name string, duration time.Duration)

	Cache() RegistryCache
}

type RegistryClient interface {
	Create(*v1alpha1.Registry) (*v1alpha1.Registry, error)
	Update(*v1alpha1.Registry) (*v1alpha1.Registry, error)
	UpdateStatus(*v1alpha1.Registry) (*v1alpha1.Registry, error)
	Delete(name string, options *metav1.DeleteOptions) error
	Get(name string, options metav1.GetOptions) (*v1alpha1.Registry, error)
	List(opts metav1.ListOptions) (*v1alpha1.RegistryList, error)
	Watch(opts metav1.ListOptions) (watch.Interface, error)
	Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1alpha1.Registry, err error)
}

type RegistryCache interface {
	Get(name string) (*v1alpha1.Registry, error)
	List(selector labels.Selector) ([]*v1alpha1.Registry, error)

	AddIndexer(indexName string, indexer RegistryIndexer)
	GetByIndex(indexName, key string) ([]*v1alpha1.Registry, error)
}

type RegistryIndexer func(obj *v1alpha1.Registry) ([]string, error)

type registryController struct {
	controllerManager *generic.ControllerManager
	clientGetter      clientset.RegistriesGetter
	informer          informers.RegistryInformer
	gvk               schema.GroupVersionKind
}

func NewRegistryController(gvk schema.GroupVersionKind, controllerManager *generic.ControllerManager, clientGetter clientset.RegistriesGetter, informer informers.RegistryInformer) RegistryController {
	return &registryController{
		controllerManager: controllerManager,
		clientGetter:      clientGetter,
		informer:          informer,
		gvk:               gvk,
	}
}

func FromRegistryHandlerToHandler(sync RegistryHandler) generic.Handler {
	return func(key string, obj runtime.Object) (ret runtime.Object, err error) {
		var v *v1alpha1.Registry
		if obj == nil {
			v, err = sync(key, nil)
		} else {
			v, err = sync(key, obj.(*v1alpha1.Registry))
		}
		if v == nil {
			return nil, err
		}
		return v, err
	}
}

func (c *registryController) Updater() generic.Updater {
	return func(obj runtime.Object) (runtime.Object, error) {
		newObj, err := c.Update(obj.(*v1alpha1.Registry))
		if newObj == nil {
			return nil, err
		}
		return newObj, err
	}
}

func UpdateRegistryDeepCopyOnChange(client RegistryClient, obj *v1alpha1.Registry, handler func(obj *v1alpha1.Registry) (*v1alpha1.Registry, error)) (*v1alpha1.Registry, error) {
	if obj == nil {
		return obj, nil
	}

	copyObj := obj.DeepCopy()
	newObj, err := handler(copyObj)
	if newObj != nil {
		copyObj = newObj
	}
	if obj.ResourceVersion == copyObj.ResourceVersion && !equality.Semantic.DeepEqual(obj, copyObj) {
		return client.Update(copyObj)
	}

	return copyObj, err
}

func (c *registryController) AddGenericHandler(ctx context.Context, name string, handler generic.Handler) {
	c.controllerManager.AddHandler(ctx, c.gvk, c.informer.Informer(), name, handler)
}

func (c *registryController) AddGenericRemoveHandler(ctx context.Context, name string, handler generic.Handler) {
	removeHandler := generic.NewRemoveHandler(name, c.Updater(), handler)
	c.controllerManager.AddHandler(ctx, c.gvk, c.informer.Informer(), name, removeHandler)
}

func (c *registryController) OnChange(ctx context.Context, name string, sync RegistryHandler) {
	c.AddGenericHandler(ctx, name, FromRegistryHandlerToHandler(sync))
}

func (c *registryController) OnRemove(ctx context.Context, name string, sync RegistryHandler) {
	removeHandler := generic.NewRemoveHandler(name, c.Updater(), FromRegistryHandlerToHandler(sync))
	c.AddGenericHandler(ctx, name, removeHandler)
}

func (c *registryController) Enqueue(name string) {
	c.controllerManager.Enqueue(c.gvk, c.informer.Informer(), "", name)
}

func (c *registryController) EnqueueAfter(name string, duration time.Duration) {
	c.controllerManager.EnqueueAfter(c.gvk, c.informer.Informer(), "", name, duration)
}

func (c *registryController) Informer() cache.SharedIndexInformer {
	return c.informer.Informer()
}

func (c *registryController) GroupVersionKind() schema.GroupVersionKind {
	return c.gvk
}

func (c *registryController) Cache() RegistryCache {
	return &registryCache{
		lister:  c.informer.Lister(),
		indexer: c.informer.Informer().GetIndexer(),
	}
}

func (c *registryController) Create(obj *v1alpha1.Registry) (*v1alpha1.Registry, error) {
	return c.clientGetter.Registries().Create(context.TODO(), obj, metav1.CreateOptions{})
}

func (c *registryController) Update(obj *v1alpha1.Registry) (*v1alpha1.Registry, error) {
	return c.clientGetter.Registries().Update(context.TODO(), obj, metav1.UpdateOptions{})
}

func (c *registryController) UpdateStatus(obj *v1alpha1.Registry) (*v1alpha1.Registry, error) {
	return c.clientGetter.Registries().UpdateStatus(context.TODO(), obj, metav1.UpdateOptions{})
}

func (c *registryController) Delete(name string, options *metav1.DeleteOptions) error {
	if options == nil {
		options = &metav1.DeleteOptions{}
	}
	return c.clientGetter.Registries().Delete(context.TODO(), name, *options)
}

func (c *registryController) Get(name string, options metav1.GetOptions) (*v1alpha1.Registry, error) {
	return c.clientGetter.Registries().Get(context.TODO(), name, options)
}

func (c *registryController) List(opts metav1.ListOptions) (*v1alpha1.RegistryList, error) {
	return c.clientGetter.Registries().List(context.TODO(), opts)
}

func (c *registryController) Watch(opts metav1.ListOptions) (watch.Interface, error) {
	return c.clientGetter.Registries().Watch(context.TODO(), opts)
}

func (c *registryController) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1alpha1.Registry, err error) {
	return c.clientGetter.Registries().Patch(context.TODO(), name, pt, data, metav1.PatchOptions{}, subresources...)
}

type registryCache struct {
	lister  listers.RegistryLister
	indexer cache.Indexer
}

func (c *registryCache) Get(name string) (*v1alpha1.Registry, error) {
	return c.lister.Get(name)
}

func (c *registryCache) List(selector labels.Selector) ([]*v1alpha1.Registry, error) {
	return c.lister.List(selector)
}

func (c *registryCache) AddIndexer(indexName string, indexer RegistryIndexer) {
	utilruntime.Must(c.indexer.AddIndexers(map[string]cache.IndexFunc{
		indexName: func(obj interface{}) (strings []string, e error) {
			return indexer(obj.(*v1alpha1.Registry))
		},
	}))
}

func (c *registryCache) GetByIndex(indexName, key string) (result []*v1alpha1.Registry, err error) {
	objs, err := c.indexer.ByIndex(indexName, key)
	if err != nil {
		return nil, err
	}
	result = make([]*v1alpha1.Registry, 0, len(objs))
	for _, obj := range objs {
		result = append(result, obj.(*v1alpha1.Registry))
	}
	return result, nil
}

type RegistryStatusHandler func(obj *v1alpha1.Registry, status v1alpha1.RegistryStatus) (v1alpha1.RegistryStatus, error)

type RegistryGeneratingHandler func(obj *v1alpha1.Registry, status v1alpha1.RegistryStatus) ([]runtime.Object, v1alpha1.RegistryStatus, error)

func RegisterRegistryStatusHandler(ctx context.Context, controller RegistryController, condition condition.Cond, name string, handler RegistryStatusHandler) {
	statusHandler := &registryStatusHandler{
		client:    controller,
		condition: condition,
		handler:   handler,
	}
	controller.AddGenericHandler(ctx, name, FromRegistryHandlerToHandler(statusHandler.sync))
}

func RegisterRegistryGeneratingHandler(ctx context.Context, controller RegistryController, apply apply.Apply,
	condition condition.Cond, name string, handler RegistryGeneratingHandler, opts *generic.GeneratingHandlerOptions) {
	statusHandler := &registryGeneratingHandler{
		RegistryGeneratingHandler: handler,
		apply:                     apply,
		name:                      name,
		gvk:                       controller.GroupVersionKind(),
	}
	if opts != nil {
		statusHandler.opts = *opts
	}
	controller.OnChange(ctx, name, statusHandler.Remove)
	RegisterRegistryStatusHandler(ctx, controller, condition, name, statusHandler.Handle)
}

type registryStatusHandler struct {
	client    RegistryClient
	condition condition.Cond
	handler   RegistryStatusHandler
}

func (a *registryStatusHandler) sync(key string, obj *v1alpha1.Registry) (*v1alpha1.Registry, error) {
	if obj == nil {
		return obj, nil
	}

	origStatus := obj.Status.DeepCopy()
	obj = obj.DeepCopy()
	newStatus, err := a.handler(obj, obj.Status)
	if err != nil {
		// Revert to old status on error
		newStatus = *origStatus.DeepCopy()
	}

	if a.condition != "" {
		if errors.IsConflict(err) {
			a.condition.SetError(&newStatus, "", nil)
		} else {
			a.condition.SetError(&newStatus, "", err)
		}
	}
	if !equality.Semantic.DeepEqual(origStatus, &newStatus) {
		var newErr error
		obj.Status = newStatus
		obj, newErr = a.client.UpdateStatus(obj)
		if err == nil {
			err = newErr
		}
	}
	return obj, err
}

type registryGeneratingHandler struct {
	RegistryGeneratingHandler
	apply apply.Apply
	opts  generic.GeneratingHandlerOptions
	gvk   schema.GroupVersionKind
	name  string
}

func (a *registryGeneratingHandler) Remove(key string, obj *v1alpha1.Registry) (*v1alpha1.Registry, error) {
	if obj != nil {
		return obj, nil
	}

	obj = &v1alpha1.Registry{}
	obj.Namespace, obj.Name = kv.RSplit(key, "/")
	obj.SetGroupVersionKind(a.gvk)

	return nil, generic.ConfigureApplyForObject(a.apply, obj, &a.opts).
		WithOwner(obj).
		WithSetID(a.name).
		ApplyObjects()
}

func (a *registryGeneratingHandler) Handle(obj *v1alpha1.Registry, status v1alpha1.RegistryStatus) (v1alpha1.RegistryStatus, error) {
	objs, newStatus, err := a.RegistryGeneratingHandler(obj, status)
	if err != nil {
		return newStatus, err
	}

	return newStatus, generic.ConfigureApplyForObject(a.apply, obj, &a.opts).
		WithOwner(obj).
		WithSetID(a.name).
		ApplyObjects(objs...)
}
