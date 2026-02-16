package merger

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/jeremy-kr/ccfg/internal/model"
	"github.com/jeremy-kr/ccfg/internal/parser"
)

// SourcedValue represents a value along with its origin scope.
type SourcedValue struct {
	Key   string      // Dot-notation path (e.g., "permissions.allow")
	Value any         // Actual value
	Scope model.Scope // Scope the value originates from
}

// MergedConfig holds the result of merging settings from all scopes.
type MergedConfig struct {
	Values []SourcedValue // Flat list of key-value pairs
}

// Merge merges JSON config files from a ScanResult according to priority.
// Priority: Project > User > Managed.
func Merge(result *model.ScanResult) *MergedConfig {
	merged := make(map[string]SourcedValue)

	// Apply from lowest priority first (later entries overwrite earlier ones)
	applyScope(merged, result.Managed, model.ScopeManaged)
	applyScope(merged, result.User, model.ScopeUser)
	applyScope(merged, result.Project, model.ScopeProject)

	// Sort by key
	values := make([]SourcedValue, 0, len(merged))
	for _, v := range merged {
		values = append(values, v)
	}
	sort.Slice(values, func(i, j int) bool {
		return values[i].Key < values[j].Key
	})

	return &MergedConfig{Values: values}
}

func applyScope(merged map[string]SourcedValue, files []model.ConfigFile, scope model.Scope) {
	for _, f := range files {
		if !f.Exists || f.FileType != model.FileTypeJSON || f.IsDir {
			continue
		}
		if f.Category != model.CategorySettings {
			continue
		}

		data, err := os.ReadFile(f.Path)
		if err != nil {
			continue
		}

		cleaned := parser.StripJSONC(string(data))
		var obj map[string]any
		if err := json.Unmarshal([]byte(cleaned), &obj); err != nil {
			continue
		}

		flatten("", obj, scope, merged)
	}
}

func flatten(prefix string, obj map[string]any, scope model.Scope, out map[string]SourcedValue) {
	for k, v := range obj {
		key := k
		if prefix != "" {
			key = prefix + "." + k
		}
		switch val := v.(type) {
		case map[string]any:
			flatten(key, val, scope, out)
		default:
			out[key] = SourcedValue{Key: key, Value: v, Scope: scope}
			_ = val
		}
	}
}

// Render formats the merged config into a human-readable string.
func (mc *MergedConfig) Render() string {
	if len(mc.Values) == 0 {
		return "(no settings to merge)"
	}

	var b strings.Builder
	b.WriteString("Merged Settings (Project > User > Managed)\n")
	b.WriteString(strings.Repeat("─", 50) + "\n\n")

	for _, v := range mc.Values {
		valStr := fmt.Sprintf("%v", v.Value)
		if len(valStr) > 60 {
			valStr = valStr[:57] + "..."
		}
		b.WriteString(fmt.Sprintf("  %-35s = %-20s [%s]\n", v.Key, valStr, v.Scope))
	}

	return b.String()
}
