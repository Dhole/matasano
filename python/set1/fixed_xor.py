#! /usr/bin/env python3

import sys

def xor(a, b):
    x = b''
    for n in range(len(a)):
        x += bytes([a[n] ^ b[n]])
    return x

a = sys.stdin.readline()
b = sys.stdin.readline()

print(xor(bytes.fromhex(a[:-1]), bytes.fromhex(b[:-1])).decode('utf-8'))
