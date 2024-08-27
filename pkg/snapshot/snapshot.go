package snapshot

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	v1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/vmware-tanzu/velero/pkg/util/kube"
)

type Internal interface {
	CreateSnapshot(volume *v1.PersistentVolume, uuid string) (mount string, err error)
	DeleteSnapshot(volume *v1.PersistentVolume, uuid string) error
}

type Snapshot interface {
	CreateSnapshot() (mount string, err error)
	DeleteSnapshot() error
}

func NewSnapshotWithName(Client client.Client, pod *v1.Pod, volumeName string, log logrus.FieldLogger) Snapshot {
	_, volume, _, err := kube.GetPodPVCVolume(context.Background(), log, pod, volumeName, Client)
	if err != nil {
		return &ErrSnapshot{err: err}
	}
	return NewSnapshot(volume)
}

func NewSnapshot(volume *v1.PersistentVolume) Snapshot {
	switch volume.Spec.StorageClassName {
	case StorageClassLoopDevice:
		return NewLoopDevice(volume)
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
	UUID   string
	Volume *v1.PersistentVolume

	impl Internal
}

func NewLoopDevice(volume *v1.PersistentVolume) *LoopDevice {
	return &LoopDevice{
		Volume: volume,
		UUID:   uuid.New().String(),
		impl:   &loopDevice{},
	}
}

func (ld *LoopDevice) CreateSnapshot() (mount string, err error) {
	return ld.impl.CreateSnapshot(ld.Volume, ld.UUID)
}

func (ld *LoopDevice) DeleteSnapshot() error {
	return ld.impl.DeleteSnapshot(ld.Volume, ld.UUID)
}
