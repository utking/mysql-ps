# MysQL Processlist Viewer

The tool is created to have a more convenient way to see MySQL processlist, with time-bases descending order, ability to stop/resume refreshing the list, and a way to save some queries to a SQL log file to analyze them later.

## Requirements

- Go v1.20+

## Build

You can run the following command to build the ./bin/mysql-ps binary
```
make build
```
### Prepare a distribution archive

Run the following to build and archive the program and its .env
```
make tar
```

## How-To

### Pre-flight configuration

To run the compiled binary, you need to have proper MySQL DSN and credentials. If you have a ~/.my.cnf, the application will try to use it. Alternatively, you can configure the following env variables (or put them in the .env file):

- MYSQL_DSN - `tcp(HOST:MYSQL-PORT)` OR `@unix(/path/to/sock/file)`, default - `tcp(127.0.0.1:3306)`
- MYSQL_USER
- MYSQL_PASSWORD
- REFRESH_INTERVAL - process-list refresh rate, in seconds; default - 1. Any non-positive or non-numeric values will be reset to 1
- USE_MOUSE - `true` OR `1` will enable mouse navigation

Instead of putting these variables in .env, you can

### Runtime

There a several shurtcuts to use while using the application - hints will be shown on the top ot the progream screen

- P - pause/resume refreshing the process-list (when paused, the statusbar border color will change to yellow)
- Q - exit the program (no additional questions asked)
- Enter - when pressed on a process list line, will show the full SQL query in the bottom preview area
- Esc - when the preview area is visible, pressing Esc will hide it
- ? - press to show the shortcut hints
- L - when viewing the full SQL query, press L to switch back to the process-list area (without hiding the preview)
- V - open it and switch to the preview area
- Ctrl+S - clean the saved-queries file (queries.sql) and save the currently selected process-list line to it
- Ctrl+A - append the currently selected process-list line to the existing log file, or create the file if it is missing
