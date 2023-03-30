package handlers

import (
	"net"
	"net/http"

	httplib "github.com/adykaaa/online-notes/lib/http"
	"github.com/google/uuid"
)

// conn, _, _, err := ws.UpgradeHTTP(r, w)
type Room struct {
	ID        uuid.UUID
	Name      string
	CreatedBy string
	Clients   map[net.Conn]bool
}

func CreateRoom() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l, ctx, cancel := httplib.SetupHandler(w, r.Context())
		defer cancel()

		req := struct {
			Name     string `json:"roomName"`
			Password string `json:""`
		}{}
	}
}
