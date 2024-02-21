package ws

import (
	"context"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
	"net/http"
	"time"
	"upbit/internal/config"
	"upbit/internal/domain"
	"upbit/internal/ws/token"
	"upbit/pkg/log"
	"upbit/pkg/rabbitmq"
)

func WebsocketConnect(url string, t domain.Token) (*websocket.Conn, error) {
	header := http.Header{}
	header.Add("Authorization", "Bearer "+token.CreateToken(t))

	dialer := *websocket.DefaultDialer
	dialer.HandshakeTimeout = 10 * time.Second
	dialer.WriteBufferSize = 1024
	dialer.ReadBufferSize = 1024

	ws, _, err := websocket.DefaultDialer.Dial(url, header)

	return ws, err
}

type ConnectionManager struct {
	Ctx       context.Context
	Cancel    context.CancelFunc
	WsURL     string
	Token     domain.Token
	Platform  string
	WebSocket *websocket.Conn
	Cfg       *config.Config
}

func NewConnectionManager(ctx context.Context, url string, platform string, cfg *config.Config) *ConnectionManager {
	return &ConnectionManager{
		Ctx:      ctx,
		WsURL:    url,
		Platform: platform,
		Cfg:      cfg,
	}
}

func (cm *ConnectionManager) StartManager(ctx context.Context, wsURL string, token domain.Token, platform string, dataType string, restartChan chan<- string) {
	cm.Ctx = ctx
	cm.WsURL = wsURL
	cm.Token = token
	cm.Platform = platform
	cm.startConnection(restartChan, dataType)
}

func (cm *ConnectionManager) WebSocketIsConnected() bool {
	return cm.WebSocket != nil && cm.WebSocket.UnderlyingConn() != nil && cm.WebSocket.UnderlyingConn().RemoteAddr() != nil
}

func (cm *ConnectionManager) startConnection(restartChan chan<- string, dataType string) {
	for {
		select {
		case <-cm.Ctx.Done():
			return
		default:
			// If connection closed, sending the signal to reconnect
			cm.connectAndHandle(restartChan, dataType)
		}
	}
}

func (cm *ConnectionManager) connectAndHandle(restartChan chan<- string, dataType string) {
	backoff := 1
	maxBackoff := 120
	var queue string
	if "ticker" == dataType {
		queue = "ticker_queue"
	} else if "trade" == dataType {
		queue = "trade_queue"
	} else {
		log.Logger.Info("Unknown data type: " + dataType)
		return
	}

	for {
		select {
		case <-cm.Ctx.Done():
			log.Logger.Info("Context cancelled, stopping connection attempts")
			if cm.WebSocket != nil {
				if err := cm.WebSocket.Close(); err != nil {
					log.Logger.Error("Error closing WebSocket", zap.Error(err))
				}
			}
			return
		default:
			ws, err := WebsocketConnect(cm.WsURL, cm.Token)
			if err != nil {
				log.Logger.Error("Failed to connect: retrying...", zap.Error(err))
				if backoff < maxBackoff {
					backoff *= 2
				}
				continue
			}
			cm.WebSocket = ws
			cm.sendRequest(dataType, cm.Platform)
			cm.handleMessages(ws, queue)
			// Reset backoff after a successful connection
			backoff = 1
			// If connection closed, sending the signal to reconnect
			log.Logger.Info("Connection closed, signaling for reconnect")
			restartChan <- dataType
		}
	}
}

func (cm *ConnectionManager) handleMessages(ws *websocket.Conn, queue string) {
	producer, err := rabbitmq.NewProducer(cm.Cfg)
	if err != nil {
		log.Logger.Error("Could not create producer", zap.Error(err))
		return
	}
	defer producer.Close()

	for {
		_, message, err := ws.ReadMessage()
		if err != nil {
			log.Logger.Error("Failed to read a message", zap.Error(err))
			break // or handle the error as needed
		}

		// Вывод полученного сообщения в консоль
		//fmt.Printf("Received message for queue %s: %s\n", queue, string(message))

		if producer != nil {
			if err := producer.SendMessage(queue, string(message)); err != nil {
				log.Logger.Error("Failed to send message to queue "+queue, zap.Error(err))
			}
		} else {
			log.Logger.Error("Producer is nil")
		}

		select {
		case <-cm.Ctx.Done():
			log.Logger.Info("Context cancelled, stopping message handling")
			return
		default:
		}
	}
}

func (cm *ConnectionManager) sendRequest(dataType string, platform string) {
	if cm.WebSocket == nil {
		log.Logger.Info("WebSocket connection is nil")
		return
	}
	upBitMarkets := []string{"KRW-BTC", "KRW-ETH", "KRW-NEO", "KRW-MTL", "KRW-XRP", "KRW-ETC", "KRW-SNT", "KRW-WAVES", "KRW-XEM", "KRW-QTUM", "KRW-LSK", "KRW-STEEM", "KRW-XLM", "KRW-ARDR", "KRW-ARK", "KRW-STORJ", "KRW-GRS", "KRW-ADA", "KRW-SBD", "KRW-POWR", "KRW-BTG", "KRW-ICX", "KRW-EOS", "KRW-TRX", "KRW-SC", "KRW-ONT", "KRW-ZIL", "KRW-POLYX", "KRW-ZRX", "KRW-LOOM", "KRW-BCH", "KRW-BAT", "KRW-IOST", "KRW-CVC", "KRW-IQ", "KRW-IOTA", "KRW-HIFI", "KRW-ONG", "KRW-GAS", "KRW-UPP", "KRW-ELF", "KRW-KNC", "KRW-BSV", "KRW-THETA", "KRW-QKC", "KRW-BTT", "KRW-MOC", "KRW-TFUEL", "KRW-MANA", "KRW-ANKR", "KRW-AERGO", "KRW-ATOM", "KRW-TT", "KRW-CRE", "KRW-MBL", "KRW-WAXP", "KRW-HBAR", "KRW-MED", "KRW-MLK", "KRW-STPT", "KRW-ORBS", "KRW-VET", "KRW-CHZ", "KRW-STMX", "KRW-DKA", "KRW-HIVE", "KRW-KAVA", "KRW-AHT", "KRW-LINK", "KRW-XTZ", "KRW-BORA", "KRW-JST", "KRW-CRO", "KRW-TON", "KRW-SXP", "KRW-HUNT", "KRW-PLA", "KRW-DOT", "KRW-MVL", "KRW-STRAX", "KRW-AQT", "KRW-GLM", "KRW-SSX", "KRW-META", "KRW-FCT2", "KRW-CBK", "KRW-SAND", "KRW-HPO", "KRW-DOGE", "KRW-STRK", "KRW-PUNDIX", "KRW-FLOW", "KRW-AXS", "KRW-STX", "KRW-XEC", "KRW-SOL", "KRW-MATIC", "KRW-AAVE", "KRW-1INCH", "KRW-ALGO", "KRW-NEAR", "KRW-AVAX", "KRW-T", "KRW-CELO", "KRW-GMT", "KRW-APT", "KRW-SHIB", "KRW-MASK", "KRW-ARB", "KRW-EGLD", "KRW-SUI", "KRW-GRT", "KRW-BLUR", "KRW-IMX", "KRW-SEI", "KRW-MINA"}
	var request []byte
	switch dataType {
	case "ticker":
		request = tickerRequest(platform, upBitMarkets)
	case "trade":
		request = tradeRequest(platform, upBitMarkets)
	}
	if platform == "binance" {
		return
	}
	err := cm.WebSocket.WriteMessage(websocket.TextMessage, request)
	if err != nil {
		log.Logger.Error("Failed to write to websocket", zap.Error(err))
		return
	}
}
