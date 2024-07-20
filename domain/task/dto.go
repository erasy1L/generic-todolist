package task

import (
	"errors"
	"net/http"
)

type Request struct {
	Title    string `json:"title"`
	ActiveAt string `json:"activeAt"`
	Status   string `json:"status,omitempty"`
}

const dateLayout = "2006-01-02"

func (i *Request) Bind(r *http.Request) error {
	if i.Title == "" {
		return errors.New("title is required")
	}

	if len(i.Title) >= 200 {
		return errors.New("title must be less than 200 characters")
	}

	if i.Status != "active" && i.Status != "done" {
		if i.Status == "" {
			i.Status = "active"
			return nil
		}

		return errors.New("invalid status")
	}

	return nil
}

type Response struct {
	ID       string `json:"id"`
	Title    string `json:"title"`
	ActiveAt string `json:"activeAt"`
	Status   string `json:"status,omitempty"`
}

func ParseFromEntity(e Entity) Response {
	return Response{
		ID:       e.ID,
		Title:    e.Title,
		ActiveAt: string(e.ActiveAt),
		Status:   e.Status,
	}
}

func ParseFromEntities(e Entity) Entity {
	return Entity{
		ID:       e.ID,
		Title:    e.Title,
		ActiveAt: e.ActiveAt,
		Status:   e.Status,
	}
}
