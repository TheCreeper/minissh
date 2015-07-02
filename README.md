# MiniSSH

MiniSSH is an SSH server written in Go for portability and shares many of the goals of TinySSH.

## Features

- easy to audit.
- simple configuration.
- reusing code - MiniSSH is reusing libraries many of the Go core libraries and no third party code is included.
- limited amount of features - MiniSSH doesn't have features such: SSH1 protocol, compression, scp, sftp, ...
- no older cryptographic primitives - rsa, dsa, classic diffie-hellman, hmac-md5, hmac-sha1, 3des, arcfour, ...
- free software ([see the license](LICENSE)).
