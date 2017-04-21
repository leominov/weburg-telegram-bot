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
	PullsTotalCounter = prometheus.NewCounter(prometheus.CounterOpts{
		Name: prometheus.BuildFQName(namespace, "", "pulls_total_count"),
		Help: "How many pulls.",
	})

	PullsFailCounter = prometheus.NewCounter(prometheus.CounterOpts{
		Name: prometheus.BuildFQName(namespace, "", "pulls_fail_count"),
		Help: "How many failed pulls.",
	})

	MessagesTotalCounter = prometheus.NewCounter(prometheus.CounterOpts{
		Name: prometheus.BuildFQName(namespace, "", "messages_total_count"),
		Help: "How many messages are sent.",
	})

	MessagesFailCounter = prometheus.NewCounter(prometheus.CounterOpts{
		Name: prometheus.BuildFQName(namespace, "", "messages_fail_count"),
		Help: "How many messages are fail.",
	})
)

func InitMetrics() {
	prometheus.MustRegister(version.NewCollector(namespace))

	prometheus.MustRegister(PullsTotalCounter)
	prometheus.MustRegister(PullsFailCounter)
	prometheus.MustRegister(MessagesTotalCounter)
	prometheus.MustRegister(MessagesFailCounter)
}

func (w *WeburgBot) ServeMetrics() {
	http.Handle(w.MetricsPath, promhttp.Handler())
	logrus.Infof("Listen address: %s", w.ListenAddr)
	logrus.Fatal(http.ListenAndServe(w.ListenAddr, nil))
}
