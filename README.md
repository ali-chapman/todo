# Todo CLI

A simple command-line todo application built with Go and Cobra that helps you manage your tasks efficiently.

## Features

- Add, edit, and remove todo items
- List active and completed todos
- Automatic timestamping of completed tasks
- Customizable file locations via environment variables
- Cross-platform support (Linux, macOS, Windows)

## Installation

### Download from Releases

Download the latest binary for your platform from the [releases page](https://github.com/your-username/todo/releases).

#### Linux/macOS
```bash
# Download and install (replace with actual download URL)
wget https://github.com/your-username/todo/releases/latest/download/todo-linux-amd64
chmod +x todo-linux-amd64
sudo mv todo-linux-amd64 /usr/local/bin/todo
```

#### Windows
Download `todo-windows-amd64.exe` and place it in your PATH.

### Build from Source

```bash
git clone https://github.com/your-username/todo.git
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
todo "Buy groceries"
todo "Finish project documentation"
```

#### Remove a todo (mark as complete)
```bash
todo -r 0  # Remove todo at index 0
todo --remove 1  # Remove todo at index 1
```

#### Edit a todo
```bash
todo -e 0 "Updated todo text"
todo --edit 1 "New description"
```

#### List completed todos
```bash
todo -d
todo --done
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

By default, todos are stored in:
- Active todos: `~/.config/.todo`
- Completed todos: `~/.config/.done`

You can customize these locations using environment variables:

```bash
export TODO_PATH="/path/to/your/todo/file"
export DONE_PATH="/path/to/your/done/file"
```

## Command Reference

| Command | Flag | Description |
|---------|------|-------------|
| `todo` | | List all active todos |
| `todo "item"` | | Add a new todo item |
| `todo -r <index>` | `--remove` | Mark todo at index as complete |
| `todo -e <index> "text"` | `--edit` | Edit todo at index |
| `todo -d` | `--done` | List completed todos with timestamps |
