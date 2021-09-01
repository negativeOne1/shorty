package shorty

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"

	"git.6summits.net/srv/shorty/pkg/random"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"goji.io/pat"
)

var (
	supportedScheme = []string{"http", "https", "ftp", "ftps", "mailto", "mms", "rtmp", "rtmpt", "ed2k", "pop", "imap", "nntp", "news", "ldap", "gopher", "dict", "dns"}
)

func inStringList(v []string, s string) bool {
	for _, v := range v {
		if v == s {
			return true
		}
	}
	return false
}

func (s *Shorty) getShort(w http.ResponseWriter, r *http.Request) {
	shortURL := pat.Param(r, "short")
	ctx := context.Background()
	re := &Record{}

	if err := s.r.FindOne(ctx, primitive.M{"short_url": shortURL}, re); err != nil {
		log.Error().Err(err).Msg("")
		http.NotFound(w, r)
		return
	}
	if re == nil {
		http.NotFound(w, r)
		return
	}

	http.Redirect(w, r, re.LongURL, http.StatusSeeOther)
}

func (s *Shorty) postLong(w http.ResponseWriter, r *http.Request) {
	longURLEncoded := pat.Param(r, "long")
	ctx := context.Background()

	longURL, err := url.ParseRequestURI(longURLEncoded)
	if err != nil {
		log.Error().Err(err).Msg("")
		http.Error(w, "that doesn't look like a valid URL", http.StatusBadRequest)
		return
	}

	if !inStringList(supportedScheme, longURL.Scheme) {
		log.Error().Str("URL", longURLEncoded).Msg("scheme not supported")
		http.Error(w, "scheme not supported", http.StatusBadRequest)
		return
	}

	token, err := random.GetToken(3)
	if err != nil {
		log.Error().Err(err).Msg("")
		return
	}

	re := Record{
		ID:       primitive.NewObjectID(),
		LongURL:  longURLEncoded,
		ShortURL: token,
	}

	_, err = s.r.CreateDocument(ctx, re)
	if err != nil {
		log.Error().Err(err).Msg("")
		return
	}

	b, err := json.MarshalIndent(re, "", "  ")

	if err != nil {
		log.Error().Err(err).Msg("")
		return
	}

	if _, err := w.Write(b); err != nil {
		log.Error().Err(err).Msg("")
		return
	}
}

func (s *Shorty) getUI(http.ResponseWriter, *http.Request) {
}
