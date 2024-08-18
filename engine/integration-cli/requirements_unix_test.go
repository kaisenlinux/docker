//go:build !windows

package main

import (
<<<<<<< HEAD
	"os"
=======
	"bytes"
	"io/ioutil"
	"os/exec"
>>>>>>> parent of ea55db5 (Import the 20.10.24 version)
	"strings"

	"github.com/docker/docker/pkg/sysinfo"
)

var sysInfo *sysinfo.SysInfo

func setupLocalInfo() {
	sysInfo = sysinfo.New()
}

func cpuCfsPeriod() bool {
	return testEnv.DaemonInfo.CPUCfsPeriod
}

func cpuCfsQuota() bool {
	return testEnv.DaemonInfo.CPUCfsQuota
}

func cpuShare() bool {
	return testEnv.DaemonInfo.CPUShares
}

func oomControl() bool {
	return testEnv.DaemonInfo.OomKillDisable
}

func pidsLimit() bool {
	return sysInfo.PidsLimit
}

func memoryLimitSupport() bool {
	return testEnv.DaemonInfo.MemoryLimit
}

func memoryReservationSupport() bool {
	return sysInfo.MemoryReservation
}

func swapMemorySupport() bool {
	return testEnv.DaemonInfo.SwapLimit
}

func memorySwappinessSupport() bool {
	return testEnv.IsLocalDaemon() && sysInfo.MemorySwappiness
}

func blkioWeight() bool {
	return testEnv.IsLocalDaemon() && sysInfo.BlkioWeight
}

func cgroupCpuset() bool {
	return testEnv.DaemonInfo.CPUSet
}

func seccompEnabled() bool {
	return sysInfo.Seccomp
}

func bridgeNfIptables() bool {
	return !sysInfo.BridgeNFCallIPTablesDisabled
}

func unprivilegedUsernsClone() bool {
	content, err := ioutil.ReadFile("/proc/sys/kernel/unprivileged_userns_clone")
	return err != nil || !strings.Contains(string(content), "0")
}
