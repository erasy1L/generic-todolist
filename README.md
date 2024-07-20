# Todo List

Simple todo list

## Features

- Add new tasks with a title and active date
- Mark tasks as done
- Edit existing tasks
- Delete tasks
- Filter tasks based on their completion status

## Installation & Usage

1. Clone the repository: `git clone https://github.com/erazr/todo-list.git`.
2. Navigate to the project directory: `cd todo-list`.
3. Rename .env.example to .env and change variables accordingly.
4. Start the docker containers: `make up`.
5. Navigate to swagger docs at http://localhost:8080/api/docs/index.htm.

## Libraries

1. [go-chi](https://github.com/go-chi/chi) as router
2. [zerolog](https://github.com/rs/zerolog) as logger
3. [golang-migrate](https://github.com/golang-migrate/migrate) for migrating the database
