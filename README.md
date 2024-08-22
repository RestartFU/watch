### Installation

To install the `watch` tool, follow these steps:

1. Clone the repository:
    ```sh
    git clone https://github.com/restartfu/watch
    ```

2. Change directory to the cloned repository:
    ```sh
    cd watch
    ```

3. Build the tool:
    ```sh
    go build -o watch .
    ```

4. Make the binary executable:
    ```sh
    chmod +x watch
    ```

5. Move the binary to a directory in your `PATH`:
    ```sh
    mv watch /usr/bin/watch
    ```

Now, you can run the `watch` command from anywhere in your terminal.
# Deployment Language Documentation

## Introduction

The deployment language provides a simple way to automate the deployment of applications from source code repositories. The language is designed to handle common tasks such as cloning repositories, setting variables, running commands, and extracting files.

## Syntax Overview

The deployment language consists of the following directives:

1. **CLONE**: Clones a repository from a given URL and branch/tag.
2. **SET**: Sets variables for use in subsequent directives.
3. **RUN**: Executes a command in the deployment environment.
4. **EXTRACT**: Copies files from one location to another.

## Directives

### CLONE
```CLONE (repository_url@branch) AS (alias)```

- **repository_url**: The URL of the repository to clone.
- **branch**: The branch or tag to checkout.
- **alias**: A name to refer to the cloned repository path in subsequent directives.

**Example:**
```CLONE (github.com/go-training/helloworld@master) AS (STABLE)```

### SET
```SET (VARIABLE_NAME=value)```

- **VARIABLE_NAME**: The name of the variable to set.
- **value**: The value of the variable.

**Example:**
```
SET (BINARY_NAME=sample_go)
SET (BINARY_PATH=./../../bin/$[BINARY_NAME])
```

### RUN
```RUN (command)```

- **command**: The command to execute in the deployment environment.

**Example:**
```
RUN (cd $[STABLE] && go build -o $[BINARY_NAME] main.go)
RUN ($[BINARY_PATH])
```

### EXTRACT
```EXTRACT (source_path destination_path)```


- **source_path**: The path of the files to copy.
- **destination_path**: The path where the files will be copied.

**Example:**
```EXTRACT ($[STABLE]/$[BINARY_NAME] $[BINARY_PATH])```


## Example Configurations

### Python Deployment Example

**Watchfile:**
```
CLONE (github.com/dbarnett/python-helloworld@main) AS (STABLE)
SET (MAIN_NAME=helloworld.py)
SET (PATH=./../../bin/py)

RUN (mkdir $[PATH])
EXTRACT ($[STABLE]/* $[PATH])
RUN (python $[PATH]/$[MAIN_NAME])
```


**Explanation:**

1. Clone the `python-helloworld` repository from GitHub and name it `STABLE`.
2. Set `MAIN_NAME` to `helloworld.py` and `PATH` to `./../../bin/py`.
3. Create the directory specified by `PATH`.
4. Extract all files from the cloned repository to the created directory.
5. Run the Python script located at `$[PATH]/$[MAIN_NAME]`.

### Go Deployment Example

**Watchfile:**
```
CLONE (github.com/go-training/helloworld@master) AS (STABLE)
SET (BINARY_NAME=sample_go)
SET (BINARY_PATH=./../../bin/$[BINARY_NAME])

RUN (cd $[STABLE] && go build -o $[BINARY_NAME] main.go)
EXTRACT ($[STABLE]/$[BINARY_NAME] $[BINARY_PATH])
RUN ($[BINARY_PATH])
```


**Explanation:**

1. Clone the `helloworld` repository from GitHub and name it `STABLE`.
2. Set `BINARY_NAME` to `sample_go` and `BINARY_PATH` to `./../../bin/$[BINARY_NAME]`.
3. Change directory to `STABLE` and build the Go application, outputting the binary as `sample_go`.
4. Extract the built binary from the cloned repository to the specified `BINARY_PATH`.
5. Run the binary located at `$[BINARY_PATH]`.

## Custom Deployment Language Example

### Node.js Deployment Example

**Watchfile:**
```
CLONE (github.com/user/node-app@main) AS (STABLE)
SET (MAIN_SCRIPT=index.js)
SET (INSTALL_DIR=./../../node_modules)

RUN (mkdir -p $[INSTALL_DIR])
RUN (cd $[STABLE] && npm install --prefix $[INSTALL_DIR])
EXTRACT ($[STABLE]/* $[INSTALL_DIR])
RUN (node $[INSTALL_DIR]/$[MAIN_SCRIPT])
```

**Explanation:**

1. Clone the `node-app` repository from GitHub and name it `STABLE`.
2. Set `MAIN_SCRIPT` to `index.js` and `INSTALL_DIR` to `./../../node_modules`.
3. Create the directory specified by `INSTALL_DIR`.
4. Change directory to `STABLE` and run `npm install` to install dependencies into `INSTALL_DIR`.
5. Extract all files from the cloned repository to the `INSTALL_DIR`.
6. Run the Node.js script located at `$[INSTALL_DIR]/$[MAIN_SCRIPT]`.

