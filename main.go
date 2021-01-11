package main

import (
	"fmt"
	"os"

	"github.com/sensu-community/sensu-plugin-sdk/sensu"
	"github.com/sensu/sensu-go/types"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/load"
)

// Config represents the check plugin config.
type Config struct {
	sensu.PluginConfig
	WarningMultiplier   float64
	CriticalMultiplier  float64
	CountLogicalCPU     bool
	CompareAllIntervals bool
}

var (
	nCPUPhysical int
	nCPULogical  int
	plugin       = Config{
		PluginConfig: sensu.PluginConfig{
			Name:     "check-load",
			Short:    "Sensu Load Average Check",
			Keyspace: "sensu.io/plugins/check-load/config",
		},
	}

	options = []*sensu.PluginConfigOption{
		{
			Path:      "warning-multiplier",
			Argument:  "warning-multiplier",
			Shorthand: "w",
			Default:   float64(1.5),
			Usage:     "The warning threshold multiplier (# CPUs x multiplier)",
			Value:     &plugin.WarningMultiplier,
		},
		{
			Path:      "critical-multiplier",
			Argument:  "critical-multiplier",
			Shorthand: "c",
			Default:   float64(2),
			Usage:     "The critical threshold multiplier (# CPUs x multiplier)",
			Value:     &plugin.CriticalMultiplier,
		},
		{
			Path:      "count-logical-cpu",
			Argument:  "count-logical-cpu",
			Shorthand: "l",
			Default:   false,
			Usage:     "Include Logical CPUs (e.g. hyperthreading) in factoring thresholds",
			Value:     &plugin.CountLogicalCPU,
		},
		{
			Path:      "compare-all-intervals",
			Argument:  "compare-all-intervals",
			Shorthand: "a",
			Default:   false,
			Usage:     "Compare thresholds to all (1m, 5m, 15m) load averages",
			Value:     &plugin.CompareAllIntervals,
		},
	}
)

func main() {
	var err error
	nCPUPhysical, err = cpu.Counts(false)
	if err != nil {
		fmt.Printf("%s CRITICAL: failed to get number of CPUs, error: %v\n", plugin.PluginConfig.Name, err)
		os.Exit(sensu.CheckStateCritical)
	}
	nCPULogical, err = cpu.Counts(true)
	if err != nil {
		fmt.Printf("%s CRITICAL: failed to get number of CPUs, error: %v\n", plugin.PluginConfig.Name, err)
		os.Exit(sensu.CheckStateCritical)
	}
	check := sensu.NewGoCheck(&plugin.PluginConfig, options, checkArgs, executeCheck, false)
	check.Execute()
}

func checkArgs(event *types.Event) (int, error) {
	if plugin.CriticalMultiplier <= float64(0) {
		return sensu.CheckStateWarning, fmt.Errorf("--critical-multiplier must be greater than 0")
	}
	if plugin.WarningMultiplier <= float64(0) {
		return sensu.CheckStateWarning, fmt.Errorf("--warning-multiplier must be greater than 0")
	}
	if plugin.CriticalMultiplier <= plugin.WarningMultiplier {
		return sensu.CheckStateWarning, fmt.Errorf("--critical-multiplier must be greater than--warning-multiplier")
	}
	return sensu.CheckStateOK, nil
}

func executeCheck(event *types.Event) (int, error) {
	var (
		warningThreshold  float64
		criticalThreshold float64
	)

	if plugin.CountLogicalCPU {
		warningThreshold = float64(nCPULogical) * plugin.WarningMultiplier
		criticalThreshold = float64(nCPULogical) * plugin.CriticalMultiplier
	} else {
		warningThreshold = float64(nCPUPhysical) * plugin.WarningMultiplier
		criticalThreshold = float64(nCPUPhysical) * plugin.CriticalMultiplier
	}

	avg, err := load.Avg()
	if err != nil {
		fmt.Printf("%s CRITICAL: failed to get load average, error: %v\n", plugin.PluginConfig.Name, err)
		return sensu.CheckStateCritical, nil
	}

	perfData := fmt.Sprintf("load1=%.2f, load5=%.2f, load15=%.2f", avg.Load1, avg.Load5, avg.Load15)

	if plugin.CompareAllIntervals {
		if avg.Load1 >= criticalThreshold && avg.Load5 >= criticalThreshold && avg.Load15 >= criticalThreshold {
			fmt.Printf("%s CRITICAL: load avg %.2f, %.2f, %.2f | %s\n", plugin.PluginConfig.Name, avg.Load1, avg.Load5, avg.Load15, perfData)
			return sensu.CheckStateCritical, nil
		} else if avg.Load1 >= warningThreshold && avg.Load5 >= warningThreshold && avg.Load15 >= warningThreshold {
			fmt.Printf("%s WARNING: load avg %.2f, %.2f, %.2f | %s\n", plugin.PluginConfig.Name, avg.Load1, avg.Load5, avg.Load15, perfData)
			return sensu.CheckStateWarning, nil
		}
		fmt.Printf("%s OK: load avg %.2f, %.2f, %.2f | %s\n", plugin.PluginConfig.Name, avg.Load1, avg.Load5, avg.Load15, perfData)
		return sensu.CheckStateOK, nil
	}

	if avg.Load1 >= criticalThreshold {
		fmt.Printf("%s CRITICAL: 1m load avg %.2f | %s\n", plugin.PluginConfig.Name, avg.Load1, perfData)
		return sensu.CheckStateCritical, nil
	} else if avg.Load1 >= warningThreshold {
		fmt.Printf("%s WARNING: 1m load avg %.2f | %s\n", plugin.PluginConfig.Name, avg.Load1, perfData)
		return sensu.CheckStateWarning, nil
	}
	fmt.Printf("%s OK: 1m load avg %.2f | %s\n", plugin.PluginConfig.Name, avg.Load1, perfData)
	return sensu.CheckStateOK, nil
}
