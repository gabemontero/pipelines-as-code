package informers

import (
	"context"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	kruntime "k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/tools/cache"
	"k8s.io/klog/v2"

	"knative.dev/pkg/controller"
	"knative.dev/pkg/injection"

	kcpsharedinformer "github.com/kcp-dev/apimachinery/third_party/informers"

	"github.com/openshift-pipelines/pipelines-as-code/pkg/params/clients"

	apispipelinev1beta1 "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1beta1"
	versioned2 "github.com/tektoncd/pipeline/pkg/client/clientset/versioned"
	pipelineruninformer "github.com/tektoncd/pipeline/pkg/client/injection/informers/pipeline/v1beta1/pipelinerun"
	pipelinev1beta1 "github.com/tektoncd/pipeline/pkg/client/listers/pipeline/v1beta1"
)

type Informers struct {
	Clients             clients.Clients
	PipelineRunInformer *prwrapper
	//TODO add webhook certificate informer see webhook main
	//TODO add validating webhook informer see pkg/webhook/controller.go
	//TODO add secret informer see pkg/webhook/controller.go
	//TODO add repository informer see pkg/webhook/controller.go
}

func (i *Informers) NewInformers(ctx context.Context) {
	// access the default tekton injector  package so inits run, and then we override
	_ = pipelineruninformer.Key{}
	i.PipelineRunInformer = &prwrapper{client: i.Clients.Tekton, resourceVersion: injection.GetResourceVersion(ctx)}
	klog.Infof("GGM NewInformers override informer %#v", i.PipelineRunInformer)
	injectionInformer := func(ctx context.Context) (context.Context, controller.Informer) {
		klog.Infof("GGM new injection informer func called")
		return context.WithValue(ctx, pipelineruninformer.Key{}, i.PipelineRunInformer), i.PipelineRunInformer.Informer()
	}
	injection.Default.RegisterInformer(injectionInformer)
}

type prwrapper struct {
	client versioned2.Interface

	namespace string

	resourceVersion string
}

func (w *prwrapper) Informer() cache.SharedIndexInformer {
	//FYI - the reset of the controllers are still doing v1beta1, but I picked v1 to help distinguish.
	klog.Infof("GGM prwrapper Informer returning kcp NewSharedIndexInformer with version specific listerwatcher")
	listerWatcher := &cache.ListWatch{
		ListFunc: func(opts metav1.ListOptions) (kruntime.Object, error) {
			klog.Infof("GGMGGM informer calling our lister")
			return w.client.TektonV1().PipelineRuns(w.namespace).List(context.Background(), opts)
		},
		WatchFunc: func(opts metav1.ListOptions) (watch.Interface, error) {
			klog.Infof("GGMGGM informer calling our watcher")
			return w.client.TektonV1().PipelineRuns(w.namespace).Watch(context.Background(), opts)
		},
	}
	return kcpsharedinformer.NewSharedIndexInformer(listerWatcher, &apispipelinev1beta1.PipelineRun{}, 0, nil)
}

func (w *prwrapper) Lister() pipelinev1beta1.PipelineRunLister {
	klog.Infof("GGM prwrapper global Lister")
	return w
}

func (w *prwrapper) PipelineRuns(namespace string) pipelinev1beta1.PipelineRunNamespaceLister {
	klog.Infof("GGM prwrapper namespaced Lister")
	return &prwrapper{client: w.client, namespace: namespace, resourceVersion: w.resourceVersion}
}

// SetResourceVersion allows consumers to adjust the minimum resourceVersion
// used by the underlying client.  It is not accessible via the standard
// lister interface, but can be accessed through a user-defined interface and
// an implementation check e.g. rvs, ok := foo.(ResourceVersionSetter)
func (w *prwrapper) SetResourceVersion(resourceVersion string) {
	w.resourceVersion = resourceVersion
}

func (w *prwrapper) List(selector labels.Selector) (ret []*apispipelinev1beta1.PipelineRun, err error) {
	klog.Infof("GGM PR List")
	lo, err := w.client.TektonV1beta1().PipelineRuns(w.namespace).List(context.TODO(), metav1.ListOptions{
		LabelSelector:   selector.String(),
		ResourceVersion: w.resourceVersion,
	})
	if err != nil {
		return nil, err
	}
	for idx := range lo.Items {
		ret = append(ret, &lo.Items[idx])
	}
	return ret, nil
}

func (w *prwrapper) Get(name string) (*apispipelinev1beta1.PipelineRun, error) {
	klog.Infof("GGM PR Get")
	return w.client.TektonV1beta1().PipelineRuns(w.namespace).Get(context.TODO(), name, metav1.GetOptions{
		ResourceVersion: w.resourceVersion,
	})
}
