# meta system Command Spec

| Metadata | Value |
|----------|-------|
| **ADR** | N/A (extends existing pattern) |
| **Status** | Draft |
| **Issue** | #53 |

## Command

```bash
ado meta system [flags]
```

## Purpose

Display comprehensive system diagnostic information (OS, CPU, GPU, NPU, memory, storage) for troubleshooting, bug reports, and requirements validation. Follows the existing `meta` command pattern (`info`, `env`, `features`).

## Usage Examples

```bash
# Example 1: Basic usage (text output)
ado meta system
# Expected output:
# OS: darwin
# Platform: macOS 14.2
# Architecture: arm64
# CPUs: 10 (Apple M2 Pro)
# Memory: 16384 MB total, 8192 MB available (50.0% used)
# Storage:
#   /: 494 GB total, 123 GB used (25.0%)
# GPU:
#   Apple M2 Pro (integrated)
# NPU: Apple Neural Engine (detected)

# Example 2: JSON output (for bug reports)
ado meta system --output json
# Expected output:
# {
#   "os": "darwin",
#   "platform": "macOS 14.2",
#   "kernel": "Darwin 23.2.0",
#   "architecture": "arm64",
#   "cpu": {
#     "model": "Apple M2 Pro",
#     "vendor": "Apple",
#     "cores": 10,
#     "physical_cores": 10,
#     "frequency_mhz": 0
#   },
#   "memory": {
#     "total_mb": 16384,
#     "available_mb": 8192,
#     "used_mb": 8192,
#     "used_percent": 50.0,
#     "swap_total_mb": 0,
#     "swap_used_mb": 0
#   },
#   "storage": [
#     {
#       "device": "/dev/disk3s1s1",
#       "mountpoint": "/",
#       "filesystem": "apfs",
#       "total_gb": 494,
#       "used_gb": 123,
#       "free_gb": 371,
#       "used_percent": 25.0
#     }
#   ],
#   "gpu": [
#     {
#       "vendor": "Apple",
#       "model": "Apple M2 Pro",
#       "type": "integrated"
#     }
#   ],
#   "npu": {
#     "detected": true,
#     "type": "Apple Neural Engine",
#     "inference_method": "cpu_model"
#   }
# }

# Example 3: YAML output
ado meta system --output yaml
# Expected output:
# os: darwin
# platform: macOS 14.2
# kernel: Darwin 23.2.0
# architecture: arm64
# cpu:
#   model: Apple M2 Pro
#   vendor: Apple
#   cores: 10
#   physical_cores: 10
#   frequency_mhz: 0
# memory:
#   total_mb: 16384
#   available_mb: 8192
#   used_mb: 8192
#   used_percent: 50.0
#   swap_total_mb: 0
#   swap_used_mb: 0
# storage:
#   - device: /dev/disk3s1s1
#     mountpoint: /
#     filesystem: apfs
#     total_gb: 494
#     used_gb: 123
#     free_gb: 371
#     used_percent: 25.0
# gpu:
#   - vendor: Apple
#     model: Apple M2 Pro
#     type: integrated
# npu:
#   detected: true
#   type: Apple Neural Engine
#   inference_method: cpu_model

# Example 4: Piping JSON to jq for specific field extraction
ado meta system --output json | jq '.memory.used_percent'
# Expected output: 50.0
```

## Arguments

This command takes no positional arguments.

## Flags

### Command-Specific Flags

| Flag | Short | Type | Default | Description |
|------|-------|------|---------|-------------|
| `--output` | `-o` | enum | `text` | Output format: text, json, yaml |

### Inherited Global Flags

All commands inherit these flags from the root command:

- `--config PATH` - Config file path (default: auto-detected)
- `--log-level LEVEL` - Log level: debug, info, warn, error (default: info)
- `--help, -h` - Show help for command

## Behavior

### Operation Flow

1. **Initialization**: Create default SystemInfo struct with "unknown" values
2. **OS Detection**: Collect OS name, platform, kernel version, architecture
3. **CPU Detection**: Collect CPU model, vendor, core count, frequency
4. **Memory Detection**: Collect total, available, used memory and swap info
5. **Storage Detection**: Enumerate mounted volumes with usage statistics
6. **GPU Detection**: Attempt to detect GPU vendor and model (best-effort)
7. **NPU Detection**: Infer NPU presence from CPU model (best-effort)
8. **Format Output**: Render as text, JSON, or YAML based on `--output` flag
9. **Return**: Always exit 0 (diagnostic tool, not validation tool)

### Graceful Degradation Strategy

**Design Principle**: Always succeed, never fail due to missing system info.

- If OS info unavailable → show `os: "unknown"`
- If CPU info unavailable → show `cpu_model: "unknown"`, `cores: 0`
- If memory info unavailable → show `memory: null` (JSON) or omit section (text)
- If storage info unavailable → show `storage: []` (empty array)
- If GPU not detectable → show `gpu: null` or omit section
- If NPU not detectable → show `npu: null` or `detected: false`

**Exit Codes:**
- Success (always): `0`
- Invalid flag (e.g., `--output xml`): `1`

### Output Formats

#### Text (default)

Human-readable, sectioned output:

```
OS: darwin
Platform: macOS 14.2
Kernel: Darwin 23.2.0
Architecture: arm64

CPU:
  Model: Apple M2 Pro
  Vendor: Apple
  Cores: 10 (10 physical)
  Frequency: unknown

Memory:
  Total: 16384 MB
  Available: 8192 MB
  Used: 8192 MB (50.0%)
  Swap: 0 MB total, 0 MB used

Storage:
  /: 494 GB total, 123 GB used (25.0%)

GPU:
  Apple M2 Pro (integrated)

NPU:
  Type: Apple Neural Engine
  Detection Method: Inferred from CPU model
```

#### JSON (`--output json`)

Structured JSON with stable keys (suitable for parsing):

```json
{
  "os": "darwin",
  "platform": "macOS 14.2",
  "kernel": "Darwin 23.2.0",
  "architecture": "arm64",
  "cpu": {
    "model": "Apple M2 Pro",
    "vendor": "Apple",
    "cores": 10,
    "physical_cores": 10,
    "frequency_mhz": 0
  },
  "memory": {
    "total_mb": 16384,
    "available_mb": 8192,
    "used_mb": 8192,
    "used_percent": 50.0,
    "swap_total_mb": 0,
    "swap_used_mb": 0
  },
  "storage": [
    {
      "device": "/dev/disk3s1s1",
      "mountpoint": "/",
      "filesystem": "apfs",
      "total_gb": 494,
      "used_gb": 123,
      "free_gb": 371,
      "used_percent": 25.0
    }
  ],
  "gpu": [
    {
      "vendor": "Apple",
      "model": "Apple M2 Pro",
      "type": "integrated"
    }
  ],
  "npu": {
    "detected": true,
    "type": "Apple Neural Engine",
    "inference_method": "cpu_model"
  }
}
```

**JSON Schema Notes:**
- All top-level fields always present (never omitted)
- `gpu` and `npu` may be `null` if not detectable
- `storage` is always an array (may be empty)
- Numeric values use appropriate types (int, float64)

#### YAML (`--output yaml`)

Same structure as JSON, YAML format:

```yaml
os: darwin
platform: macOS 14.2
kernel: Darwin 23.2.0
architecture: arm64
cpu:
  model: Apple M2 Pro
  vendor: Apple
  cores: 10
  physical_cores: 10
  frequency_mhz: 0
memory:
  total_mb: 16384
  available_mb: 8192
  used_mb: 8192
  used_percent: 50.0
  swap_total_mb: 0
  swap_used_mb: 0
storage:
  - device: /dev/disk3s1s1
    mountpoint: /
    filesystem: apfs
    total_gb: 494
    used_gb: 123
    free_gb: 371
    used_percent: 25.0
gpu:
  - vendor: Apple
    model: Apple M2 Pro
    type: integrated
npu:
  detected: true
  type: Apple Neural Engine
  inference_method: cpu_model
```

## Error Cases

| Condition | Exit Code | Error Message |
|-----------|-----------|---------------|
| Invalid output format | 1 | `Error: invalid output format "xml". Valid formats: text, json, yaml` |
| Unknown flag | 1 | `Error: unknown flag: --invalid-flag` |

**Non-Errors (Graceful Degradation):**
- GPU not detectable → Show `gpu: null`, exit 0
- NPU not detectable → Show `npu: null` or `detected: false`, exit 0
- Memory info unavailable → Show partial info with "unknown" fields, exit 0

## Implementation

### File Structure

| Purpose | Path |
|---------|------|
| Command registration | `cmd/ado/meta/meta.go` (modify: add `newSystemCommand()`) |
| System info collector | `internal/meta/system.go` (new) |
| Command tests | `cmd/ado/meta/meta_test.go` (modify: add `TestSystemCommand`) |
| Collector tests | `internal/meta/system_test.go` (new) |

### Dependencies

#### Required Dependencies

```bash
go get github.com/shirou/gopsutil/v4@latest
```

**Justification:**
- **gopsutil/v4**: Cross-platform system info library (10.5k stars)
  - Covers: OS, CPU, Memory, Storage
  - No CGO required (easy cross-compilation)
  - Battle-tested in production systems

#### Optional Dependencies (for GPU/NPU enhancement)

```bash
go get github.com/jaypipes/ghw@latest
```

**Justification:**
- **ghw**: GPU detection library (1.8k stars)
  - Only Go library with GPU support
  - No CGO required
  - Graceful fallback if GPU not detectable

**Alternative**: Implement GPU/NPU detection without ghw using platform-specific commands (Linux: `lspci`, macOS: `system_profiler`, Windows: WMI). Trade-off: More code, CGO-free.

**Recommendation**: Start with gopsutil only (covers 80% of use cases), add ghw if GPU detection is critical for initial release.

### Implementation Outline

#### `internal/meta/system.go`

```go
package meta

import (
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
	GPU          []GPUInfo     `json:"gpu,omitempty" yaml:"gpu,omitempty"`
	NPU          *NPUInfo      `json:"npu,omitempty" yaml:"npu,omitempty"`
}

type CPUInfo struct {
	Model         string  `json:"model" yaml:"model"`
	Vendor        string  `json:"vendor" yaml:"vendor"`
	Cores         int     `json:"cores" yaml:"cores"`
	PhysicalCores int     `json:"physical_cores" yaml:"physical_cores"`
	FrequencyMHz  float64 `json:"frequency_mhz" yaml:"frequency_mhz"`
}

type MemoryInfo struct {
	TotalMB       uint64  `json:"total_mb" yaml:"total_mb"`
	AvailableMB   uint64  `json:"available_mb" yaml:"available_mb"`
	UsedMB        uint64  `json:"used_mb" yaml:"used_mb"`
	UsedPercent   float64 `json:"used_percent" yaml:"used_percent"`
	SwapTotalMB   uint64  `json:"swap_total_mb" yaml:"swap_total_mb"`
	SwapUsedMB    uint64  `json:"swap_used_mb" yaml:"swap_used_mb"`
}

type StorageInfo struct {
	Device      string  `json:"device" yaml:"device"`
	Mountpoint  string  `json:"mountpoint" yaml:"mountpoint"`
	Filesystem  string  `json:"filesystem" yaml:"filesystem"`
	TotalGB     uint64  `json:"total_gb" yaml:"total_gb"`
	UsedGB      uint64  `json:"used_gb" yaml:"used_gb"`
	FreeGB      uint64  `json:"free_gb" yaml:"free_gb"`
	UsedPercent float64 `json:"used_percent" yaml:"used_percent"`
}

type GPUInfo struct {
	Vendor string `json:"vendor" yaml:"vendor"`
	Model  string `json:"model" yaml:"model"`
	Type   string `json:"type" yaml:"type"` // integrated, discrete, unknown
}

type NPUInfo struct {
	Detected        bool   `json:"detected" yaml:"detected"`
	Type            string `json:"type" yaml:"type"` // Apple Neural Engine, Intel AI Boost, AMD Ryzen AI, unknown
	InferenceMethod string `json:"inference_method" yaml:"inference_method"` // cpu_model, platform_api, unknown
}

// CollectSystemInfo gathers system diagnostic information.
// Returns partial information if some detection fails (graceful degradation).
// Never returns an error (diagnostic tool, not validation tool).
func CollectSystemInfo() SystemInfo {
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
	if hostInfo, err := host.Info(); err == nil {
		info.OS = hostInfo.OS
		info.Platform = hostInfo.Platform + " " + hostInfo.PlatformVersion
		info.Kernel = hostInfo.KernelVersion
		info.Architecture = hostInfo.KernelArch
	}

	// CPU info (graceful degradation)
	if cpuInfos, err := cpu.Info(); err == nil && len(cpuInfos) > 0 {
		first := cpuInfos[0]
		info.CPU = CPUInfo{
			Model:         first.ModelName,
			Vendor:        first.VendorID,
			Cores:         int(first.Cores),
			PhysicalCores: int(first.Cores), // gopsutil doesn't distinguish, use same value
			FrequencyMHz:  first.Mhz,
		}
	}

	// Memory info (graceful degradation)
	if memInfo, err := mem.VirtualMemory(); err == nil {
		info.Memory.TotalMB = memInfo.Total / 1024 / 1024
		info.Memory.AvailableMB = memInfo.Available / 1024 / 1024
		info.Memory.UsedMB = memInfo.Used / 1024 / 1024
		info.Memory.UsedPercent = memInfo.UsedPercent
	}

	// Swap info (graceful degradation)
	if swapInfo, err := mem.SwapMemory(); err == nil {
		info.Memory.SwapTotalMB = swapInfo.Total / 1024 / 1024
		info.Memory.SwapUsedMB = swapInfo.Used / 1024 / 1024
	}

	// Storage info (graceful degradation)
	if partitions, err := disk.Partitions(false); err == nil {
		for _, partition := range partitions {
			if usage, err := disk.Usage(partition.Mountpoint); err == nil {
				info.Storage = append(info.Storage, StorageInfo{
					Device:      partition.Device,
					Mountpoint:  partition.Mountpoint,
					Filesystem:  partition.Fstype,
					TotalGB:     usage.Total / 1024 / 1024 / 1024,
					UsedGB:      usage.Used / 1024 / 1024 / 1024,
					FreeGB:      usage.Free / 1024 / 1024 / 1024,
					UsedPercent: usage.UsedPercent,
				})
			}
		}
	}

	// GPU detection (best-effort, using ghw if available)
	// Implementation: detectGPU() helper function
	info.GPU = detectGPU()

	// NPU detection (best-effort, CPU model-based inference)
	// Implementation: detectNPU() helper function
	info.NPU = detectNPU(info.CPU.Model, info.OS)

	return info
}

// detectGPU attempts to detect GPU information.
// Returns empty slice if detection fails (graceful degradation).
func detectGPU() []GPUInfo {
	// TODO: Implement using ghw or platform-specific commands
	// For initial implementation, return empty slice
	// Future: Use github.com/jaypipes/ghw for cross-platform GPU detection
	return []GPUInfo{}
}

// detectNPU attempts to infer NPU presence from CPU model.
// Returns nil if NPU not detected (graceful degradation).
func detectNPU(cpuModel, os string) *NPUInfo {
	// Best-effort detection based on CPU model keywords
	// Apple Silicon: M1, M2, M3, M4 → Apple Neural Engine
	// Intel Core Ultra: "Ultra" → Intel AI Boost
	// AMD Ryzen AI: "Ryzen AI" → AMD Ryzen AI

	// TODO: Implement keyword-based detection
	// For initial implementation, return nil
	return nil
}
```

#### `cmd/ado/meta/meta.go` (modifications)

```go
// In NewCommand(), add:
cmd.AddCommand(
	newInfoCommand(buildInfo),
	newEnvCommand(),
	newFeaturesCommand(),
	newSystemCommand(), // NEW
)

// Add new function:
func newSystemCommand() *cobra.Command {
	var output string

	cmd := &cobra.Command{
		Use:   "system",
		Short: "Show system diagnostic information",
		Long: `Display system-level diagnostic information including OS, CPU, GPU, NPU, memory, and storage.

Useful for:
  - Troubleshooting environment-specific issues
  - Sharing system information in bug reports
  - Validating system requirements for ado commands
  - Capturing system state in CI/CD pipelines

Output formats:
  - text (default): Human-readable sectioned output
  - json: Structured JSON for parsing/automation
  - yaml: Structured YAML for parsing/automation

Examples:
  # Show system info in human-readable format
  ado meta system

  # Export as JSON for bug report
  ado meta system --output json

  # Extract specific field with jq
  ado meta system --output json | jq '.memory.used_percent'`,
		RunE: func(cmd *cobra.Command, args []string) error {
			info := internalmeta.CollectSystemInfo()
			format, err := ui.ParseOutputFormat(output)
			if err != nil {
				return err
			}

			return ui.PrintOutput(cmd.OutOrStdout(), format, info, func() (string, error) {
				return formatSystemInfo(info), nil
			})
		},
	}

	cmd.Flags().StringVarP(&output, "output", "o", "text", "Output format: text, json, yaml")
	return cmd
}

func formatSystemInfo(info internalmeta.SystemInfo) string {
	var b strings.Builder

	// OS Section
	fmt.Fprintf(&b, "OS: %s\n", info.OS)
	fmt.Fprintf(&b, "Platform: %s\n", info.Platform)
	fmt.Fprintf(&b, "Kernel: %s\n", info.Kernel)
	fmt.Fprintf(&b, "Architecture: %s\n", info.Architecture)
	fmt.Fprintln(&b)

	// CPU Section
	fmt.Fprintln(&b, "CPU:")
	fmt.Fprintf(&b, "  Model: %s\n", info.CPU.Model)
	fmt.Fprintf(&b, "  Vendor: %s\n", info.CPU.Vendor)
	fmt.Fprintf(&b, "  Cores: %d (%d physical)\n", info.CPU.Cores, info.CPU.PhysicalCores)
	if info.CPU.FrequencyMHz > 0 {
		fmt.Fprintf(&b, "  Frequency: %.0f MHz\n", info.CPU.FrequencyMHz)
	} else {
		fmt.Fprintln(&b, "  Frequency: unknown")
	}
	fmt.Fprintln(&b)

	// Memory Section
	fmt.Fprintln(&b, "Memory:")
	fmt.Fprintf(&b, "  Total: %d MB\n", info.Memory.TotalMB)
	fmt.Fprintf(&b, "  Available: %d MB\n", info.Memory.AvailableMB)
	fmt.Fprintf(&b, "  Used: %d MB (%.1f%%)\n", info.Memory.UsedMB, info.Memory.UsedPercent)
	if info.Memory.SwapTotalMB > 0 {
		fmt.Fprintf(&b, "  Swap: %d MB total, %d MB used\n", info.Memory.SwapTotalMB, info.Memory.SwapUsedMB)
	}
	fmt.Fprintln(&b)

	// Storage Section
	if len(info.Storage) > 0 {
		fmt.Fprintln(&b, "Storage:")
		for _, storage := range info.Storage {
			fmt.Fprintf(&b, "  %s: %d GB total, %d GB used (%.1f%%)\n",
				storage.Mountpoint, storage.TotalGB, storage.UsedGB, storage.UsedPercent)
		}
		fmt.Fprintln(&b)
	}

	// GPU Section
	if len(info.GPU) > 0 {
		fmt.Fprintln(&b, "GPU:")
		for _, gpu := range info.GPU {
			fmt.Fprintf(&b, "  %s %s (%s)\n", gpu.Vendor, gpu.Model, gpu.Type)
		}
		fmt.Fprintln(&b)
	}

	// NPU Section
	if info.NPU != nil && info.NPU.Detected {
		fmt.Fprintln(&b, "NPU:")
		fmt.Fprintf(&b, "  Type: %s\n", info.NPU.Type)
		fmt.Fprintf(&b, "  Detection Method: %s\n", info.NPU.InferenceMethod)
	}

	return b.String()
}
```

### Implementation Notes

1. **Error Handling**: Use graceful degradation pattern from `internal/meta/env.go`:
   - Initialize with defaults
   - Try to collect info, ignore errors
   - Return partial info (never fail)

2. **Testing Strategy**: Table-driven tests with subtests:
   - Test all output formats (text, json, yaml)
   - Test graceful degradation (mock missing system info)
   - Test cross-platform (CI matrix: Linux, macOS, Windows)

3. **Cross-Platform Considerations**:
   - gopsutil handles OS differences automatically
   - GPU detection may fail on some platforms → show `gpu: null`
   - NPU detection is best-effort → show `detected: false` if not found

4. **Performance**: All info collected in &lt;100ms (gopsutil is fast)

5. **Security**: No elevated permissions required (read-only system info)

## Testing Checklist

### Unit Tests

- [ ] `TestSystemCommand` with table-driven tests:
  - [ ] Text output contains expected fields
  - [ ] JSON output is valid JSON with correct schema
  - [ ] YAML output is valid YAML with correct schema
  - [ ] Invalid `--output` format returns error
  - [ ] Help text (`--help`) is accurate

- [ ] `TestCollectSystemInfo`:
  - [ ] Returns non-nil SystemInfo
  - [ ] OS/Platform/Architecture populated (or "unknown")
  - [ ] CPU info populated (or defaults)
  - [ ] Memory info populated
  - [ ] Storage is array (may be empty)
  - [ ] GPU is array (may be empty, graceful degradation)
  - [ ] NPU may be nil (graceful degradation)

- [ ] `TestFormatSystemInfo`:
  - [ ] Text output formatted correctly
  - [ ] All sections present when data available
  - [ ] Sections omitted gracefully when data unavailable

### Integration Tests

- [ ] Run on CI matrix (Linux, macOS, Windows)
- [ ] Verify command executes without error on all platforms
- [ ] Verify JSON output parseable on all platforms
- [ ] Verify YAML output parseable on all platforms

### Coverage Requirements

- [ ] 80%+ coverage for `internal/meta/system.go`
- [ ] 80%+ coverage for `cmd/ado/meta/meta.go` (system command section)
- [ ] Coverage check passes in CI (`make go.test.cover.check`)

### Manual Testing

- [ ] Test on Linux (CI: ubuntu-latest)
- [ ] Test on macOS (CI: macos-latest)
- [ ] Test on Windows (CI: windows-latest)
- [ ] Test with GPU (if available)
- [ ] Test without GPU (graceful degradation)
- [ ] Test NPU detection (Apple Silicon, Intel Core Ultra, AMD Ryzen AI)

## Cross-Platform Support Matrix

| Feature | Linux | macOS | Windows | Method |
|---------|-------|-------|---------|--------|
| **OS Info** | ✅ Full | ✅ Full | ✅ Full | gopsutil/host |
| **CPU Info** | ✅ Full | ✅ Full | ✅ Full | gopsutil/cpu |
| **Memory** | ✅ Full | ✅ Full | ✅ Full | gopsutil/mem |
| **Storage** | ✅ Full | ✅ Full | ✅ Full | gopsutil/disk |
| **GPU** | ⚠️ Best-effort | ⚠️ Best-effort | ⚠️ Best-effort | ghw/gpu or platform cmds |
| **NPU** | ⚠️ Best-effort | ⚠️ Best-effort | ⚠️ Best-effort | CPU model inference |

**Legend:**
- ✅ Full: Guaranteed to work on all systems
- ⚠️ Best-effort: Works on most systems, graceful degradation if unavailable

## Related Commands

- `ado meta info` - Show ado build metadata (version, commit, etc.)
- `ado meta env` - Show ado environment configuration (config paths, env vars)
- `ado meta features` - Show enabled/disabled ado features

## Dependencies

### Required

```go
require (
	github.com/shirou/gopsutil/v4 v4.24.0 // Cross-platform system info (OS, CPU, Memory, Storage)
)
```

### Optional (Future Enhancement)

```go
require (
	github.com/jaypipes/ghw v0.12.0 // GPU detection (optional, for enhanced GPU support)
)
```

## Implementation Phases

### Phase 1: Core System Info (Required for Initial Release)

**Scope:**
- OS, Platform, Kernel, Architecture
- CPU (model, vendor, cores, frequency)
- Memory (total, available, used, swap)
- Storage (mounted volumes, usage)
- Text/JSON/YAML output formats
- 80%+ test coverage
- Cross-platform CI validation (Linux/macOS/Windows)

**Dependencies:**
- `gopsutil/v4` only

**Deliverables:**
- `internal/meta/system.go` (SystemInfo, CollectSystemInfo, no GPU/NPU)
- `cmd/ado/meta/meta.go` (newSystemCommand, formatSystemInfo)
- `internal/meta/system_test.go` (unit tests)
- `cmd/ado/meta/meta_test.go` (command tests)
- `docs/commands/05-meta-system.md` (this spec)

**Success Criteria:**
- ✅ `ado meta system` works on Linux/macOS/Windows
- ✅ All output formats (text/json/yaml) work
- ✅ 80%+ coverage in CI
- ✅ No breaking changes

### Phase 2: GPU Detection (Optional Enhancement)

**Scope:**
- Add `GPUInfo` to `SystemInfo`
- Implement `detectGPU()` using ghw or platform commands
- Update tests for GPU section
- Graceful degradation if GPU not detectable

**Dependencies:**
- `ghw` (optional) or platform-specific commands

**Deliverables:**
- Update `internal/meta/system.go` (add GPU detection)
- Update tests (add GPU test cases)

**Success Criteria:**
- ✅ GPU detected on systems with GPU
- ✅ Graceful degradation on systems without GPU
- ✅ Tests cover GPU present/absent scenarios

### Phase 3: NPU Detection (Optional Enhancement)

**Scope:**
- Add `NPUInfo` to `SystemInfo`
- Implement `detectNPU()` using CPU model inference
- Support Apple Neural Engine, Intel AI Boost, AMD Ryzen AI
- Update tests for NPU section

**Dependencies:**
- None (CPU model-based inference)

**Deliverables:**
- Update `internal/meta/system.go` (add NPU detection)
- Update tests (add NPU test cases)

**Success Criteria:**
- ✅ NPU detected on Apple Silicon (M1/M2/M3/M4)
- ✅ NPU detected on Intel Core Ultra
- ✅ NPU detected on AMD Ryzen AI
- ✅ Graceful degradation on systems without NPU

## Security Considerations

- **No Elevated Permissions**: All system info collection works without sudo/admin
- **Read-Only**: No modifications to system state
- **Graceful Failures**: Missing info does not prevent command from succeeding
- **No Sensitive Data**: Only collects hardware/OS info, no user data or credentials

## Performance Requirements

- **Execution Time**: &lt;200ms total (gopsutil is fast, &lt;100ms for all info)
- **Memory Usage**: &lt;10MB (small structs, no large allocations)
- **CPU Usage**: Minimal (single syscall per info type)

## Documentation Requirements

- [ ] Update `docs/commands/03-meta.md` to reference new `system` subcommand
- [ ] Add `docs/commands/05-meta-system.md` (this spec)
- [ ] Update `README.md` examples if relevant
- [ ] Update `docs/workflow.md` demo examples (if used for 30-min demo)

## Acceptance Criteria

**Spec Ready for Implementation When:**
- ✅ All sections filled in
- ✅ Usage examples cover common use cases
- ✅ Output formats clearly defined with examples
- ✅ Error cases documented
- ✅ Implementation outline complete
- ✅ Testing strategy defined with 80%+ coverage requirement
- ✅ Cross-platform support documented
- ✅ Dependencies justified
- ✅ Implementation phases clearly separated

**Implementation PR Ready for Merge When:**
- ✅ All Phase 1 features implemented
- ✅ All tests pass on Linux/macOS/Windows (CI)
- ✅ 80%+ test coverage achieved
- ✅ No breaking changes to existing `meta` commands
- ✅ Documentation updated
- ✅ Conventional commit format (`feat(meta): add system subcommand`)
- ✅ CI checks pass (lint, test, coverage, build)
- ✅ Code owner approval obtained

## Notes

- **Phase 2 & 3 (GPU/NPU)**: Optional enhancements, can be deferred to follow-up PRs
- **Demo Compatibility**: This spec supports 30-minute demo workflow (Issue → Spec → Implementation → Release)
- **Extensibility**: Future phases could add network, battery, temperature, etc.
