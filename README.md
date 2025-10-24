# tomashevi.ch
## My personal website btw
using `net/http` 1.24.4

---

### Features
1. Infinity procedural generated fishes background: get visitor soul and insert into database with unique seed (uuidv7)
2. Pixelbattle in header word `tomashevich`: every soul can paint 10 pixels

---

### How to Run/Deploy
1. `cp config.example.json config.json` make config file
2. `go run .`/`go build .` => `./tomashevich` start the server

or use docker

1. `cp config.example.json config.json` make config file
2. `docker-compose up --build -d` build docker and start the server

---

### About proxy
Enable `server.is_behind_proxy` in `config.json`, we use `X-Forwarded-For` ONLY if r.RemoteAddr is LOCAL. (dont use cloudflare proxy 4 example)

---

### Embeding files!
directory `./static` and `config.json` building IN binary bc i love embeding btw

---

### Middlewares
1. `cache` => set caching headers for browser
2. `compress` => set and encode content to br/zstd/gzip or none (if unsupported by client)
3. `helheim` => get your soul and generate seed for fish
4. `rate_limiter` => no dosing pls

---

### API
1. `GET /fishes?page=N` => returning fishes seeds
2. `GET /fishes/me` => returning your seed
3. `GET /pixels` => returning all pixels from pixelbattle
4. `POST /pixels:paint` `{"x": int, "y": int, "color": string}` => no content return
5. `POST /pixels:register` `{"pixels": [{"x": int, "y": int}]}` => no content return

---

### Known issue
1. no logs (idc)
2. `/pixels:register` can be poisoned by wrong invalid coords and broke pixelbattle
