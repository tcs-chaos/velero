package snapshot

import (
	"context"
	"fmt"

	"github.com/sirupsen/logrus"
	v1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/vmware-tanzu/velero/pkg/util/kube"
)

type Internal interface {
	CreateSnapshot(volume *v1.PersistentVolume, pvdNamespacedName string,
		log logrus.FieldLogger) (mount string, err error)
	DeleteSnapshot(volume *v1.PersistentVolume, pvdNamespacedName string,
		log logrus.FieldLogger) error
}

type Snapshot interface {
	CreateSnapshot() (mount string, err error)
	DeleteSnapshot() error
}

func NewSnapshotWithName(client client.Client, pod *v1.Pod, pvdNamespacedName, volumeName string,
	log logrus.FieldLogger) Snapshot {
	_, volume, _, err := kube.GetPodPVCVolume(context.Background(), log, pod, volumeName, client)
	if err != nil {
		return &ErrSnapshot{err: err}
	}
	return NewSnapshot(pvdNamespacedName, volume, log)
}

func NewSnapshot(pvdNamespacedName string, volume *v1.PersistentVolume, log logrus.FieldLogger) Snapshot {
	switch volume.Spec.StorageClassName {
	case StorageClassLoopDevice:
		return NewLoopDevice(pvdNamespacedName, volume, log)
	}
	return &ErrSnapshot{err: fmt.Errorf("unsupported storage class %s", volume)}
}

type ErrSnapshot struct {
	err error
}

func (es *ErrSnapshot) CreateSnapshot() (mount string, err error) {
	return "", fmt.Errorf("error snapshot, %v", es.err)
}

func (es *ErrSnapshot) DeleteSnapshot() error {
	return fmt.Errorf("error snapshot, %v", es.err)
}

type LoopDevice struct {
	PvdNamespacedName string
	Volume            *v1.PersistentVolume

	log  logrus.FieldLogger
	impl Internal
}

func NewLoopDevice(pvdNamespacedName string, volume *v1.PersistentVolume, log logrus.FieldLogger) *LoopDevice {
	return &LoopDevice{
		Volume:            volume,
		PvdNamespacedName: pvdNamespacedName,
		impl:              &loopDevice{},
		log:               log,
	}
}

func (ld *LoopDevice) CreateSnapshot() (mount string, err error) {
	return ld.impl.CreateSnapshot(ld.Volume, ld.PvdNamespacedName, ld.log)
}

func (ld *LoopDevice) DeleteSnapshot() error {
	return ld.impl.DeleteSnapshot(ld.Volume, ld.PvdNamespacedName, ld.log)
}
