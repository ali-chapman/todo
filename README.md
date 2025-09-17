# Todo CLI

A simple todo CLI built with Go and Sqlite.

*Why did you use Sqlite instead of just text files?*

Honestly just because I wanted to. It also makes it a bit snappier to do the
filtering. You can use the `--format` flag to customize the output to
csv, json, or simple text file which you could store in version control.

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

## Configuration

### File Locations

By default, the sqlite database file is stored in: `~/.config/.todo.db`

You can customize this using an environment variable:

```bash
export TODO_PATH="/path/to/your/todo/database_file"
```

### Todo tags

You can add tags to a todo by appending `@` to a word within the new todo, for example:
```bash
todo @project1 do something relating to project1 @urgent
```
When listing the todo you will see `[project1, urgent]` as the tags on this todo.

You can filter using the `-T` flag:
```bash
todo -T urgent
```

If you set the environment variable TODO_TAG, then this will be added as a tag to all new
todos, and when listing the todos will be filtered by this tag.

## Command Reference

| Command | Flag | Description |
|---------|------|-------------|
| `todo` | | List all active todos |
| `todo "item"` | | Add a new todo item |
| `todo "item @tag1 @tag2"` | | Add a new todo item with tags |
| `todo -d <index>` | `--done` | Mark todo at index as complete |
| `todo -e <index> "text"` | `--edit` | Edit todo at index |
| `todo -f <string>` | `--format` | Output format (table, json, csv, txt), defaults to "table"
| `todo -x <index>` | `--remove` | Remove todo at index from database, THIS CANNOT BE UNDONE |
| `todo -sd` | `--status done` | List completed todos |
| `todo -sa` | `--status all` | List all todos |
| `todo -c` | `--created` | Display relative creation time |
| `todo -C` | `--completed` | Display relative completion time |
| `todo -n` | `--hide-tags` | Do not display tags for each todo |
| `todo -T mytag` | `--tag mytag` | Filter todos by tag. If not set then uses TODO_TAG environment variable |
