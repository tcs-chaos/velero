/*
Copyright 2024 The Kubernetes Authors.

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

package v1

import (
	"context"
	"time"

	v1 "github.com/kubernetes-csi/external-snapshotter/client/v7/apis/volumesnapshot/v1"
	scheme "github.com/kubernetes-csi/external-snapshotter/client/v7/clientset/versioned/scheme"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	rest "k8s.io/client-go/rest"
)

// VolumeSnapshotsGetter has a method to return a VolumeSnapshotInterface.
// A group's client should implement this interface.
type VolumeSnapshotsGetter interface {
	VolumeSnapshots(namespace string) VolumeSnapshotInterface
}

// VolumeSnapshotInterface has methods to work with VolumeSnapshot resources.
type VolumeSnapshotInterface interface {
	Create(ctx context.Context, volumeSnapshot *v1.VolumeSnapshot, opts metav1.CreateOptions) (*v1.VolumeSnapshot, error)
	Update(ctx context.Context, volumeSnapshot *v1.VolumeSnapshot, opts metav1.UpdateOptions) (*v1.VolumeSnapshot, error)
	UpdateStatus(ctx context.Context, volumeSnapshot *v1.VolumeSnapshot, opts metav1.UpdateOptions) (*v1.VolumeSnapshot, error)
	Delete(ctx context.Context, name string, opts metav1.DeleteOptions) error
	DeleteCollection(ctx context.Context, opts metav1.DeleteOptions, listOpts metav1.ListOptions) error
	Get(ctx context.Context, name string, opts metav1.GetOptions) (*v1.VolumeSnapshot, error)
	List(ctx context.Context, opts metav1.ListOptions) (*v1.VolumeSnapshotList, error)
	Watch(ctx context.Context, opts metav1.ListOptions) (watch.Interface, error)
	Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts metav1.PatchOptions, subresources ...string) (result *v1.VolumeSnapshot, err error)
	VolumeSnapshotExpansion
}

// volumeSnapshots implements VolumeSnapshotInterface
type volumeSnapshots struct {
	client rest.Interface
	ns     string
}

// newVolumeSnapshots returns a VolumeSnapshots
func newVolumeSnapshots(c *SnapshotV1Client, namespace string) *volumeSnapshots {
	return &volumeSnapshots{
		client: c.RESTClient(),
		ns:     namespace,
	}
}

// Get takes name of the volumeSnapshot, and returns the corresponding volumeSnapshot object, and an error if there is any.
func (c *volumeSnapshots) Get(ctx context.Context, name string, options metav1.GetOptions) (result *v1.VolumeSnapshot, err error) {
	result = &v1.VolumeSnapshot{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("volumesnapshots").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do(ctx).
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of VolumeSnapshots that match those selectors.
func (c *volumeSnapshots) List(ctx context.Context, opts metav1.ListOptions) (result *v1.VolumeSnapshotList, err error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	result = &v1.VolumeSnapshotList{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("volumesnapshots").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Do(ctx).
		Into(result)
	return
}

// Watch returns a watch.Interface that watches the requested volumeSnapshots.
func (c *volumeSnapshots) Watch(ctx context.Context, opts metav1.ListOptions) (watch.Interface, error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	opts.Watch = true
	return c.client.Get().
		Namespace(c.ns).
		Resource("volumesnapshots").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Watch(ctx)
}

// Create takes the representation of a volumeSnapshot and creates it.  Returns the server's representation of the volumeSnapshot, and an error, if there is any.
func (c *volumeSnapshots) Create(ctx context.Context, volumeSnapshot *v1.VolumeSnapshot, opts metav1.CreateOptions) (result *v1.VolumeSnapshot, err error) {
	result = &v1.VolumeSnapshot{}
	err = c.client.Post().
		Namespace(c.ns).
		Resource("volumesnapshots").
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(volumeSnapshot).
		Do(ctx).
		Into(result)
	return
}

// Update takes the representation of a volumeSnapshot and updates it. Returns the server's representation of the volumeSnapshot, and an error, if there is any.
func (c *volumeSnapshots) Update(ctx context.Context, volumeSnapshot *v1.VolumeSnapshot, opts metav1.UpdateOptions) (result *v1.VolumeSnapshot, err error) {
	result = &v1.VolumeSnapshot{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("volumesnapshots").
		Name(volumeSnapshot.Name).
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(volumeSnapshot).
		Do(ctx).
		Into(result)
	return
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().
func (c *volumeSnapshots) UpdateStatus(ctx context.Context, volumeSnapshot *v1.VolumeSnapshot, opts metav1.UpdateOptions) (result *v1.VolumeSnapshot, err error) {
	result = &v1.VolumeSnapshot{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("volumesnapshots").
		Name(volumeSnapshot.Name).
		SubResource("status").
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(volumeSnapshot).
		Do(ctx).
		Into(result)
	return
}

// Delete takes name of the volumeSnapshot and deletes it. Returns an error if one occurs.
func (c *volumeSnapshots) Delete(ctx context.Context, name string, opts metav1.DeleteOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("volumesnapshots").
		Name(name).
		Body(&opts).
		Do(ctx).
		Error()
}

// DeleteCollection deletes a collection of objects.
func (c *volumeSnapshots) DeleteCollection(ctx context.Context, opts metav1.DeleteOptions, listOpts metav1.ListOptions) error {
	var timeout time.Duration
	if listOpts.TimeoutSeconds != nil {
		timeout = time.Duration(*listOpts.TimeoutSeconds) * time.Second
	}
	return c.client.Delete().
		Namespace(c.ns).
		Resource("volumesnapshots").
		VersionedParams(&listOpts, scheme.ParameterCodec).
		Timeout(timeout).
		Body(&opts).
		Do(ctx).
		Error()
}

// Patch applies the patch and returns the patched volumeSnapshot.
func (c *volumeSnapshots) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts metav1.PatchOptions, subresources ...string) (result *v1.VolumeSnapshot, err error) {
	result = &v1.VolumeSnapshot{}
	err = c.client.Patch(pt).
		Namespace(c.ns).
		Resource("volumesnapshots").
		Name(name).
		SubResource(subresources...).
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(data).
		Do(ctx).
		Into(result)
	return
}
