package config

type (
	LogingConfig struct {
		Observability Observability `koanf:"observability"`
	}

	Observability struct {
		Logging Logging `koanf:"logging"`
	}

	Logging struct {
		Level    string   `koanf:"level" `
		Logstash Logstash `koanf:"logstash"`
	}

	Logstash struct {
		Enabled bool   `koanf:"enabled"`
		Address string `koanf:"address"`
	}

	Server struct {
		HTTP      HTTPServer `koanf:"http"`
		DomainURL string     `koanf:"domainurl"`
	}

	HTTPServer struct {
		Port       int    `koanf:"port"`
		Host       string `koanf:"host"`
		Production bool   `koanf:"production"`
		BasePath   string `koanf:"basepath"`
	}

	Security struct {
		Secret string `koanf:"secret"`
	}

	RabbitMQ struct {
		Host          string `koanf:"host"`
		EventExchange string `koanf:"exchange"`
		EventQueue    string `koanf:"queue"`
	}

	MongoDB struct {
		Host   string `koanf:"host"`
		DBName string `koanf:"database"`
	}
)
