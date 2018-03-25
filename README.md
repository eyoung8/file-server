# file-server
Simple file server intended for local network use

Create a local `files` directory adjacent to your `file-server` binary.

Put anything you want in this `files` directory, images, text files, other directories, etc.

Run `./file-server` and access the file server on your local network via the URI `<machine ip>:port`.
`port` is overrideable using the `-p` input flag like `./file-server -p 3000` but defaults to `8080`.

To make accessing the file server easier, it is recommended that you set up forwarding in your
router settings and run the server on a machine that is always up. 
