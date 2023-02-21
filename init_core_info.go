package horizon

import "github.com/SafeRE-IT/horizon/log"

func initCoreInfo(app *App) {
	err := app.UpdateCoreInfo()
	if err != nil {
		log.WithField("service", "core-info").WithError(err).Panic("Failed to init core info")
	}
}

func init() {
	appInit.Add("core-info", initCoreInfo, "core_connector", "app-context", "log")
}
