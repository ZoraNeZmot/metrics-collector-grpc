# Metrics Collector

A client-server application written in **Go** for collecting metrics from client machines and aggregating them on a central server.

## Architecture

The system consists of two main components:

- **Client** – collects metrics from the host machine. Implemented as a Windows service and a Linux daemon.
- **Server** – receives metrics from clients, processes them, and stores the collected data.

## Technologies

- **Go** – primary programming language.
- **gRPC** – communication between clients and the server.
- **Windows Services / Linux Daemons** – client deployment.

## Configuration

Application configuration is loaded at startup from:
- Environment variables.
- `.ini` configuration file.

## Features

- Cross-platform client support (Windows and Linux).
- Client-server communication using gRPC.
- Centralized metrics collection and aggregation.
- Configurable application settings.