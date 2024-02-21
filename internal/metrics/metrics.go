package metrics

import (
	"github.com/gorilla/websocket"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/mem"
	"github.com/streadway/amqp"
	"time"
)

var (
	webSocketClient *websocket.Conn
	rabbitMQConn    *amqp.Connection
	cpuUsageGauge   = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "cpu_usage_wsdr",
		Help: "Current CPU usage percentage",
	})
	diskUsageGauge = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "disk_usage_wsdr",
		Help: "Current disk usage percentage",
	})
	ramUsageGauge = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "ram_usage_wsdr",
		Help: "Current RAM usage percentage",
	})
	webSocketConnectionGauge = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "websocket_connection_status",
		Help: "Current WebSocket connection status (1: connected, 0: disconnected)",
	})
	rabbitMQConnectionGauge = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "rabbitmq_connection_status",
		Help: "Current RabbitMQ connection status (1: connected, 0: disconnected)",
	})
)

func init() {
	prometheus.MustRegister(cpuUsageGauge, diskUsageGauge, ramUsageGauge)
}

func UpdateResourceUsageMetrics(sec int) {
	for {
		cpuPercent, _ := cpu.Percent(0, false)
		if len(cpuPercent) > 0 {
			cpuUsageGauge.Set(cpuPercent[0])
		}

		diskUsage, _ := disk.Usage("/")
		diskUsageGauge.Set(diskUsage.UsedPercent)

		ramStats, _ := mem.VirtualMemory()
		ramUsageGauge.Set(ramStats.UsedPercent)

		if webSocketIsConnected() {
			webSocketConnectionGauge.Set(1) // Подключено
		} else {
			webSocketConnectionGauge.Set(0) // Отключено
		}

		// Проверяем состояние подключения к RabbitMQ
		if rabbitMQIsConnected() {
			rabbitMQConnectionGauge.Set(1) // Подключено
		} else {
			rabbitMQConnectionGauge.Set(0) // Отключено
		}

		time.Sleep(time.Duration(sec) * time.Second)
	}
}

//func UpdateConnectionMetrics(webSocketClient *websocket.Conn, rabbitMQConn *amqp.Connection) {
//	for {
//		// Проверяем состояние подключения к WebSocket
//		if webSocketIsConnected(webSocketClient) {
//			webSocketConnectionGauge.Set(1) // Подключено
//		} else {
//			webSocketConnectionGauge.Set(0) // Отключено
//		}
//
//		// Проверяем состояние подключения к RabbitMQ
//		if rabbitMQIsConnected(rabbitMQConn) {
//			rabbitMQConnectionGauge.Set(1) // Подключено
//		} else {
//			rabbitMQConnectionGauge.Set(0) // Отключено
//		}
//
//		time.Sleep(10 * time.Second) // Например, проверяем каждые 10 секунд
//	}
//}

func webSocketIsConnected() bool {
	return webSocketClient != nil && webSocketClient.UnderlyingConn() != nil && webSocketClient.UnderlyingConn().RemoteAddr() != nil
}

func rabbitMQIsConnected() bool {
	return rabbitMQConn != nil && !rabbitMQConn.IsClosed()
}
