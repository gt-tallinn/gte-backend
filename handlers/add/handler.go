package add

import (
	"net/http"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"encoding/json"
	"context"
)

type StorageWriter interface {
	Write(context.Context, *Request) error
}

type Handler struct {
	logger  *logrus.Entry
	storage StorageWriter
}

func New(logger *logrus.Entry, storage StorageWriter) *Handler {
	return &Handler{
		logger:  logger,
		storage: storage,
	}
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	bytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		h.logger.WithError(err).Error("Can't read body")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	h.logger.WithField("request", "/add").Debug(string(bytes))
	r.Body.Close()
	var req Request
	err = json.Unmarshal(bytes, &req)
	if err != nil {
		h.logger.WithError(err).Error("Can't unmarshal request")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = h.storage.Write(context.Background(), &req)
	if err != nil {
		h.logger.WithError(err).Error("Can't write data to storage")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	return
}

type Request struct {
	ID      string `json:"id" bson:"id"`
	Type    string `json:"type" bson:"type"`
	Service string `json:"service" bson:"service"`
	Context string `json:"context" bson:"context"`
	Start   int64  `json:"start" bson:"start"`
	Finish  int64  `json:"finish" bson:"finish"`
}
