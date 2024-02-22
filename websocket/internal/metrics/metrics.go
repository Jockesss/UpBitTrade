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
	prometheus.MustRegister(cpuUsageGauge, diskUsageGauge, ramUsageGauge, webSocketConnectionGauge, rabbitMQConnectionGauge)
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

		//if webSocketIsConnected(webSocketClient) {
		//	webSocketConnectionGauge.Set(1) // Подключено
		//} else {
		//	webSocketConnectionGauge.Set(0) // Отключено
		//}

		//// Проверяем состояние подключения к RabbitMQ
		//if rabbitMQIsConnected() {
		//	rabbitMQConnectionGauge.Set(1) // Подключено
		//} else {
		//	rabbitMQConnectionGauge.Set(0) // Отключено
		//}

		time.Sleep(time.Duration(sec) * time.Second)
	}
}

//func webSocketIsConnected(ws *websocket.Conn) bool {
//	// Не блокирующее чтение
//	err := ws.SetReadDeadline(time.Now().Add(1 * time.Second))
//	if err != nil {
//		return false
//	}
//	defer func(ws *websocket.Conn, t time.Time) {
//		err := ws.SetReadDeadline(t)
//		if err != nil {
//
//		}
//	}(ws, time.Time{}) // Сброс после использования
//
//	if err := ws.WriteMessage(websocket.PingMessage, nil); err != nil {
//		return false // Если не можем отправить пинг, соединение мертво
//	}
//
//	_, _, err = ws.ReadMessage()
//	if err != nil {
//		return false // Если не получаем сообщение, соединение мертво
//	}
//
//	return true
//}

//func rabbitMQIsConnected() bool {
//	conn, err := rabbitmq.GetConnection(cfg)
//	if err != nil {
//		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
//	}
//	rabbitMQConnectionGauge.Set(0) // Предполагаем, что подключение разорвано
//	if conn, _ := rabbitmq.GetConnection(cfg); conn != nil && !conn.IsClosed() {
//		rabbitMQConnectionGauge.Set(1) // Подключение активно
//	}
//	return rabbitMQConn != nil && !rabbitMQConn.IsClosed()
//}
