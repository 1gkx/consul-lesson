#!/bin/sh

/bin/consul agent -recursor=8.8.8.8 -config-dir=/consul/config
# /usr/app/main