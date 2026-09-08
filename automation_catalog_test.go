package gomeassistant

import (
	"context"
	"path/filepath"
	"regexp"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestKeyIsAvailableOnEveryFinalBuilderStage(t *testing.T) {
	schedule := NewDailySchedule().Call(func(*Service, State) {}).At("12:00").Key(" first ").Key("daily").Build()
	interval := NewInterval().Call(func(*Service, State) {}).Every("1h").Key("interval").Build()
	entity := NewEntityListener().EntityIDs("light.test").Call(func(*Service, State, EntityData) {}).Key("entity").Build()
	event := NewEventListener().EventTypes("test").TriggerWhen(func(EventData) bool { return true }).Call(func(*Service, State, EventData) {}).Key("event").Build()

	require.Equal(t, "daily", schedule.key)
	require.Equal(t, "interval", interval.key)
	require.Equal(t, "entity", entity.key)
	require.Equal(t, "event", event.key)

	require.PanicsWithValue(t, ErrInvalidArgs, func() {
		NewDailySchedule().Call(func(*Service, State) {}).At("12:00").Key("  ")
	})
	require.PanicsWithValue(t, ErrInvalidArgs, func() {
		NewInterval().Call(func(*Service, State) {}).Every("1h").Key("")
	})
	require.PanicsWithValue(t, ErrInvalidArgs, func() {
		NewEntityListener().EntityIDs("light.test").Call(func(*Service, State, EntityData) {}).Key("\t")
	})
	require.PanicsWithValue(t, ErrInvalidArgs, func() {
		NewEventListener().EventTypes("test").Call(func(*Service, State, EventData) {}).Key("\n")
	})
}

func TestDisableMetricsIsAvailableOnEveryFinalBuilderStage(t *testing.T) {
	schedule := NewDailySchedule().Call(func(*Service, State) {}).At("12:00").DisableMetrics().Key("daily").Build()
	interval := NewInterval().Call(func(*Service, State) {}).Every("1h").Key("interval").DisableMetrics().Build()
	entity := NewEntityListener().EntityIDs("light.test").Call(func(*Service, State, EntityData) {}).DisableMetrics().Key("entity").Build()
	event := NewEventListener().EventTypes("test").TriggerWhen(func(EventData) bool { return true }).Call(func(*Service, State, EventData) {}).Key("event").DisableMetrics().Build()

	require.True(t, schedule.metricsDisabled)
	require.True(t, interval.metricsDisabled)
	require.True(t, entity.metricsDisabled)
	require.True(t, event.metricsDisabled)
}

func TestMetricsDisabledAutomationsAreExcludedFromCatalogAndKeyValidation(t *testing.T) {
	app := newConfigurationTestApp(t)
	app.RegisterSchedules(
		NewDailySchedule().Call(func(*Service, State) {}).At("12:00").Key("shared").Build(),
		NewDailySchedule().Call(func(*Service, State) {}).At("13:00").Key("shared").DisableMetrics().Build(),
		NewDailySchedule().Call(nil).At("14:00").DisableMetrics().Build(),
	)
	app.RegisterIntervals(NewInterval().Call(func(*Service, State) {}).Every("5s").Key("shared").DisableMetrics().Build())
	app.RegisterEntityListeners(NewEntityListener().EntityIDs("light.test").Call(func(*Service, State, EntityData) {}).Key("shared").DisableMetrics().Build())
	app.RegisterEventListeners(NewEventListener().EventTypes("test").Call(func(*Service, State, EventData) {}).Key("shared").DisableMetrics().Build())
	require.NoError(t, app.EnableMetricsCollection(filepath.Join(t.TempDir(), "metrics.db")))

	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	require.ErrorIs(t, app.Start(ctx), context.Canceled)
	require.Len(t, app.catalog, 1)
	require.Equal(t, "shared", app.catalog[0].Key)
	require.Equal(t, "daily", app.catalog[0].Type)
	require.Len(t, app.metricsUIAutomations(), 1)

	for _, action := range app.scheduledActions {
		switch automation := action.(type) {
		case *DailySchedule:
			if automation.metricsDisabled {
				require.Empty(t, automation.automationKey)
			}
		case *Interval:
			require.True(t, automation.metricsDisabled)
			require.Empty(t, automation.automationKey)
		}
	}
}

func TestCatalogDescribesEventTriggerPredicate(t *testing.T) {
	app := newConfigurationTestApp(t)
	app.RegisterEventListeners(NewEventListener().
		EventTypes("zha_event").
		TriggerWhen(func(EventData) bool { return true }).
		Call(func(*Service, State, EventData) {}).
		Key("filtered-event").
		Build())
	require.NoError(t, app.EnableMetricsCollection(filepath.Join(t.TempDir(), "metrics.db")))

	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	require.ErrorIs(t, app.Start(ctx), context.Canceled)
	require.Len(t, app.catalog, 1)
	require.Equal(t, true, app.catalog[0].Configuration["trigger_when"])
}

func TestCatalogUsesCallbackSourceFallbackAndDeduplicatesMultiEntityListener(t *testing.T) {
	app := newConfigurationTestApp(t)
	listener := NewEntityListener().EntityIDs("light.one", "light.two").Call(func(*Service, State, EntityData) {}).Build()
	app.RegisterEntityListeners(listener)
	require.NoError(t, app.EnableMetricsCollection(filepath.Join(t.TempDir(), "metrics.db")))

	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	require.ErrorIs(t, app.Start(ctx), context.Canceled)
	require.Len(t, app.catalog, 1)
	require.Regexp(t, regexp.MustCompile(`^automation_catalog_test\.go:\d+$`), app.catalog[0].Key)
	require.Equal(t, app.catalog[0].Key, app.catalog[0].Source)
	require.Equal(t, "entity", app.catalog[0].Type)
	require.ElementsMatch(t, []string{"light.one", "light.two"}, app.catalog[0].Configuration["entity_ids"])
}

func TestCatalogRejectsGlobalExplicitAndFallbackCollisions(t *testing.T) {
	t.Run("explicit across types", func(t *testing.T) {
		app := newConfigurationTestApp(t)
		app.RegisterSchedules(NewDailySchedule().Call(func(*Service, State) {}).At("12:00").Key("same").Build())
		app.RegisterEventListeners(NewEventListener().EventTypes("test").Call(func(*Service, State, EventData) {}).Key("same").Build())
		require.NoError(t, app.EnableMetricsCollection(filepath.Join(t.TempDir(), "metrics.db")))
		err := app.Start(context.Background())
		require.ErrorContains(t, err, `duplicate automation key "same"`)
		require.ErrorContains(t, err, ".Key() builder method")
	})

	t.Run("same callback fallback", func(t *testing.T) {
		callback := func(*Service, State, EntityData) {}
		app := newConfigurationTestApp(t)
		app.RegisterEntityListeners(
			NewEntityListener().EntityIDs("light.one").Call(callback).Build(),
			NewEntityListener().EntityIDs("light.two").Call(callback).Build(),
		)
		require.NoError(t, app.EnableMetricsCollection(filepath.Join(t.TempDir(), "metrics.db")))
		err := app.Start(context.Background())
		require.ErrorContains(t, err, "duplicate automation key")
		require.ErrorContains(t, err, "callback source locations are used as keys by default")
		require.ErrorContains(t, err, ".Key() builder method")
	})
}

func TestExplicitKeyAllowsUnavailableCallbackMetadata(t *testing.T) {
	app := newConfigurationTestApp(t)
	app.RegisterSchedules(NewDailySchedule().Call(nil).At("12:00").Key("explicit").Build())
	require.NoError(t, app.EnableMetricsCollection(filepath.Join(t.TempDir(), "metrics.db")))
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	require.ErrorIs(t, app.Start(ctx), context.Canceled)
	require.Equal(t, "explicit", app.catalog[0].Key)
}

func TestFallbackRequiresCallbackMetadata(t *testing.T) {
	app := newConfigurationTestApp(t)
	app.RegisterSchedules(NewDailySchedule().Call(nil).At("12:00").Build())
	require.NoError(t, app.EnableMetricsCollection(filepath.Join(t.TempDir(), "metrics.db")))
	require.ErrorContains(t, app.Start(context.Background()), "supply an explicit Key")
}

func TestSolarInitializationRefreshesFrozenCatalogNextDue(t *testing.T) {
	app := newConfigurationTestApp(t)
	app.RegisterSchedules(NewDailySchedule().Call(func(*Service, State) {}).Sunrise().Key("sunrise").Build())
	require.NoError(t, app.EnableMetricsCollection(filepath.Join(t.TempDir(), "metrics.db")))

	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	require.ErrorIs(t, app.Start(ctx), context.Canceled)
	require.True(t, app.catalog[0].NextDue.IsZero())

	app.state.latitude = 40.7128
	app.state.longitude = -74.0060
	app.initializeSolarSchedules()

	catalog := app.automationCatalogSnapshot()
	require.Len(t, catalog, 1)
	require.False(t, catalog[0].NextDue.IsZero())

	uiCatalog := app.metricsUIAutomations()
	require.Len(t, uiCatalog, 1)
	require.Equal(t, "sunrise", uiCatalog[0].Key)
	require.Equal(t, catalog[0].NextDue, uiCatalog[0].NextDue)
}
