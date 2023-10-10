#!/bin/sh
make rebuild build-linux-amd64 && scp build/nsxctl_linux_amd64 lin01:~ && echo 'VMware1!' | ssh lin01 sudo -S mv nsxctl_linux_amd64 /usr/local/bin/nsxctl
