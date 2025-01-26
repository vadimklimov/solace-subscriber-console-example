> [!IMPORTANT]
> This repository has been archived and is no longer maintained.

# Solace subscriber console application example

An example Solace subscriber application that was used in the demo during the session _Architecture and Use Cases for Event Mesh_ at SAP Inside Track Netherlands 2023.

## Usage

For the application to function correctly, the following environment variables must be set:

| Environment variable | Description                       | Example               |
| -------------------- | --------------------------------- | --------------------- |
| SOLACE_HOST          | Event broker host                 | tcp://localhost:55555 |
| SOLACE_VPN           | Message VPN name                  | default               |
| SOLACE_USERNAME      | Username for basic authentication | example_user          |
| SOLACE_PASSWORD      | Password for basic authentication | example_password      |
| SOLACE_QUEUE         | Queue name                        | examples              |
