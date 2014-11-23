#! /usr/bin/env python3

import sys, operator

LETTERS = 'ABCDEFGHIJKLMNOPQRSTUVWXYZ 1234567890\',.?!-"$#=*<>/'
ETAOIN = 'ETAOIN SHRDLCUMWFGYPBVKJXQZ1234567890\',.?!"$#=*<>/'

def xor(a, b):
    x = b''
    for n in range(len(a)):
        x += bytes([a[n] ^ b[n]])
    return x

def letterCount(msg):
    count = {}
    for l in LETTERS:
        count[l] = 0

    for l in msg.upper():
        if l in LETTERS:
            count[l] += 1

    return count
    
line = sys.stdin.readline()[:-1]

scores = {}
for n in range(128):
    xored = xor(bytes.fromhex(line), (chr(n) * len(line)).encode())
    try:
        count = letterCount(xored.decode('utf-8'))
        sorted_count = sorted(count.items(), key=operator.itemgetter(1), reverse=True)
        sorted_vals = [x[0] for x in sorted_count]
        s = 0
        for p in range(len(ETAOIN)):
            s += abs(p - sorted_vals.index(ETAOIN[p]))
            #print(p, sorted_vals.index(ETAOIN[p]))

        scores[n] = s
        #print(sorted_count)
    except UnicodeError:
        pass

k = chr(min(scores.items(), key=operator.itemgetter(1))[0])
print(xor(bytes.fromhex(line), (k * len(line)).encode()).decode('utf-8'))
    
