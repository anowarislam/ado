package meta

import (
	"context"
	"log/slog"
	"strings"

	"github.com/jaypipes/ghw"
	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/disk"
	"github.com/shirou/gopsutil/v4/host"
	"github.com/shirou/gopsutil/v4/mem"
)

// SystemInfo represents comprehensive system diagnostic information.
type SystemInfo struct {
	OS           string        `json:"os" yaml:"os"`
	Platform     string        `json:"platform" yaml:"platform"`
	Kernel       string        `json:"kernel" yaml:"kernel"`
	Architecture string        `json:"architecture" yaml:"architecture"`
	CPU          CPUInfo       `json:"cpu" yaml:"cpu"`
	Memory       MemoryInfo    `json:"memory" yaml:"memory"`
	Storage      []StorageInfo `json:"storage" yaml:"storage"`
	GPU          []GPUInfo     `json:"gpu" yaml:"gpu"`
	NPU          *NPUInfo      `json:"npu" yaml:"npu"`
}

// CPUInfo represents CPU information.
type CPUInfo struct {
	Model        string  `json:"model" yaml:"model"`
	Vendor       string  `json:"vendor" yaml:"vendor"`
	Cores        int32   `json:"cores" yaml:"cores"`
	FrequencyMHz float64 `json:"frequency_mhz" yaml:"frequency_mhz"`
}

// MemoryInfo represents memory and swap information.
type MemoryInfo struct {
	TotalMB     uint64  `json:"total_mb" yaml:"total_mb"`
	AvailableMB uint64  `json:"available_mb" yaml:"available_mb"`
	UsedMB      uint64  `json:"used_mb" yaml:"used_mb"`
	UsedPercent float64 `json:"used_percent" yaml:"used_percent"`
	SwapTotalMB uint64  `json:"swap_total_mb" yaml:"swap_total_mb"`
	SwapUsedMB  uint64  `json:"swap_used_mb" yaml:"swap_used_mb"`
}

// StorageInfo represents storage volume information.
type StorageInfo struct {
	Device      string  `json:"device" yaml:"device"`
	Mountpoint  string  `json:"mountpoint" yaml:"mountpoint"`
	Filesystem  string  `json:"filesystem" yaml:"filesystem"`
	TotalMB     uint64  `json:"total_mb" yaml:"total_mb"`
	UsedMB      uint64  `json:"used_mb" yaml:"used_mb"`
	FreeMB      uint64  `json:"free_mb" yaml:"free_mb"`
	UsedPercent float64 `json:"used_percent" yaml:"used_percent"`
}

// GPUInfo represents GPU information.
type GPUInfo struct {
	Vendor string `json:"vendor" yaml:"vendor"`
	Model  string `json:"model" yaml:"model"`
	Type   string `json:"type" yaml:"type"` // integrated, discrete, unknown
}

// NPUInfo represents NPU (Neural Processing Unit) information.
type NPUInfo struct {
	Detected        bool   `json:"detected" yaml:"detected"`
	Type            string `json:"type" yaml:"type"`              // Apple Neural Engine, Intel AI Boost, AMD Ryzen AI, unknown
	InferenceMethod string `json:"inference_method" yaml:"inference_method"` // cpu_model, platform_api, unknown
}

// CollectSystemInfo gathers system diagnostic information.
// Returns partial information if some detection fails (graceful degradation).
// Detection failures are logged via slog at debug level.
// Never returns an error (diagnostic tool, not validation tool).
//
// Zero values indicate "unknown" or "not detected":
// - Cores: 0 = unknown
// - FrequencyMHz: 0.0 = unknown (common on Apple Silicon)
// - TotalMB/UsedMB: 0 = detection failed
func CollectSystemInfo(ctx context.Context) SystemInfo {
	info := SystemInfo{
		OS:           "unknown",
		Platform:     "unknown",
		Kernel:       "unknown",
		Architecture: "unknown",
		CPU: CPUInfo{
			Model:  "unknown",
			Vendor: "unknown",
			Cores:  0,
		},
		Memory:  MemoryInfo{},
		Storage: []StorageInfo{},
		GPU:     []GPUInfo{},
	}

	// OS and host info (graceful degradation)
	if hostInfo, err := host.InfoWithContext(ctx); err == nil {
		info.OS = hostInfo.OS
		info.Platform = hostInfo.Platform + " " + hostInfo.PlatformVersion
		info.Kernel = hostInfo.KernelVersion
		info.Architecture = hostInfo.KernelArch
	} else {
		slog.DebugContext(ctx, "Host info detection failed", "error", err)
	}

	// CPU info (graceful degradation)
	if cpuInfos, err := cpu.InfoWithContext(ctx); err == nil && len(cpuInfos) > 0 {
		first := cpuInfos[0]
		info.CPU = CPUInfo{
			Model:        first.ModelName,
			Vendor:       first.VendorID,
			Cores:        int32(first.Cores),
			FrequencyMHz: first.Mhz,
		}
	} else if err != nil {
		slog.DebugContext(ctx, "CPU detection failed", "error", err)
	}

	// Memory info (graceful degradation)
	if memInfo, err := mem.VirtualMemoryWithContext(ctx); err == nil {
		info.Memory.TotalMB = memInfo.Total / 1024 / 1024
		info.Memory.AvailableMB = memInfo.Available / 1024 / 1024
		info.Memory.UsedMB = memInfo.Used / 1024 / 1024
		info.Memory.UsedPercent = memInfo.UsedPercent
	} else {
		slog.DebugContext(ctx, "Memory detection failed", "error", err)
	}

	// Swap info (graceful degradation)
	if swapInfo, err := mem.SwapMemoryWithContext(ctx); err == nil {
		info.Memory.SwapTotalMB = swapInfo.Total / 1024 / 1024
		info.Memory.SwapUsedMB = swapInfo.Used / 1024 / 1024
	} else {
		slog.DebugContext(ctx, "Swap detection failed", "error", err)
	}

	// Storage info (graceful degradation)
	if partitions, err := disk.PartitionsWithContext(ctx, false); err == nil {
		// Filter out pseudo-filesystems (Linux /proc, /sys, etc.)
		skipFsTypes := map[string]bool{
			"sysfs": true, "proc": true, "devtmpfs": true, "tmpfs": true,
			"devpts": true, "cgroup": true, "cgroup2": true, "overlay": true,
		}

		for _, partition := range partitions {
			// Skip pseudo-filesystems
			if skipFsTypes[partition.Fstype] {
				continue
			}

			if usage, err := disk.UsageWithContext(ctx, partition.Mountpoint); err == nil {
				info.Storage = append(info.Storage, StorageInfo{
					Device:      partition.Device,
					Mountpoint:  partition.Mountpoint,
					Filesystem:  partition.Fstype,
					TotalMB:     usage.Total / 1024 / 1024,
					UsedMB:      usage.Used / 1024 / 1024,
					FreeMB:      usage.Free / 1024 / 1024,
					UsedPercent: usage.UsedPercent,
				})
			}
		}
	} else {
		slog.DebugContext(ctx, "Storage detection failed", "error", err)
	}

	// Phase 2: GPU detection (best-effort)
	info.GPU = detectGPU(ctx)

	// Phase 3: NPU detection (best-effort, CPU model-based inference)
	info.NPU = detectNPU(ctx, info.CPU.Model, info.OS)

	return info
}

// detectGPU attempts to detect GPU information using hardware-level detection.
// Returns empty slice if detection fails (graceful degradation).
// Logs detection failures via slog at debug level.
//
// Phase 2 implementation: Cross-platform GPU detection using ghw.
// Detects NVIDIA, AMD, Intel, Apple, and other GPUs on Linux, Windows, and macOS.
func detectGPU(ctx context.Context) []GPUInfo {
	gpus := []GPUInfo{}

	// Use ghw for hardware-level GPU detection
	gpu, err := ghw.GPU()
	if err != nil {
		slog.DebugContext(ctx, "GPU detection failed", "error", err)
		return gpus
	}

	if gpu == nil || len(gpu.GraphicsCards) == 0 {
		slog.DebugContext(ctx, "No GPUs detected")
		return gpus
	}

	for _, card := range gpu.GraphicsCards {
		if card.DeviceInfo == nil {
			continue
		}

		// Determine GPU vendor from device info
		vendor := "Unknown"
		gpuType := "unknown"

		// Normalize vendor name
		vendorLower := strings.ToLower(card.DeviceInfo.Vendor.Name)
		if strings.Contains(vendorLower, "nvidia") {
			vendor = "NVIDIA"
			gpuType = "discrete"
		} else if strings.Contains(vendorLower, "amd") || strings.Contains(vendorLower, "advanced micro devices") {
			vendor = "AMD"
			gpuType = "discrete"
		} else if strings.Contains(vendorLower, "intel") {
			vendor = "Intel"
			// Intel GPUs can be integrated or discrete
			if strings.Contains(strings.ToLower(card.DeviceInfo.Product.Name), "arc") {
				gpuType = "discrete"
			} else {
				gpuType = "integrated"
			}
		} else if strings.Contains(vendorLower, "apple") {
			vendor = "Apple"
			gpuType = "integrated"
		} else {
			vendor = card.DeviceInfo.Vendor.Name
		}

		model := card.DeviceInfo.Product.Name
		if model == "" {
			model = "Unknown Model"
		}

		gpus = append(gpus, GPUInfo{
			Vendor: vendor,
			Model:  model,
			Type:   gpuType,
		})

		slog.DebugContext(ctx, "Detected GPU", "vendor", vendor, "model", model, "type", gpuType)
	}

	return gpus
}

// detectNPU attempts to infer NPU presence from CPU model.
// Returns nil if NPU not detected (graceful degradation).
// Logs detection attempts via slog at debug level.
//
// Phase 3 implementation: Keyword-based NPU detection from CPU model.
// Supports:
//   - Apple Neural Engine (M1, M2, M3, M4 series)
//   - Intel AI Boost (Core Ultra series)
//   - AMD Ryzen AI (Ryzen AI series)
func detectNPU(ctx context.Context, cpuModel, os string) *NPUInfo {
	cpuLower := strings.ToLower(cpuModel)

	// Apple Silicon: M1, M2, M3, M4 → Apple Neural Engine
	if strings.Contains(cpuLower, "apple m1") ||
		strings.Contains(cpuLower, "apple m2") ||
		strings.Contains(cpuLower, "apple m3") ||
		strings.Contains(cpuLower, "apple m4") {
		slog.DebugContext(ctx, "Detected Apple Neural Engine", "cpu_model", cpuModel)
		return &NPUInfo{
			Detected:        true,
			Type:            "Apple Neural Engine",
			InferenceMethod: "cpu_model",
		}
	}

	// Intel Core Ultra: "Ultra" → Intel AI Boost
	if strings.Contains(cpuLower, "intel") && strings.Contains(cpuLower, "ultra") {
		slog.DebugContext(ctx, "Detected Intel AI Boost", "cpu_model", cpuModel)
		return &NPUInfo{
			Detected:        true,
			Type:            "Intel AI Boost",
			InferenceMethod: "cpu_model",
		}
	}

	// AMD Ryzen AI: "Ryzen AI" or specific AI models
	if strings.Contains(cpuLower, "ryzen") && strings.Contains(cpuLower, "ai") {
		slog.DebugContext(ctx, "Detected AMD Ryzen AI", "cpu_model", cpuModel)
		return &NPUInfo{
			Detected:        true,
			Type:            "AMD Ryzen AI",
			InferenceMethod: "cpu_model",
		}
	}

	// No NPU detected
	slog.DebugContext(ctx, "No NPU detected", "cpu_model", cpuModel, "os", os)
	return nil
}
