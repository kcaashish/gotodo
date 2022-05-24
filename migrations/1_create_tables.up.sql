CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY,
    username VARCHAR(50) NOT NULL,
    first_name VARCHAR(50) NOT NULL,
    last_name VARCHAR(50) NOT NULL,
    email VARCHAR(300) NOT NULL,
    password VARCHAR(50) NOT NULL
);

CREATE TABLE IF NOT EXISTS todo_list (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL,
    title TEXT NOT NULL,
    description TEXT NOT NULL,
    created_date TIMESTAMPTZ NOT NULL,
    updated_date TIMESTAMPTZ NOT NULL,
    due_date TIMESTAMPTZ NOT NULL,
    completed BOOLEAN NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS todo_entry (
    id UUID PRIMARY KEY,
    todolist_id UUID NOT NULL,
    content TEXT NOT NULL,
    created_date TIMESTAMPTZ NOT NULL,
    updated_date TIMESTAMPTZ NOT NULL,
    due_date TIMESTAMPTZ NOT NULL,
    completed BOOLEAN NOT NULL,
    FOREIGN KEY (todolist_id) REFERENCES todo_list (id) ON DELETE CASCADE
);