# 图书<<自己动手写docker>> 源码 喻成改编版

This project is forked from https://github.com/xianlubird/mydocker which is a educational
container engine for book "write Docker from scratch". This book and project provides guidence
to readers about how a build a simplified "Docker Engine".


Cheng Yu (s3341458, chengyu0316@gmail.com) modified it on each chapter branch on
{{ original branch }}-cheng for following reasons:

1. original code has bugs and obsolete code due to library updates, so go build will fail
2. original code has wrong imports
3. original code works on a old Ubuntu 14.04 linux kernel 3.13, I am using archlinux with much newer kernel
4. original code is not ideal in terms of the code of itself (but at OK level)

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

# for the code of Chapter 5.5 and whats next
almost the same as the original code, just fixed some tiny import issues and function name issue

# for the code of Chapter 5.6 and whats next
almost the same as the original code, add a clean command to remove all the containers related files for testing reason

# for the code of Chapter 5.7 and whats next
almost the same as the original code, however since the implementation of this commit is too far away from the real docker so really did not get motivated.

In real docker world commit are implemented by copy files of container writable layer to image read-only layers with proper re-indexing.
To implemented things to that really will take too much time.

# for the code of Chapter 5.8 and whats next
almost the same as the original code, during the time test this stage I found issues in command line arg parsing(when environment and container name, it will mistakely take the later one as the image name), however such thing is not the main point of this project(understand docker working mechnism). SoI do not want spend extra time for this problem for now.

# for the code of Chapter 6.5 (last chapter)
Almost the same as the original code, fix some bugs manage the code in a slightly more elegant way.
All networks and ipam information will be put into networks folder in mydocker directory
Due to a command args parsing bug, you need to specify image by --image and command by --command, this is different with how Docker looks in real world.

# P.S
branch {{branch}}-cheng is my modifications on {{branch}} which makes it runable on my machine

sometimes the code in this book can not properly distinguish the whether the args are more  mydocker or for the wrapped command.
better use it in sudo ./mydocker -args argv "command command-args"
