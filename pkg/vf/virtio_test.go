package vf

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/Code-Hex/vz/v3"
	"github.com/crc-org/vfkit/pkg/config"
	"github.com/stretchr/testify/require"
)

func TestImageRootSyncModeKeepsCachedIO(t *testing.T) {
	path := filepath.Join(t.TempDir(), "root.raw")
	require.NoError(t, os.WriteFile(path, make([]byte, 512), 0600))
	for _, tc := range []struct {
		name   string
		option string
		want   vz.DiskImageSynchronizationMode
	}{
		{"default", "", vz.DiskImageSynchronizationModeFsync},
		{"ignore guest flushes", ",sync=none", vz.DiskImageSynchronizationModeNone},
	} {
		t.Run(tc.name, func(t *testing.T) {
			vm := config.NewVirtualMachine(1, 512, nil)
			require.NoError(t, vm.AddDevicesFromCmdLine([]string{"virtio-blk,path=" + path + tc.option}))
			root := vm.Devices[0].(*config.VirtioBlk)
			require.Equal(t, vz.DiskImageCachingModeCached, toVzCachingMode(root.CachingMode))
			require.Equal(t, tc.want, toVzSyncMode(root.SynchronizationMode))
		})
	}
}
