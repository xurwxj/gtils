package sys

type SysInfoObj struct {
	Soft     OSSoftExtObj `json:"soft"`
	Hardware OSHdExtObj   `json:"hardware"`
}

type OSSoftExtObj struct {
	OSName    string `form:"osName" json:"osName" binding:"required"`
	HostName  string `form:"hostName" json:"hostName" binding:"required"`
	OSType    string `form:"osType" json:"osType" binding:"required"`
	OSVersion string `form:"osVersion" json:"osVersion" binding:"required"`
	// get from soft
	Lang       string `form:"lang" json:"lang"`
	Resolution string `form:"resolution" json:"resolution" binding:"required"`
	// get from soft
	InstallPath string `json:"installPath"`
	// can't get from os for now
	BrowserObj `json:"browser"`
}

type BrowserObj struct {
	Name string `form:"name" json:"name"`
	Ver  string `form:"ver" json:"ver"`
}

type OSHdExtObj struct {
	CpuObjs       []CpuObj `json:"cpu" binding:"required"`
	GpuObjs       []GpuObj `json:"gpu" binding:"required"`
	MemSize       int64    `form:"memSize" json:"memSize" binding:"required"`
	HDDSize       int64    `form:"hddSize" json:"hddSize" binding:"required"`
	HDDRemainSize int64    `form:"hddRemainSize" json:"hddRemainSize" binding:"required"`
}

type GpuObj struct {
	GPUModel   string `form:"gpuModel" json:"gpuModel" binding:"required"`
	GPUMemSize int64  `form:"gpuMemSize" json:"gpuMemSize" binding:"required"`
	GPUVendor  string `form:"gpuVendor" json:"gpuVendor" binding:"required"`
}

type CpuObj struct {
	CPUCores  int64  `form:"cpuCores" json:"cpuCores" binding:"required"`
	CPUModel  string `form:"cpuModel" json:"cpuModel" binding:"required"`
	CPUVendor string `form:"cpuVendor" json:"cpuVendor" binding:"required"`
}
