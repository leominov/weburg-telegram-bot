package metrics

import (
	"gotel/watcher"
	"net/http"

	"github.com/Sirupsen/logrus"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/common/version"
)

const (
	namespace = "weburg"
)

var (
	MessagesTotalCounter = prometheus.NewCounter(prometheus.CounterOpts{
		Name: prometheus.BuildFQName(namespace, "", "messages_total_count"),
		Help: "How many messages are sent",
	})

	MessagesFailCounter = prometheus.NewCounter(prometheus.CounterOpts{
		Name: prometheus.BuildFQName(namespace, "", "messages_fail_count"),
		Help: "How many messages are fail",
	})

	PullsTotalCounter = prometheus.NewCounter(prometheus.CounterOpts{
		Name: prometheus.BuildFQName(namespace, "", "pulls_total_count"),
		Help: "How many pulls.",
	})

	PullsFailCounter = prometheus.NewCounter(prometheus.CounterOpts{
		Name: prometheus.BuildFQName(namespace, "", "pulls_fail_count"),
		Help: "How many failed pulls.",
	})

	MessagesTotalCounters = map[string]prometheus.Counter{}
	MessagesFailCounters  = map[string]prometheus.Counter{}
	PullsTotalCounters    = map[string]prometheus.Counter{}
	PullsFailCounters     = map[string]prometheus.Counter{}
)

func InitMetrics() {
	prometheus.MustRegister(version.NewCollector(namespace))

	prometheus.MustRegister(MessagesTotalCounter)
	prometheus.MustRegister(MessagesFailCounter)
	prometheus.MustRegister(PullsTotalCounter)
	prometheus.MustRegister(PullsFailCounter)

	for _, c := range watcher.RssAgentsCollection {
		MessagesTotalCounters[c.Type] = prometheus.NewCounter(prometheus.CounterOpts{
			Name:        prometheus.BuildFQName(namespace, "", "messages_feed_total_count"),
			Help:        "How many messages are sent",
			ConstLabels: prometheus.Labels{"feed": c.Type},
		})
		prometheus.MustRegister(MessagesTotalCounters[c.Type])
		MessagesFailCounters[c.Type] = prometheus.NewCounter(prometheus.CounterOpts{
			Name:        prometheus.BuildFQName(namespace, "", "messages_feed_fail_count"),
			Help:        "How many messages are fail",
			ConstLabels: prometheus.Labels{"feed": c.Type},
		})
		prometheus.MustRegister(MessagesFailCounters[c.Type])
		PullsTotalCounters[c.Type] = prometheus.NewCounter(prometheus.CounterOpts{
			Name:        prometheus.BuildFQName(namespace, "", "pulls_feed_total_count"),
			Help:        "How many pulls.",
			ConstLabels: prometheus.Labels{"feed": c.Type},
		})
		prometheus.MustRegister(PullsTotalCounters[c.Type])
		PullsFailCounters[c.Type] = prometheus.NewCounter(prometheus.CounterOpts{
			Name:        prometheus.BuildFQName(namespace, "", "pulls_feed_fail_count"),
			Help:        "How many failed pulls.",
			ConstLabels: prometheus.Labels{"feed": c.Type},
		})
		prometheus.MustRegister(PullsFailCounters[c.Type])
	}
}

func ServeMetrics(listenAddr, metricsPath string) {
	http.Handle(metricsPath, promhttp.Handler())
	logrus.Infof("Listen address: %s", listenAddr)
	logrus.Fatal(http.ListenAndServe(listenAddr, nil))
}
