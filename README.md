# <<自己动手写docker>> 源码

Cheng Yu (s3341458, chengyu0316@gmail.com) modified it in order to makes it work on my machine

difference with original code base

# for the code of Chapter 3.1 and what's next
I fixed the issue that old code tried to imported the wrong package from github

# for the code of Chapter 3.2 and what's next
Using syscall.Mount("", "/", "", syscall.MS_PRIVATE|syscall.MS_REC, "") makes the mount namespace works properly on my archlinux computer, systemd made "/" mounted as shared by default

# for the code of Chapter 3.3 and what's next
Change the logs functions makes so I can find out why pivot was not working properly

# for the code of Chapter 4.1
It is almost the same as code-3.3-cheng due the issue of code of the book code-3.2 is exactly the same as code-4.1
add Cmd.Dir which gives a static image layer

# for the code of Chapter 4.2 and whats next
Using overlayfs instead of aufs to implement container writable layer and base readble image layer

# for the code of Chapter 4.3 and whats next
Using bind mount instead of aufs to implement docker volume

# for the code of Chapter 4.4 and whats next
the code 4.4 actually implemented a mini version of "docker export"
instead of "docker commit", so I just changed name of command and
functions to "export"

# for the code of Chapter 5.1
exactly the same as original book
I fix the issue that I miss remove the cgroup management in previouse commit

# for the code of Chapter 5.2 and whats next
each docker container will have its own writeable layer based on their names
fixed one import problem
all container information will be stored in a local directory instead in /var/run/mydocker

# for the code of Chapter 5.3 and whats next
all logs will be stored in current directory logs, each container will have a different log file based on its name

fix the issue where logs and info are not removed in delete workspace

# for the code of Chapter 5.4 and whats next
almost the same as the original code, just fixed some tiny import issues

# for the code of Chapter 5.4 and whats next
almost the same as the original code, just fixed some tiny import issues and function name issue

# P.S
branch {{branch}}-cheng is my modifications on {{branch}} which makes it runable on my machine


sometimes the code in this book can not properly distinguish the whether the args are more  mydocker or for the wrapped command.

better use it in sudo ./mydocker -args argv "command command-args"
