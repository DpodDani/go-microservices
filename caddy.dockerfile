FROM caddy:2.4.6-alpine

# basically specifying our own version of a Caddy image
COPY Caddyfile /etc/caddy/Caddyfile
