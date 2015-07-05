# MiniSSH

MiniSSH is an SSH server written in pure Go.

## Limitations

- Go lacks the ability to change the UID/GID of a child process.
- os/user has no facility to find the groups a user is apart of with the exception of the primary group.
