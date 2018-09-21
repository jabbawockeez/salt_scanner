#!/usr/bin/python

import json
import socket
import salt.client
import datetime

target = "127.0.0.1"
port = 6463


def send(minions):
    client = socket.socket(socket.AF_INET, socket.SOCK_STREAM)

    client.connect((target, port))

    hostname = socket.gethostname()
    
    try:
        client.sendall(json.dumps({hostname: minions}))
    except:
        raise
    finally:
        client.close()


def get_minions():
    c = salt.client.LocalClient()

    result = c.cmd_iter("*", "test.ping", timeout = 1)

    minions = []

    for m in result:
        minions.extend(m.keys())
        print minions

    return minions


if __name__ == "__main__":
    minions = get_minions()
    #send(minions)
