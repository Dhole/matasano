#! /usr/bin/env python3 

import sys, base64
from Crypto.Cipher import AES

def read_b64file(filename):
    with open (filename, 'r') as f:
        data_b64 = f.read()
        data = base64.b64decode(data_b64)
    return data

filename = sys.argv[1]
key = sys.argv[2]
data = read_b64file(filename)
cipher = AES.new(key, AES.MODE_ECB)
plaintext = cipher.decrypt(data)
print(plaintext.decode('ASCII'))
