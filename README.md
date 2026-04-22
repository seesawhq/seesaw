# Seesaw: Lightweight feature management system

Seesaw is lightweight feature management system written in Golang with few external dependencies which makes it easier to deploy on prem and maintain.
The whole idea behind it is to start with single server and scale it as needed. By default seesaw stores data in sqlite and uses in-memory cache.
If you need you can configure postgres and redis to make it more scale as per your need.

---

### Development setup

1. clone the code `git clone https://github.com/seesawhq/seesaw.git`
2. cd into directory `cd seesaw`
3. install dependencies `go mod download`
4. Run the service `go tool air`
