package sauron_test

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/Bowbaq/sauron"
	"github.com/Bowbaq/sauron/model"
	"github.com/Bowbaq/sauron/store"
	"github.com/Bowbaq/sauron/watcher"
)

var (
	watchSauron = model.WatchOptions{
		Repository: model.Repository{Owner: "Bowbaq", Name: "sauron"},
	}
	watchSauronFeature = model.WatchOptions{
		Repository: model.Repository{Owner: "Bowbaq", Name: "sauron"},
		Branch:     "feature",
	}
	watchSauronReadme = model.WatchOptions{
		Repository: model.Repository{Owner: "Bowbaq", Name: "sauron"},
		Path:       "README.md",
	}
)

func TestRequiresOwner(t *testing.T) {
	s, _, _ := setup()

	require.Error(t, s.Watch(model.WatchOptions{}), "Expected 'Repository owner required' error")
}

func TestRequiresRepository(t *testing.T) {
	s, _, _ := setup()

	require.Error(t, s.Watch(model.WatchOptions{
		Repository: model.Repository{
			Owner: "owner",
		},
	}), "Expected 'Repository name required' error")
}

func TestDoesntNotifyOnFirstRun(t *testing.T) {
	s, _, notifications := setup()

	require.Empty(t, notifications, "Shouldn't have any notifications before first run")

	require.NoError(t, s.Watch(watchSauron))
	require.Empty(t, notifications, "Shouldn't have any notifications after first run")
}

func TestDoesntNotifyWhenUnchanged(t *testing.T) {
	s, _, notifications := setup()

	require.Empty(t, notifications, "Shouldn't have any notifications before first run")

	require.NoError(t, s.Watch(watchSauron))
	require.Empty(t, notifications, "Shouldn't have any notifications after first run")

	require.NoError(t, s.Watch(watchSauron))
	require.Empty(t, notifications, "Shouldn't have any notifications after second run (no changes)")
}

func TestNotifiesOnUpdate(t *testing.T) {
	s, _, notifications := setup()

	s.SetWatcher(withChanges(map[store.WatchKey][]model.Update{
		store.Key(watchSauron): {{Timestamp: time.Now()}, {Timestamp: time.Now().Add(1 * time.Hour)}},
	}))

	require.Empty(t, notifications, "Shouldn't have any notifications before first run")

	require.NoError(t, s.Watch(watchSauron))
	require.Empty(t, notifications, "Shouldn't have any notifications after first run")

	require.NoError(t, s.Watch(watchSauron))
	require.Len(t, notifications, 1, "Should have one notification after second run (changed)")
}

func TestFiltersOnBranch(t *testing.T) {
	s, _, notifications := setup()

	now := time.Now()
	s.SetWatcher(withChanges(map[store.WatchKey][]model.Update{
		store.Key(watchSauron):        {{Timestamp: now}, {Timestamp: now.Add(1 * time.Hour)}}, // Master changed
		store.Key(watchSauronFeature): {{Timestamp: now}, {Timestamp: now}},                    // Feature branch didn't
	}))

	require.Empty(t, notifications, "Shouldn't have any notifications before first run")

	sauronKey := store.Key(watchSauron)
	require.NoError(t, s.Watch(watchSauron))
	require.Empty(t, notifications[sauronKey], "Shouldn't have any notifications after first run")

	require.NoError(t, s.Watch(watchSauron))
	require.Len(t, notifications[sauronKey], 1, "Should have one notification after second run (changed)")

	sauronFeatureKey := store.Key(watchSauronFeature)
	require.NoError(t, s.Watch(watchSauronFeature))
	require.Empty(t, notifications[sauronFeatureKey], "Shouldn't have any notifications after first run")

	require.NoError(t, s.Watch(watchSauronFeature))
	require.Empty(t, notifications[sauronFeatureKey], "Shouldn't have any notifications after second run (no changes)")
}

func TestFiltersOnPath(t *testing.T) {
	s, _, notifications := setup()

	now := time.Now()
	s.SetWatcher(withChanges(map[store.WatchKey][]model.Update{
		store.Key(watchSauron):       {{Timestamp: now}, {Timestamp: now.Add(1 * time.Hour)}}, // Master changed
		store.Key(watchSauronReadme): {{Timestamp: now}, {Timestamp: now}},                    // README.md didn't
	}))

	require.Empty(t, notifications, "Shouldn't have any notifications before first run")

	sauronKey := store.Key(watchSauron)
	require.NoError(t, s.Watch(watchSauron))
	require.Empty(t, notifications[sauronKey], "Shouldn't have any notifications after first run")

	require.NoError(t, s.Watch(watchSauron))
	require.Len(t, notifications[sauronKey], 1, "Should have one notification after second run (changed)")

	sauronReadmeKey := store.Key(watchSauronReadme)
	require.NoError(t, s.Watch(watchSauronReadme))
	require.Empty(t, notifications[sauronReadmeKey], "Shouldn't have any notifications after first run")

	require.NoError(t, s.Watch(watchSauronReadme))
	require.Empty(t, notifications[sauronReadmeKey], "Shouldn't have any notifications after second run (no changes)")
}

type testNotification struct {
	previous, current model.Update
}

type testNotifier map[store.WatchKey][]testNotification

func (tn testNotifier) Notify(opts model.WatchOptions, lastUpdate, update model.Update) error {
	tn[store.Key(opts)] = append(tn[store.Key(opts)], testNotification{lastUpdate, update})

	return nil
}

func setup() (*sauron.Sauron, sauron.Options, testNotifier) {
	os.Remove(".sauron")

	options := sauron.Options{}
	notifier := testNotifier{}

	s := sauron.New(options)
	s.SetNotifier(notifier)

	return s, options, notifier
}

type changingWatcher struct {
	changes map[store.WatchKey][]model.Update
	hits    map[store.WatchKey]int
}

func (w *changingWatcher) CheckForUpdate(opts model.WatchOptions, _ time.Time) (model.Update, error) {
	key := store.Key(opts)
	u := w.changes[key][w.hits[key]]
	w.hits[key]++

	return u, nil
}

func withChanges(changes map[store.WatchKey][]model.Update) watcher.Watcher {
	return &changingWatcher{
		changes: changes,

		hits: make(map[store.WatchKey]int),
	}
}
