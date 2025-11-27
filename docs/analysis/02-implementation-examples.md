# Implementation Examples & Code Patterns

## 1. Data Structure Design Pattern

### Graceful Degradation with Capability Flags

```go
package system

import "time"

// SystemInfo represents complete system information
// All fields that can fail are pointers (nil = unavailable)
type SystemInfo struct {
    OS           OSInfo        `json:"os"`
    CPU          CPUInfo       `json:"cpu"`
    GPU          *GPUInfo      `json:"gpu,omitempty"`      // Optional
    NPU          *NPUInfo      `json:"npu,omitempty"`      // Optional
    Memory       MemoryInfo    `json:"memory"`
    Storage      []StorageInfo `json:"storage"`
    Capabilities Capabilities `json:"capabilities"`
}

// Capabilities tracks what information was successfully collected
type Capabilities struct {
    CPU            bool `json:"cpu"`              // Always true (Go stdlib)
    Memory         bool `json:"memory"`           // Always true
    Storage        bool `json:"storage"`          // Always true
    OS             bool `json:"os"`               // Always true
    GPU            bool `json:"gpu"`              // False if no GPU found
    GPUDetails     bool `json:"gpu_details"`      // False if nvidia-smi/rocm-smi not found
    GPUTemperature bool `json:"gpu_temperature"` // False if can't read temps
    NPU            bool `json:"npu"`              // False if not detected
    StorageSmartData bool `json:"storage_smart"` // False if smartctl not available
}

// OSInfo represents operating system information
type OSInfo struct {
    Name       string `json:"name"`        // "Linux", "macOS", "Windows"
    Version    string `json:"version"`     // "14.1", "22.04", "11"
    Build      string `json:"build,omitempty"` // macOS: "23.1.0", Windows: Build number
    Kernel     string `json:"kernel"`     // Kernel version
    Architecture string `json:"architecture"` // "amd64", "arm64"
}

// CPUInfo represents CPU information
type CPUInfo struct {
    Model          string   `json:"model"`           // "Apple M3 Max", "Intel Core i9"
    LogicalCores   int      `json:"logical_cores"`
    PhysicalCores  int      `json:"physical_cores,omitempty"`
    ThreadsPerCore int      `json:"threads_per_core,omitempty"`
    MaxFrequency   string   `json:"max_frequency,omitempty"` // "4.2 GHz"
    Features       []string `json:"features,omitempty"` // ["AVX", "AVX2", "SVE"]
}

// GPUInfo represents GPU information
type GPUInfo struct {
    Devices []GPUDevice `json:"devices"`
}

type GPUDevice struct {
    Vendor        string `json:"vendor"`      // "NVIDIA", "AMD", "Intel", "Apple"
    Model         string `json:"model"`       // Full device name
    DriverVersion string `json:"driver_version,omitempty"`
    Memory        int64  `json:"memory_bytes,omitempty"` // In bytes, 0 if shared (Apple Silicon)
    Shared        bool   `json:"shared_memory"` // For Apple Silicon: true (uses system RAM)
}

// NPUInfo represents Neural Processing Unit information
type NPUInfo struct {
    Available      bool   `json:"available"`
    Type           string `json:"type"`      // "Apple Neural Engine", "Intel AI Boost"
    Details        string `json:"details,omitempty"` // Limitation notes
}

// MemoryInfo represents memory statistics
type MemoryInfo struct {
    TotalBytes     int64  `json:"total_bytes"`
    AvailableBytes int64  `json:"available_bytes"`
    UsedBytes      int64  `json:"used_bytes"`
    SwapTotal      int64  `json:"swap_total_bytes,omitempty"` // 0 on macOS/Windows
    SwapAvailable  int64  `json:"swap_available_bytes,omitempty"`
}

// StorageInfo represents a mounted filesystem
type StorageInfo struct {
    Mountpoint     string `json:"mountpoint"`      // "/", "C:\", etc.
    Device         string `json:"device"`          // "/dev/sda1", "C:", etc.
    Type           string `json:"type"`            // "ext4", "APFS", "NTFS"
    TotalBytes     int64  `json:"total_bytes"`
    AvailableBytes int64  `json:"available_bytes"`
    UsedBytes      int64  `json:"used_bytes"`
    IsReadOnly     bool   `json:"is_readonly"`
    IsSSD          bool   `json:"is_ssd,omitempty"` // Unknown if false
}

func (s StorageInfo) UsagePercent() float64 {
    if s.TotalBytes == 0 {
        return 0
    }
    return float64(s.UsedBytes) / float64(s.TotalBytes) * 100
}
```

---

## 2. Collection Interface Pattern

### Main Entry Point with Context-Aware Collection

```go
package system

import (
    "context"
    "fmt"
)

// Collector coordinates all system information gathering
type Collector struct {
    // Option flags for optional features
    IncludeGPU      bool
    IncludeNPU      bool
    IncludeSmartData bool
    
    // Timeout for external tool execution
    CommandTimeout  time.Duration
}

// NewCollector creates a new system information collector
func NewCollector() *Collector {
    return &Collector{
        IncludeGPU:     true,
        IncludeNPU:     true,
        IncludeSmartData: false, // Skip by default (requires tools/privileges)
        CommandTimeout: 5 * time.Second,
    }
}

// Collect gathers all system information
// Never returns error - gracefully degrades if components fail
func (c *Collector) Collect(ctx context.Context) SystemInfo {
    info := SystemInfo{
        Capabilities: Capabilities{
            CPU:      true,    // Always available (Go stdlib)
            Memory:   true,    // Always available
            Storage:  true,    // Always available
            OS:       true,    // Always available
        },
    }
    
    // Always-available info (non-optional)
    info.OS = collectOSInfo(ctx)
    info.CPU = collectCPUInfo(ctx)
    info.Memory = collectMemoryInfo(ctx)
    info.Storage = collectStorageInfo(ctx)
    
    // Optional info with graceful degradation
    if c.IncludeGPU {
        if gpu, err := collectGPUInfo(ctx, c.CommandTimeout); err == nil {
            info.GPU = &gpu
            info.Capabilities.GPU = true
        }
    }
    
    if c.IncludeNPU {
        if npu, err := collectNPUInfo(ctx); err == nil {
            info.NPU = &npu
            info.Capabilities.NPU = true
        }
    }
    
    return info
}
```

---

## 3. Platform-Specific Implementation Pattern

### OS-Agnostic Interface with Platform Implementations

```go
// internal/system/cpu_info.go
package system

import (
    "context"
    "runtime"
)

func collectCPUInfo(ctx context.Context) CPUInfo {
    info := CPUInfo{
        LogicalCores: runtime.NumCPU(),
    }
    
    // Platform-specific details
    switch runtime.GOOS {
    case "linux":
        info = collectCPUInfoLinux(ctx, info)
    case "darwin":
        info = collectCPUInfoDarwin(ctx, info)
    case "windows":
        info = collectCPUInfoWindows(ctx, info)
    }
    
    return info
}
```

```go
// internal/system/cpu_darwin.go
// +build darwin

package system

import (
    "context"
    "os/exec"
    "strings"
)

func collectCPUInfoDarwin(ctx context.Context, info CPUInfo) CPUInfo {
    // Get CPU model
    if model, err := sysctlString("hw.model"); err == nil {
        info.Model = model
    }
    
    // Get physical cores
    if cores, err := sysctlInt("hw.physicalcpu"); err == nil {
        info.PhysicalCores = cores
    }
    
    // Get max frequency
    if freq, err := sysctlInt("hw.cpufrequency_max"); err == nil {
        info.MaxFrequency = formatFrequency(freq)
    }
    
    // Get CPU features
    info.Features = detectCPUFeatures()
    
    return info
}

func sysctlString(key string) (string, error) {
    cmd := exec.CommandContext(context.Background(), "sysctl", "-n", key)
    out, err := cmd.Output()
    if err != nil {
        return "", err
    }
    return strings.TrimSpace(string(out)), nil
}

func sysctlInt(key string) (int, error) {
    val, err := sysctlString(key)
    if err != nil {
        return 0, err
    }
    return strconv.Atoi(val)
}

func detectCPUFeatures() []string {
    // Use golang.org/x/sys/cpu for portable feature detection
    var features []string
    
    // Platform-specific intrinsics check
    // This is complex and platform-specific
    // For now, return basic info
    
    return features
}

func formatFrequency(hz int) string {
    if hz >= 1e9 {
        return fmt.Sprintf("%.1f GHz", float64(hz)/1e9)
    }
    return fmt.Sprintf("%d MHz", hz/1e6)
}
```

```go
// internal/system/cpu_linux.go
// +build linux

package system

import (
    "context"
    "os"
    "strings"
)

func collectCPUInfoLinux(ctx context.Context, info CPUInfo) CPUInfo {
    // Parse /proc/cpuinfo
    cpuinfo, err := os.ReadFile("/proc/cpuinfo")
    if err != nil {
        return info
    }
    
    for _, line := range strings.Split(string(cpuinfo), "\n") {
        if strings.HasPrefix(line, "model name") {
            parts := strings.SplitN(line, ":", 2)
            if len(parts) == 2 {
                info.Model = strings.TrimSpace(parts[1])
                break // All cores have same model
            }
        }
        
        if strings.HasPrefix(line, "flags") {
            parts := strings.SplitN(line, ":", 2)
            if len(parts) == 2 {
                info.Features = strings.Fields(parts[1])
            }
        }
    }
    
    return info
}
```

```go
// internal/system/cpu_windows.go
// +build windows

package system

import (
    "context"
    "syscall"
)

func collectCPUInfoWindows(ctx context.Context, info CPUInfo) CPUInfo {
    // Use WMI or Windows Registry
    // This requires Windows-specific libraries
    // For now, return basic info from stdlib
    
    // TODO: Implement WMI query for Win32_Processor
    // This requires additional dependencies or cgo
    
    return info
}
```

---

## 4. GPU Detection with Graceful Fallback

### Tool-Based Detection Pattern

```go
// internal/system/gpu_info.go
package system

import (
    "context"
    "os/exec"
    "strings"
    "time"
)

func collectGPUInfo(ctx context.Context, timeout time.Duration) (GPUInfo, error) {
    info := GPUInfo{
        Devices: []GPUDevice{},
    }
    
    // Try NVIDIA
    if devices, err := detectNVIDIA(ctx, timeout); err == nil && len(devices) > 0 {
        info.Devices = append(info.Devices, devices...)
    }
    
    // Try AMD ROCm
    if devices, err := detectAMD(ctx, timeout); err == nil && len(devices) > 0 {
        info.Devices = append(info.Devices, devices...)
    }
    
    // Try Intel
    if devices, err := detectIntel(ctx, timeout); err == nil && len(devices) > 0 {
        info.Devices = append(info.Devices, devices...)
    }
    
    // Platform-specific detection
    switch runtime.GOOS {
    case "darwin":
        if devices, err := detectAppleSilicon(ctx, timeout); err == nil && len(devices) > 0 {
            info.Devices = append(info.Devices, devices...)
        }
    case "windows":
        if devices, err := detectWindowsGPU(ctx, timeout); err == nil && len(devices) > 0 {
            info.Devices = append(info.Devices, devices...)
        }
    }
    
    if len(info.Devices) == 0 {
        return info, fmt.Errorf("no GPU detected")
    }
    
    return info, nil
}

// detectNVIDIA tries to query nvidia-smi
func detectNVIDIA(ctx context.Context, timeout time.Duration) ([]GPUDevice, error) {
    ctx, cancel := context.WithTimeout(ctx, timeout)
    defer cancel()
    
    // Check if nvidia-smi exists
    cmd := exec.CommandContext(ctx, "which", "nvidia-smi")
    if err := cmd.Run(); err != nil {
        return nil, fmt.Errorf("nvidia-smi not found")
    }
    
    // Query GPU info
    cmd = exec.CommandContext(ctx,
        "nvidia-smi",
        "--query-gpu=name,driver_version,memory.total",
        "--format=csv,noheader,nounits",
    )
    
    output, err := cmd.Output()
    if err != nil {
        return nil, err
    }
    
    var devices []GPUDevice
    for _, line := range strings.Split(strings.TrimSpace(string(output)), "\n") {
        if line == "" {
            continue
        }
        
        parts := strings.Split(line, ",")
        if len(parts) >= 3 {
            device := GPUDevice{
                Vendor:        "NVIDIA",
                Model:         strings.TrimSpace(parts[0]),
                DriverVersion: strings.TrimSpace(parts[1]),
                Memory:        parseMemory(strings.TrimSpace(parts[2])),
            }
            devices = append(devices, device)
        }
    }
    
    return devices, nil
}

// detectAppleSilicon detects Apple Metal GPU (M1+)
func detectAppleSilicon(ctx context.Context, timeout time.Duration) ([]GPUDevice, error) {
    // Check if running on Apple Silicon (M1+)
    // This is available in runtime or via sysctl hw.model
    
    model, _ := sysctlString("hw.model")
    
    // Apple Silicon models: MacBookPro18,x, MacBookAir11,x, Mac14,x, etc.
    if !strings.HasPrefix(model, "MacBook") && !strings.HasPrefix(model, "Mac") {
        return nil, fmt.Errorf("not Apple Silicon")
    }
    
    // Parse system_profiler output for detailed GPU info
    ctx, cancel := context.WithTimeout(ctx, timeout)
    defer cancel()
    
    cmd := exec.CommandContext(ctx, "system_profiler", "SPDisplaysDataType")
    output, err := cmd.Output()
    if err != nil {
        return nil, err
    }
    
    // Parse GPU info from output
    // This is fragile; consider it best-effort
    devices := parseAppleMetalGPU(string(output))
    if len(devices) == 0 {
        // Fallback: we know Apple Silicon exists, just can't get details
        devices = []GPUDevice{
            {
                Vendor: "Apple",
                Model:  "Integrated GPU",
                Shared: true,
            },
        }
    }
    
    return devices, nil
}

func parseMemory(sizeStr string) int64 {
    // Parse "15360 MB" or "8 GB"
    // This is a simplified version
    return 0 // TODO: Implement proper parsing
}

func parseAppleMetalGPU(output string) []GPUDevice {
    // Parse system_profiler SPDisplaysDataType output
    // This is complex and fragile; return empty on parse errors
    return []GPUDevice{}
}
```

---

## 5. Text Output Formatting with Graceful Degradation

### Formatter Pattern

```go
// internal/system/format.go
package system

import (
    "fmt"
    "strings"
)

// FormatText returns human-readable system information
func (s *SystemInfo) FormatText() string {
    var b strings.Builder
    
    fmt.Fprintln(&b, "System Information")
    fmt.Fprintln(&b, strings.Repeat("=", 50))
    fmt.Fprintln(&b)
    
    // OS Section
    fmt.Fprintln(&b, "OS")
    fmt.Fprintf(&b, "  Name: %s\n", s.OS.Name)
    fmt.Fprintf(&b, "  Version: %s\n", s.OS.Version)
    if s.OS.Build != "" {
        fmt.Fprintf(&b, "  Build: %s\n", s.OS.Build)
    }
    fmt.Fprintf(&b, "  Kernel: %s\n", s.OS.Kernel)
    fmt.Fprintf(&b, "  Architecture: %s\n", s.OS.Architecture)
    fmt.Fprintln(&b)
    
    // CPU Section
    fmt.Fprintln(&b, "CPU")
    fmt.Fprintf(&b, "  Model: %s\n", s.CPU.Model)
    fmt.Fprintf(&b, "  Logical Cores: %d\n", s.CPU.LogicalCores)
    if s.CPU.PhysicalCores > 0 {
        fmt.Fprintf(&b, "  Physical Cores: %d\n", s.CPU.PhysicalCores)
    }
    if s.CPU.MaxFrequency != "" {
        fmt.Fprintf(&b, "  Max Frequency: %s\n", s.CPU.MaxFrequency)
    }
    if len(s.CPU.Features) > 0 {
        fmt.Fprintf(&b, "  Features: %s\n", strings.Join(s.CPU.Features[:min(5, len(s.CPU.Features))], ", "))
        if len(s.CPU.Features) > 5 {
            fmt.Fprintf(&b, "    + %d more\n", len(s.CPU.Features)-5)
        }
    }
    fmt.Fprintln(&b)
    
    // Memory Section
    fmt.Fprintln(&b, "Memory")
    fmt.Fprintf(&b, "  Total: %s\n", formatBytes(s.Memory.TotalBytes))
    fmt.Fprintf(&b, "  Available: %s\n", formatBytes(s.Memory.AvailableBytes))
    fmt.Fprintf(&b, "  Used: %s (%.1f%%)\n", 
        formatBytes(s.Memory.UsedBytes),
        float64(s.Memory.UsedBytes)/float64(s.Memory.TotalBytes)*100)
    fmt.Fprintln(&b)
    
    // Storage Section
    if len(s.Storage) > 0 {
        fmt.Fprintln(&b, "Storage")
        for _, vol := range s.Storage {
            fmt.Fprintf(&b, "  %s (%s)\n", vol.Mountpoint, vol.Device)
            fmt.Fprintf(&b, "    Type: %s", vol.Type)
            if vol.IsSSD {
                fmt.Fprintf(&b, " (SSD)")
            }
            fmt.Fprintln(&b)
            fmt.Fprintf(&b, "    Total: %s\n", formatBytes(vol.TotalBytes))
            fmt.Fprintf(&b, "    Used: %s (%.1f%%)\n", 
                formatBytes(vol.UsedBytes),
                vol.UsagePercent())
            fmt.Fprintf(&b, "    Available: %s\n", formatBytes(vol.AvailableBytes))
        }
        fmt.Fprintln(&b)
    }
    
    // GPU Section (Optional)
    if s.GPU != nil && s.Capabilities.GPU {
        fmt.Fprintln(&b, "GPU")
        for i, device := range s.GPU.Devices {
            if i > 0 {
                fmt.Fprintln(&b)
            }
            fmt.Fprintf(&b, "  [%d] %s %s\n", i, device.Vendor, device.Model)
            if device.DriverVersion != "" {
                fmt.Fprintf(&b, "      Driver: %s\n", device.DriverVersion)
            }
            if device.Memory > 0 {
                fmt.Fprintf(&b, "      Memory: %s\n", formatBytes(device.Memory))
            } else if device.Shared {
                fmt.Fprintf(&b, "      Memory: Shared (system RAM)\n")
            }
        }
        fmt.Fprintln(&b)
    } else if !s.Capabilities.GPU {
        fmt.Fprintln(&b, "GPU")
        fmt.Fprintln(&b, "  Not detected (or nvidia-smi/rocm-smi not available)")
        fmt.Fprintln(&b)
    }
    
    // NPU Section (Optional)
    if s.NPU != nil && s.Capabilities.NPU {
        fmt.Fprintln(&b, "NPU (Neural Processing Unit)")
        fmt.Fprintf(&b, "  Type: %s\n", s.NPU.Type)
        fmt.Fprintf(&b, "  Available: %v\n", s.NPU.Available)
        if s.NPU.Details != "" {
            fmt.Fprintf(&b, "  Note: %s\n", s.NPU.Details)
        }
        fmt.Fprintln(&b)
    }
    
    // Capabilities Section
    fmt.Fprintln(&b, "Capabilities")
    fmt.Fprintf(&b, "  GPU Details: %v\n", s.Capabilities.GPUDetails)
    fmt.Fprintf(&b, "  NPU Info: %v\n", s.Capabilities.NPU)
    fmt.Fprintf(&b, "  SMART Data: %v\n", s.Capabilities.StorageSmartData)
    
    return b.String()
}

func formatBytes(bytes int64) string {
    units := []string{"B", "KB", "MB", "GB", "TB"}
    value := float64(bytes)
    
    for _, unit := range units {
        if value < 1024 {
            return fmt.Sprintf("%.1f %s", value, unit)
        }
        value /= 1024
    }
    
    return fmt.Sprintf("%.1f PB", value)
}

func min(a, b int) int {
    if a < b {
        return a
    }
    return b
}
```

---

## 6. Error Handling Pattern

### Context-Aware Error Handling

```go
// internal/system/errors.go
package system

import (
    "fmt"
)

// CollectionError represents a partial collection failure
// The partial data is still returned; this error indicates what failed
type CollectionError struct {
    Component string // "GPU", "NPU", "SMART", etc.
    Err       error
    IsFatal   bool // If true, entire collection should fail
}

func (e CollectionError) Error() string {
    if e.IsFatal {
        return fmt.Sprintf("fatal error collecting %s: %v", e.Component, e.Err)
    }
    return fmt.Sprintf("optional component %s unavailable: %v", e.Component, e.Err)
}

// CollectionErrors represents multiple non-fatal collection failures
type CollectionErrors struct {
    Errors []CollectionError
}

func (e CollectionErrors) Error() string {
    if len(e.Errors) == 0 {
        return "no errors"
    }
    
    msg := fmt.Sprintf("%d collection issues:\n", len(e.Errors))
    for _, err := range e.Errors {
        if err.IsFatal {
            msg += fmt.Sprintf("  [FATAL] %s: %v\n", err.Component, err.Err)
        } else {
            msg += fmt.Sprintf("  [WARN] %s: %v\n", err.Component, err.Err)
        }
    }
    return msg
}

// Example usage in collector:
func (c *Collector) CollectWithLogging(ctx context.Context) (SystemInfo, []CollectionError) {
    var errs []CollectionError
    info := c.Collect(ctx)
    
    // Log collection errors without failing
    // Return both complete info and error details
    
    return info, errs
}
```

---

## 7. Integration with `ado meta system` Command

### Command Implementation

```go
// cmd/ado/meta/system.go
package meta

import (
    "fmt"
    "github.com/spf13/cobra"
    "github.com/anowarislam/ado/internal/system"
    "github.com/anowarislam/ado/internal/ui"
)

func newSystemCommand() *cobra.Command {
    var output string
    var includeOptional bool
    
    cmd := &cobra.Command{
        Use:   "system",
        Short: "Display system hardware and OS information",
        Long: `Displays detailed system information including CPU, GPU, memory, and storage.

GPU and NPU information is optional and gracefully degraded if tools are unavailable.
`,
        RunE: func(cmd *cobra.Command, args []string) error {
            collector := system.NewCollector()
            collector.IncludeGPU = true
            collector.IncludeNPU = true
            collector.IncludeSmartData = includeOptional
            
            info := collector.Collect(cmd.Context())
            
            format, err := ui.ParseOutputFormat(output)
            if err != nil {
                return err
            }
            
            return ui.PrintOutput(cmd.OutOrStdout(), format, info, func() (string, error) {
                return info.FormatText(), nil
            })
        },
    }
    
    cmd.Flags().StringVarP(&output, "output", "o", "text", "Output format: text, json, yaml")
    cmd.Flags().BoolVar(&includeOptional, "include-optional", false, "Include SMART data and other optional info")
    
    return cmd
}

// Update root command to register the new subcommand
func NewCommand(buildInfo internalmeta.BuildInfo) *cobra.Command {
    cmd := &cobra.Command{
        Use:   "meta",
        Short: "Introspect the ado binary and its environment",
    }

    cmd.AddCommand(
        newInfoCommand(buildInfo),
        newEnvCommand(),
        newFeaturesCommand(),
        newSystemCommand(),  // Add this line
    )

    return cmd
}
```

---

## 8. Testing Pattern

### Platform-Agnostic Testing

```go
// internal/system/system_test.go
package system

import (
    "context"
    "testing"
)

func TestCollectorAlwaysReturnsBaseInfo(t *testing.T) {
    collector := NewCollector()
    ctx := context.Background()
    
    info := collector.Collect(ctx)
    
    // Core info always available
    if info.CPU.LogicalCores == 0 {
        t.Error("CPU logical cores should always be available")
    }
    
    if info.Memory.TotalBytes == 0 {
        t.Error("Memory total should always be available")
    }
    
    if info.OS.Architecture == "" {
        t.Error("Architecture should always be available")
    }
    
    // Capabilities always set correctly
    if !info.Capabilities.CPU || !info.Capabilities.Memory || !info.Capabilities.OS {
        t.Error("Core capabilities should always be true")
    }
}

func TestGracefulDegradation(t *testing.T) {
    collector := NewCollector()
    collector.CommandTimeout = 100 * time.Millisecond // Aggressive timeout
    
    ctx := context.Background()
    info := collector.Collect(ctx)
    
    // GPU may or may not be available, but that's OK
    // We should still have all base info
    if info.CPU.LogicalCores == 0 {
        t.Error("Should have CPU info despite GPU failure")
    }
    
    // Capabilities should reflect what succeeded
    // (GPU may be false if tools not available)
}

func TestStorageUsageCalculation(t *testing.T) {
    s := StorageInfo{
        TotalBytes:     1000,
        AvailableBytes: 300,
        UsedBytes:      700,
    }
    
    expected := 70.0
    if usage := s.UsagePercent(); usage != expected {
        t.Errorf("Expected %.1f%%, got %.1f%%", expected, usage)
    }
}

func TestFormatBytes(t *testing.T) {
    tests := []struct {
        bytes    int64
        expected string
    }{
        {1024, "1.0 KB"},
        {1024 * 1024, "1.0 MB"},
        {1024 * 1024 * 1024, "1.0 GB"},
    }
    
    for _, tt := range tests {
        if result := formatBytes(tt.bytes); result != tt.expected {
            t.Errorf("formatBytes(%d): expected %q, got %q", tt.bytes, tt.expected, result)
        }
    }
}
```

