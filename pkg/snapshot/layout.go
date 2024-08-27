package snapshot

import (
	"path/filepath"

	v1 "k8s.io/api/core/v1"
)

// /host/disk/{pv path}/
// /host/disk/{pv path}/{pv.name}.img
// /host/disk/{pv path}/{pv.name}/velero/{snapshot.id}/snapshot.img
// /host/disk/{pv path}/{pv.name}/velero/{snapshot.id}/snapshot

func Disk(disk string) layoutDisk {
	return defaultLayout().Disk(disk)
}

func defaultLayout() layout {
	return "/host"
}

type layout string

func (l layout) Path() string {
	return string(l)
}

func (l layout) Disk(disk string) layoutDisk {
	if disk == "" {
		disk = LayOutDefaultDisk
	}
	return layoutDisk(filepath.Join(string(l), disk))
}

type layoutDisk string

func (l layoutDisk) Path() string {
	return string(l)
}

func (l layoutDisk) PV(pv *v1.PersistentVolume) layoutPV {
	return layoutPV(filepath.Join(l.Path(), LayOutDefaultPV, pv.Name))
}

type layoutPV string

func (l layoutPV) Path() string {
	return string(l)
}

func (l layoutPV) Image() string {
	return string(l) + ".img"
}

func (l layoutPV) Snapshot(snapshotID string) layoutSnapshot {
	return layoutSnapshot(filepath.Join(string(l), "velero", snapshotID))
}

type layoutSnapshot string

func (l layoutSnapshot) Path() string {
	return string(l)
}

func (l layoutSnapshot) Image() string {
	return filepath.Join(string(l), "snapshot.img")
}

func (l layoutSnapshot) MountPoint() string {
	return filepath.Join(string(l), "snapshot")
}
