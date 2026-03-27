package plugin

import (
	"fmt"

	"github.com/openagent-md/paralleltestctx/pkg/paralleltestctx"
	"github.com/golangci/plugin-module-register/register"
	"golang.org/x/tools/go/analysis"
)

func init() {
	register.Plugin("paralleltestctx", New)
}

func New(settings any) (register.LinterPlugin, error) {
	s, ok := settings.(map[string]any)
	if !ok {
		return nil, fmt.Errorf("expect paralleltestctx's configurations to be a map[string]string, got %T", settings)
	}
	conf := make(map[string]string, len(s))
	for k, v := range s {
		vStr, ok := v.(string)
		if !ok {
			return nil, fmt.Errorf("expect paralleltestctx's configuration values for %q to be strings, got %T", k, v)
		}
		conf[k] = vStr
	}
	return &parallelTestCtxPlugin{conf: conf}, nil
}

type parallelTestCtxPlugin struct {
	conf map[string]string
}

func (p *parallelTestCtxPlugin) BuildAnalyzers() ([]*analysis.Analyzer, error) {
	analyzer := paralleltestctx.Analyzer()
	for k, v := range p.conf {
		if err := analyzer.Flags.Set(k, v); err != nil {
			return nil, fmt.Errorf("set config flag %s with %s: %w", k, v, err)
		}
	}
	return []*analysis.Analyzer{analyzer}, nil
}

func (p *parallelTestCtxPlugin) GetLoadMode() string { return register.LoadModeTypesInfo }
