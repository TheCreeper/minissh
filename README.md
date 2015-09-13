# MiniSSH

This repository contains the small and portable MiniSSH server. This SSH server
is written in Go and attempts to be usable on most platforms that the Go
compiler supports (Windows, Linux, FreeBSD).

## Features

- Highly portable. Should run on Windows/Linux/BSD.
- Use only modern and reliable encryption methods (RSA/ECDSA).
- Uses only SSH keys. No passwords.
- Uses only the Go standard library.

## Limitations of the Standard Library

- Package os/exec has no facility to change the UID/GID of a child process.
- Package os/user has no facility to find the default shell of a user.