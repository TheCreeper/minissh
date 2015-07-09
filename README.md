# MiniSSH

MiniSSH is an SSH server written in Go.

## Limitations

- os/exec has no facility to change the UID/GID of a child process.
- os/user has no facility to find the groups a user is apart of with the exception of the primary group.
- os/user has no facility to find the default shell of a user.
