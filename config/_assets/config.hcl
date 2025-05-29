node "example" {
  debug = true
  listen "tcp6" {
    addr = "[::1]:80"
  }
  listen "unix" {
    addr = "/run/service/unix.sock"
  }

  reachable_address = "example.com:666"
}

env = [
  "OTEL_EXPORTER_OTLP_ENDPOINT=http://localhost:4317",
]
