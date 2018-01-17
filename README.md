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

# P.S
branch {{branch}}-cheng is my modifications on {{branch}} which makes it runable on my machine


sometimes the code in this book can not properly distinguish the whether the args are more  mydocker or for the wrapped command.

better use it in sudo ./mydocker -args argv "command command-args"
