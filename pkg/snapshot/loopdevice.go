package snapshot

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/pkg/errors"
	v1 "k8s.io/api/core/v1"
	"k8s.io/klog/v2"
)

type loopDevice struct {
}

func (ld *loopDevice) CreateSnapshot(volume *v1.PersistentVolume, uuid string) (path string, err error) {
	return createSnapshot(volume, uuid)
}

func (ld *loopDevice) DeleteSnapshot(volume *v1.PersistentVolume, uuid string) error {
	return deleteSnapshot(volume, uuid)
}

func deleteSnapshot(volume *v1.PersistentVolume, uuid string) error {
	pv := Disk("").PV(volume)
	snapshot := pv.Snapshot(uuid)
	klog.Infof("delete snapshot %s", snapshot.Path())
	if err := Unmount(snapshot.MountPoint()); err != nil {
		klog.Errorf("unmount %s failed: %v", snapshot.MountPoint(), err)
		return errors.WithStack(err)
	}
	if err := os.RemoveAll(snapshot.Path()); err != nil {
		klog.Errorf("remove %s failed: %v", snapshot.Path(), err)
		return errors.WithStack(err)
	}
	return nil
}

func createSnapshot(volume *v1.PersistentVolume, uuid string) (path string, err error) {
	pv := Disk("").PV(volume)
	if err = EnsureDir(pv.Path()); err != nil {
		return "", errors.WithStack(err)
	}

	// create snapshot image
	var src, dst string
	snapshot := pv.Snapshot(uuid)
	src, dst = pv.Image(), snapshot.Image()
	if err = EnsureDir(snapshot.Path()); err != nil {
		return "", errors.WithStack(err)
	}
	if err = ReFlink(src, dst); err != nil {
		return "", errors.WithStack(err)
	}

	// create mount path
	src, dst = snapshot.Image(), snapshot.MountPoint()
	if err = EnsureDir(dst); err != nil {
		return "", errors.WithStack(err)
	}
	if err = MountWithNoUUID(src, dst); err != nil {
		return "", errors.WithStack(err)
	}
	return dst, nil
}

func EnsureDir(dir string) error {
	if _, err := os.Stat(dir); err != nil && os.IsNotExist(err) {
		if err = os.MkdirAll(dir, os.ModePerm); err != nil {
			return errors.WithStack(err)
		}
		if err = os.Chmod(dir, os.ModePerm); err != nil {
			return errors.WithStack(err)
		}
	}
	return nil
}

func ReFlink(src, dst string, opts ...string) (err error) {
	klog.V(4).Infof("reflink %s to %s", src, dst)
	if err = exec.Command("cp", append([]string{"--reflink=always", src, dst}, opts...)...).Run(); err != nil {
		klog.V(4).Infof("reflink %s to %s failed: %v", src, dst, err)
		return errors.Wrapf(err, "reflink %s to %s failed", src, dst)
	}
	return nil
}

func MountWithNoUUID(src, dst string) error {
	return Mount(src, dst, "nouuid")
}

func Mount(src, dst string, options ...string) error {
	klog.V(4).Infof("mount %s to %s", src, dst)

	args := []string{src, dst}

	if len(options) > 0 {
		args = append(args, "-o")
		args = append(args, strings.Join(options, ","))
	}

	cmd := exec.Command("mount", args...)
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("mount failed: %v, output: %s", err, string(output))
	}
	return nil
}

func Unmount(path string) error {
	klog.V(4).Infof("umount %s", path)
	if err := exec.Command("umount", path).Run(); err != nil {
		return errors.Wrapf(err, "umount %s failed", path)
	}
	return nil
}
