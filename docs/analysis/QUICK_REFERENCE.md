# Cross-Platform System Information - Quick Reference Card

## Feature Availability at a Glance

```
Feature              | Linux | macOS | Windows | Tools/Methods
---------------------|-------|-------|---------|-----------------------------------
CPU Model Name       |  ✓    |  ✓    |  ✓      | /proc/cpuinfo, sysctl, WMI
CPU Logical Cores    |  ✓    |  ✓    |  ✓      | runtime.NumCPU()
CPU Physical Cores   |  ✓    |  ✓    |  ✓      | /proc/cpuinfo, sysctl
CPU Frequency        |  ✓    |  ✓    |  ✓      | sysctl, WMI
CPU Features (AVX)   |  ✓    |  ✓    |  ✓      | /proc/cpuinfo, golang.org/x/sys/cpu
Total Memory         |  ✓    |  ✓    |  ✓      | syscall.Sysinfo, sysctl
Available Memory     |  ✓    |  ✓    |  ✓      | /proc/meminfo, vm_stat, WMI
Memory Swap          |  ✓    |  ✗    |  ✗      | /proc/meminfo only
Storage Total        |  ✓    |  ✓    |  ✓      | statfs syscall
Storage Used         |  ✓    |  ✓    |  ✓      | statfs syscall
Mount Points         |  ✓    |  ✓    |  ✓      | /proc/mounts, mount, WMI
Filesystem Type      |  ✓    |  ✓    |  ✓      | statfs f_type
SSD Detection        |  ✓*   |  ✓*   |  ✓*     | /sys/block/*/queue/rotational
OS Version           |  ✓    |  ✓    |  ✓      | /etc/os-release, sw_vers, Registry
Kernel Version       |  ✓    |  ✓    |  ✓      | uname, sysctl, API
Architecture         |  ✓    |  ✓    |  ✓      | runtime.GOARCH
Uptime               |  ✓    |  ✓    |  ✓      | /proc/uptime, uptime cmd, WMI
---------------------|-------|-------|---------|-----------------------------------
NVIDIA GPU           |  ✓    |  ✓    |  ✓      | nvidia-smi (if installed)
AMD GPU              |  ✓    |  ✗    |  ✓      | rocm-smi (if installed)
Intel GPU            |  ✓    |  ✓    |  ✓      | system_profiler, tools
Apple Metal GPU      |  ✗    |  ✓    |  ✗      | system_profiler, IOKit
Apple Neural Engine  |  ✗    |  ✓**  |  ✗      | Detection only (M1+)
Intel AI Boost       |  ✓**  |  ✓**  |  ✓**    | CPU model detection only
---------------------|-------|-------|---------|-----------------------------------

Legend:
✓   = Available and reliable
✗   = Not available
✓*  = Available but unreliable (SSD detection varies)
✓** = Detection only (no detailed metrics available)
```

## No Dependencies Required? YES!

**For Phase 1 (MVP):**
- runtime (Go stdlib) - CPU cores, architecture
- golang.org/x/sys/unix - Memory and storage via syscall
- os/exec - Shell commands for system info

No external Go packages needed beyond what's already in project!

## Go Code Patterns

### Always Available Info (Never Fail)
```go
cpu.LogicalCores = runtime.NumCPU()          // Always works
memory.Total = getPhysicalMemoryBytes()      // Use syscall fallback
storage := enumerateMountPoints()            // Use statfs
```

### Optional Info (Graceful Degradation)
```go
// Try GPU detection, don't fail if it doesn't work
if gpu, err := detectGPU(); err == nil {
    info.GPU = &gpu
    info.Capabilities.GPU = true
}
```

### Output Flexibility
```go
// Always output capability flags
// Tools can check what was actually collected
type Capabilities struct {
    GPU      bool `json:"gpu"`       // Was GPU detection attempted?
    NPU      bool `json:"npu"`       // Was NPU detection attempted?
    SmartData bool `json:"smart"`    // Can we read SMART data?
}
```

## Platform-Specific Entry Points

### Core Info Collection (Platform-Agnostic)
```
collectOSInfo()     → Works on all platforms
collectCPUInfo()    → Works on all platforms (with platform-specific enhancements)
collectMemoryInfo() → Works on all platforms
collectStorageInfo()→ Works on all platforms
```

### Optional Info Collection
```
collectGPUInfo()    → Optional (tools may not exist)
collectNPUInfo()    → Optional (NPU may not exist)
```

## Testing Approach

### What to Test Per Platform
```
Linux:
  - /proc/cpuinfo parsing
  - /proc/meminfo parsing
  - /proc/mounts parsing
  - nvidia-smi output parsing
  - rocm-smi output parsing

macOS:
  - sysctl command execution
  - system_profiler parsing
  - Darwin-specific code paths
  - Apple Silicon (M1+) detection

Windows:
  - WMI queries (or file-based fallback)
  - Registry reading
  - Path handling (\\ vs /)
  - admin vs non-admin contexts
```

### Graceful Degradation Tests
```go
// Test 1: Missing tools don't break collection
if err := testMissingNvidiaSmi(); err != nil {
    t.Fatal("Should not fail when nvidia-smi missing")
}

// Test 2: Parsing errors don't break collection
if err := testCorruptedProcCpuinfo(); err != nil {
    t.Fatal("Should not fail on malformed /proc/cpuinfo")
}

// Test 3: Core info always present
info := collector.Collect(ctx)
if info.CPU.LogicalCores == 0 {
    t.Fatal("CPU info should always be available")
}
```

## Security Checklist

```
NO Special Privileges Needed For:
  CPU Info              ✓
  Memory Stats          ✓
  Storage Stats         ✓
  OS Version            ✓
  GPU Detection         ✓ (nvidia-smi, rocm-smi work unprivileged)
  NPU Detection         ✓

MAY Need Elevated Privileges For:
  SMART Data            (smartctl requires sudo)
  Thermal Info          (varies by platform)
  Detailed Memory Info  (some metrics restricted)

NEVER Request Elevation - Just Gracefully Skip
```

## Implementation Priority

### Week 1 (Must Have)
- OS info (name, version, kernel, arch)
- CPU info (model, cores, frequency)
- Memory info (total, available, used)
- Storage info (mount points, usage)

### Week 2 (Should Have)
- GPU detection (nvidia-smi, rocm-smi, system_profiler)
- NPU detection (M1+, Meteor Lake+)

### Week 3+ (Nice to Have)
- SSD/HDD detection
- Uptime
- SMART data
- Thermal info

## File Structure

```
internal/system/
├── info.go              (Data structures)
├── collector.go         (Main logic)
├── os_*.go             (OS-specific implementations)
├── cpu_*.go            (CPU-specific implementations)
├── memory_*.go         (Memory-specific implementations)
├── storage_*.go        (Storage-specific implementations)
├── gpu_*.go            (GPU detection)
├── npu_*.go            (NPU detection)
└── format.go           (Text/JSON/YAML output)

cmd/ado/meta/
└── system.go           (Command wiring)
```

## Expected Output

### Text (Default)
```
System Information
==================================================

OS
  Name: macOS
  Version: 14.1
  Kernel: Darwin 23.1.0
  Architecture: arm64

CPU
  Model: Apple M3 Max
  Logical Cores: 12
  Physical Cores: 8
  Max Frequency: 4.2 GHz

Memory
  Total: 36.0 GiB
  Available: 12.5 GiB
  Used: 23.5 GiB (65.3%)

Storage
  / (APFS, SSD)
    Total: 494.4 GiB
    Used: 250.0 GiB (50.6%)
    Available: 244.4 GiB

GPU
  [0] Apple Integrated GPU
      Memory: Shared (system RAM)

NPU (Neural Processing Unit)
  Type: Apple Neural Engine
  Available: true

Capabilities
  GPU Details: true
  NPU Info: true
  SMART Data: false
```

### JSON (--output json)
```json
{
  "os": { "name": "macOS", "version": "14.1", ... },
  "cpu": { "model": "Apple M3 Max", "logical_cores": 12, ... },
  "memory": { "total_bytes": 38654705664, ... },
  "storage": [ { "mountpoint": "/", "device": "/dev/disk0s2", ... } ],
  "gpu": { "devices": [ { "vendor": "Apple", "model": "Integrated GPU", ... } ] },
  "npu": { "available": true, "type": "Apple Neural Engine" },
  "capabilities": { "cpu": true, "memory": true, "storage": true, ... }
}
```

## Common Pitfalls to Avoid

1. **Assuming GPU exists** - Always check capability flags
2. **Using gopsutil** - Too heavy; use syscalls instead
3. **Requesting sudo** - Never do this; gracefully degrade
4. **Ignoring platform differences** - Use build tags and platform-specific files
5. **Not testing graceful degradation** - Mock failing tools in tests
6. **Hardcoding paths** - Use runtime detection or standard locations
7. **Forgetting Windows paths** - Use filepath.Join, not string concatenation
8. **Ignoring sandboxed environments** - Some metrics unavailable in containers/WSL

## Key Decision: Shell Commands Over Libraries

### Why Use `nvidia-smi` Instead of NVIDIA Go Bindings?
- Already installed on target systems
- No additional Go dependencies
- Works in more environments
- Easier to test (mock command output)
- More portable

### Pattern
```go
// Good: Shell command wrapper
func detectNVIDIA() ([]GPU, error) {
    cmd := exec.Command("nvidia-smi", "--query-gpu=...")
    // Parse output
}

// Avoid: Heavy library dependency
import "github.com/NVIDIA/gpu-monitoring-tools/..."
```

## Contact Point with Existing Code

### Updates Required:
1. `cmd/ado/meta/meta.go` - Add `newSystemCommand()` and register in `NewCommand()`
2. `docs/commands/03-meta.md` - Add `ado meta system` section
3. `CLAUDE.md` - Add new command to Quick Commands

### No Changes Needed:
- go.mod (no new external dependencies)
- internal/ui (output formatting already generic)
- internal/meta (separate from system info)

---

**Total Analysis Pages:** 3 documents (~4000 lines)
**Ready to Implement:** Yes (Phase 1 can start immediately)
**Dependencies:** Zero external Go packages
**Estimated Implementation Time:** 3-4 weeks (Phase 1+2+3)
