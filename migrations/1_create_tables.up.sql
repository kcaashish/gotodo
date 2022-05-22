CREATE TABLE users (
    id UUID PRIMARY KEY,
    username TEXT NOT NULL,
    first_name TEXT NOT NULL,
    last_name TEXT NOT NULL,
    email TEXT NOT NULL,
    password TEXT NOT NULL
)

CREATE TABLE todo_list (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL,
    title TEXT NOT NULL,
    description TEXT NOT NULL,
    created_date TIMESTAMP NOT NULL,
    updated_date TIMESTAMP NOT NULL,
    due_date TIMESTAMP NOT NULL,
    completed BOOLEAN NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
)

CREATE TABLE todo_entry (
    id UUID PRIMARY KEY,
    todolist_id UUID NOT NULL,
    content TEXT NOT NULL,
    created_date TIMESTAMP NOT NULL,
    updated_date TIMESTAMP NOT NULL,
    due_date TIMESTAMP NOT NULL,
    completed BOOLEAN NOT NULL,
    FOREIGN KEY (todo_list_id) REFERENCES todo_list (id) ON DELETE CASCADE
)