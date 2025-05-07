Bitget Trade Collector
Bitget Trade Collector is a Go-based application that collects real-time trade data from the Bitget cryptocurrency exchange via WebSocket and stores it in SQLite databases. Each trading pair (e.g., BTCUSDT, ETHUSDT) runs as an independent systemd service instance, ensuring fault isolation and scalability. The project is packaged as a DEB package for easy deployment on Debian/Ubuntu systems.
Features

Collects trade data (trade ID, symbol, price, size, side, timestamp) for specified trading pairs.
Stores data in per-pair SQLite databases (e.g., trades_BTCUSDT.db) using Write-Ahead Logging (WAL) for non-blocking access.
Supports multiple trading pairs via systemd template services (collector@<pair>.service).
Includes optional debug logging for troubleshooting.
Automated user and directory setup via DEB package installation scripts.
Automatically starts a default BTCUSDT collector service upon installation.

Prerequisites

Go 1.18 or later (for building from source).
Debian/Ubuntu system for DEB package installation.
sqlite3 for querying the database.
dpkg for building and installing the DEB package.
Internet access to connect to the Bitget WebSocket API.

Installation
Building from Source

Clone the repository:git clone https://github.com/magf/bitget-collector.git
cd bitget-collector


Initialize Go modules and fetch dependencies:go mod init bitget-collector
go get github.com/gorilla/websocket
go get github.com/mattn/go-sqlite3


Build the binary:make build



Building DEB Package

Ensure dpkg-deb is installed:
sudo apt install dpkg


Build the DEB package:
make deb


Install the DEB package:
sudo dpkg -i bitget-collector.deb

The installation process will:

Create the bitget system user.
Set up the /var/lib/bitget-collector directory with appropriate permissions.
Install the systemd template service.
Automatically enable and start the collector@BTCUSDT service.



Running the Service

Start a service instance for a specific trading pair (e.g., ETHUSDT):sudo systemctl start collector@ETHUSDT
sudo systemctl enable collector@ETHUSDT


Check the service status:sudo systemctl status collector@ETHUSDT


View collected data:sqlite3 /var/lib/bitget-collector/trades_ETHUSDT.db "SELECT * FROM trades LIMIT 10;"



Usage

Run manually for testing with debug output:./collector -pair=BTCUSDT -debug


Use make run to run with the default pair (BTCUSDT): make run


Data is stored in /var/lib/bitget-collector/trades_<pair>.db.

Project Structure
bitget-collector/
├── cmd/
│   └── collector/        # Application entry point
├── pkg/                  # Placeholder for future Go packages
├── debian/               # DEB package configuration (control, systemd service, install scripts)
├── scripts/              # Build scripts for DEB packaging
├── LICENSE               # MIT License
├── README.md             # English documentation
├── README-ru.md          # Russian documentation
├── Makefile              # Build automation
├── go.mod                # Go module dependencies
├── go.sum                # Go dependency checksums
└── .gitignore            # Git ignore rules

Contributing
Contributions are welcome! Please submit issues or pull requests to the repository. Ensure any changes are tested and follow the project's coding style.
License
This project is licensed under the MIT License. See the LICENSE file for details.
Author
Maxim Gajdaj maxim.gajdaj@gmail.com
