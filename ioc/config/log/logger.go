package log

import (
	"fmt"
	"io"
	"strconv"
	"strings"
	"sync"

	// oteltrace "go.opentelemetry.io/otel/trace"
	"github.com/kade-chen/library/ioc"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/pkgerrors"
	"gopkg.in/natefinch/lumberjack.v2"
)

func init() {
	ioc.Config().Registry(&Config{
		CallerDeep: 3,
		Level:      zerolog.DebugLevel,
		TraceFiled: "trace_id",
		Console: Console{
			Enable:  true,
			NoColor: false,
		},
		File: File{
			Enable:     false,
			MaxSize:    100,
			MaxBackups: 6,
		},
		loggers: make(map[string]*zerolog.Logger),
	})
}

type Config struct {
	// 0 为打印日志全路径, 默认打印2层路径
	CallerDeep int `toml:"caller_deep" json:"caller_deep" yaml:"caller_deep"  env:"LOG_CALLER_DEEP"`
	// 日志的级别, 默认Debug
	Level zerolog.Level `toml:"level" json:"level" yaml:"level"  env:"LOG_LEVEL"`
	// 开启Trace时, 记录的TraceId名称, 默认trace_id
	TraceFiled string `toml:"trace_filed" json:"trace_filed" yaml:"trace_filed"  env:"LOG_TRACE_FILED"`

	// 控制台日志配置
	Console Console `toml:"console" json:"console" yaml:"console" envPrefix:"LOG_CONSOLE_"`
	// 日志文件配置
	File File `toml:"file" json:"file" yaml:"file" envPrefix:"LOG_FILE_"`

	ioc.ObjectImpl
	root    *zerolog.Logger
	lock    sync.Mutex
	loggers map[string]*zerolog.Logger
}

// 注册时版本初始化覆盖ioc.ObjectImpl
func (c *Config) Init() error {
	var writers []io.Writer
	if c.Console.Enable {
		writers = append(writers, c.Console.ConsoleWriter())
	}
	if c.File.Enable {
		writers = append(writers, c.File.FileWriter())
	}

	if len(writers) == 0 {
		return nil
	}
	//当使用 zerolog 记录错误日志时，它可以处理 error 类型的堆栈信息，将其格式化并记录下来。
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
	//2024-01-25T22:35:09+08:00
	// zerolog.TimeFieldFormat = "01/02 03:04:05PM '06 -0700"
	// logger := zerolog.New(os.Stdout).With().Timestamp().Logger()
	//MultiWriter函数创建一个多重写入器，将多个写入器（writers 切片中的 io.Writer）组合在一起
	// root := zerolog.New(io.MultiWriter(writers...)).With().Timestamp().Str("SUB_LOGGER_KEY", "ccc")
	root := zerolog.New(io.MultiWriter(writers...)).With().Timestamp()
	// 0 为打印日志全路径, 默认打印2层路径
	if c.CallerDeep > 0 {
		//这个方法会在日志记录中包含调用者信息。这样，root 就变成了一个包含了调用者信息的新的日志记录器
		root = root.Caller()
		//这行代码设置了 zerolog 全局变量 CallerMarshalFunc，这是一个函数，用于格式化调用者信息。
		//通过将配置对象 c 中的 CallerMarshalFunc 赋值给全局变量，可以自定义调用者信息的格式化方式。
		zerolog.CallerMarshalFunc = c.CallerMarshalFunc
	}
	//第一种方式
	// cc := root.Logger().Level(c.Level)
	// c.root = &cc
	//第二种
	c.SetRoot(root.Logger().Level(c.Level))
	return nil
}

// 这是一个函数，用于格式化调用者信息。
// 通过将配置对象 c 中的 CallerMarshalFunc 赋值给全局变量，可以自定义调用者信息的格式化方式
func (c *Config) CallerMarshalFunc(pc uintptr, file string, line int) string {
	if c.CallerDeep == 0 {
		return file
	}

	short := file
	count := 0
	for i := len(file) - 1; i > 0; i-- {
		if file[i] == '/' {
			short = file[i+1:]
			count++
		}
		if count >= c.CallerDeep {
			break
		}
	}
	file = short
	return file + ":" + strconv.Itoa(line)
}

func (c *Config) SetRoot(r zerolog.Logger) {
	c.root = &r
}

func (c *Config) Name() string {
	return AppName
}

// 这个可以在init里面做
func (c *Config) Logger(name string) *zerolog.Logger {
	c.lock.Lock()
	defer c.lock.Unlock()

	if _, ok := c.loggers[name]; !ok {
		//logger := zerolog.New(os.Stdout).With().Timestamp().Logger()
		l := c.root.With().Str(SUB_LOGGER_KEY, name).Logger()
		c.loggers[name] = &l
	}
	return c.loggers[name]
}

// 开启Trace时, 记录的TraceId名称
func (c *Config) TraceLogger(name string) *TraceLogger {
	return &TraceLogger{
		l:          c.Logger(name),
		traceFiled: c.TraceFiled, //"level":"debug",
	}
}

type TraceLogger struct {
	l          *zerolog.Logger
	traceFiled string
}

// 将上下文中的跟踪信息添加到日志记录中
// func (t *TraceLogger) Trace(ctx context.Context) *zerolog.Logger {
// 	traceId := oteltrace.SpanFromContext(ctx).SpanContext().TraceID().String()
// 	if traceId == "" {
// 		return t.l
// 	}
// 	l := t.l.With().Str(t.traceFiled, traceId).Logger()
// 	return &l
// }

type File struct {
	// 是否开启文件记录
	Enable bool `toml:"enable" json:"enable" yaml:"enable"  env:"ENABLE"`
	// 文件的路径
	FilePath string `toml:"file_path" json:"file_path" yaml:"file_path"  env:"PATH"`
	// 单位M, 默认100M
	MaxSize int `toml:"max_size" json:"max_size" yaml:"max_size"  env:"MAX_SIZE"`
	// 默认保存 6个文件
	MaxBackups int `toml:"max_backups" json:"max_backups" yaml:"max_backups"  env:"MAX_BACKUPS"`
	// 保存多久
	MaxAge int `toml:"max_age" json:"max_age" yaml:"max_age"  env:"MAX_AGE"`
	// 是否压缩
	Compress bool `toml:"compress" json:"compress" yaml:"compress"  env:"COMPRESS"`
}

func (f *File) FileWriter() io.Writer {
	return &lumberjack.Logger{
		Filename:   f.FilePath,
		MaxSize:    f.MaxSize,
		MaxAge:     f.MaxAge,
		MaxBackups: f.MaxBackups,
		Compress:   f.Compress,
	}
}

type Console struct {
	Enable  bool `toml:"enable" json:"enable" yaml:"enable"  env:"ENABLE"`
	NoColor bool `toml:"no_color" json:"no_color" yaml:"no_color"  env:"NO_COLOR"`
}

// 设置日志打印在console上
// 在这里，i interface{} 是一个空接口类型，表示可以接受任意类型的值。
// 在 FormatLevel、FormatMessage、FormatFieldName、FormatFieldValue 四个函数中
// i 表示相应的字段值，但具体的类型会根据日志库的实现和日志记录时传递的值而变化。这种设计允许日志库接受各种类型的字段，并使用用户提供的函数来格式化这些字段的输出
func (c *Console) ConsoleWriter() io.Writer {
	output := zerolog.NewConsoleWriter(func(w *zerolog.ConsoleWriter) {
		w.NoColor = c.NoColor
		//设置时间格式为 RFC3339
		// w.TimeFormat = time.RFC3339
		w.TimeFormat = " 2006-01-02 15:04:05 PM '06 -0700"
	})
	// output.TimeFormat = "2006 01/02 15:04:05 PM '06 -0700"
	// 设置日志级别的格式。在这里，将日志级别转为大写，并左对齐，占用 6 个字符的宽度
	output.FormatLevel = func(i interface{}) string {
		return strings.ToUpper(fmt.Sprintf("%-6s", i))
	}
	//设置日志消息的格式
	output.FormatMessage = func(i interface{}) string {
		return fmt.Sprintf("%s", i)
	}
	//设置字段名的格式。在这里，字段名后面加上冒号 :
	output.FormatFieldName = func(i interface{}) string {
		return fmt.Sprintf("%s:", i)
	}
	//设置字段值的格式。在这里，将字段值转为大写
	output.FormatFieldValue = func(i interface{}) string {
		return strings.ToUpper(fmt.Sprintf("%s", i))
	}
	return output
}

func (i *Config) Priority() int {
	return 10000
}
