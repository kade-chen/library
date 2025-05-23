package apidoc

const (
	AppName    = "apidoc-restful"
	AppNamegin = "apidoc-gin"
)

type ApiDoc struct {
	// Swagger API Doc URL路径
	Path string `json:"path" yaml:"path" toml:"path" env:"HTTP_API_DOC_PATH"`
}
