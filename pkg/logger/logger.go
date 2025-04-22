package logger

import (
	"fmt"
	"github.com/mellgit/shorturl/internal/config"
	"github.com/natefinch/lumberjack"
	log "github.com/sirupsen/logrus"
	"io"
	"os"
	"path/filepath"
)

func SetUpLogger(cfg config.Config) error {

	level, err := log.ParseLevel(cfg.Logging.Level)
	if err != nil {
		return fmt.Errorf("invalid logging level %s", cfg.Logging.Level)
	}
	log.SetLevel(level)

	// Setting log formatter
	formatter := cfg.Logging.Formatter
	switch formatter {
	case "text":
		log.SetFormatter(&log.TextFormatter{
			DisableLevelTruncation: true,
			FullTimestamp:          true,
			DisableColors:          true,
		})
	case "json":
		log.SetFormatter(&log.JSONFormatter{})
	default:
		return fmt.Errorf("unsupported log formatter: %s", formatter)
	}

	// Setting log handler
	var handler io.Writer
	switch cfg.Logging.Handler {
	case "file":
		handler = &lumberjack.Logger{
			Filename: filepath.Join(cfg.Logging.Path, "entry.log"),
			MaxSize:  128,
		}
	case "console":
		handler = os.Stdout
	default:
		return fmt.Errorf("unsupported log handler type %s", cfg.Logging.Handler)
	}
	log.SetOutput(handler)

	return nil

}
