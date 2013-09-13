# gocat

Socket activated transparent SSL proxy written in Go.
The goal is to make it easy to write a simple unit file that exposes a Unix socket to the internet securely

cool-service-ssl.socket

```
[Unit]
Description=Cool Service Internet Proxy

[Socket]
ListenStream=1234
```

cool-service-ssl.service

```
[Unit]
Description=Proxy Cool Service to the Internet

[Service]
Type=simple
ExecStart=/usr/bin/gocat -key <path to key> -cert <path to cert pem> /var/run/cool-service/service.socket
```

## Goals

- Simple transparent SSL proxy
- Socket activated by default using systemd's socket activation protocol
- Support for SSL client certificates
- Simple command line interface
