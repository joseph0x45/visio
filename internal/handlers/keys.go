package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/jmoiron/sqlx"
	"github.com/oklog/ulid/v2"
	"log/slog"
	"math/rand"
	"net/http"
	"time"
	"visio/internal/store"
	"visio/internal/types"
	"visio/pkg"
)

type KeyHandler struct {
	db       *sqlx.DB
	keys     *store.Keys
	logger   *slog.Logger
	sessions *store.Sessions
}

func NewKeyHandler(db *sqlx.DB, keysStore *store.Keys, sessionsStore *store.Sessions, logger *slog.Logger) *KeyHandler {
	return &KeyHandler{
		db:       db,
		keys:     keysStore,
		sessions: sessionsStore,
		logger:   logger,
	}
}

func generateRandomString(length int) string {
	const CHARACTER_POOL = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz1234567890-_=+:;?><|"
	key := ""
	for i := 0; i < length; i++ {
		idx := rand.Intn(len(CHARACTER_POOL))
		key += string(CHARACTER_POOL[idx])
	}
	return key
}

func (h *KeyHandler) GetNew(w http.ResponseWriter, r *http.Request) {
	log := h.logger.With("requestid", chiMiddleware.GetReqID(r.Context()))
	currentUser, ok := r.Context().Value("currentUser").(*types.User)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	tx, err := h.db.Beginx()
	if err != nil {
		log.Error(fmt.Sprintf("Error while starting transaction: %s", err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	err = h.keys.Delete(tx, currentUser.Id)
	if err != nil {
		log.Error(err.Error())
		err = tx.Rollback()
		if err != nil {
			log.Error(fmt.Sprintf("Error while rolling back transaction: %s", err.Error()))
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	prefix := generateRandomString(7)
	suffix := generateRandomString(23)
	key := fmt.Sprintf("%s.%s", prefix, suffix)
	keyHash, err := pkg.Hash(suffix)
	if err != nil {
		log.Error(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	newKey := &types.Key{
		Id:           ulid.Make().String(),
		UserId:       currentUser.Id,
		Prefix:       prefix,
		KeyHash:      keyHash,
		CreationDate: time.Now().UTC().Format("January, 2 2006"),
	}
	err = h.keys.Insert(newKey, tx)
	if err != nil {
		if errors.Is(err, types.ErrDuplicatePrefix) {
			log.Debug("A duplicate prefix error occured")
		}
		log.Error(err.Error())
		err = tx.Rollback()
		if err != nil {
			log.Error(fmt.Sprintf("Error while rolling back transaction: %s", err.Error()))
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	data, err := json.Marshal(
		map[string]interface{}{
			"data": map[string]interface{}{
				"naked_key":  key,
				"key_object": newKey,
			},
		},
	)
	if err != nil {
		log.Error(err.Error())
		err = tx.Rollback()
		if err != nil {
			log.Error(fmt.Sprintf("Error while rolling back transaction: %s", err.Error()))
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	err = tx.Commit()
	if err != nil {
		log.Error(fmt.Sprintf("Error while commiting transaction: %s", err.Error()))
		err = tx.Rollback()
		if err != nil {
			log.Error(fmt.Sprintf("Error while rolling back transaction: %s", err.Error()))
		}
	}
	w.WriteHeader(http.StatusCreated)
	w.Write(data)
	return
}
