proxy
-----

A proxy server with a per route override mechanism.

# Configuration

See overrides.json.example for example of config file.

Create your own and use `-c` to point to it. See Usage for example.

# Usage

`proxy -port 8001 -c overrides.json`

Default port is 8000
There is no default config file. It is required.

By default `proxy` accepts only localhost connections.

Use `-a` to accept *all* requested connections.
