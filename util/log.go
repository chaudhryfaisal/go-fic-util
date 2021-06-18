package util

import (
	"github.com/chaudhryfaisal/elogrus"
	"github.com/olivere/elastic"
	"github.com/sirupsen/logrus"
	"os"
)

var Log = logrus.New()

func SetupLog() {
	Log.SetOutput(os.Stdout)
	// Create elastic client
	LogEsToken := PropS("LOG_ES_TOKEN", "")
	if len(LogEsToken) > 0 {
		LogEsUrl := PropS("LOG_ES_URL", "https://cloud.humio.com:443/api/v1/ingest/elastic-bulk")
		LogEsUser := PropS("LOG_ES_USER", "LOG_ES_USER")
		LogEsHost := PropS("LOG_ES_HOST", "LOG_ES_HOST")
		client, err := elastic.NewClient(
			elastic.SetURL(LogEsUrl),
			elastic.SetBasicAuth(LogEsUser, LogEsToken),
			elastic.SetHealthcheck(false),
			elastic.SetSniff(false),
			elastic.SetGzip(true),
			//elastic.SetErrorLog(logrus.New()),
			//elastic.SetTraceLog(logrus.New()),
		)
		if err != nil {
			Log.WithError(err).Fatal("Failed to construct elasticsearch client")
		}
		// Create logger with 15 seconds flush interval
		Log.Info("NewBulkProcessorElasticHook")
		hook, err := elogrus.NewBulkProcessorElasticHook(client, LogEsHost, logrus.DebugLevel, "none")
		if err != nil {
			Log.WithError(err).Error("Failed to create elasticsearch hook for logger")
		}
		Log.Hooks.Add(hook)
	} else {
		Log.Warnf("LOG_ES_SETUP Skipped!, set LOG_ES_TOKEN to enable")
	}
}
