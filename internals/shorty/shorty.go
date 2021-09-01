package shorty

import (
	"net/http"

	"git.6summits.net/srv/shorty/pkg/logging"
	"git.6summits.net/srv/shorty/pkg/mongodb"
	"goji.io"
	"goji.io/pat"
)

type Shorty struct {
	r *mongodb.Repository
	m *goji.Mux
}

func New(db *mongodb.Client) *Shorty {
	s := &Shorty{
		r: mongodb.NewRepository(db, "records"),
		m: goji.NewMux(),
	}

	return s
}

func (s *Shorty) registerHandler() {
	s.m.HandleFunc(pat.Get("/"), s.getUI)
	s.m.HandleFunc(pat.Get("/:short"), s.getShort)
	s.m.HandleFunc(pat.Post("/:long"), s.postLong)
}

func (s *Shorty) ListenAndServe(addr string) error {
	s.registerHandler()

	h := logging.Middleware(s.m)

	return http.ListenAndServe(addr, h)
}
