package gomeassistant

import (
	"context"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func newConfigurationTestApp(t *testing.T) *App {
	t.Helper()
	app, err := NewApp(context.Background(), NewAppRequest{
		URL:         "http://127.0.0.1:1",
		HAAuthToken: "token",
	})
	require.NoError(t, err)
	return app
}

func TestMetricsAndWebConfigurationValidationAndReplacement(t *testing.T) {
	app := newConfigurationTestApp(t)
	databaseOne := filepath.Join(t.TempDir(), "one.db")
	databaseTwo := filepath.Join(t.TempDir(), "two.db")

	for _, databasePath := range []string{"", "relative.db", "~", "~/gomeassistant/metrics.db", `~\gomeassistant\metrics.db`, "~someone/metrics.db"} {
		err := app.EnableMetricsCollection(databasePath)
		require.ErrorIs(t, err, ErrInvalidArgs)
		require.ErrorContains(t, err, "metrics database path must be absolute")
	}
	require.NoError(t, app.EnableMetricsCollection("  "+databaseOne+"  "))
	require.Equal(t, databaseOne, app.metricsDatabasePath)
	require.NoError(t, app.EnableMetricsCollection(databaseTwo))
	require.Equal(t, databaseTwo, app.metricsDatabasePath)

	require.NoError(t, app.EnableWebUI(false))
	require.Equal(t, localWebUIAddress, app.webUIAddress)
	require.NoError(t, app.EnableWebUI(true))
	require.Equal(t, exposedWebUIAddress, app.webUIAddress)
}

func TestFirstStartFreezesConfiguration(t *testing.T) {
	app := newConfigurationTestApp(t)
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	require.ErrorIs(t, app.Start(ctx), context.Canceled)

	require.ErrorIs(t, app.EnableMetricsCollection(filepath.Join(t.TempDir(), "metrics.db")), ErrConfigurationFrozen)
	require.ErrorIs(t, app.EnableWebUI(false), ErrConfigurationFrozen)
}

func TestWebUIRequiresMetricsAtStartRegardlessOfCallOrder(t *testing.T) {
	app := newConfigurationTestApp(t)
	require.NoError(t, app.EnableWebUI(false))
	require.ErrorContains(t, app.Start(context.Background()), "requires metrics")
	require.ErrorIs(t, app.EnableMetricsCollection(filepath.Join(t.TempDir(), "metrics.db")), ErrConfigurationFrozen)
}

func TestAutomationRegistrationFreezesOnlyWhenMetricsEnabled(t *testing.T) {
	t.Run("metrics enabled", func(t *testing.T) {
		app := newConfigurationTestApp(t)
		require.NoError(t, app.EnableMetricsCollection(filepath.Join(t.TempDir(), "metrics.db")))
		ctx, cancel := context.WithCancel(context.Background())
		cancel()
		require.ErrorIs(t, app.Start(ctx), context.Canceled)

		require.PanicsWithValue(t, ErrConfigurationFrozen, func() {
			app.RegisterSchedules(NewDailySchedule().Call(func(*Service, State) {}).At("12:00").Build())
		})
		require.PanicsWithValue(t, ErrConfigurationFrozen, func() {
			app.RegisterIntervals(NewInterval().Call(func(*Service, State) {}).Every("1h").Build())
		})
		require.PanicsWithValue(t, ErrConfigurationFrozen, func() {
			app.RegisterEntityListeners(NewEntityListener().EntityIDs("light.test").Call(func(*Service, State, EntityData) {}).Build())
		})
		require.PanicsWithValue(t, ErrConfigurationFrozen, func() {
			app.RegisterEventListeners(NewEventListener().EventTypes("test").Call(func(*Service, State, EventData) {}).Build())
		})
	})

	t.Run("metrics disabled", func(t *testing.T) {
		app := newConfigurationTestApp(t)
		ctx, cancel := context.WithCancel(context.Background())
		cancel()
		require.ErrorIs(t, app.Start(ctx), context.Canceled)
		require.NotPanics(t, func() {
			app.RegisterSchedules(NewDailySchedule().Call(func(*Service, State) {}).At("12:00").Build())
		})
	})
}
