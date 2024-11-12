package ioc

import "context"

const (
	DEFAULT_VERSION = "v1"
)

type ObjectImpl struct {
}

func (i *ObjectImpl) Init() error {
	return nil
}

func (i *ObjectImpl) Name() string {
	return ""
}

func (i *ObjectImpl) Hello1111() string {
	return "hello shhshsahbsakbsajkb"
}

func (i *ObjectImpl) Version() string {
	return ""
}

func (i *ObjectImpl) Priority() int {
	return 0
}

func (i *ObjectImpl) AllowOverwrite() bool {
	return true
}

func (i *ObjectImpl) Close(ctx context.Context) error {
	return nil
}

// go-restful
func (i *ObjectImpl) Meta() ObjectMeta {
	return DefaultObjectMeta()
}

func DefaultObjectMeta() ObjectMeta {
	return ObjectMeta{
		//		CustomPathPrefix: "/s", 必须要/号 http://127.0.0.1:8080/s
		CustomPathPrefix: "",
		// CustomPathPrefix: "/s",
		Extra: map[string]string{},
	}
}

type ObjectMeta struct {
	CustomPathPrefix string
	Extra            map[string]string
}
