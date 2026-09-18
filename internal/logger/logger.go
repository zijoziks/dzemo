package logger

//TODO I WANT TO MAKE LOGGING LOOK MORE NICER, NOT JUST RAW OUTPUT

import (
	"log/slog"
	"os"
)

var Logger = slog.New(slog.NewTextHandler(os.Stdout, nil))
