package logger

const (
	DebugLevel = "debug"
	ErrorLevel = "error"
	InfoLevel  = "info"
	WarnLevel  = "warn"
)

const (
	DefaultOutputPath = "stdout"
)

type Config struct {
	Debug       bool     `envconfig:"debug"`
	Level       string   `envconfig:"level" default:"info"`
	Output      []string `envconfig:"output"`
	TimeEncoder string   `envconfig:"time_encoder" default:"epoch"`
}

func (c Config) GetOutput() []string {
	if len(c.Output) == 0 {
		return []string{DefaultOutputPath}
	}
	return c.Output
}
