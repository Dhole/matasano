#!/usr/bin/env python3

import sys, base64

def hex2base64(hex_line):
    data = bytes.fromhex(hex_line)
    return base64.b64encode(data)

for line in sys.stdin:
    print(hex2base64(line[:-1]).decode('utf-8'))
