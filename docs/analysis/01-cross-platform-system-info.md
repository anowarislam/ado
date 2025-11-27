# Cross-Platform System Information Analysis for `ado meta system`

## 1. Feature Availability Matrix

### CPU Information
| Feature | Linux | macOS | Windows | Notes |
|---------|-------|-------|---------|-------|
| CPU Count | ✓ | ✓ | ✓ | `runtime.NumCPU()` - Go stdlib |
| CPU Model Name | ✓ | ✓ | ✓ | Available via stdlib; macOS needs `sysctl hw.model` |
| CPU Frequency (current) | ✓ | ✓ | ✓ | Linux: `/proc/cpuinfo`; macOS: `sysctl hw.cpufrequency`; Windows: WMI |
| CPU Temperature | ✓ (partial) | ✓ (partial) | ✓ (partial) | Requires system-specific access; may need privileges |
| CPU Cores (logical/physical) | ✓ | ✓ | ✓ | Go stdlib `runtime.NumCPU()`; physical cores require OS calls |
| CPU Features (AVX, SSE, etc) | ✓ | ✓ | ✓ | `golang.org/x/sys/cpu` - built-in intrinsics |
| Hypervisor detection | ✓ | ✓ | ✓ | Possible but unreliable; hypervisor-specific |

### GPU Information
| Feature | Linux | macOS | Windows | Notes |
|---------|-------|-------|---------|-------|
| GPU Count | ✓ (NVIDIA/AMD) | ✓ (Metal) | ✓ (DXGI) | No universal Go stdlib; vendor-specific |
| GPU Model | ✓ (NVIDIA/AMD) | ✓ (Metal) | ✓ (DXGI) | Requires specific APIs per GPU type |
| VRAM | ✓ (NVIDIA/AMD) | ✓ (Metal) | ✓ (DXGI) | No universal access; may need tools |
| GPU Driver Version | ✓ (NVIDIA/AMD) | ✓ (IOKit) | ✓ (Windows Registry) | OS-specific; tools often required |
| Apple Silicon GPU | ✗ | ✓ | ✗ | macOS-only; Metal API + IOKit |
| NVIDIA GPU | ✓ | ✓ (CUDA) | ✓ | Requires nvidia-smi or CUDA toolkit |
| AMD GPU | ✓ | ✗ | ✓ | ROCm or vendor tools required |
| Intel GPU | ✓ | ✓ | ✓ | Limited support; Intel Arc/iGPU |

### NPU (Neural Processing Unit) Information
| Feature | Linux | macOS | Windows | Notes |
|---------|-------|-------|---------|-------|
| Apple Neural Engine | ✗ | ✓ (M1+) | ✗ | macOS-only; no official public API; private frameworks |
| Intel AI Boost (Meteor Lake+) | ✓ (partial) | ? | ✓ (partial) | Limited public documentation; IPEX recommended |
| Qualcomm Hexagon (Android) | ✗ | ✗ | ✗ | Not relevant for desktop/server; would be Android |
| MediaTek APU | ✗ | ✗ | ✗ | Not relevant for desktop/server |
| Generic NPU detection | ✗ | ✗ | ✗ | No universal mechanism; vendor-specific |

### Memory Statistics
| Feature | Linux | macOS | Windows | Notes |
|---------|-------|-------|---------|-------|
| Total Physical RAM | ✓ | ✓ | ✓ | `syscall.Sysconf` (Linux); `sysctlbyname` (macOS); WMI (Windows) |
| Available Memory | ✓ | ✓ | ✓ | `/proc/meminfo` (Linux); `vm_stat` (macOS); WMI (Windows) |
| Used Memory | ✓ | ✓ | ✓ | Derived from above |
| Memory Swap | ✓ | ✗ | ✗ | Linux only; macOS uses compressed memory; Windows uses pagefile |
| Memory Pressure | ✓ | ✓ | ✓ | Linux: % utilization; macOS: `memory_pressure` metric; Windows: WMI |
| Virtual Memory/Pagefile Size | ✓ | ✓ | ✓ | OS-specific metrics |
| NUMA info | ✓ | ✗ | ✗ | Linux server systems only |

### Storage/Volume Information
| Feature | Linux | macOS | Windows | Notes |
|---------|-------|-------|---------|-------|
| Total Disk Size | ✓ | ✓ | ✓ | `syscall.Statfs` (Unix); WMI (Windows) |
| Available Disk Space | ✓ | ✓ | ✓ | `syscall.Statfs` (Unix); WMI (Windows) |
| Used Disk Space | ✓ | ✓ | ✓ | Derived from above |
| Filesystem Type | ✓ | ✓ | ✓ | `statfs` f_type field (Linux); FSVolumeInfo (macOS); WMI (Windows) |
| Mount Points | ✓ | ✓ | ✓ | `/proc/mounts` (Linux); `mount` output (macOS); WMI (Windows) |
| Disk I/O Stats | ✓ | ✓ (partial) | ✓ | `/proc/diskstats` (Linux); iostat (macOS); WMI (Windows) |
| SSD vs HDD Detection | ✓ (partial) | ✓ (partial) | ✓ | `/sys/block/*/queue/rotational` (Linux); `VENDOR_ID` (macOS); WMI (Windows) |
| SMART Status | ✓ (partial) | ✓ (partial) | ✓ (partial) | Requires tool access (smartctl); may need sudo |

### OS Version Details
| Feature | Linux | macOS | Windows | Notes |
|---------|-------|-------|---------|-------|
| OS Name | ✓ | ✓ | ✓ | `runtime.GOOS` (basic); detailed via OS-specific |
| OS Version/Release | ✓ | ✓ | ✓ | `/etc/os-release` (Linux); `sw_vers` (macOS); `GetVersionEx` (Windows) |
| Kernel Version | ✓ | ✓ | ✓ | `uname -r` (Unix); `uname -v` (macOS); `RtlGetVersion` (Windows) |
| Linux Distribution | ✓ | ✗ | ✗ | `/etc/os-release`, `/etc/lsb-release`, etc. |
| macOS Version/Build | ✗ | ✓ | ✗ | `sw_vers` or `system_profiler` |
| Windows Edition/Build | ✗ | ✗ | ✓ | `RtlGetVersion`, Windows Registry |
| Architecture | ✓ | ✓ | ✓ | `runtime.GOARCH` |
| Uptime | ✓ | ✓ | ✓ | `/proc/uptime` (Linux); `uptime` (macOS); WMI (Windows) |

---

## 2. Platform-Specific Implementation Approaches

### CPU Information
**Universal (Go stdlib):**
```go
runtime.NumCPU()              // Logical CPU count
runtime.GOARCH               // CPU architecture
```

**Linux-specific:**
- `/proc/cpuinfo` - model name, flags, frequency
- `sysconf(_SC_PHYS_PAGES)` - physical memory
- `sysconf(_SC_CLK_TCK)` - clock ticks

**macOS-specific:**
- `sysctl hw.model` - CPU model
- `sysctl hw.cpufrequency_max` - max frequency
- `sysctl hw.physicalcpu` - physical CPU count
- `sysctl hw.logicalcpu` - logical CPU count

**Windows-specific:**
- WMI `Win32_Processor`
- Registry `HKEY_LOCAL_MACHINE\HARDWARE\DESCRIPTION\System\CentralProcessor`

### GPU Information
**Option 1: Shell Command Wrappers (Recommended for broad compatibility)**
```
NVIDIA:  nvidia-smi --query-gpu=... --format=csv
AMD:     rocm-smi --showid --showtemp
Intel:   Intel Graphics Monitor (Windows) / ioreg (macOS)
Apple:   system_profiler SPDisplaysDataType
```

**Option 2: Native API Calls (Most reliable but OS-specific)**
- **Linux (NVIDIA):** NVIDIA Management Library (NVML)
- **Linux (AMD):** ROCm tools or read `/sys/class/kfd/devices/`
- **macOS (Metal):** Metal Framework via IOKit
- **macOS (Apple Silicon):** IOKit (private frameworks; undocumented)
- **Windows (DXGI):** Direct3D/DXGI
- **Windows (NVIDIA):** NVIDIA NVML

**Practical approach:** Try tool-based detection first (nvidia-smi, rocm-smi, etc.), fall back gracefully if unavailable.

### NPU Information
**Apple Neural Engine (macOS M1+):**
- No public API available
- Private frameworks: `com.apple.CoreML`, `com.apple.foundation`
- Detection: System Profiler shows "Neural Engine"
- Fallback: `/usr/sbin/system_profiler SPHardwareDataType`

**Intel AI Boost (Meteor Lake+):**
- Requires IPEX (Intel Python Extension for PyTorch)
- Limited Linux/Windows documentation
- Detection: CPU model name contains "Meteor Lake"
- Fallback: Model detection + flag checking

**General approach:** Platform-specific detection + graceful "not available" response

### Memory Statistics
**Universal (Go stdlib):**
```go
import "golang.org/x/sys/unix"
var info unix.Sysinfo_t
unix.Sysinfo(&info)
// info.Totalram, info.Freeram available
```

**Linux:** `/proc/meminfo` for detailed breakdown
**macOS:** `vm_stat`, `sysctlbyname("hw.memsize")`
**Windows:** WMI `Win32_LogicalMemoryConfiguration`, `Win32_OperatingSystem`

### Storage Information
**Universal (Go stdlib):**
```go
import "golang.org/x/sys/unix"
var stat unix.Statfs_t
unix.Statfs(path, &stat)
// stat.Blocks, stat.Bavail, stat.Namemax available
```

**Linux:**
- `/proc/mounts` for mount points
- `/sys/block/*/queue/rotational` for SSD detection
- `/proc/diskstats` for I/O metrics

**macOS:**
- `mount` output for mount points
- `diskutil info` or FSVolumeInfo for metadata
- System Profiler for SSD detection

**Windows:**
- WMI `Win32_LogicalDisk`, `Win32_Volume`
- Registry for mount points
- WMI `Win32_DiskDrive` for physical drive info

### OS Version Details
**Universal (Go stdlib):**
```go
runtime.GOOS   // "linux", "darwin", "windows"
runtime.GOARCH // "amd64", "arm64", etc.
```

**Linux:**
- `/etc/os-release` (modern standard)
- `/etc/lsb-release` (Ubuntu/Debian fallback)
- `/etc/redhat-release` (RHEL-based fallback)
- `uname -r` for kernel version

**macOS:**
- `sw_vers` command
- `system_profiler SPSoftwareDataType`
- Registry key: `com.apple.system.version.components.ProductVersion`

**Windows:**
- `RtlGetVersion` API (most reliable)
- Registry: `HKEY_LOCAL_MACHINE\SOFTWARE\Microsoft\Windows NT\CurrentVersion`
- WMI `Win32_OperatingSystem`

---

## 3. Graceful Degradation Strategy

### Recommended Approach: Capability-Based Feature Flags

```go
type SystemInfo struct {
    OS         OSInfo         `json:"os"`
    CPU        CPUInfo        `json:"cpu"`
    GPU        *GPUInfo       `json:"gpu,omitempty"`           // Nullable
    NPU        *NPUInfo       `json:"npu,omitempty"`           // Nullable
    Memory     MemoryInfo     `json:"memory"`
    Storage    []StorageInfo  `json:"storage"`
    Capabilities CapabilitySet `json:"capabilities"`           // Track what was collected
}

type CapabilitySet struct {
    CPU              bool `json:"cpu"`              // Always true
    Memory           bool `json:"memory"`           // Always true
    Storage          bool `json:"storage"`          // Always true
    OS               bool `json:"os"`               // Always true
    GPU              bool `json:"gpu"`              // False if no GPU detected/accessible
    NPU              bool `json:"npu"`              // False if not available on platform
    GPUDetails       bool `json:"gpu_details"`      // False if GPU info incomplete
    GPUTemperature   bool `json:"gpu_temperature"` // False if tools not available
    StorageSmartData bool `json:"storage_smart"`   // False if smartctl not available
}
```

### Error Handling Pattern
```go
// Don't fail the entire command if one feature is unavailable
// Instead, return partial data + capability flags

func collectSystemInfo(ctx context.Context) (SystemInfo, error) {
    info := SystemInfo{}
    
    // Always-available info (non-nullable)
    info.OS = collectOSInfo()       // Never fails
    info.CPU = collectCPUInfo()     // Never fails
    info.Memory = collectMemoryInfo() // May have partial data
    info.Storage = collectStorageInfo() // Never fails
    
    // Optional info (graceful degradation)
    gpu, err := collectGPUInfo(ctx)
    if err == nil {
        info.GPU = &gpu
        info.Capabilities.GPU = true
    }
    
    npu, err := collectNPUInfo(ctx)
    if err == nil {
        info.NPU = &npu
        info.Capabilities.NPU = true
    }
    
    return info, nil // Return what we could collect
}
```

### Text Output with Graceful Degradation
```
System Information
==================

OS
  Name: macOS
  Version: 14.1 (Sonoma)
  Kernel: Darwin 23.1.0
  Arch: arm64

CPU
  Model: Apple M3 Max
  Cores: 12 logical, 8 performance / 4 efficiency
  Max Frequency: 4.2 GHz
  Features: AVX, AVX2, SVE

Memory
  Total: 36.0 GiB
  Available: 12.5 GiB
  Used: 23.5 GiB

Storage
  /dev/disk0s2 (APFS)
    Mount: /
    Total: 494.4 GiB
    Used: 250.0 GiB
    Available: 244.4 GiB
    Type: SSD

GPU
  Model: Apple M3 Max (Integrated)
  Memory: Not available (shared system memory)
  [NOTE: Detailed GPU metrics require system-profiler access]

NPU
  Apple Neural Engine
  Status: Available (M3+)
  [NOTE: Detailed NPU metrics not publicly available]
```

---

## 4. Security & Permission Requirements

### CPU Information
- **Linux:** No special privileges needed
- **macOS:** No special privileges needed
- **Windows:** No special privileges needed
- **Exception:** Reading `/proc/cpuinfo` requires read access (usually available to all users)

### GPU Information
**NVIDIA:**
- `nvidia-smi` output: Readable without sudo on most systems
- Direct NVML API: May require additional privileges
- CUDA toolkit: Usually requires installation but not elevated privileges to query
- Recommendation: Call `nvidia-smi` if available; document that it may not be accessible in sandboxed environments

**AMD (ROCm):**
- `rocm-smi` output: May require elevated privileges on some systems
- `/sys/class/kfd/devices/`: Read-only access needed
- Recommendation: Try unprivileged access first, gracefully degrade

**Intel GPU:**
- Limited public APIs; typically embedded in system tools
- No special privileges usually needed

**macOS (Metal/Apple Silicon):**
- IOKit queries: Usually require elevated privileges for detailed metrics
- System Profiler: Available without sudo but data is limited
- Recommendation: Use `system_profiler` (no privileges) as fallback

**Windows (DXGI):**
- Requires Direct3D SDK but not elevated privileges
- WMI queries: Usually accessible to standard users

### NPU Information
**Apple Neural Engine:**
- No public API; reverse-engineered info from system_profiler
- `system_profiler SPHardwareDataType`: Available without sudo
- Recommendation: Use system_profiler; document that info is limited

**Intel AI Boost:**
- IPEX (PyTorch extension): Not applicable to general system query
- Detection via CPUID: No special privileges needed
- Recommendation: Detect via CPU model name; feature info gracefully degraded

### Memory Statistics
- **Linux:** `/proc/meminfo` is world-readable
- **macOS:** `vm_stat` available to all users
- **Windows:** WMI queries usually available to standard users
- **Recommendation:** No elevation needed for basic stats; detailed breakdown (swap, pressure) may have restrictions

### Storage Information
- **Statfs calls:** Available without privileges for mounted filesystems readable by user
- **smartctl (SMART data):** Requires root/admin on all platforms
- **WMI/Registry (Windows):** Standard user access usually sufficient for basic info
- **Recommendation:** Never require sudo/admin for basic stats; optional SMART info gracefully fails

### OS Version Details
- **All platforms:** Available without special privileges
- **Windows Registry:** Standard user reads available for version info
- **Recommendation:** No elevation needed

---

## 5. Implementation Roadmap & Dependencies

### Phase 1: Core (No external dependencies)
- OS version details (GOOS, kernel version)
- CPU info (count, model, architecture)
- Memory stats (basic via syscall)
- Storage info (basic statfs)

**Dependencies:** None beyond Go stdlib + `golang.org/x/sys`

### Phase 2: Enhanced (Optional tool-based)
- GPU detection (nvidia-smi, rocm-smi, system_profiler)
- Extended CPU features (CPUID flags)
- Memory pressure metrics
- Storage type detection (SSD/HDD)

**Dependencies:** External tools (shell commands); gracefully skip if unavailable

### Phase 3: Advanced (Optional specialized libraries)
- Detailed GPU metrics (if library available)
- SMART data collection
- Thermal metrics
- Hypervisor detection

**Dependencies:** Optional libraries with graceful fallback

### Recommended Dependency Strategy:
- **Zero hard dependencies** - only Go stdlib + cobra (already in use)
- **Optional shell command detection** - check if tools exist before using
- **Fail gracefully** - missing tools/info → omit field or set capability flag false

---

## 6. Cross-Platform Testing Checklist

### Linux (x86_64 + ARM)
- [ ] No GPU system
- [ ] NVIDIA GPU (with nvidia-smi)
- [ ] NVIDIA GPU (without nvidia-smi)
- [ ] AMD GPU (with rocm-smi)
- [ ] Intel Arc GPU
- [ ] Multiple GPUs
- [ ] Root vs non-root access differences

### macOS (Intel + Apple Silicon M1/M2/M3+)
- [ ] Intel Mac with integrated GPU
- [ ] Intel Mac with discrete GPU (eGPU)
- [ ] Apple Silicon M1/M2/M3 with Neural Engine
- [ ] Multiple physical disks (Fusion Drive behavior)
- [ ] SSD detection
- [ ] APFS filesystem metrics

### Windows (x86_64 + ARM)
- [ ] No GPU system
- [ ] NVIDIA GPU
- [ ] Intel integrated GPU
- [ ] AMD GPU
- [ ] Multiple GPU vendors
- [ ] Sandboxed/WSL2 context
- [ ] Different disk types (SATA, NVMe, USB)

### Special Scenarios
- [ ] Docker/container detection
- [ ] WSL2 vs native Windows
- [ ] Virtual machine detection (affects GPU info)
- [ ] Limited user permissions (non-admin Windows)
- [ ] Sandboxed environments (macOS App Sandbox)

---

## 7. Recommended File Structure for Implementation

```
internal/
  system/
    system.go         # Main collection logic
    info.go           # Data structures
    
    os_info.go
    os_unix.go        # Linux + macOS shared
    os_darwin.go      # macOS-specific
    os_windows.go     # Windows-specific
    
    cpu_info.go
    cpu_unix.go
    cpu_darwin.go
    cpu_windows.go
    
    gpu_info.go       # Detection + tool wrappers
    memory_info.go
    memory_unix.go
    memory_windows.go
    
    storage_info.go
    storage_unix.go
    storage_windows.go
    
    npu_info.go       # Special handling, always graceful
```

---

## 8. Key Implementation Decisions

### Decision 1: GPU Detection Approach
**Recommendation:** Shell command detection (priority order)
1. `nvidia-smi --query-gpu=name,driver_version,memory.total --format=csv,noheader`
2. `rocm-smi --showid --showtemp --csv`
3. macOS: `system_profiler SPDisplaysDataType`
4. Windows: WMI query (if Go bindings available)

**Rationale:**
- Works across Linux/macOS/Windows
- Gracefully skips if tool unavailable
- No large external dependencies
- Users familiar with these tools

### Decision 2: NPU Information
**Recommendation:** Detection only, no enumeration
- Apple Neural Engine: Detect from CPU model (M1+), document as available
- Intel AI Boost: Detect from CPU model (Meteor Lake+), document as available
- Full capability info: Mark as "requires vendor tools" in text output

**Rationale:**
- No public APIs available
- Avoids private framework reverse-engineering
- Honest about limitations
- Useful for user awareness

### Decision 3: Error Handling Philosophy
**Recommendation:** "Missing ≠ Error"
- Graceful degradation is default
- Capability flags enable structured output to convey what was collected
- Never fail the entire command due to missing GPU/NPU/SMART info
- Return 0 exit code always (unless I/O error)

**Rationale:**
- User can still get useful baseline system info
- Structured output (JSON/YAML) can indicate capabilities
- No surprises in automation/CI contexts

### Decision 4: Tool Dependency vs Library Dependency
**Recommendation:** Prefer tools (shell commands) over Go libraries
- `nvidia-smi` vs NVIDIA Go bindings
- `rocm-smi` vs ROCm Go SDK
- WMI (Windows) vs external WMI library

**Rationale:**
- Tools already installed on target systems
- Simpler dependency management
- Better cross-platform compatibility
- Easier testing (can mock commands)

---

## Summary Table: Recommended Scope

| Feature | Impl | Notes |
|---------|------|-------|
| **OS Version** | Full | Always available |
| **CPU Info** | Full | Count, model, architecture always available |
| **CPU Features** | Partial | Basic info; full details gracefully degraded |
| **Memory** | Full | Basic stats always; pressure metrics where available |
| **Storage** | Full | Mountpoint, size, usage always; SMART optional |
| **GPU Count/Model** | Full (with tools) | Via tool detection; optional without |
| **GPU VRAM** | Partial | Tool-dependent; may fail gracefully |
| **GPU Driver** | Partial | Tool-dependent; may fail gracefully |
| **Apple Neural Engine** | Detection only | Mark "Available (M1+)" with no detailed metrics |
| **Intel AI Boost** | Detection only | Mark "Available (Meteor Lake+)" with no detailed metrics |
| **Thermal/SMART** | Skip v1 | Complex dependencies; can add as optional feature later |

