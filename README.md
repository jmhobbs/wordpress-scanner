[![Build Status](https://travis-ci.org/jmhobbs/wordpress-scanner.svg?branch=master)](https://travis-ci.org/jmhobbs/wordpress-scanner) [![codecov](https://codecov.io/gh/jmhobbs/wordpress-scanner/branch/master/graph/badge.svg)](https://codecov.io/gh/jmhobbs/wordpress-scanner)

This is an experimental server which downloads plugins from WordPress.org on demand, and hashes their contents.

The idea is that a client could check the hashes against their existing files to quickly check if the plugin has been hacked or otherwise corrupted.

# Endpoints

  * `GET /plugin/{name}/{version}` - Get hashes for a plugin from wordpress.org
  * `POST /plugin/{name}/{version}/diff` - Compare a client hash against a wordpress.org hash (Not Implemented)
  * `GET /plugin` - List of plugins we have hashed versions of
  * `GET /plugin/{name}` - List of versions we have hashed

# Binary Encoding

I wrote a custom binary encoding of the Scan struct for storage and wire xfer.  A scan of bbpress 2.3 (PHP files only) compares as such:

| Bytes   | JSON  | Binary |
|---------|-------|--------|
| Plain   | 11684 | 7973   |
| gzipped | 2153  | 1985   |

You don't gain much after gzip, but it's still interesting, and decoding should be faster than JSON.

If we move to a prefix tree, I think we could easily go even smaller.

# Ideas

  * https://github.com/d4l3k/messagediff - Diff the output of client/server
  * Optionally use protobufs
  * Make sure gzip is on
  * HTTP/2?
  * TLS & Auth
  * Use a tree structure (radix tree?) and binary encoding for xfer
  * Hash at the directory level (sorted filenames + hashes)

