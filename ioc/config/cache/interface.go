package cache

import (
	"fmt"

	"github.com/dgraph-io/ristretto"
	"github.com/kade-chen/library/exception"
	"github.com/kade-chen/library/ioc"
	"github.com/kade-chen/library/ioc/config/log"
	"github.com/rs/zerolog"
)

const (
	AppName = "cache"
)

func init() {
	ioc.Config().Registry(&Cache{
		Enabled: true,
		//访问频率计数器数量（不是缓存条目数）
		NumCounters: 1e6,
		// 缓存能用的“总成本”上限
		MaxCost: 512 << 20, // 256MB
		//Get 请求的异步缓冲队列长度
		BufferItems: 64,
		// 是否记录命中率、reject、evict 等指标
		Metrics: true,
		//TTL 扫描频率（不是精确到期时间）
		TtlTickerDurationInSec: 60,
		//     OnReject: func(item *ristretto.Item) {
		//         log.Warn().Msg("cache reject")
		//     },
	})
}

type Cache struct {
	Enabled                bool  `toml:"enabled" json:"enabled" yaml:"enabled"  env:"ENABLED"`
	NumCounters            int64 `toml:"num_counters" json:"num_counters" yaml:"num_counters"  env:"NUM_COUNTERS"`
	MaxCost                int64 `toml:"max_cost" json:"max_cost" yaml:"max_cost"  env:"MAX_COST"`
	BufferItems            int64 `toml:"buffer_items" json:"buffer_items" yaml:"buffer_items"  env:"BUFFER_ITEMS"`
	Metrics                bool  `toml:"metrics" json:"metrics" yaml:"metrics"  env:"METRICS"`
	TtlTickerDurationInSec int64 `toml:"ttl_ticker_duration_in_sec" json:"ttl_ticker_duration_in_sec" yaml:"ttl_ticker_duration_in_sec"  env:"TTL_TICKER_DURATION_IN_SEC"`
	Ristretto              *ristretto.Cache
	log                    *zerolog.Logger
	ioc.ObjectImpl
}

func (m *Cache) Name() string {
	return AppName
}

func (m *Cache) Init() error {
	m.log = log.Sub("cache")
	ristretto, err := ristretto.NewCache(&ristretto.Config{
		NumCounters:            m.NumCounters,
		MaxCost:                m.MaxCost,
		BufferItems:            m.BufferItems,
		Metrics:                m.Metrics,
		TtlTickerDurationInSec: m.TtlTickerDurationInSec,
		//“我宁愿不存，也不会把自己撑爆”
		OnReject: func(item *ristretto.Item) {
			m.log.Warn().Str("key", fmt.Sprint(item.Key)).
				Msg("CACHE REJECT")
		},
		// 	已经 成功进 cache
		// •	后来因为：
		// •	空间不够
		// •	更热的 key 进来
		// •	被策略踢掉
		OnEvict: func(item *ristretto.Item) {
			m.log.Warn().Str("key", fmt.Sprint(item.Key)).
				Int64("cost", item.Cost).
				Msg("CACHE EVICTED")
		},
		// 	任何离开 cache 的路径都会触发
		// •	eviction
		// •	reject
		// •	过期
		// •	Del()
		// •	Clear()
		OnExit: func(val interface{}) {
			m.log.Warn().Msg("CACHE VALUE EXIT")
		},
	})
	if err != nil {
		return exception.NewIocRegisterFailed("Cache Ioc Register Failed, ERROR: %v", err)
	}
	m.Ristretto = ristretto
	m.log.Debug().Msgf("%v Ioc Register Is Successsful", m.Name())
	return nil
}

func (i *Cache) Priority() int {
	return 0
}
