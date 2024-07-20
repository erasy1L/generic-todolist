package server

import (
	"errors"
	"net/http"
	"todo-list/domain/task"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

func HandlerFromMux(s *HTTPServer, router *chi.Mux) http.Handler {
	router.Route("/todo-list/tasks", func(r chi.Router) {
		r.Post("/", s.create)
		r.Get("/", s.list)
		r.Route("/{id}", func(r chi.Router) {
			r.Put("/", s.update)
			r.Delete("/", s.delete)
			r.Put("/done", s.done)
		})
	})

	return router
}

func (s *HTTPServer) create(w http.ResponseWriter, r *http.Request) {
	var t task.Request
	err := render.Bind(r, &t)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	id, err := s.taskService.CreateTask(r.Context(), t)
	if err != nil {
		if errors.Is(err, task.ErrExists) {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(id))
}

func (s *HTTPServer) update(w http.ResponseWriter, r *http.Request) {
	var t task.Request
	err := render.Bind(r, &t)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	id := chi.URLParam(r, "id")

	err = s.taskService.UpdateTask(r.Context(), id, t)
	if err != nil {
		if errors.Is(err, task.ErrNotFound) {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (s *HTTPServer) delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	err := s.taskService.DeleteTask(r.Context(), id)
	if err != nil {
		if errors.Is(err, task.ErrNotFound) {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (s *HTTPServer) done(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	err := s.taskService.DoneTask(r.Context(), id)
	if err != nil {
		if errors.Is(err, task.ErrNotFound) {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (s *HTTPServer) list(w http.ResponseWriter, r *http.Request) {
	status := chi.URLParam(r, "status")

	tasks, err := s.taskService.ListTasks(r.Context(), status)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	render.JSON(w, r, tasks)
}
