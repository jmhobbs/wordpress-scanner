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

# Binary Encoding

I wrote a custom binary encoding of the Scan struct for storage and wire xfer.  A scan of bbpress 2.3 (PHP files only) compares as such:

| Bytes   | JSON  | Binary |
|---------|-------|--------|
| Plain   | 11684 | 7973   |
| gzipped | 2153  | 1985   |

You don't gain much after gzip, but it's still interesting, and decoding should be faster than JSON.

If we move to a prefix tree, I think we could easily go even smaller.
