package statsd

import (
	"runtime"
	"strings"
	"time"

	"github.com/atlassian/gostatsd"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"golang.org/x/time/rate"
)

// DefaultBackends is the list of default backends' names.
var DefaultBackends = []string{"graphite"}

// DefaultMaxReaders is the default number of socket reading goroutines.
var DefaultMaxReaders = runtime.NumCPU()

// DefaultMaxWorkers is the default number of goroutines that aggregate metrics.
var DefaultMaxWorkers = runtime.NumCPU()

// DefaultPercentThreshold is the default list of applied percentiles.
var DefaultPercentThreshold = []float64{90}

// DefaultTags is the default list of additional tags.
var DefaultTags = gostatsd.Tags{}

const (
	// DefaultMaxCloudRequests is the maximum number of cloud provider requests per second.
	DefaultMaxCloudRequests = 10
	// DefaultBurstCloudRequests is the burst number of cloud provider requests per second.
	DefaultBurstCloudRequests = DefaultMaxCloudRequests + 5
	// DefaultExpiryInterval is the default expiry interval for metrics.
	DefaultExpiryInterval = 5 * time.Minute
	// DefaultFlushInterval is the default metrics flush interval.
	DefaultFlushInterval = 1 * time.Second
	// DefaultIgnoreHost is the default value for whether the source should be used as the host
	DefaultIgnoreHost = false
	// DefaultMetricsAddr is the default address on which to listen for metrics.
	DefaultMetricsAddr = ":8125"
	// DefaultMaxQueueSize is the default maximum number of buffered metrics per worker.
	DefaultMaxQueueSize = 10000 // arbitrary
	// DefaultMaxConcurrentEvents is the default maximum number of events sent concurrently.
	DefaultMaxConcurrentEvents = 1024 // arbitrary
	// DefaultCacheRefreshPeriod is the default cache refresh period.
	DefaultCacheRefreshPeriod = 1 * time.Minute
	// DefaultCacheEvictAfterIdlePeriod is the default idle cache eviction period.
	DefaultCacheEvictAfterIdlePeriod = 10 * time.Minute
	// DefaultCacheTTL is the default cache TTL for successful lookups.
	DefaultCacheTTL = 30 * time.Minute
	// DefaultCacheNegativeTTL is the default cache TTL for failed lookups (errors or when instance was not found).
	DefaultCacheNegativeTTL = 1 * time.Minute
)

const (
	// ParamBackends is the name of parameter with backends.
	ParamBackends = "backends"
	// ParamCloudProvider is the name of parameter with the name of cloud provider.
	ParamCloudProvider = "cloud-provider"
	// ParamMaxCloudRequests is the name of parameter with maximum number of cloud provider requests per second.
	ParamMaxCloudRequests = "max-cloud-requests"
	// ParamBurstCloudRequests is the name of parameter with burst number of cloud provider requests per second.
	ParamBurstCloudRequests = "burst-cloud-requests"
	// ParamDefaultTags is the name of parameter with the list of additional tags.
	ParamDefaultTags = "default-tags"
	// ParamExpiryInterval is the name of parameter with expiry interval for metrics.
	ParamExpiryInterval = "expiry-interval"
	// ParamFlushInterval is the name of parameter with metrics flush interval.
	ParamFlushInterval = "flush-interval"
	// ParamIgnoreHost is the name of parameter indicating if the source should be used as the host
	ParamIgnoreHost = "ignore-host"
	// ParamMaxReaders is the name of parameter with number of socket readers.
	ParamMaxReaders = "max-readers"
	// ParamMaxWorkers is the name of parameter with number of goroutines that aggregate metrics.
	ParamMaxWorkers = "max-workers"
	// ParamMaxQueueSize is the name of parameter with maximum number of buffered metrics per worker.
	ParamMaxQueueSize = "max-queue-size"
	// ParamMaxConcurrentEvents is the name of parameter with maximum number of events sent concurrently.
	ParamMaxConcurrentEvents = "max-concurrent-events"
	// ParamCacheRefreshPeriod is the name of parameter with cache refresh period.
	ParamCacheRefreshPeriod = "cloud-cache-refresh-period"
	// ParamCacheEvictAfterIdlePeriod is the name of parameter with idle cache eviction period.
	ParamCacheEvictAfterIdlePeriod = "cloud-cache-evict-after-idle-period"
	// ParamCacheTTL is the name of parameter with cache TTL for successful lookups.
	ParamCacheTTL = "cloud-cache-ttl"
	// ParamCacheNegativeTTL is the name of parameter with cache TTL for failed lookups (errors or when instance was not found).
	ParamCacheNegativeTTL = "cloud-cache-negative-ttl"
	// ParamMetricsAddr is the name of parameter with address on which to listen for metrics.
	ParamMetricsAddr = "metrics-addr"
	// ParamNamespace is the name of parameter with namespace for all metrics.
	ParamNamespace = "namespace"
	// ParamPercentThreshold is the name of parameter with list of applied percentiles.
	ParamPercentThreshold = "percent-threshold"
)

// NewServer will create a new Server with the default configuration.
func NewServer() *Server {
	return &Server{
		Limiter:             rate.NewLimiter(DefaultMaxCloudRequests, DefaultBurstCloudRequests),
		DefaultTags:         DefaultTags,
		ExpiryInterval:      DefaultExpiryInterval,
		FlushInterval:       DefaultFlushInterval,
		IgnoreHost:          DefaultIgnoreHost,
		MaxReaders:          DefaultMaxReaders,
		MaxWorkers:          DefaultMaxWorkers,
		MaxQueueSize:        DefaultMaxQueueSize,
		MaxConcurrentEvents: DefaultMaxConcurrentEvents,
		MetricsAddr:         DefaultMetricsAddr,
		PercentThreshold:    DefaultPercentThreshold,
		CacheOptions: CacheOptions{
			CacheRefreshPeriod:        DefaultCacheRefreshPeriod,
			CacheEvictAfterIdlePeriod: DefaultCacheEvictAfterIdlePeriod,
			CacheTTL:                  DefaultCacheTTL,
			CacheNegativeTTL:          DefaultCacheNegativeTTL,
		},
		Viper: viper.New(),
	}
}

// AddFlags adds flags to the specified FlagSet.
func AddFlags(fs *pflag.FlagSet) {
	fs.String(ParamCloudProvider, "", "If set, use the cloud provider to retrieve metadata about the sender")
	fs.Duration(ParamExpiryInterval, DefaultExpiryInterval, "After how long do we expire metrics (0 to disable)")
	fs.Duration(ParamFlushInterval, DefaultFlushInterval, "How often to flush metrics to the backends")
	fs.Bool(ParamIgnoreHost, DefaultIgnoreHost, "Ignore the source for populating the hostname field of metrics")
	fs.Int(ParamMaxReaders, DefaultMaxReaders, "Maximum number of socket readers")
	fs.Int(ParamMaxWorkers, DefaultMaxWorkers, "Maximum number of workers to process metrics")
	fs.Int(ParamMaxQueueSize, DefaultMaxQueueSize, "Maximum number of buffered metrics per worker")
	fs.Int(ParamMaxConcurrentEvents, DefaultMaxConcurrentEvents, "Maximum number of events sent concurrently")
	fs.Duration(ParamCacheRefreshPeriod, DefaultCacheRefreshPeriod, "Cloud cache refresh period")
	fs.Duration(ParamCacheEvictAfterIdlePeriod, DefaultCacheEvictAfterIdlePeriod, "Idle cloud cache eviction period")
	fs.Duration(ParamCacheTTL, DefaultCacheTTL, "Cloud cache TTL for successful lookups")
	fs.Duration(ParamCacheNegativeTTL, DefaultCacheNegativeTTL, "Cloud cache TTL for failed lookups")
	fs.String(ParamMetricsAddr, DefaultMetricsAddr, "Address on which to listen for metrics")
	fs.String(ParamNamespace, "", "Namespace all metrics")
	//TODO Remove workaround when https://github.com/spf13/viper/issues/112 is fixed
	// https://github.com/spf13/viper/issues/200
	fs.String(ParamBackends, strings.Join(DefaultBackends, ","), "Comma-separated list of backends")
	fs.Int(ParamMaxCloudRequests, DefaultMaxCloudRequests, "Maximum number of cloud provider requests per second")
	fs.Int(ParamBurstCloudRequests, DefaultBurstCloudRequests, "Burst number of cloud provider requests per second")
	fs.String(ParamDefaultTags, strings.Join(DefaultTags, ","), "Comma-separated list of tags to add to all metrics")
	fs.String(ParamPercentThreshold, strings.Join(toStringSlice(DefaultPercentThreshold), ","), "Comma-separated list of percentiles")
}
