# Deployer

Auto deploy with GitHub webhook

## Possibility

- Run `make build` from directory for configuration `repository + branch`
- Send build status to specific chat in telegram

## Handlers

- `/ping` - response with status OK
- `/hook/github` - webhook for GitHub Webhooks
    - response with status Accepted if deploy has been started
    - response with status NonAcceptable if deploy impossible to start (invalid payload, not push event and another)

## Important

- For deploy will be created only one job, other deploy jobs will be canceled
- If push event accepted but repository or branch not configured - job will be skipped

## Requirements

- `Go` v1.21+
- `make` in system

## Installation

1. Clone this repository `git clone https://github.com/toppi-me/deployer.git`
2. Create .env file `cp .env.example .env`
3. Create config file `cp .config.json.example .config.json`
4. Setup `.env` and `config.json` files with your data
5. Run `go mod download` and `go build .`
6. Run `chmod +x deployer`

## Create systemd

- `nano` or `vim` `/etc/systemd/system/deployer.service`

Example of `deployer.service` file
```ini
[Unit]
Description=Deployer
After=network.target

[Service]
User=root
Group=root
ExecStart=/home/root/deployer/deployer
Restart=always

[Install]
WantedBy=multi-user.target
```

- Need to change `User`, `Group` and path to binary `ExecStart`
- Start service `sudo systemctl start deployer.service`

## Contributions:
Any contributions are welcome, also if there are problems in the process, then create an issue
