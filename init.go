package horizon

import (
	"github.com/SafeRE-IT/horizon/log"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

// InitFn is a function that contributes to the initialization of an App struct
type InitFn func(*App)

type initializer struct {
	Name string
	Fn   InitFn
	Deps []string
}

type initializerSet []initializer

var appInit initializerSet

// Add adds a new initializer into the chain
func (is *initializerSet) Add(name string, fn InitFn, deps ...string) {
	*is = append(*is, initializer{
		Name: name,
		Fn:   fn,
		Deps: deps,
	})
}

// Run initializes the provided application, but running every Initializer
func (is *initializerSet) Run(app *App) {
	err := is.checkDuplicates()
	if err != nil {
		log.WithField("err", err).Fatal("failed to init initializer")
	}
	init := *is
	alreadyRun := make(map[string]bool)

	for {
		ranInitializer := false
		for _, i := range init {
			// if we've already been run, skip
			if ok := alreadyRun[i.Name]; ok {
				continue
			}

			// if any of our dependencies haven't been run, skip
			isReadyToRun := true
			for _, d := range i.Deps {
				if ok := alreadyRun[d]; !ok {
					alreadyRun[d] = false
					isReadyToRun = false
					break
				}
			}

			if !isReadyToRun {
				alreadyRun[i.Name] = false
				continue
			}

			log.WithField("init_name", i.Name).Debug("running initializer")
			i.Fn(app)
			alreadyRun[i.Name] = true
			ranInitializer = true
		}
		// If, after a full loop through the initializers we ran nothing
		// we are done
		if !ranInitializer {
			break
		}
	}

	// if we didn't get to run all initializers, we have a cycle
	if len(alreadyRun) != len(init) {
		var failedToRun []string
		for name, isStarted := range alreadyRun {
			if !isStarted {
				failedToRun = append(failedToRun, name)
			}
		}
		log.WithField("failedToRun", failedToRun).Panic("initializer cycle detected")
	}
}

func (is *initializerSet) checkDuplicates() error {
	init := *is
	unique := map[string]struct{}{}
	for _, runner := range init {
		if _, exists := unique[runner.Name]; exists {
			return errors.From(errors.New("duplicated initializer"), logan.F{"name": runner.Name})
		}
		unique[runner.Name] = struct{}{}
	}

	return nil
}
