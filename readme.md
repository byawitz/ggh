<p align="center">
    <img width="80" height="70" src="./assets/ggh.png" alt="Elk logo">
</p>
<h1 align="center"/>GGH</h1>

<p align="center"><i>Recall your SSH sessions</i></p>

## Install

Run one of the following script, or download the binary from the [latest release](https://github.com/byawitz/ggh/releases) page.

```shell

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
```