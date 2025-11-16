package config

type (
	LogingConfig struct {
		Observability Observability `koanf:"observability"`
	}

	Observability struct {
		Logging Logging `koanf:"logging"`
		Tracing Tracing `koanf:"tracing"`
		Metrics Metrics `koanf:"metrics"`
	}

	Logging struct {
		Level               string `koanf:"level"`
		AuthorizationHeader string `koanf:"authorizationheader"`
		Organization        string `koanf:"organization"`
		StreamName          string `koanf:"streamName"`
	}

	Tracing struct {
		Enabled             bool   `koanf:"enabled"`
		AuthorizationHeader string `koanf:"authorizationheader"`
		Organization        string `koanf:"organization"`
		StreamName          string `koanf:"streamName"`
	}

	Metrics struct {
		Enabled             bool   `koanf:"enabled"`
		AuthorizationHeader string `koanf:"authorizationheader"`
		Organization        string `koanf:"organization"`
		StreamName          string `koanf:"streamName"`
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

	Oauth struct {
		ClientID     string `koanf:"clientid"`
		ClientSecret string `koanf:"clientsecret"`
		ProviderURL  string `koanf:"providerurl"`
	}
)
