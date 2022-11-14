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

// Code generated by client-gen. DO NOT EDIT.

package fake

import (
	"context"

	v1alpha1 "github.com/openshift-pipelines/pipelines-as-code/pkg/apis/pipelinesascode/v1alpha1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	labels "k8s.io/apimachinery/pkg/labels"
	schema "k8s.io/apimachinery/pkg/runtime/schema"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	testing "k8s.io/client-go/testing"
)

// FakeRepositories implements RepositoryInterface
type FakeRepositories struct {
	Fake *FakePipelinesascodeV1alpha1
	ns   string
}

var repositoriesResource = schema.GroupVersionResource{Group: "pipelinesascode.tekton.dev", Version: "v1alpha1", Resource: "repositories"}

var repositoriesKind = schema.GroupVersionKind{Group: "pipelinesascode.tekton.dev", Version: "v1alpha1", Kind: "Repository"}

// Get takes name of the repository, and returns the corresponding repository object, and an error if there is any.
func (c *FakeRepositories) Get(ctx context.Context, name string, options v1.GetOptions) (result *v1alpha1.Repository, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewGetAction(repositoriesResource, c.ns, name), &v1alpha1.Repository{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.Repository), err
}

// List takes label and field selectors, and returns the list of Repositories that match those selectors.
func (c *FakeRepositories) List(ctx context.Context, opts v1.ListOptions) (result *v1alpha1.RepositoryList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewListAction(repositoriesResource, repositoriesKind, c.ns, opts), &v1alpha1.RepositoryList{})

	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &v1alpha1.RepositoryList{ListMeta: obj.(*v1alpha1.RepositoryList).ListMeta}
	for _, item := range obj.(*v1alpha1.RepositoryList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested repositories.
func (c *FakeRepositories) Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewWatchAction(repositoriesResource, c.ns, opts))

}

// Create takes the representation of a repository and creates it.  Returns the server's representation of the repository, and an error, if there is any.
func (c *FakeRepositories) Create(ctx context.Context, repository *v1alpha1.Repository, opts v1.CreateOptions) (result *v1alpha1.Repository, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewCreateAction(repositoriesResource, c.ns, repository), &v1alpha1.Repository{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.Repository), err
}

// Update takes the representation of a repository and updates it. Returns the server's representation of the repository, and an error, if there is any.
func (c *FakeRepositories) Update(ctx context.Context, repository *v1alpha1.Repository, opts v1.UpdateOptions) (result *v1alpha1.Repository, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateAction(repositoriesResource, c.ns, repository), &v1alpha1.Repository{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.Repository), err
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().
func (c *FakeRepositories) UpdateStatus(ctx context.Context, repository *v1alpha1.Repository, opts v1.UpdateOptions) (*v1alpha1.Repository, error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateSubresourceAction(repositoriesResource, "status", c.ns, repository), &v1alpha1.Repository{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.Repository), err
}

// Delete takes name of the repository and deletes it. Returns an error if one occurs.
func (c *FakeRepositories) Delete(ctx context.Context, name string, opts v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewDeleteActionWithOptions(repositoriesResource, c.ns, name, opts), &v1alpha1.Repository{})

	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeRepositories) DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error {
	action := testing.NewDeleteCollectionAction(repositoriesResource, c.ns, listOpts)

	_, err := c.Fake.Invokes(action, &v1alpha1.RepositoryList{})
	return err
}

// Patch applies the patch and returns the patched repository.
func (c *FakeRepositories) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v1alpha1.Repository, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewPatchSubresourceAction(repositoriesResource, c.ns, name, pt, data, subresources...), &v1alpha1.Repository{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.Repository), err
}
