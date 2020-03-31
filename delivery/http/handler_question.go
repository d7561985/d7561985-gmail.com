package http

import (
	"net/http"

	"github.com/d7561985/questions/internal/tr"
	"github.com/d7561985/questions/model"
	"github.com/go-chi/chi"
)

type questionHandler struct {
	*factory
}

func (t *questionHandler) route() chi.Router {
	r := chi.NewRouter()
	r.Get("/", t.list)
	r.Post("/", t.create)

	return r
}

func (t *questionHandler) list(w http.ResponseWriter, r *http.Request) {
	l := tr.FromCtx(r.Context())

	res, err := t.uc.QuestionList(r.Context(), r.URL.Query().Get("lang"))
	if err != nil {
		Send(w, r.Header.Get("Content-Type"), nil, http.StatusInternalServerError, l, err)
		return
	}

	Send(w, r.Header.Get("Content-Type"), res, http.StatusOK, l, nil)
}

func (t *questionHandler) create(w http.ResponseWriter, r *http.Request) {
	l, m := tr.FromCtx(r.Context()), model.Question{}

	if err := Read(r, &m); err != nil {
		Send(w, r.Header.Get("Content-Type"), nil, http.StatusInternalServerError, l, err)
		return
	}

	res, err := t.uc.AddQuestion(r.Context(), m)
	if err != nil {
		Send(w, r.Header.Get("Content-Type"), nil, http.StatusInternalServerError, l, err)
		return
	}

	Send(w, r.Header.Get("Content-Type"), res, http.StatusOK, l, nil)
}
