# Todo CLI

A simple command-line todo application built with Go and Cobra that helps you manage your tasks efficiently.
A simple todo CLI built with Go and Sqlite. 

## Features

- Add, edit, complete and remove todo items
- List active and completed todos
- Automatic timestamping of completed tasks
- Customizable database file locations via environment variables
- Cross-platform support (Linux, macOS, Windows)

## Installation

### Download from Releases

Download the latest binary for your platform from the [releases page](https://github.com/ali-chapman/todo/releases).

#### Linux/macOS
```bash
# Download and install (replace with actual download URL)
wget https://github.com/ali-chapman/todo/releases/latest/download/todo-linux-amd64
chmod +x todo-linux-amd64
sudo mv todo-linux-amd64 /usr/local/bin/todo
```

#### Windows
Download `todo-windows-amd64.exe` and place it in your PATH.

### Build from Source

```bash
git clone https://github.com/ali-chapman/todo.git
cd todo
go build -o todo .
```

## Usage

### Basic Commands

#### List todos
```bash
todo
```

#### Add a todo
```bash
todo Buy groceries
todo Finish project documentation
```

#### Add a todo with tags
```bash
todo Write tests for @project1
```
All words that start with a '@' symbol will be stored as tags on the todo.

#### Complete a todo
```bash
todo -d 0  # Remove todo at index 0
todo --done 1  # Remove todo at index 1
```

#### Edit a todo
```bash
todo -e 0 "Updated todo text"
todo --edit 1 "New description"
```

#### List completed todos
```bash
todo -sd
todo --status done
```

#### List all todos
```bash
todo -sa
todo --status all
```

#### Show creation date
```bash
todo -c
todo --created
```

#### Show completion date
```bash
todo -C
todo --completed
```

#### Hide tags
```bash
todo -n
todo --hide-tags
```

### Examples

```bash
# Add some todos
$ todo "Learn Go programming"
Adding todo: Learn Go programming
[0]: Learn Go programming

$ todo "Write unit tests"
Adding todo: Write unit tests
[0]: Learn Go programming
[1]: Write unit tests

# Complete a todo
$ todo -r 0
Completed: Learn Go programming
[0]: Write unit tests

# View completed todos with timestamps
$ todo -d
[0]: 2024-01-15 10:30:25 - Learn Go programming

# Edit a todo
$ todo -e 0 "Write comprehensive unit tests"
[0]: Write comprehensive unit tests
```

## Configuration

### File Locations

By default, the sqlite database file is stored in: `~/.config/.todo.db`

You can customize this using an environment variable:

```bash
export TODO_PATH="/path/to/your/todo/database_file"
```

## Command Reference

| Command | Flag | Description |
|---------|------|-------------|
| `todo` | | List all active todos |
| `todo "item"` | | Add a new todo item |
| `todo "item @tag1 @tag2"` | | Add a new todo item with tags |
| `todo -d <index>` | `--done` | Mark todo at index as complete |
| `todo -e <index> "text"` | `--edit` | Edit todo at index |
| `todo -x <index>` | `--remove` | Remove todo at index from database, THIS CANNOT BE UNDONE |
| `todo -sd` | `--status done` | List completed todos |
| `todo -sa` | `--status all` | List all todos |
| `todo -c` | `--created` | Display relative creation time |
| `todo -C` | `--completed` | Display relative completion time |
| `todo -n` | `--hide-tags` | Do not display tags for each todo |
| `todo -T mytag` | `--tag mytag` | Filter todos by tag |
