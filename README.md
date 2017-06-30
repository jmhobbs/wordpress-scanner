# Ideas

  * https://github.com/d4l3k/messagediff - Diff the output of client/server
  * Optionally use protobufs
  * Make sure gzip is on
  * HTTP/2?
  * TLS & Auth
  * Use a tree structure (radix tree?) and binary encoding for xfer
  * Hash at the directory level (sorted filenames + hashes)

# Endpoints

  * `GET /plugin/{name}/{version}` - Get hashes for a plugin from wordpress.org
  * `POST /plugin/{name}/{version}/diff` - Compare a client hash against a wordpress.org hash (Not Implemented)
  * `GET /plugin` - List of plugins we have hashed versions of (Not Implemented)
  * `GET /plugin/{name}` - List if versions we have hashed (Not Implemented)
