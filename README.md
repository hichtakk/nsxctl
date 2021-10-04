nsxctl
======

`nsxctl` acts as command-line client for NSX-T.

## Usage
First of all, you need to configure NSX-T endpoint and credential for subsequent operation.

```
nsxctl config set-site ${NSX-T-SITE} --endpoint https://${MANAGER-IP} --user ${USER} --password ${PASSWORD} --init
```

This configuration is stored `~/.config/nsxctl.json` by default. You can change the path and file name with `-c/--config` option. Sensitive data such as user credential will be base64 encoded.

Then you could run nsxctl to get and create NSX-T inventories.

```
nsxctl show gateway
```

## Options

```
nsxctl 
modern NSX-T client

Usage:
  nsxctl [command]

Available Commands:
  completion  generate the autocompletion script for the specified shell
  config      configuration
  help        Help about any command
  show        Show resources

Flags:
  -c, --config string   path to nsxctl config file (default "~/.config/nsxctl.json")
      --debug           enable debug mode
  -h, --help            help for nsxctl

Use "nsxctl [command] --help" for more information about a command.
```