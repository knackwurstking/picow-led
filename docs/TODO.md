# TODO

## Project

- [ ] Think about separating the home handler in to more handlers, one for devices and one for groups?
- [ ] Improve error handling and create a errors package

## Features

- [ ] Change group power on/off to toggle, just like devices
- [ ] OOB render group power buttons
- [ ] Find a good way to show devices overview, show power on/off, color, pins

## Performance

- **Unnecessary Database Operations**: The `services/device-controls.go` file has a method called `GetPins` that calls `control.GetPins` which seems to be a separate function. This could be optimized by caching or pre-fetching pins.

- **Caching**: There's no caching mechanism for frequently accessed data, which could lead to performance issues.

- **Database Connection Management**: The `services/registry.go` file creates a single database connection that's shared across all services. This could become a bottleneck under high load, as all services share the same connection.

### Reddit:

Here are my general tips regarding mattn's driver, which my team has used to build SQLite-backed microservices:

1. Set the journal mode to WAL and synchronous to Normal.

2. Use two connections, one read-only with max open connections set to some large number, and one read-write set to a maximum of 1 open connection.

3. Set the transaction locking mode to IMMEDIATE and use transactions for any multi-query methods.

4. Set the busy timeout to some large value, like 5000. I'm not sure why this is necessary, since I figured the pool size of 1 would obviate the need for this, but it seems necessary (otherwise you can get database is locked errors).

With these few settings, we get good performance for our use case (>2K mid-size writes/sec, 30K reads per second on 2 vCPU and an SSD). I'd also recommend using Litestream to perform WAL shipping to S3.

> [SQLite performance tuning](https://phiresky.github.io/blog/2020/sqlite-performance-tuning/)

Also see here: [Database Connection Bottleneck in Services](./Database_Connection_Bottleneck_in_Services.md)
