package meta

import (
	"context"
	"testing"
)

func TestCollectSystemInfo(t *testing.T) {
	ctx := context.Background()
	info := CollectSystemInfo(ctx)

	t.Run("OS fields populated", func(t *testing.T) {
		if info.OS == "" {
			t.Error("OS should not be empty")
		}
		if info.Platform == "" {
			t.Error("Platform should not be empty")
		}
		if info.Architecture == "" {
			t.Error("Architecture should not be empty")
		}
		// Kernel may be empty on some platforms, so we don't enforce it
	})

	t.Run("CPU fields populated", func(t *testing.T) {
		if info.CPU.Model == "" {
			t.Error("CPU Model should not be empty")
		}
		// Cores should be > 0 on any real system
		if info.CPU.Cores <= 0 {
			t.Error("CPU Cores should be > 0")
		}
		// FrequencyMHz may be 0 on Apple Silicon, so we don't enforce it
	})

	t.Run("Memory fields populated", func(t *testing.T) {
		if info.Memory.TotalMB == 0 {
			t.Error("Memory TotalMB should be > 0")
		}
		if info.Memory.UsedPercent < 0 || info.Memory.UsedPercent > 100 {
			t.Errorf("Memory UsedPercent should be 0-100, got %.2f", info.Memory.UsedPercent)
		}
	})

	t.Run("Storage is array", func(t *testing.T) {
		// Storage should always be an array (may be empty)
		if info.Storage == nil {
			t.Error("Storage should not be nil")
		}
		// On real systems, we should have at least one storage device
		if len(info.Storage) == 0 {
			t.Log("Warning: No storage devices detected (this may be normal in containers)")
		}
	})

	t.Run("GPU is array", func(t *testing.T) {
		// GPU should always be an array (may be empty)
		if info.GPU == nil {
			t.Error("GPU should not be nil")
		}
		// GPU detection is best-effort, so empty array is acceptable
	})

	t.Run("NPU may be nil", func(t *testing.T) {
		// NPU is optional and may be nil if not detected
		// This is acceptable - no error
		t.Logf("NPU detected: %v", info.NPU != nil)
	})
}

func TestDetectGPU(t *testing.T) {
	ctx := context.Background()

	// Test that detectGPU returns a valid slice (not nil)
	gpus := detectGPU(ctx)

	if gpus == nil {
		t.Error("detectGPU() returned nil, should return empty slice if no GPUs detected")
	}

	// GPU detection is hardware-dependent, so we can't assert specific GPUs
	// Just verify the return type and that each GPU has required fields
	for i, gpu := range gpus {
		if gpu.Vendor == "" {
			t.Errorf("GPU[%d].Vendor is empty", i)
		}
		if gpu.Model == "" {
			t.Errorf("GPU[%d].Model is empty", i)
		}
		if gpu.Type == "" {
			t.Errorf("GPU[%d].Type is empty", i)
		}
		// Type should be one of the valid values
		validTypes := map[string]bool{"integrated": true, "discrete": true, "unknown": true}
		if !validTypes[gpu.Type] {
			t.Errorf("GPU[%d].Type = %s, want one of: integrated, discrete, unknown", i, gpu.Type)
		}
	}
}

func TestDetectNPU(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name         string
		cpuModel     string
		os           string
		wantDetected bool
		wantType     string
	}{
		{
			name:         "Apple M1",
			cpuModel:     "Apple M1",
			os:           "darwin",
			wantDetected: true,
			wantType:     "Apple Neural Engine",
		},
		{
			name:         "Apple M2 Pro",
			cpuModel:     "Apple M2 Pro",
			os:           "darwin",
			wantDetected: true,
			wantType:     "Apple Neural Engine",
		},
		{
			name:         "Apple M3 Max",
			cpuModel:     "Apple M3 Max",
			os:           "darwin",
			wantDetected: true,
			wantType:     "Apple Neural Engine",
		},
		{
			name:         "Apple M4",
			cpuModel:     "Apple M4",
			os:           "darwin",
			wantDetected: true,
			wantType:     "Apple Neural Engine",
		},
		{
			name:         "Intel Core Ultra",
			cpuModel:     "Intel(R) Core(TM) Ultra 7 155H",
			os:           "windows",
			wantDetected: true,
			wantType:     "Intel AI Boost",
		},
		{
			name:         "AMD Ryzen AI",
			cpuModel:     "AMD Ryzen AI 9 HX 370",
			os:           "windows",
			wantDetected: true,
			wantType:     "AMD Ryzen AI",
		},
		{
			name:         "Intel Core i7 (no NPU)",
			cpuModel:     "Intel(R) Core(TM) i7-9750H CPU @ 2.60GHz",
			os:           "linux",
			wantDetected: false,
			wantType:     "",
		},
		{
			name:         "AMD Ryzen 9 (no NPU)",
			cpuModel:     "AMD Ryzen 9 5900X",
			os:           "linux",
			wantDetected: false,
			wantType:     "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			npu := detectNPU(ctx, tt.cpuModel, tt.os)

			if tt.wantDetected {
				if npu == nil {
					t.Error("detectNPU() returned nil, want NPU detected")
					return
				}
				if !npu.Detected {
					t.Error("detectNPU() returned Detected=false, want true")
				}
				if npu.Type != tt.wantType {
					t.Errorf("detectNPU() returned Type=%s, want %s", npu.Type, tt.wantType)
				}
				if npu.InferenceMethod != "cpu_model" {
					t.Errorf("detectNPU() returned InferenceMethod=%s, want cpu_model", npu.InferenceMethod)
				}
			} else {
				if npu != nil {
					t.Errorf("detectNPU() returned %+v, want nil", npu)
				}
			}
		})
	}
}

func TestVerifyAnnotations(t *testing.T) {
	t.Fatal("Intentional failure - testing annotation display")
}
