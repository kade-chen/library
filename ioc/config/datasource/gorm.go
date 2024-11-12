package datasource

import (
	"context"
	"fmt"

	"github.com/kade-chen/library/ioc"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func init() {
	ioc.Config().Registry(&dataSource{
		Host:        "35.220.136.110",
		Port:        3307,
		DB:          "sso",
		Username:    "root",
		Password:    "cys000522",
		Debug:       false,
		EnableTrace: true,
	})
}

type dataSource struct {
	// Provider    PROVIDER `json:"provider" yaml:"provider" toml:"provider" env:"DATASOURCE_PROVIDER"`
	Host        string `json:"host" yaml:"host" toml:"host" env:"DATASOURCE_HOST"`
	Port        int    `json:"port" yaml:"port" toml:"port" env:"DATASOURCE_PORT"`
	DB          string `json:"database" yaml:"database" toml:"database" env:"DATASOURCE_DB"`
	Username    string `json:"username" yaml:"username" toml:"username" env:"DATASOURCE_USERNAME"`
	Password    string `json:"password" yaml:"password" toml:"password" env:"DATASOURCE_PASSWORD"`
	Debug       bool   `json:"debug" yaml:"debug" toml:"debug" env:"DATASOURCE_DEBUG"`
	EnableTrace bool   `toml:"enable_trace" json:"enable_trace" yaml:"enable_trace"  env:"DATASOURCE_ENABLE_TRACE"`

	db *gorm.DB
	ioc.ObjectImpl
}

//toml，配置的名字datasource
func (m *dataSource) Name() string {
	return AppName
}

func (m *dataSource) Init() error {
	db, err := gorm.Open(mysql.Open(m.DSN()), &gorm.Config{})
	if err != nil {
		return err
	}

	if m.Debug {
		db = db.Debug()
	}
	m.db = db
	fmt.Println("mysql init succeessful")
	return nil
}

func (m *dataSource) DSN() string {
	//dsn := "user:pass@tcp(127.0.0.1:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local"
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		m.Username,
		m.Password,
		m.Host,
		m.Port,
		m.DB,
	)
}

// 关闭数据库连接
func (m *dataSource) Close(ctx context.Context) error {
	if m.db == nil {
		return nil
	}

	d, err := m.db.DB()
	if err != nil {
		return err
	}
	return d.Close()
}

func (i *dataSource) Version() string {
	return ""
}

func (i *dataSource) Priority() int {
	return 99
}

func (i *dataSource) AllowOverwrite() bool {
	return false
}
