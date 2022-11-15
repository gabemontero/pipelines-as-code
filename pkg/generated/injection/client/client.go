/*
Copyright Red Hat

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

// Code generated by injection-gen. DO NOT EDIT.

package client

import (
	context "context"
	json "encoding/json"
	errors "errors"
	fmt "fmt"
	"k8s.io/klog/v2"

	v1alpha1 "github.com/openshift-pipelines/pipelines-as-code/pkg/apis/pipelinesascode/v1alpha1"
	versioned "github.com/openshift-pipelines/pipelines-as-code/pkg/generated/clientset/versioned"
	typedpipelinesascodev1alpha1 "github.com/openshift-pipelines/pipelines-as-code/pkg/generated/clientset/versioned/typed/pipelinesascode/v1alpha1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	unstructured "k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	runtime "k8s.io/apimachinery/pkg/runtime"
	schema "k8s.io/apimachinery/pkg/runtime/schema"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	discovery "k8s.io/client-go/discovery"
	dynamic "k8s.io/client-go/dynamic"
	rest "k8s.io/client-go/rest"
	injection "knative.dev/pkg/injection"
	dynamicclient "knative.dev/pkg/injection/clients/dynamicclient"
	logging "knative.dev/pkg/logging"
)

func init() {
	injection.Default.RegisterClient(withClientFromConfig)
	injection.Default.RegisterClientFetcher(func(ctx context.Context) interface{} {
		return Get(ctx)
	})
	injection.Dynamic.RegisterDynamicClient(withClientFromDynamic)
}

// Key is used as the key for associating information with a context.Context.
type Key struct{}

func withClientFromConfig(ctx context.Context, cfg *rest.Config) context.Context {
	return context.WithValue(ctx, Key{}, versioned.NewForConfigOrDie(cfg))
}

func withClientFromDynamic(ctx context.Context) context.Context {
	return context.WithValue(ctx, Key{}, &wrapClient{dyn: dynamicclient.Get(ctx)})
}

// Get extracts the versioned.Interface client from the context.
func Get(ctx context.Context) versioned.Interface {
	klog.Infof("GGM pac client injection getter called")
	untyped := ctx.Value(Key{})
	if untyped == nil {
		if injection.GetConfig(ctx) == nil {
			logging.FromContext(ctx).Panic(
				"Unable to fetch github.com/openshift-pipelines/pipelines-as-code/pkg/generated/clientset/versioned.Interface from context. This context is not the application context (which is typically given to constructors via sharedmain).")
		} else {
			logging.FromContext(ctx).Panic(
				"Unable to fetch github.com/openshift-pipelines/pipelines-as-code/pkg/generated/clientset/versioned.Interface from context.")
		}
	}
	return untyped.(versioned.Interface)
}

type wrapClient struct {
	dyn dynamic.Interface
}

var _ versioned.Interface = (*wrapClient)(nil)

func (w *wrapClient) Discovery() discovery.DiscoveryInterface {
	panic("Discovery called on dynamic client!")
}

func convert(from interface{}, to runtime.Object) error {
	bs, err := json.Marshal(from)
	if err != nil {
		return fmt.Errorf("Marshal() = %w", err)
	}
	if err := json.Unmarshal(bs, to); err != nil {
		return fmt.Errorf("Unmarshal() = %w", err)
	}
	return nil
}

// PipelinesascodeV1alpha1 retrieves the PipelinesascodeV1alpha1Client
func (w *wrapClient) PipelinesascodeV1alpha1() typedpipelinesascodev1alpha1.PipelinesascodeV1alpha1Interface {
	return &wrapPipelinesascodeV1alpha1{
		dyn: w.dyn,
	}
}

type wrapPipelinesascodeV1alpha1 struct {
	dyn dynamic.Interface
}

func (w *wrapPipelinesascodeV1alpha1) RESTClient() rest.Interface {
	panic("RESTClient called on dynamic client!")
}

func (w *wrapPipelinesascodeV1alpha1) Repositories(namespace string) typedpipelinesascodev1alpha1.RepositoryInterface {
	return &wrapPipelinesascodeV1alpha1RepositoryImpl{
		dyn: w.dyn.Resource(schema.GroupVersionResource{
			Group:    "pipelinesascode.tekton.dev",
			Version:  "v1alpha1",
			Resource: "repositories",
		}),

		namespace: namespace,
	}
}

type wrapPipelinesascodeV1alpha1RepositoryImpl struct {
	dyn dynamic.NamespaceableResourceInterface

	namespace string
}

var _ typedpipelinesascodev1alpha1.RepositoryInterface = (*wrapPipelinesascodeV1alpha1RepositoryImpl)(nil)

func (w *wrapPipelinesascodeV1alpha1RepositoryImpl) Create(ctx context.Context, in *v1alpha1.Repository, opts v1.CreateOptions) (*v1alpha1.Repository, error) {
	in.SetGroupVersionKind(schema.GroupVersionKind{
		Group:   "pipelinesascode.tekton.dev",
		Version: "v1alpha1",
		Kind:    "Repository",
	})
	uo := &unstructured.Unstructured{}
	if err := convert(in, uo); err != nil {
		return nil, err
	}
	uo, err := w.dyn.Namespace(w.namespace).Create(ctx, uo, opts)
	if err != nil {
		return nil, err
	}
	out := &v1alpha1.Repository{}
	if err := convert(uo, out); err != nil {
		return nil, err
	}
	return out, nil
}

func (w *wrapPipelinesascodeV1alpha1RepositoryImpl) Delete(ctx context.Context, name string, opts v1.DeleteOptions) error {
	return w.dyn.Namespace(w.namespace).Delete(ctx, name, opts)
}

func (w *wrapPipelinesascodeV1alpha1RepositoryImpl) DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error {
	return w.dyn.Namespace(w.namespace).DeleteCollection(ctx, opts, listOpts)
}

func (w *wrapPipelinesascodeV1alpha1RepositoryImpl) Get(ctx context.Context, name string, opts v1.GetOptions) (*v1alpha1.Repository, error) {
	uo, err := w.dyn.Namespace(w.namespace).Get(ctx, name, opts)
	if err != nil {
		return nil, err
	}
	out := &v1alpha1.Repository{}
	if err := convert(uo, out); err != nil {
		return nil, err
	}
	return out, nil
}

func (w *wrapPipelinesascodeV1alpha1RepositoryImpl) List(ctx context.Context, opts v1.ListOptions) (*v1alpha1.RepositoryList, error) {
	uo, err := w.dyn.Namespace(w.namespace).List(ctx, opts)
	if err != nil {
		return nil, err
	}
	out := &v1alpha1.RepositoryList{}
	if err := convert(uo, out); err != nil {
		return nil, err
	}
	return out, nil
}

func (w *wrapPipelinesascodeV1alpha1RepositoryImpl) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v1alpha1.Repository, err error) {
	uo, err := w.dyn.Namespace(w.namespace).Patch(ctx, name, pt, data, opts)
	if err != nil {
		return nil, err
	}
	out := &v1alpha1.Repository{}
	if err := convert(uo, out); err != nil {
		return nil, err
	}
	return out, nil
}

func (w *wrapPipelinesascodeV1alpha1RepositoryImpl) Update(ctx context.Context, in *v1alpha1.Repository, opts v1.UpdateOptions) (*v1alpha1.Repository, error) {
	in.SetGroupVersionKind(schema.GroupVersionKind{
		Group:   "pipelinesascode.tekton.dev",
		Version: "v1alpha1",
		Kind:    "Repository",
	})
	uo := &unstructured.Unstructured{}
	if err := convert(in, uo); err != nil {
		return nil, err
	}
	uo, err := w.dyn.Namespace(w.namespace).Update(ctx, uo, opts)
	if err != nil {
		return nil, err
	}
	out := &v1alpha1.Repository{}
	if err := convert(uo, out); err != nil {
		return nil, err
	}
	return out, nil
}

func (w *wrapPipelinesascodeV1alpha1RepositoryImpl) UpdateStatus(ctx context.Context, in *v1alpha1.Repository, opts v1.UpdateOptions) (*v1alpha1.Repository, error) {
	in.SetGroupVersionKind(schema.GroupVersionKind{
		Group:   "pipelinesascode.tekton.dev",
		Version: "v1alpha1",
		Kind:    "Repository",
	})
	uo := &unstructured.Unstructured{}
	if err := convert(in, uo); err != nil {
		return nil, err
	}
	uo, err := w.dyn.Namespace(w.namespace).UpdateStatus(ctx, uo, opts)
	if err != nil {
		return nil, err
	}
	out := &v1alpha1.Repository{}
	if err := convert(uo, out); err != nil {
		return nil, err
	}
	return out, nil
}

func (w *wrapPipelinesascodeV1alpha1RepositoryImpl) Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error) {
	return nil, errors.New("NYI: Watch")
}
