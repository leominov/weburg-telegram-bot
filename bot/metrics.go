package bot

import (
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

func (b *Bot) InitMetrics() {
	prometheus.MustRegister(version.NewCollector(namespace))

	prometheus.MustRegister(MessagesTotalCounter)
	prometheus.MustRegister(MessagesFailCounter)
	prometheus.MustRegister(PullsTotalCounter)
	prometheus.MustRegister(PullsFailCounter)

	for _, c := range b.Config.Agents {
		MessagesTotalCounters[c.Name] = prometheus.NewCounter(prometheus.CounterOpts{
			Name:        prometheus.BuildFQName(namespace, "", "messages_feed_total_count"),
			Help:        "How many messages are sent",
			ConstLabels: prometheus.Labels{"feed": c.Name},
		})
		prometheus.MustRegister(MessagesTotalCounters[c.Name])
		MessagesFailCounters[c.Name] = prometheus.NewCounter(prometheus.CounterOpts{
			Name:        prometheus.BuildFQName(namespace, "", "messages_feed_fail_count"),
			Help:        "How many messages are fail",
			ConstLabels: prometheus.Labels{"feed": c.Name},
		})
		prometheus.MustRegister(MessagesFailCounters[c.Name])
		PullsTotalCounters[c.Name] = prometheus.NewCounter(prometheus.CounterOpts{
			Name:        prometheus.BuildFQName(namespace, "", "pulls_feed_total_count"),
			Help:        "How many pulls.",
			ConstLabels: prometheus.Labels{"feed": c.Name},
		})
		prometheus.MustRegister(PullsTotalCounters[c.Name])
		PullsFailCounters[c.Name] = prometheus.NewCounter(prometheus.CounterOpts{
			Name:        prometheus.BuildFQName(namespace, "", "pulls_feed_fail_count"),
			Help:        "How many failed pulls.",
			ConstLabels: prometheus.Labels{"feed": c.Name},
		})
		prometheus.MustRegister(PullsFailCounters[c.Name])
	}
}

func (b *Bot) ServeMetrics() {
	http.Handle(b.Config.MetricsPath, promhttp.Handler())
	logrus.Infof("Listen address: %s", b.Config.ListenAddr)
	logrus.Fatal(http.ListenAndServe(b.Config.ListenAddr, nil))
}
