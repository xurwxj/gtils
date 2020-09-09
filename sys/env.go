package sys

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/cloudfoundry/jibber_jabber"
	"github.com/rs/zerolog"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/mem"
	"github.com/xurwxj/gtils/base"
)

func GetOsInfo(Log *zerolog.Logger) SysInfoObj {
	var rs SysInfoObj
	var err error
	var n *host.InfoStat
	if c, err := cpu.Info(); err == nil {
		rs.Hardware.CpuObjs = append(rs.Hardware.CpuObjs, CpuObj{CPUCores: int64(c[0].Cores), CPUModel: c[0].ModelName, CPUVendor: c[0].VendorID})
	} else {
		Log.Err(err).Msg("GetOsInfo get cpu")
	}
	if v, err := mem.VirtualMemory(); err == nil {
		rs.Hardware.MemSize = int64(v.Total)
	} else {
		Log.Err(err).Msg("GetOsInfo get memory")
	}
	if n, err = host.Info(); err == nil {
		rs.Soft.OSName = n.Platform
		rs.Soft.OSType = n.OS
		rs.Soft.OSVersion = n.PlatformVersion
		rs.Soft.HostName = n.Hostname
	} else {
		Log.Err(err).Msg("GetOsInfo get host")
	}
	// i, _ := net.Interfaces()
	path := "/"
	if n.OS != "" {
		switch n.OS {
		case "windows":
			path = "\\"
			cmd := exec.Command("wmic", "path", "win32_VideoController", "get", "/all", "/format:list")
			if out, err := cmd.Output(); err == nil {
				scanner := bufio.NewScanner(strings.NewReader(string(out)))
				var gobj GpuObj
				var h, w string
				for scanner.Scan() {
					dLineT := scanner.Text()
					lArr := strings.Split(dLineT, "=")
					if label := strings.TrimSpace(lArr[0]); len(lArr) > 1 && base.FindInStringSlice([]string{"Name", "AdapterRAM", "AdapterCompatibility", "CurrentVerticalResolution", "CurrentHorizontalResolution"}, label) {
						switch label {
						case "Name":
							gobj.GPUModel = strings.TrimSpace(lArr[1])
						case "AdapterRAM":
							if ar, err := strconv.ParseInt(strings.TrimSpace(lArr[1]), 10, 64); err == nil {
								gobj.GPUMemSize = ar
							}
						case "AdapterCompatibility":
							gobj.GPUVendor = strings.TrimSpace(lArr[1])
						case "CurrentVerticalResolution":
							h = strings.TrimSpace(lArr[1])
						case "CurrentHorizontalResolution":
							w = strings.TrimSpace(lArr[1])
						}
					}
				}
				rs.Hardware.GpuObjs = append(rs.Hardware.GpuObjs, gobj)
				if h != "" && w != "" {
					rs.Soft.Resolution = fmt.Sprintf("%sx%s", w, h)
				}
			} else {
				Log.Err(err).Msg("GetOsInfo get gpu info under window")
			}
		case "linux":
			cmd := exec.Command("lspci")
			if out, err := cmd.Output(); err == nil {
				scanner := bufio.NewScanner(strings.NewReader(string(out)))
				for scanner.Scan() {
					dLineT := scanner.Text()
					var gobj GpuObj
					lArr := strings.Split(dLineT, "controller:")
					if len(lArr) > 1 && strings.Index(lArr[0], "VGA") > -1 {
						gobj.GPUModel = strings.TrimSpace(lArr[1])
						gobj.GPUVendor = strings.Split(gobj.GPUModel, " ")[0]
						gobj.GPUMemSize = 0
						// Log.Info(strings.TrimSpace(lArr[1]))
						rs.Hardware.GpuObjs = append(rs.Hardware.GpuObjs, gobj)
					}
				}
			} else {
				Log.Err(err).Msg("GetOsInfo get gpu info under linux")
			}
			// Log.Info(err)
			cmdd := exec.Command("xdpyinfo")
			if outd, errd := cmdd.Output(); errd == nil {
				scannerd := bufio.NewScanner(strings.NewReader(string(outd)))
				for scannerd.Scan() {
					dLineTd := scannerd.Text()
					lArrd := strings.Split(dLineTd, ":")
					if len(lArrd) > 1 && strings.Index(lArrd[0], "dimensions") > -1 {
						rs.Soft.Resolution = strings.Split(strings.TrimSpace(lArrd[1]), " ")[0]
					}
				}
			} else {
				rs.Soft.Resolution = "unknown"
				Log.Err(errd).Msg("GetOsInfo get resolution info under linux")
			}
			// Log.Info(errd)
		case "darwin":
			cmd := exec.Command("system_profiler", "SPDisplaysDataType")
			if out, err := cmd.Output(); err == nil {
				scanner := bufio.NewScanner(strings.NewReader(string(out)))
				var gobj GpuObj
				for scanner.Scan() {
					dLineT := scanner.Text()
					lArr := strings.Split(dLineT, ":")
					if label := strings.TrimSpace(lArr[0]); len(lArr) > 1 && base.FindInStringSlice([]string{"Chipset Model", "Vendor", "VRAM (Dynamic, Max)", "Resolution"}, label) {
						switch label {
						case "Chipset Model":
							gobj.GPUModel = strings.TrimSpace(lArr[1])
						case "VRAM (Dynamic, Max)":
							vr := strings.TrimSpace(lArr[1])
							if strings.Index(vr, "MB") > -1 {
								if vri, err := strconv.ParseInt(strings.TrimSpace(strings.Split(vr, " ")[0]), 10, 64); err == nil {
									gobj.GPUMemSize = vri * 1024 * 1024
								}
							} else if ar, err := strconv.ParseInt(vr, 10, 64); err == nil {
								gobj.GPUMemSize = ar
							}
						case "Vendor":
							gobj.GPUVendor = strings.TrimSpace(lArr[1])
						case "Resolution":
							rs.Soft.Resolution = strings.TrimSpace(lArr[1])
						}
					}
				}
				rs.Hardware.GpuObjs = append(rs.Hardware.GpuObjs, gobj)
			} else {
				Log.Err(err).Msg("GetOsInfo get gpu info under mac")
			}
		}
		if diskStat, err := disk.Usage(path); err == nil {
			rs.Hardware.HDDSize = int64(diskStat.Total)
			rs.Hardware.HDDRemainSize = int64(diskStat.Free)
		} else {
			Log.Err(err).Msg("GetOsInfo get disk info under mac")
		}
	}
	currentDir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		Log.Err(err).Msg("GetOsInfo get install path")
	}
	rs.Soft.InstallPath = currentDir
	rs.Soft = CheckOsExtInfo(rs.Soft, Log)
	return rs
}
func GetHostID() string {
	var n *host.InfoStat
	var err error
	if n, err = host.Info(); err == nil {
		// fmt.Println(n.BootTime)
		// fmt.Println(n.Uptime)
		// fmt.Println(n.HostID)
		return n.HostID
	}
	return ""
}

func CheckOsExtInfo(OssE OSSoftExtObj, Log *zerolog.Logger) OSSoftExtObj {
	if OssE.Resolution == "" {
		OssE.Resolution = "unknown"
	}
	if OssE.Lang == "" {
		userLanguage, err := jibber_jabber.DetectLanguage()
		if err != nil {
			Log.Err(err).Interface("OssE", OssE).Msg("CheckOsExtInfo DetectLanguage")
		}
		OssE.Lang = userLanguage
	}
	if OssE.Lang == "" {
		OssE.Lang = "unknown"
	}
	return OssE
}
