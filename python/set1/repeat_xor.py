#! /usr/bin/env python3

import sys, binascii

def xor(a, b):
    x = b''
    for n in range(len(a)):
        x += bytes([a[n] ^ b[n]])
    return x

def rep_xor(data, key):
    rep_key = (key * int(len(data) / len(key) + 1)).encode()
    return xor(data, rep_key)

key = sys.argv[1]
data = sys.stdin.read()
xored = rep_xor(data.encode(), key)
print(binascii.hexlify(xored).decode('utf-8'))
