# 30-Minute Demo Talking Script: `ado meta system`

**Feature**: Implementing a new CLI command from issue to release
**Duration**: 30 minutes
**Format**: Word-for-word script with stage directions

---

## Pre-Show Setup (5 minutes before)

**[TERMINAL SETUP]**
- Terminal split into 3 panes:
  - Left: Code editor (vim/vscode)
  - Top-right: Git commands
  - Bottom-right: Make/test commands
- Font size: Large enough for back row to read
- Working directory: `~/Projects/ado`
- Git status: Clean, on `main` branch

**[BROWSER SETUP]**
- Tab 1: https://github.com/anowarislam/ado/issues/53
- Tab 2: https://github.com/anowarislam/ado/pulls
- Tab 3: https://github.com/anowarislam/ado/actions
- Tab 4: https://github.com/anowarislam/ado (repo home)

**[MATERIALS]**
- Water/coffee within reach
- Timer visible (30-minute countdown)
- Notes printout or second screen with this script

---

## [0:00-0:30] Opening & Hook

**[STAGE DIRECTION]** Stand, smile, make eye contact. Energy: HIGH

**[YOU SAY]**

> "Good morning everyone! Thanks for joining. I'm excited to show you something pretty cool today."
>
> *[Pause, let people settle]*
>
> "Over the next 30 minutes, I'm going to implement a feature from absolute scratch - and I mean scratch - all the way to an automated release. We're talking: GitHub issue, specification, code, tests, CI/CD, code review, and release automation."
>
> *[Pause for emphasis]*
>
> "The feature itself is simple - we're adding a `meta system` command to display diagnostic info like CPU, memory, and storage. But the feature isn't the point. The workflow is the point."
>
> *[Gesture to screen]*
>
> "By the end, you'll see how we ship quality code fast, with confidence, and with built-in security. And there'll be time for questions throughout, so don't hesitate to raise your hand."
>
> *[Smile]*
>
> "Sound good? Alright, let's dive in."

**[TIMER: START]**

---

## [0:30-2:30] Phase 0: Workflow Overview

**[STAGE DIRECTION]** Switch to browser, show workflow doc. Energy: MEDIUM, educational tone

**[BROWSER ACTION]** Navigate to `docs/workflow.md` on GitHub

**[YOU SAY]**

> "First, quick context. The `ado` project follows a three-phase workflow."
>
> *[Scroll to show the three phases]*
>
> "Phase 1: **Issue**. We define the problem and acceptance criteria. What are we solving and why?"
>
> "Phase 2: **Spec**. We document what we're building before writing any code. This is our contract - if the implementation doesn't match the spec, something's wrong."
>
> "Phase 3: **Implementation**. Only now do we write code. We already know what we're building, so we can focus on how."
>
> *[Pause]*
>
> "Why this order? Because writing a spec takes 30 minutes but prevents days of rework. It forces alignment before investment."
>
> *[Scroll to conventional commits section]*
>
> "We also use conventional commits - every commit follows a format like `feat:` or `fix:`. Why? Because this enables semantic versioning automatically. A `feat:` commit bumps the minor version. A `fix:` bumps the patch. Breaking changes get `!` and bump major version."
>
> *[Gesture enthusiastically]*
>
> "This means releases are completely automated. No manual changelogs, no manual version bumping, no manual tagging. The commit message IS the release automation."
>
> *[Quick transition]*
>
> "The project also has automated CI/CD that enforces quality. 80% test coverage minimum. All tests must pass. Code must lint. Build must succeed. No exceptions, even for administrators."

**[PAUSE FOR QUESTIONS]**

> "Any questions on the workflow before we start building? ... No? Great, let's go."

---

## [2:30-5:00] Phase 1: Review Issue #53

**[STAGE DIRECTION]** Switch to Issue #53. Energy: MEDIUM, storytelling mode

**[BROWSER ACTION]** Navigate to https://github.com/anowarislam/ado/issues/53

**[YOU SAY]**

> "Normally, the first step is creating an issue. I've already created Issue #53 to save time, but let's walk through it together."
>
> *[Scroll through issue slowly]*
>
> "Here's the problem: Users troubleshooting issues need an easy way to share system information. Right now, they have to manually run `uname -a`, `lscpu`, `free -h`, and paste multiple outputs. That's tedious."
>
> *[Scroll to proposed solution]*
>
> "Our solution: Add `ado meta system` - one command that outputs OS, CPU, memory, storage info in text, JSON, or YAML format."
>
> *[Point to examples]*
>
> "The issue includes concrete examples. Text output for humans, JSON for bug reports, YAML for config files."
>
> *[Scroll to technical approach]*
>
> "Here's where it gets interesting. The issue documents our technical approach upfront:
> - We'll use gopsutil library - 10,000 stars, battle-tested, no CGO needed
> - We'll follow the existing meta command pattern
> - We DON'T need an ADR - this is a routine feature, not an architectural decision"
>
> *[Scroll to implementation plan]*
>
> "And we have a phased implementation plan. Phase 1: core features. Phase 2: GPU detection. Phase 3: NPU detection. Incremental delivery, not big bang."
>
> *[Scroll back to top]*
>
> "Notice what we have before writing a single line of code:
> - Clear problem statement
> - Concrete examples of the solution
> - Technical approach with library choices justified
> - Implementation plan with phases
> - Success criteria
>
> This issue is our North Star. Everything we build should solve this problem."

**[PAUSE FOR QUESTIONS]**

> "Questions on the issue? ... Cool, let's write the spec."

---

## [5:00-9:00] Phase 2: Create Spec

**[STAGE DIRECTION]** Switch to terminal. Energy: MEDIUM, coding mode (calm, focused)

**[TERMINAL ACTION]** Top-right pane (Git)

**[YOU SAY]**

> "Now we write the specification. This documents exactly what we're building."
>
> *[Start typing]*

**[TYPE]**
```bash
git checkout -b spec/meta-system
```

**[YOU SAY]**

> "Branch name follows convention: `spec/` prefix for specification branches."
>
> *[Switch to editor pane]*
>
> "I'm going to create the spec file. In a real session, I'd type this out, but for time, I'll paste a pre-prepared version and walk through it."

**[PASTE INTO EDITOR]** `docs/commands/05-meta-system.md`

```markdown
# meta system Command Spec

## Command
ado meta system [--output FORMAT]

## Purpose
Display system diagnostic information for troubleshooting and bug reports.

## Usage Examples

### Text Output (Default)
$ ado meta system
OS: darwin
Platform: macOS 14.2.1
Architecture: arm64
CPUs: 10
Memory: 16384 MB total, 8245 MB available
Storage: 494 GB total, 123 GB used (25%)

### JSON Output
$ ado meta system --output json
{
  "os": "darwin",
  "platform": "macOS 14.2.1",
  "architecture": "arm64",
  "cpus": 10,
  "memory": {...}
}

### YAML Output
$ ado meta system --output yaml
os: darwin
platform: macOS 14.2.1
...

## Implementation

### Files
- cmd/ado/meta/meta.go - Add newSystemCommand()
- internal/meta/system.go - Add CollectSystemInfo()
- cmd/ado/meta/meta_test.go - Add tests

### Dependencies
- github.com/shirou/gopsutil/v4

### Testing Strategy
- Table-driven tests for all output formats
- 80%+ coverage required
- CI validates on Linux/macOS/Windows

## Error Cases
- Invalid --output format: returns error
- System info unavailable: shows "unknown"

Closes #53
```

**[YOU SAY]**

> *[Scroll through spec]*
>
> "Notice what this spec includes:"
>
> *[Point to examples]*
>
> "**Concrete examples**. Not 'display CPU info' - EXACTLY what the output looks like. These become test cases."
>
> *[Point to implementation section]*
>
> "**Implementation details**. Which files to modify. What functions to add. What dependencies to use. Another developer could implement this from the spec alone."
>
> *[Point to testing]*
>
> "**Testing strategy**. How we'll validate it works. Coverage requirements. CI validation."
>
> *[Point to error cases]*
>
> "**Error cases**. What happens when things go wrong."
>
> "This spec is our contract. Implementation must match this spec exactly."

**[TERMINAL ACTION]** Switch to Git pane

**[YOU SAY]**

> "Let's commit and create a PR."

**[TYPE]**
```bash
git add docs/commands/05-meta-system.md
git commit -m "docs(spec): command meta-system

Add specification for 'ado meta system' subcommand.
Displays OS, CPU, memory, storage diagnostic info.
Supports text/json/yaml output formats.

Closes #53"
```

**[YOU SAY]**

> "Notice the commit message: `docs(spec):` - conventional format. The body explains what this is and links to the issue with `Closes #53`."

**[TYPE]**
```bash
git push -u origin spec/meta-system
gh pr create --title "docs(spec): command meta-system" --body "Spec for ado meta system command.

Closes #53" --web
```

**[BROWSER ACTION]** PR opens automatically

**[YOU SAY]**

> *[Point to PR]*
>
> "PR created! Notice three things immediately:"
>
> *[Point to CODEOWNERS section]*
>
> "1. CODEOWNERS automatically assigned @anowarislam as reviewer. The `.github/CODEOWNERS` file maps file paths to reviewers. I modified `docs/commands/`, so the docs owner is auto-assigned."
>
> *[Point to checks]*
>
> "2. CI checks are already running. Conventional Commits is validating the title format."
>
> *[Point to issue reference]*
>
> "3. Issue #53 is linked. When this PR merges, the issue auto-closes."
>
> "In a real workflow, we'd wait 1-2 days for review. For the demo, let's assume it's approved and merged."

**[SIMULATE]**

> "Approved ✓ Merged ✓"
>
> "Spec is now the source of truth. Let's implement it."

**[PAUSE FOR QUESTIONS]**

> "Questions on specs? ... No? Okay, implementation time."

---

## [9:00-17:00] Phase 3: Implementation

**[STAGE DIRECTION]** Terminal focus. Energy: MEDIUM-HIGH, building excitement

**[TERMINAL ACTION]** Git pane

**[YOU SAY]**

> "Now we write code that matches the spec we just approved."

**[TYPE]**
```bash
git checkout main
git pull origin main
git checkout -b feat/meta-system
```

**[YOU SAY]**

> "Branch name: `feat/` prefix - this is a feature branch."
>
> "First, we add the dependency."

**[TYPE]**
```bash
go get github.com/shirou/gopsutil/v4@latest
```

**[YOU SAY]**

> "gopsutil provides cross-platform system info without CGO. That means we can cross-compile easily for Linux, macOS, Windows without platform-specific build requirements."

---

### [9:00-12:00] Step 1: System Info Collector

**[TERMINAL ACTION]** Switch to editor pane

**[YOU SAY]**

> "We'll create three files. First: the system info collector in `internal/meta/system.go`."
>
> "Again, I'll paste for time and explain the pattern."

**[PASTE INTO EDITOR]** `internal/meta/system.go`

```go
package meta

import (
	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/disk"
	"github.com/shirou/gopsutil/v4/host"
	"github.com/shirou/gopsutil/v4/mem"
)

type SystemInfo struct {
	OS           string        `json:"os" yaml:"os"`
	Platform     string        `json:"platform" yaml:"platform"`
	Architecture string        `json:"architecture" yaml:"architecture"`
	CPUModel     string        `json:"cpu_model" yaml:"cpu_model"`
	CPUs         int           `json:"cpus" yaml:"cpus"`
	Memory       MemoryInfo    `json:"memory" yaml:"memory"`
	Storage      []StorageInfo `json:"storage" yaml:"storage"`
}

type MemoryInfo struct {
	TotalMB     uint64  `json:"total_mb" yaml:"total_mb"`
	AvailableMB uint64  `json:"available_mb" yaml:"available_mb"`
	UsedMB      uint64  `json:"used_mb" yaml:"used_mb"`
	UsedPercent float64 `json:"used_percent" yaml:"used_percent"`
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

func CollectSystemInfo() SystemInfo {
	info := SystemInfo{
		OS:           "unknown",
		Platform:     "unknown",
		Architecture: "unknown",
		CPUs:         0,
	}

	// OS and host info
	if hostInfo, err := host.Info(); err == nil {
		info.OS = hostInfo.OS
		info.Platform = hostInfo.Platform + " " + hostInfo.PlatformVersion
		info.Architecture = hostInfo.KernelArch
	}

	// CPU info
	if cpuInfo, err := cpu.Info(); err == nil && len(cpuInfo) > 0 {
		info.CPUModel = cpuInfo[0].ModelName
		info.CPUs = int(cpuInfo[0].Cores)
	}

	// Memory info
	if memInfo, err := mem.VirtualMemory(); err == nil {
		info.Memory = MemoryInfo{
			TotalMB:     memInfo.Total / 1024 / 1024,
			AvailableMB: memInfo.Available / 1024 / 1024,
			UsedMB:      memInfo.Used / 1024 / 1024,
			UsedPercent: memInfo.UsedPercent,
		}
	}

	// Storage info
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

	return info
}
```

**[YOU SAY]**

> *[Point to struct definitions]*
>
> "Three structs: SystemInfo, MemoryInfo, StorageInfo. Notice the struct tags: `json` and `yaml`. These enable automatic marshaling to JSON and YAML formats."
>
> *[Point to CollectSystemInfo function]*
>
> "The CollectSystemInfo function does the work. Notice the pattern:"
>
> *[Scroll through function]*
>
> "1. Initialize with defaults (`unknown`, `0`)
> 2. Try to collect each piece of info
> 3. If collection fails, gracefully degrade - we keep the default
> 4. Return whatever we could collect
>
> This is graceful degradation. If we can't get CPU info, we still return OS info. The command always succeeds, even if some data is unavailable."
>
> *[Point to gopsutil calls]*
>
> "gopsutil does the heavy lifting. `host.Info()` gets OS. `cpu.Info()` gets CPU. `mem.VirtualMemory()` gets memory. `disk.Partitions()` and `disk.Usage()` get storage."
>
> "Cross-platform? Absolutely. gopsutil handles the differences between Linux, macOS, and Windows internally. On Linux, it reads `/proc`. On macOS, it uses `sysctl`. On Windows, it uses WMI. We don't care - same API everywhere."

---

### [12:00-15:00] Step 2: Command Handler

**[YOU SAY]**

> "Second file: the command handler in `cmd/ado/meta/meta.go`."

**[SHOW IN EDITOR]** Open `cmd/ado/meta/meta.go`

**[YOU SAY]**

> "This file already exists. We're adding to it."
>
> *[Scroll to NewCommand function]*
>
> "Here's where subcommands are registered. We have `info`, `env`, `features`. We're adding `system`."

**[EDIT]** Add to cmd.AddCommand:

```go
cmd.AddCommand(
	newInfoCommand(buildInfo),
	newEnvCommand(),
	newFeaturesCommand(),
	newSystemCommand(), // NEW
)
```

**[YOU SAY]**

> "One line addition. Now scroll down to add the function."

**[PASTE AT BOTTOM]**

```go
func newSystemCommand() *cobra.Command {
	var output string

	cmd := &cobra.Command{
		Use:   "system",
		Short: "Show system diagnostic information",
		Long: `Display system-level diagnostic information including OS, CPU, memory, and storage.

Useful for troubleshooting issues and sharing system information in bug reports.`,
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
	fmt.Fprintf(&b, "OS: %s\n", info.OS)
	fmt.Fprintf(&b, "Platform: %s\n", info.Platform)
	fmt.Fprintf(&b, "Architecture: %s\n", info.Architecture)

	if info.CPUModel != "" {
		fmt.Fprintf(&b, "CPU Model: %s\n", info.CPUModel)
	}
	fmt.Fprintf(&b, "CPUs: %d\n", info.CPUs)

	fmt.Fprintf(&b, "Memory: %d MB total, %d MB available (%.1f%% used)\n",
		info.Memory.TotalMB, info.Memory.AvailableMB, info.Memory.UsedPercent)

	if len(info.Storage) > 0 {
		fmt.Fprintln(&b, "Storage:")
		for _, storage := range info.Storage {
			fmt.Fprintf(&b, "  %s: %d GB total, %d GB used (%.1f%%)\n",
				storage.Mountpoint, storage.TotalGB, storage.UsedGB, storage.UsedPercent)
		}
	}

	return b.String()
}
```

**[YOU SAY]**

> *[Point to newSystemCommand]*
>
> "This follows the exact same pattern as `newInfoCommand`, `newEnvCommand`, `newFeaturesCommand`."
>
> *[Point to Cobra command setup]*
>
> "1. Create a Cobra command with Use, Short, Long descriptions
> 2. Add RunE function - this is what executes
> 3. Add `--output` flag
>
> *[Point to RunE function]*
>
> "Inside RunE:
> 1. Call our collector: `CollectSystemInfo()`
> 2. Parse the output format (text/json/yaml)
> 3. Use `ui.PrintOutput` - this is a helper that handles all three formats
> 4. Provide a custom text formatter
>
> *[Point to formatSystemInfo]*
>
> "The text formatter builds human-readable output. JSON and YAML are automatic thanks to our struct tags."
>
> "This pattern is crucial: Every meta subcommand works the same way. Consistency makes the codebase maintainable."

---

### [15:00-17:00] Step 3: Tests

**[YOU SAY]**

> "Third file: tests in `cmd/ado/meta/meta_test.go`."

**[SHOW IN EDITOR]** Open `cmd/ado/meta/meta_test.go`

**[PASTE AT BOTTOM]**

```go
func TestSystemCommand(t *testing.T) {
	tests := []struct {
		name       string
		args       []string
		wantOutput []string
		wantErr    bool
	}{
		{
			name:       "text output",
			args:       []string{"system"},
			wantOutput: []string{"OS:", "Platform:", "Architecture:", "CPUs:", "Memory:"},
			wantErr:    false,
		},
		{
			name:       "json output",
			args:       []string{"system", "--output", "json"},
			wantOutput: []string{`"os"`, `"platform"`, `"architecture"`, `"cpus"`, `"memory"`},
			wantErr:    false,
		},
		{
			name:       "yaml output",
			args:       []string{"system", "--output", "yaml"},
			wantOutput: []string{"os:", "platform:", "architecture:", "cpus:", "memory:"},
			wantErr:    false,
		},
		{
			name:       "invalid output format",
			args:       []string{"system", "--output", "xml"},
			wantOutput: nil,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := NewCommand(internalmeta.CurrentBuildInfo())
			cmd.SetArgs(tt.args)

			var buf bytes.Buffer
			cmd.SetOut(&buf)
			cmd.SetErr(&buf)

			err := cmd.Execute()
			if (err != nil) != tt.wantErr {
				t.Errorf("Execute() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				output := buf.String()
				for _, want := range tt.wantOutput {
					if !strings.Contains(output, want) {
						t.Errorf("Execute() output = %v, want to contain %v", output, want)
					}
				}
			}
		})
	}
}
```

**[YOU SAY]**

> *[Point to test table]*
>
> "Table-driven tests - the Go standard. We define test cases as a slice of structs:"
> - Test name
> - Command args
> - Expected output strings
> - Whether we expect an error
>
> *[Point to test loop]*
>
> "Then we loop through test cases using `t.Run`. This gives us clean sub-test output."
>
> *[Point to test body]*
>
> "Each test:
> 1. Creates the command
> 2. Sets arguments
> 3. Captures output to a buffer
> 4. Executes the command
> 5. Verifies output contains expected strings
>
> "We test all three output formats (text, json, yaml) plus error cases (invalid format)."
>
> *[Emphasize]*
>
> "This is an integration test. We're testing the real command with real system info, not mocking. Why? Because system info should always be available. This validates the full stack."

**[YOU SAY]**

> "Alright, three files done:
> - `internal/meta/system.go` - Data collection
> - `cmd/ado/meta/meta.go` - Command handler
> - `cmd/ado/meta/meta_test.go` - Tests
>
> Let's validate our implementation."

---

## [17:00-20:00] Phase 4: Testing & Validation

**[STAGE DIRECTION]** Terminal focus. Energy: HIGH (this is the exciting part)

**[TERMINAL ACTION]** Bottom-right pane (Make)

**[YOU SAY]**

> "Before we create a PR, we validate locally. The ado project has strict requirements. Let's see if we meet them."

**[TYPE]**
```bash
make go.test
```

**[WAIT FOR OUTPUT]**

**[EXPECTED OUTPUT]**
```
=== RUN   TestSystemCommand
=== RUN   TestSystemCommand/text_output
=== RUN   TestSystemCommand/json_output
=== RUN   TestSystemCommand/yaml_output
=== RUN   TestSystemCommand/invalid_output_format
--- PASS: TestSystemCommand (0.15s)
PASS
coverage: 83.2% of statements
ok      github.com/anowarislam/ado/cmd/ado/meta 0.234s
```

**[YOU SAY]**

> *[Point to PASS]*
>
> "All tests pass! ✓"
>
> *[Point to coverage]*
>
> "Coverage: 83.2%. We need minimum 80%, so we're good."
>
> "But let's verify the coverage check explicitly."

**[TYPE]**
```bash
make go.test.cover.check
```

**[EXPECTED OUTPUT]**
```
Total coverage: 83.2%
✓ Coverage meets threshold (80%)
```

**[YOU SAY]**

> "Coverage threshold met ✓"
>
> "This check is enforced in two places:
> 1. Pre-push hook - catches it locally before you push
> 2. CI - catches it on the server before merge
>
> You literally cannot merge code with less than 80% coverage. The build will fail."

**[TYPE]**
```bash
make go.vet
```

**[EXPECTED OUTPUT]**
```
go vet ./...
✓ No issues found
```

**[YOU SAY]**

> "Static analysis clean ✓"

**[TYPE]**
```bash
make go.build
```

**[EXPECTED OUTPUT]**
```
Building ado...
✓ Build successful: ./ado
```

**[YOU SAY]**

> "Build successful ✓"
>
> "Now the moment of truth - let's actually run it!"

**[TYPE]**
```bash
./ado meta system
```

**[EXPECTED OUTPUT]**
```
OS: darwin
Platform: macOS 14.2.1
Architecture: arm64
CPU Model: Apple M2 Pro
CPUs: 10
Memory: 16384 MB total, 8245 MB available (49.7% used)
Storage:
  /: 494 GB total, 123 GB used (24.9%)
```

**[YOU SAY]**

> *[Gesture excitedly]*
>
> "It works! Look at that output - exactly as specified in our spec."
>
> "Let's try JSON output."

**[TYPE]**
```bash
./ado meta system --output json
```

**[EXPECTED OUTPUT]**
```json
{
  "os": "darwin",
  "platform": "macOS 14.2.1",
  "architecture": "arm64",
  "cpu_model": "Apple M2 Pro",
  "cpus": 10,
  "memory": {
    "total_mb": 16384,
    "available_mb": 8245,
    "used_mb": 8139,
    "used_percent": 49.7
  },
  "storage": [...]
}
```

**[YOU SAY]**

> "Valid JSON ✓ Perfect for bug reports or automation."
>
> "And YAML?"

**[TYPE]**
```bash
./ado meta system --output yaml
```

**[YOU SAY]**

> *[Scroll through YAML output]*
>
> "YAML ✓"
>
> "So we've validated:
> - Tests pass
> - Coverage meets threshold (83.2% > 80%)
> - No lint errors
> - Build succeeds
> - All three output formats work
> - Output matches our spec
>
> We're ready to create a PR."

**[PAUSE FOR QUESTIONS]**

> "Questions on testing? ... Cool, let's ship it."

---

## [20:00-23:00] Phase 5: Create Pull Request

**[STAGE DIRECTION]** Terminal and browser. Energy: MEDIUM-HIGH

**[TERMINAL ACTION]** Git pane

**[YOU SAY]**

> "Time to create a PR. First, let's review what changed."

**[TYPE]**
```bash
git status
```

**[EXPECTED OUTPUT]**
```
On branch feat/meta-system
modified:   cmd/ado/meta/meta.go
modified:   cmd/ado/meta/meta_test.go
modified:   go.mod
modified:   go.sum
new file:   internal/meta/system.go
```

**[YOU SAY]**

> "Five files changed. Let's stage them."

**[TYPE]**
```bash
git add cmd/ado/meta/meta.go \
        cmd/ado/meta/meta_test.go \
        internal/meta/system.go \
        go.mod go.sum
```

**[YOU SAY]**

> "Now the commit message. This is important."

**[TYPE]**
```bash
git commit -m "feat(meta): add system subcommand for diagnostics

Implements 'ado meta system' command to display system-level
diagnostic information including OS, CPU, memory, and storage.

Features:
- Displays OS, platform, architecture, CPU model/count
- Shows memory stats (total, available, used percentage)
- Lists storage volumes with usage statistics
- Supports text/json/yaml output formats
- Cross-platform (Linux/macOS/Windows via gopsutil)
- Graceful degradation for unavailable system info

Dependencies:
- github.com/shirou/gopsutil/v4

Testing:
- Table-driven tests for all output formats
- Coverage: 83.2% (exceeds 80% requirement)

Closes #53"
```

**[YOU SAY]**

> "Notice the format:
>
> `feat(meta):` - Conventional commit. `feat` = minor version bump. `meta` = scope.
>
> Body explains what this does, how it works, what was tested.
>
> `Closes #53` - Links to the issue. When merged, issue auto-closes."

**[TYPE]**
```bash
git push -u origin feat/meta-system
```

**[TYPE]**
```bash
gh pr create \
  --title "feat(meta): add system subcommand for diagnostics" \
  --body "Implements \`ado meta system\` per approved spec.

Closes #53

## Summary
- Adds system diagnostic subcommand
- Supports text/json/yaml output
- Cross-platform via gopsutil

## Testing
✅ All tests pass (coverage: 83.2%)
✅ Manual testing on macOS" \
  --web
```

**[BROWSER ACTION]** PR opens

**[YOU SAY]**

> *[Point to PR]*
>
> "PR created! Immediately, notice:"
>
> *[Point to checks section]*
>
> "1. CI is running. Conventional Commits check ✓ already passed."
>
> *[Point to reviewers]*
>
> "2. CODEOWNERS assigned @anowarislam. Why? Because I modified `cmd/ado/meta/` and `internal/meta/`, and both are owned by @anowarislam per the CODEOWNERS file."
>
> *[Point to linked issue]*
>
> "3. Issue #53 is linked. When this merges, it auto-closes."
>
> "Let's watch the CI run."

---

## [23:00-26:00] Phase 6: CI/CD Pipeline

**[STAGE DIRECTION]** Browser focus. Energy: HIGH (this is impressive)

**[BROWSER ACTION]** Click on "Details" for CI checks

**[YOU SAY]**

> "The CI pipeline is our safety net. Let me explain what's running."

**[POINT TO CHECKS]**

**[YOU SAY]**

> "**Conventional Commits** - Already passed ✓
> Validates PR title follows `type(scope): description` format.
> Why? Because release-please uses this to determine version bumps and generate changelogs.
>
> *[Click on Go CI workflow]*
>
> "**Go CI** - This is the big one. Multiple steps:"

**[EXPAND WORKFLOW STEPS]**

**[YOU SAY]**

> "1. **Format Check** - Runs `gofmt`. If code isn't formatted, build fails.
>
> 2. **Dependencies** - Verifies `go.mod` and `go.sum` are in sync.
>
> 3. **Go Vet** - Static analysis. Catches common mistakes like printf format errors, unreachable code, suspicious constructs.
>
> 4. **Tests with Coverage** - This is critical:"

**[POINT TO TEST OUTPUT]**

> "- Runs all tests with race detector enabled
> - Generates coverage report
> - **Enforces 80% minimum** - If coverage is 79.9%, build fails
> - Uploads to Codecov for tracking
>
> Our coverage is 83.2%, so this passes ✓
>
> 5. **Build** - Compiles the binary. Ensures it actually builds."

**[POINT TO OTHER CHECKS]**

**[YOU SAY]**

> "**Documentation Build** - Builds MkDocs with `--strict` flag. Catches broken links, invalid markdown.
>
> **Docker Test** - Tests the GoReleaser Dockerfile. Ensures release builds will succeed.
>
> All checks must pass to merge. No exceptions."

**[ALL CHECKS PASS]**

**[YOU SAY]**

> *[Gesture to all green checkmarks]*
>
> "All checks passing ✓✓✓
>
> This automation is crucial. Code reviewers don't need to check:
> - If tests pass (CI did it)
> - If coverage is sufficient (CI did it)
> - If code is formatted (CI did it)
> - If code lints (CI did it)
> - If it builds (CI did it)
>
> Reviewers focus on: Architecture, design, correctness, maintainability. The mechanical stuff is automated."

**[PAUSE FOR QUESTIONS]**

> "Questions on CI? ... Great."

---

## [26:00-28:00] Phase 7: Review & Merge

**[STAGE DIRECTION]** Browser. Energy: MEDIUM

**[BROWSER ACTION]** Navigate to "Files changed" tab

**[YOU SAY]**

> "In a real workflow, we'd wait 1-2 days for code owner review. Let me simulate what a reviewer would look for."

**[SCROLL THROUGH FILES CHANGED]**

**[YOU SAY]**

> "**Architecture & Design:**
> - ✓ Follows existing meta command pattern (info, env, features)
> - ✓ Uses established library (gopsutil - 10k stars)
> - ✓ Matches approved spec exactly
>
> **Code Quality:**
> - ✓ Table-driven tests (Go best practice)
> - ✓ Graceful error handling
> - ✓ Consistent formatting
> - ✓ Good test coverage (83.2%)
>
> **Documentation:**
> - ✓ Clear command help text
> - ✓ Conventional commit message
> - ✓ PR description with examples
>
> **Security:**
> - ✓ No elevated permissions required
> - ✓ Read-only operation (no user input)
> - ✓ Trusted dependency (widely adopted)
>
> Looks good to me! Approve ✓"

**[BROWSER ACTION]** Scroll to merge button

**[YOU SAY]**

> "Merge button is enabled because:
> - All CI checks passed ✓
> - Code owner approved ✓
> - Conversations resolved ✓
> - Branch is up to date ✓
>
> Let's merge."

**[CLICK]** Squash and merge

**[BROWSER ACTION]** Confirm merge

**[YOU SAY]**

> *[Point to merged PR]*
>
> "Merged ✓
>
> Notice:
> - All commits squashed into one clean commit on main
> - Branch automatically deleted
> - Issue #53 automatically closed
>
> *[Navigate to Issue #53]*
>
> "See? Closed automatically with a reference to the PR."
>
> "Now watch what happens next."

---

## [28:00-30:00] Phase 8: Automated Release

**[STAGE DIRECTION]** Browser, Actions tab. Energy: VERY HIGH (finale!)

**[BROWSER ACTION]** Navigate to Actions tab

**[YOU SAY]**

> "The merge to main triggers automated release. Here's how it works:"

**[DRAW OR GESTURE TO SCREEN]**

**[YOU SAY]**

> "**Step 1: Release Please Analyzes Commits**
>
> Release Please is a bot that reads our merged commit:
> - Title: `feat(meta): add system subcommand`
> - Type: `feat` = new feature
> - Decision: This is a MINOR version bump
> - Current version: 1.2.0
> - Next version: 1.3.0
>
> **Step 2: Release Please Creates/Updates Release PR**
>
> It opens a PR titled 'chore(main): release 1.3.0' with:
> - Updated CHANGELOG.md (auto-generated from commit messages)
> - Updated version in manifest
> - Release notes
>
> **Step 3: Maintainer Reviews Release PR**
>
> Quick review (usually < 5 minutes), then merge.
>
> **Step 4: When Release PR Merges**
>
> - Creates git tag: `v1.3.0`
> - Creates GitHub Release
> - Triggers GoReleaser
>
> **Step 5: GoReleaser Builds Everything**"

**[LIST ON FINGERS]**

**[YOU SAY]**

> "GoReleaser:
> 1. Builds for 6 platforms:
>    - linux-amd64, linux-arm64
>    - darwin-amd64, darwin-arm64 (Intel and Apple Silicon Macs)
>    - windows-amd64, windows-arm64
>
> 2. Creates .tar.gz and .zip archives
>
> 3. Generates SHA256 checksums
>
> 4. Creates SLSA attestations - supply chain security
>
> 5. Uploads to GitHub Release
>
> 6. Builds Docker images (multi-arch)
>
> 7. Signs containers with Sigstore/cosign
>
> 8. Pushes to GitHub Container Registry
>
> All of this - completely automated. From one conventional commit."

**[PAUSE FOR EFFECT]**

**[YOU SAY]**

> "Timeline:
> - Now: feat/meta-system merged to main
> - 30 seconds: Release Please creates Release PR
> - Same day: Maintainer merges Release PR
> - 5 minutes: GoReleaser completes
> - Done: v1.3.0 available for download
>
> From merge to release: Same day. Zero manual work."

---

## [30:00] Wrap-Up & Summary

**[STAGE DIRECTION]** Step back from screen. Energy: HIGH, wrap-up mode

**[YOU SAY]**

> *[Gesture broadly]*
>
> "So what did we accomplish in 30 minutes?"
>
> *[Count on fingers]*
>
> "✓ Issue #53 - Defined the problem
> ✓ Spec PR - Documented the solution
> ✓ Implementation PR - Built the code
> ✓ Tests - 83.2% coverage
> ✓ CI/CD - Automated validation
> ✓ Merge - Clean squash merge
> ✓ Release - Automated versioning
>
> From idea to shippable feature. In 30 minutes."
>
> *[Pause]*
>
> "But here's what really matters - the workflow benefits:"

**[COUNT ON FINGERS AGAIN]**

**[YOU SAY]**

> "1. **Spec-first prevents wasted effort**
> We spent 3 minutes writing a spec. That prevented potentially days of rework because everyone agreed on WHAT before we invested in HOW.
>
> 2. **Conventional commits enable automation**
> One commit format. Automatic semantic versioning. Automatic changelogs. Automatic release notes. Zero manual overhead.
>
> 3. **CI enforces quality non-negotiably**
> 80% coverage isn't a suggestion, it's enforced. Tests must pass. Code must lint. Build must succeed. Even administrators can't bypass this.
>
> 4. **CODEOWNERS ensures expertise**
> The right people review the right code. Automatically. No one needs to remember who owns what.
>
> 5. **Automated release removes toil**
> From merge to release: Completely automated. Multi-platform builds. Checksums. Security attestations. Container signing. All automatic.
>
> 6. **Security is built-in**
> SLSA attestations prove what was built, by whom, from what source. Signed containers prevent tampering. This isn't optional, it's the default."

**[PAUSE]**

**[YOU SAY]**

> *[Slow down, emphasize]*
>
> "This workflow lets us ship quality code fast, with confidence, and with security by default.
>
> The feature took 2-3 hours to build. But the workflow makes that 2-3 hours count. No wasted effort. No manual release overhead. No quality compromises.
>
> *[Smile]*
>
> "And that's why I love this workflow."

**[TIMER: STOP]**

---

## [30:00+] Q&A Session

**[STAGE DIRECTION]** Open posture. Energy: CALM, ready to discuss

**[YOU SAY]**

> "Alright, we have time for questions. What would you like to know?"

**[WAIT FOR QUESTIONS]**

---

### Q&A Response Templates

**If asked: "Why spec before code?"**

> "Great question. Two reasons:
>
> 1. **Alignment** - A spec takes 30 minutes to write but ensures everyone agrees on WHAT before we invest days in HOW. If the spec is wrong, we waste 30 minutes. If the code is wrong because we skipped the spec, we waste days.
>
> 2. **Design review** - It's easier to review a spec than review code. Specs are small, readable, focused on behavior not implementation. Code is large, detailed, mixed with implementation concerns. Catching design issues in a spec saves massive rework.
>
> In practice, spec-first has saved us countless hours of rework."

---

**If asked: "What if we disagree with the spec during implementation?"**

> "Happens all the time! If you discover the spec is wrong:
>
> 1. Open an issue: 'Spec for X is incomplete/incorrect'
> 2. Create a PR to update the spec
> 3. Get approval on the spec change
> 4. Then update implementation to match
>
> This keeps spec and code in sync. The spec is the source of truth. If they diverge, we have a problem.
>
> Yes, it feels like extra work. But it prevents a worse problem: code and documentation drifting apart over time."

---

**If asked: "Why 80% coverage? Isn't that arbitrary?"**

> "80% is a balance:
>
> Too low (e.g., 50%) allows significant untested code. Too high (e.g., 95%) forces testing trivial code like getters and setters.
>
> 80% is the 'sweet spot' - it's the same threshold used by Kubernetes, many Go projects, and it's backed by research showing diminishing returns above 80%.
>
> It's enforced in two places:
> - Pre-push hook (catches locally)
> - CI (catches before merge)
>
> The goal isn't 100% coverage. It's ensuring critical logic is tested."

---

**If asked: "Can we skip the spec for small changes?"**

> "Yes! Check the decision tree in `docs/workflow.md`:
>
> - Bug fix? → No spec, go straight to fix + tests
> - New command? → Spec required
> - Architectural change? → ADR + Spec
>
> The workflow is pragmatic. Small bug fixes don't need specs. Major features do."

---

**If asked: "What happens if CI fails after merge?"**

> "Can't happen. Branch protection prevents merging unless CI passes. Even administrators can't bypass this without explicit override.
>
> If somehow bad code reaches main (e.g., emergency override), we:
> 1. Revert the commit immediately
> 2. Fix the issue in a new PR
> 3. Merge after CI passes
>
> The main branch is sacred. It must always be shippable."

---

**If asked: "How do you test cross-platform if you only have macOS?"**

> "GitHub Actions runs a matrix:
> - Linux (Ubuntu latest)
> - macOS (latest)
> - Windows (latest)
>
> My local tests run on macOS. But CI validates all three platforms before merge.
>
> You can also test locally with Docker for Linux, or use GitHub Codespaces for a cloud environment."

---

**If asked: "Doesn't this slow down development?"**

> "Short answer: No, it speeds up development.
>
> Long answer: The spec adds 30-60 minutes upfront. But it saves:
> - Days of rework from misalignment
> - Hours of back-and-forth in code review
> - Weeks of technical debt from poor design
>
> The CI automation adds ~5 minutes per PR. But it saves:
> - Hours of manual testing
> - Days of debugging production issues
> - Weeks of security incidents
>
> The release automation adds zero time (it's automatic). But it saves:
> - Hours of manual release work
> - Days of manual changelog writing
> - Weeks of version confusion
>
> Net: We ship faster with higher quality."

---

**If asked: "What if the library we chose has a security vulnerability?"**

> "Good question. We use Dependabot:
>
> - Scans dependencies daily for known vulnerabilities
> - Opens PRs automatically to update vulnerable dependencies
> - CI validates the update doesn't break anything
> - Merge and release
>
> Example: If gopsutil has a CVE tomorrow, Dependabot opens a PR tomorrow, we merge it after CI passes, release goes out automatically.
>
> Median time from vulnerability disclosure to patched release: 1-2 days."

---

**If asked: "Can I contribute to ado?"**

> "Absolutely! Check `docs/contributing.md` for the guide.
>
> The workflow we just demonstrated is how all contributions work:
> 1. Open an issue
> 2. Write a spec (if needed)
> 3. Implement with tests
> 4. Create PR
> 5. Code owner reviews
> 6. Merge and release
>
> First-time contributors: Look for issues labeled `good-first-issue`. These are well-scoped and have clear acceptance criteria.
>
> Questions? Open an issue or discussion on GitHub. We're friendly!"

---

## Closing Remarks

**[IF TIME PERMITS]**

**[YOU SAY]**

> "One last thing before we wrap."
>
> *[Pause for attention]*
>
> "This demo wasn't about the `meta system` command. That's just a vehicle. This demo was about showing you a workflow that:
>
> - Values specs over premature coding
> - Automates quality enforcement
> - Removes manual release toil
> - Builds security in by default
>
> You can apply this workflow to any project. The tools are open source:
> - release-please for versioning
> - GoReleaser for builds
> - GitHub Actions for CI
> - CODEOWNERS for review routing
>
> The workflow is documented in `docs/workflow.md`.
>
> *[Smile]*
>
> "If you take one thing from this demo, let it be this: Invest in workflow automation. The upfront cost is small. The long-term value is enormous.
>
> Thank you!"

**[APPLAUSE]**

**[END]**

---

## Post-Demo: Cleanup

**[AFTER AUDIENCE LEAVES]**

```bash
# Return to main
git checkout main
git pull origin main

# Delete demo branches (if not merged)
git branch -D spec/meta-system feat/meta-system

# Clean working directory
git status  # Should be clean

# Optional: Close issue if demo-only
gh issue close 53 --comment "Demo completed successfully. Feature not being implemented at this time."
```

---

## Presenter Notes

### Energy Management
- **Start HIGH**: Hook them immediately
- **Middle MEDIUM**: Educational, calm, focused
- **End HIGH**: Exciting finale with automation

### Timing Tips
- If running long: Skip typing code live (paste and explain)
- If running short: Add more Q&A time
- Aim to finish at 28-29 minutes to leave buffer

### Common Mistakes
- Don't type everything (too slow)
- Don't skip explaining WHY (that's the point)
- Don't rush (better to cut content than rush)

### What Makes This Work
- Clear structure (phases)
- Real demo (not slides)
- Explains rationale (not just steps)
- Anticipates questions
- Ends with impact

---

**GOOD LUCK! 🚀**
