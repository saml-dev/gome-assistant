package gomeassistant

import (
	"fmt"
	"path/filepath"
	"strings"
)

const (
	defaultWebUIPort    = 52189
	localWebUIAddress   = "127.0.0.1:52189"
	exposedWebUIAddress = "0.0.0.0:52189"
)

// EnableMetricsCollection configures durable metrics collection. Resources
// are opened by Start, not by this method.
func (app *App) EnableMetricsCollection(databasePath string) error {
	app.configurationMu.Lock()
	defer app.configurationMu.Unlock()

	if app.configurationFrozen {
		return ErrConfigurationFrozen
	}
	databasePath = strings.TrimSpace(databasePath)
	if databasePath == "" || !filepath.IsAbs(databasePath) {
		return fmt.Errorf("%w: metrics database path must be absolute", ErrInvalidArgs)
	}
	app.metricsDatabasePath = databasePath
	return nil
}

// EnableWebUI configures the read-only embedded web server. When expose is
// false, it listens only on loopback; when true, it listens on all IPv4
// interfaces. Start tries port 52189 and increments it until one is free.
// Resources are not bound until Start.
func (app *App) EnableWebUI(expose bool) error {
	app.configurationMu.Lock()
	defer app.configurationMu.Unlock()

	if app.configurationFrozen {
		return ErrConfigurationFrozen
	}
	app.webUIAddress = localWebUIAddress
	if expose {
		app.webUIAddress = exposedWebUIAddress
	}
	return nil
}

func (app *App) requireRegistrationOpen() {
	if app.registryFrozen {
		panic(ErrConfigurationFrozen)
	}
}

// freezeConfiguration is intentionally the first operation performed by
// Start. Once metrics are enabled it freezes both configuration and the
// automation catalog, even when that Start attempt later fails.
func (app *App) freezeConfiguration() error {
	app.configurationMu.Lock()
	defer app.configurationMu.Unlock()

	if !app.configurationFrozen {
		app.configurationFrozen = true
		app.registryFrozen = app.metricsDatabasePath != ""
		if app.registryFrozen {
			app.catalog, app.catalogErr = app.buildAutomationCatalog()
		}
	}

	if app.webUIAddress != "" && app.metricsDatabasePath == "" {
		return fmt.Errorf("web UI requires metrics collection")
	}
	return app.catalogErr
}
