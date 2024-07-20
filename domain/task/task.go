package task

import (
	"fmt"
	"time"
)

type Entity struct {
	ID       string
	Title    string
	ActiveAt OnlyDate `db:"active_at"`
	Status   string
}

var (
	ErrExists   = &TaskError{"task already exists"}
	ErrNotFound = &TaskError{"task not found"}
)

type TaskError struct {
	message string
}

func (e *TaskError) Error() string {
	return e.message
}

func (e *TaskError) Is(err error) bool {
	return e == err
}

type OnlyDate string

// method of [driver.Valuer] interface
func (o *OnlyDate) Value() (interface{}, error) {
	nt, err := time.Parse(dateLayout, string(*o))
	if err != nil {
		return nil, err
	}

	return nt, nil
}

// method of [sql.Scanner] interface
func (o *OnlyDate) Scan(val interface{}) error {
	nt, ok := val.(time.Time)
	if !ok {
		return fmt.Errorf("expected time.Time, got %T", val)
	}

	*o = OnlyDate(nt.Format(dateLayout))

	return nil
}

func (t *Entity) TransformByWeekday() (*Entity, error) {
	nt, err := t.ActiveAt.Value()
	if err != nil {
		return nil, err
	}

	date, ok := nt.(time.Time)
	if !ok {
		return nil, err
	}

	if date.Weekday() == time.Saturday || date.Weekday() == time.Sunday {
		t.Title = "ВЫХОДНОЙ - " + t.Title
		return t, nil
	}

	return t, nil
}
