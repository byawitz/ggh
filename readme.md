<p align="center">
    <img width="80" height="70" src="./assets/ggh.png" alt="GGH logo">
</p>
<h1 align="center"/>GGH</h1>

<p align="center"><i>Recall your SSH sessions</i></p>

<p align="center"><img src="./assets/ggh.gif" alt="GGH Demo"></p>


## Install

Run one of the following script, or download the latest binary from the [releases](https://github.com/byawitz/ggh/releases) page.

```shell
# Unix based
curl https://raw.githubusercontent.com/byawitz/ggh/master/install/unix.sh | sh

# Windows 
powershell -c "irm https://raw.githubusercontent.com/byawitz/ggh/master/install/windows.ps1 | iex"

# Go
go install github.com/byawitz/ggh@latest
```

## Usages

```shell
# Use it just like you're using SSH
ggh root@server.com
ggh root@server.com -p2440

# Run it with no arguments to get interactive list of the previous sessions
ggh

# Run it with - to get interactive list of all of your ~/.ssh/config listing
ggh - 

# Run it with - STRING to get interactive filtered list of your ~/.ssh/config listing
ggh - stage
ggh - meta-servers

# To get non-interactive list of history and config, run
ggh --config
ggh --history
```

### GGH is NOT replacing SSH

In fact, GGH won't work if SSH is not installed or isn't available in your system's path.

GGH is meant to act as a lightweight, fast wrapper around your SSH commands.