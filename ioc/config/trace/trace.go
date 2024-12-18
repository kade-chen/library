package trace

// import (
// 	"github.com/kade-chen/library/ioc"
// 	"go.opentelemetry.io/otel"
// 	"go.opentelemetry.io/otel/exporters/jaeger"
// 	"go.opentelemetry.io/otel/propagation"
// 	"go.opentelemetry.io/otel/sdk/resource"
// 	sdktrace "go.opentelemetry.io/otel/sdk/trace"
// )

// func init() {
// 	ioc.Config().Registry(&Trace{
// 		Provider: TRACE_PROVIDER_JAEGER,
// 		Endpoint: "http://localhost:14268/api/traces",
// 		Enable:   false,
// 	})
// }

// type Trace struct {
// 	ioc.ObjectImpl

// 	Enable   bool           `json:"enable" yaml:"enable" toml:"enable" env:"TRACE_ENABLE"`
// 	Provider TRACE_PROVIDER `toml:"provider" json:"provider" yaml:"provider" env:"TRACE_PROVIDER"`
// 	Endpoint string         `toml:"endpoint" json:"endpoint" yaml:"endpoint" env:"TRACE_ENDPOINT"`
// }

// // 优先初始化, 以供后面的组件使用
// func (i *Trace) Priority() int {
// 	return 9998
// }

// func (i *Trace) Name() string {
// 	return AppName
// }

// func (t *Trace) Init() error {
// 	exporter, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(t.Endpoint)))
// 	if err != nil {
// 		return err
// 	}

// 	tp := sdktrace.NewTracerProvider(
// 		sdktrace.WithSampler(sdktrace.AlwaysSample()),
// 		sdktrace.WithBatcher(exporter),
// 		sdktrace.WithResource(resource.Default()),
// 	)
// 	otel.SetTracerProvider(tp)
// 	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.Baggage{}, propagation.TraceContext{}))
// 	return nil
// }
