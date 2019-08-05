# bcast

listen for N connections then broadcast stdin to them all.

# Examples

Watch a movie with someone

```bash
# Server (listen)
tee < ~/Movie.mp4 >(bcast -l :1337) | mpv -

# Client (connect)
nc 127.0.0.1 1337 | mpv -
```
