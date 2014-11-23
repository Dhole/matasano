#! /usr/bin/env python3

import sys, base64
from Crypto.Cipher import AES

blk_len = 16

def xor(a, b):
    x = b''
    for n in range(len(a)):
        x += bytes([a[n] ^ b[n]])
    return x

def read_b64file(filename):
    with open (filename, 'r') as f:
        data_b64 = f.read()
        data = base64.b64decode(data_b64)
    return data

def aes_cbc_encrypt(iv, data, key):
    last = iv
    cipher = AES.new(key, AES.MODE_ECB)
    c_data = b''
    for n in range(int(len(data) / blk_len)):
        blk = data[blk_len * n: blk_len * (n + 1)]
        xored = xor(last, blk)
        c_data += cipher.encrypt(xored)
        last = c_data
    return c_data

def aes_cbc_decrypt(iv, data, key):
    cipher = AES.new(key, AES.MODE_ECB)
    prev = iv
    d_data = b''
    for n in range(int(len(data) / blk_len)):
        blk = data[blk_len * n: blk_len * (n + 1)]
        xored = cipher.decrypt(blk)
        d_data += xor(xored, prev)
        prev = blk
    return d_data
    

data = read_b64file(sys.argv[1])

key = "YELLOW SUBMARINE".encode()
iv = ("\0" * 16).encode()
d_data = aes_cbc_decrypt(iv, data, key)
print(d_data.decode('ASCII'))
