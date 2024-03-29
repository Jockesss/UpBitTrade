package v1

import (
	"common/config"
	"common/pkg/log"
	"common/pkg/rabbitmq"
	"context"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
	"upbit/internal/domain"
	"upbit/internal/ws"
)

type Handler struct {
	cmMap          map[string]map[string]*HandlerEntry
	cm             *config.Config
	rabbitProducer *rabbitmq.Producer
}

type HandlerEntry struct {
	ws     *ws.ConnectionManager
	cancel context.CancelFunc
}

func NewHandler(config *config.Config, rabbitProducer *rabbitmq.Producer) *Handler {
	return &Handler{
		cmMap:          make(map[string]map[string]*HandlerEntry),
		cm:             config,
		rabbitProducer: rabbitProducer,
	}
}

func (h *Handler) startHandler(w http.ResponseWriter, r *http.Request) {
	var (
		UpBitToken = domain.Token{
			AccessKey: h.cm.UpBit.AccessKey,
			SecretKey: h.cm.UpBit.SecretKey,
		}
	)
	platform := chi.URLParam(r, "platform")
	dataType := chi.URLParam(r, "dataType")
	log.Logger.Info(fmt.Sprintf("Starting connection manager for %s with dataType %s", platform, dataType))
	restartChan := make(chan string, 10)

	if h.cmMap[platform] == nil {
		h.cmMap[platform] = make(map[string]*HandlerEntry)
	}

	if _, ok := h.cmMap[platform][dataType]; ok {
		log.Logger.Info(fmt.Sprintf("Connection manager for platform %s and dataType %s is already started", platform, dataType))
		fmt.Fprintf(w, "Connection manager for platform %s and dataType %s is already started", platform, dataType)
		return
	}

	ctx, cancel := context.WithCancel(context.Background())
	connManager := ws.NewConnectionManager(ctx, h.cm.UpBit.WsURL, platform, h.cm, h.rabbitProducer)

	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Logger.Info(fmt.Sprintf("Recovered in startManager for %s: %v", platform, r))
			}
		}()
		connManager.StartManager(ctx, h.cm.UpBit.WsURL, UpBitToken, platform, dataType, restartChan)
	}()

	h.cmMap[platform][dataType] = &HandlerEntry{
		ws:     connManager,
		cancel: cancel,
	}

	log.Logger.Info(fmt.Sprintf("Connection manager for platform %s with dataType %s started successfully", platform, dataType))
	fmt.Fprintf(w, "Connection manager for platform %s with dataType %s started successfully", platform, dataType)
}

func (h *Handler) stopHandler(w http.ResponseWriter, r *http.Request) {
	platform := chi.URLParam(r, "platform")
	dataType := chi.URLParam(r, "dataType")

	if dataTypeMap, ok := h.cmMap[platform]; ok {
		if entry, ok := dataTypeMap[dataType]; ok {
			entry.cancel()
			delete(dataTypeMap, dataType)
			log.Logger.Info(fmt.Sprintf("Connection manager for platform %s with dataType %s stopped successfully", platform, dataType))
			fmt.Fprintf(w, "Connection manager for platform %s with dataType %s stopped successfully", platform, dataType)
			return
		}
	}

	log.Logger.Info(fmt.Sprintf("Connection manager for platform %s with dataType %s not found or already stopped", platform, dataType))
	fmt.Fprintf(w, "Connection manager for platform %s with dataType %s not found or already stopped", platform, dataType)
}

func (h *Handler) Routes() *chi.Mux {
	router := chi.NewRouter()
	router.Get("/start/{platform}/{dataType}", h.startHandler)
	router.Get("/stop/{platform}/{dataType}", h.stopHandler)
	router.Handle("/metrics", promhttp.Handler())
	return router
}
