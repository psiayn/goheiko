# Heiko
Heiko rewritten in go! 
Heiko is a lightweight distributed node manager ( at least aims to be that ).

## Installation

### Directly from github
```
go get github.com/psiayn/heiko
```

### Install it from source
```
git clone https://github.com/psiayn/heiko.git
cd heiko
go install .
```

## Usage
General overview.
```
Usage:
  heiko [command]

Available Commands:
  help        Help about any command
  init        Runs initialization of Jobs
  start       Start a new heiko job
  stop        Stops a running heiko daemon

Flags:
      --config string   config file (default is $PWD/.heiko/config.yaml)
  -h, --help            help for heiko
  -n, --name string     Unique name to give (or given) to this heiko job

Use "heiko [command] --help" for more information about a command.
```

Heiko uses a `config.yml` to store info about jobs and nodes of the cluster.
A sample config has been provided in `examples/sample-config.yml`.
The default path for the config is at `.heiko/config.yml` in the current directory where you would like to start heiko from.
You can also specify config manually.
```
heiko start/init --config path/to/config
```

You can initialize heiko, which for now runs the init jobs from your `config.yml`. More about the config can be found in [Wiki](https://github.com/psiayn/heiko/wiki).
```
heiko init -n <name you want to give>
```

Starting heiko in normal mode
```
heiko start -n <name you want to give>
```

Starting heiko in daemon mode
```
heiko start -n <name you want to give> -d
```

Once your in daemon mode, you can stop the daemon as follows.
```
heiko stop -n <name of the daemon you gave earlier>
```
