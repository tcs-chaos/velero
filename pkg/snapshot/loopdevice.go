package snapshot

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	v1 "k8s.io/api/core/v1"
)

type loopDevice struct {
}

func (ld *loopDevice) CreateSnapshot(volume *v1.PersistentVolume, pvdNamespacedName string,
	log logrus.FieldLogger) (path string, err error) {
	disk := ""
	if volume.Spec.CSI.VolumeAttributes != nil {
		disk = volume.Spec.CSI.VolumeAttributes[DiskSelector]
	}
	pv := Disk(disk).PV(volume)
	if err = EnsureDir(pv.Path()); err != nil {
		return "", errors.WithStack(err)
	}

	// create snapshot image
	var src, dst string
	snapshot := pv.Snapshot(pvdNamespacedName)
	src, dst = pv.Image(), snapshot.Image()
	if err = EnsureDir(snapshot.Path()); err != nil {
		log.Errorf("ensure dir %s failed: %v", snapshot.Path(), err)
		return "", errors.WithStack(err)
	}
	if err = ReFlink(src, dst, log); err != nil {
		log.Errorf("reflink %s to %s failed: %v", src, dst, err)
		return "", errors.WithStack(err)
	}

	// create a mount path
	src, dst = snapshot.Image(), snapshot.MountPoint()
	if err = EnsureDir(dst); err != nil {
		log.Errorf("ensure dir %s failed: %v", snapshot.Path(), err)
		return "", errors.WithStack(err)
	}
	if err = MountWithNoUUID(src, dst, log); err != nil {
		return "", errors.WithStack(err)
	}
	return dst, nil
}

func (ld *loopDevice) DeleteSnapshot(volume *v1.PersistentVolume, pvdNamespacedName string,
	log logrus.FieldLogger) error {
	disk := ""
	if volume.Spec.CSI.VolumeAttributes != nil {
		disk = volume.Spec.CSI.VolumeAttributes[DiskSelector]
	}
	pv := Disk(disk).PV(volume)
	snapshot := pv.Snapshot(pvdNamespacedName)
	log.Infof("delete snapshot %s", snapshot.Path())
	// we can get here even if the snapshot image is not mounted, so we need to umount it anyway.
	// so that we can delete the snapshot image
	if err := Unmount(snapshot.MountPoint(), log); err != nil {
		log.Warnf("unmount %s failed: %v", snapshot.MountPoint(), err)
	}
	if err := os.RemoveAll(snapshot.Path()); err != nil {
		log.Errorf("remove %s failed: %v", snapshot.Path(), err)
		// if the path does not exist (we may meet this condition), RemoveAll returns nil (no error).
		// os return error here
		return errors.WithStack(err)
	}
	log.Infof("delete snapshot %s success", snapshot.Path())
	return nil
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

func ReFlink(src, dst string, log logrus.FieldLogger, opts ...string) (err error) {
	log.Infof("reflink %s to %s", src, dst)
	if err = exec.Command("cp", append([]string{"--reflink=always", src, dst}, opts...)...).Run(); err != nil {
		log.Infof("reflink %s to %s failed: %v", src, dst, err)
		return errors.Wrapf(err, "reflink %s to %s failed", src, dst)
	}
	return nil
}

func MountWithNoUUID(src, dst string, log logrus.FieldLogger) error {
	return Mount(src, dst, log, "nouuid")
}

func Mount(src, dst string, log logrus.FieldLogger, options ...string) error {
	log.Infof("mount %s to %s", src, dst)

	args := []string{src, dst}

	if len(options) > 0 {
		args = append(args, "-o")
		args = append(args, strings.Join(options, ","))
	}

	const maxRetries = 10
	const retryInterval = time.Second

	var lastErr error
	for i := 0; i < maxRetries; i++ {
		cmd := exec.Command("mount", args...)
		if output, err := cmd.CombinedOutput(); err != nil {
			lastErr = fmt.Errorf("mount failed: %v, output: %s", err, string(output))
			log.Errorf("attempt to mount %s to %s, this is round %d, err: %s", src, dst, i+1, lastErr)
			time.Sleep(retryInterval)
			continue
		}
		return nil
	}
	return lastErr
}

func Unmount(path string, log logrus.FieldLogger) error {
	log.Infof("umount %s", path)
	if err := exec.Command("umount", path).Run(); err != nil {
		return errors.Wrapf(err, "umount %s failed", path)
	}
	return nil
}
