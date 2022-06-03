<h1 id="goprom-top" align="center">GoProm</h1>

<div align="center">
    <a href="#"><img src="logo.png" width="120px" alt="Hive"></a>
    <p align="center">
        Promotions API
        <br />
        <br />
        <a href="docs/adrs"><strong>ADRs</strong></a> · 
        <a href="https://github.com/vitorsalgado/goprom/actions/workflows/ci.yml"><strong>CI</strong></a> 
    </p>
    <div>
      <a href="https://github.com/vitorsalgado/goprom/actions/workflows/ci.yml">
        <img src="https://github.com/vitorsalgado/goprom/actions/workflows/ci.yml/badge.svg" alt="CI Status" />
      </a>
      <a href="#">
        <img src="https://img.shields.io/badge/go-1.18-blue" alt="Go 1.18" />
      </a>
      <a href="https://conventionalcommits.org">
        <img src="https://img.shields.io/badge/Conventional%20Commits-1.0.0-yellow.svg" alt="Conventional Commits"/>
      </a>
    </div>
</div>

## Overview

GoProm periodically load promotions into Redis and expose the data via an HTTP API. It's composed by two applications
that run separately: **api** and **loader**.  
The `api` simple expose an endpoint to query promotions its identifier.  
The `loader` reads a **csv** file in a specific path, by default **data**, and sends all data to Redis. It's triggered
by a `cron`.
The `loader`workflow is the following:

- a `crontab` executes the `loader` every 1 minute.
- the `loader` checks if there's any new promotions.csv file in a specified directory.
- no promotions file, it will exit.
- found a promotions file, now it'll read the promotions and generate a file with Redis commands for every entry.
- after generating the commands file with every promotion and their respective expiration, it will send the commands to
  Redis CLI and the CLI will finally send all data to the server.

## Getting Started

### Prerequisites

- Go 1.18
- Docker
- Docker Compose

### Configuration

This project uses environment variables for configuration. See [.env.sample](.env.sample) for more details.

### Running

To execute a local environment with both the `loader`, the `api` and `Redis`, execute:

```
make up
```

### Processing Promotions

With the application running, in order to process new promotions, add the **promotions.csv** file to the directory
**data** at the project root dir. After some time, the **loader** application will identify this file and start
the process automatically.

## Development

Check the [Makefile](Makefile) for more details.

### Tools

- Node.js
- Air
- Adr Tools
- Staticcheck
- Husky
- Commit Lint

## Tests

To execute all unit tests, run:

```
make test
```

## Built With

Main libraries and tools used in this project:

- Docker
- Docker Compose
- Redis
- [godotenv](https://github.com/joho/godotenv)
- [Zerolog](https://github.com/rs/zerolog)
- [go-env](https://github.com/Netflix/go-env)
- [go-redis](https://github.com/go-redis/redis)

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md) for more details.

---

## Considerations and Thoughts

- To improve the CSV processing performance, a group of goroutines are used, instead of just one in the previous
  implementation. Each goroutine writes in its own bulk load commands file and several loads are being done in parallel
  now. The performance improvement wasn't so big in the end. Need more thoughts on this.
- Promotions will be removed by Redis using the Expiration feature. For every promotion imported, a configurable
  expiration time is set as well.
- For an ideal solution, the application should have:
    - centralized log solution
    - monitoring solution
    - metrics
    - deployment should be made using solutions like Kubernetes, which facilitates configuration and scaling when needed
    - ideally, the Redis should be operated as a cluster to horizontally scale it and increase availability

<p align="center"><a href="#goprom-top">back to top</a></p>
