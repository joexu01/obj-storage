#!/bin/bash

# 注意：你可能需要手动停止进程

sudo ip addr del 10.29.1.1/16 dev lo
sudo ip addr del 10.29.1.2/16 dev lo
sudo ip addr del 10.29.1.3/16 dev lo
sudo ip addr del 10.29.1.4/16 dev lo
sudo ip addr del 10.29.1.5/16 dev lo
sudo ip addr del 10.29.1.6/16 dev lo
sudo ip addr del 10.29.2.1/16 dev lo
sudo ip addr del 10.29.2.2/16 dev lo

sudo ip addr del 10.29.102.173/16 dev lo