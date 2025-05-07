Changelog
All notable changes to the Bitget Trade Collector project will be documented in this file.
The format is based on Keep a Changelog, and this project adheres to Semantic Versioning.
[1.0] - 2025-05-07
Added

Initial release of Bitget Trade Collector.
WebSocket-based trade data collection from Bitget exchange for specified trading pairs.
Storage of trade data (trade ID, symbol, price, size, side, timestamp) in per-pair SQLite databases with Write-Ahead Logging (WAL).
Support for multiple trading pairs via systemd template services (collector@<pair>.service).
Optional debug logging via -debug flag.
DEB package support with preinst and postinst scripts for automated setup:
Creates bitget system user and /var/lib/bitget-collector directory.
Automatically enables and starts collector@BTCUSDT service on installation.


Project structure with Go modules, Makefile, and scripts for building DEB packages.
MIT License and bilingual documentation (README.md, README-ru.md).

Changed

N/A

Fixed

N/A

[Unreleased]
Added

(Placeholder for version 2.0 features, e.g., pair management, status monitoring)

Changed

(Placeholder for future changes)

Fixed

(Placeholder for future fixes)

