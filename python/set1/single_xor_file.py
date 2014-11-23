#! /usr/bin/env python3

import sys, operator

#LETTERS = 'ABCDEFGHIJKLMNOPQRSTUVWXYZ 1234567890\',.?!:;-"$#=*<>/'
#ETAOIN =  'ETAOIN SHRDLCUMWFGYPBVKJXQZ1234567890\',.?!:;-"$#=*<>/'
ETAOIN = """ etaoinsrhldcumfgpyw\nb,.vk-"_'x)(;0j1q=2:z/*!?$35>{}49[]867\+|&<%@#^`~"""
ETAOIN = ETAOIN.upper()

min_s = 9999
min_k = ''
min_m = ''

def xor(a, b):
    x = b''
    for n in range(len(a)):
        x += bytes([a[n] ^ b[n]])
    return x

def letterCount(msg):
    count = {}
    for l in ETAOIN:
        count[l] = 0

    for l in msg.upper():
        if l in ETAOIN:
            count[l] += 1

    return count

def getScore(msg):
    count = letterCount(msg)
    sorted_count = sorted(count.items(), key=operator.itemgetter(1), reverse=True)
    sorted_vals = [x[0] for x in sorted_count]
    s = 0
    for p in range(len(ETAOIN)):
        s += abs(p - sorted_vals.index(ETAOIN[p]))
        #print(p, sorted_vals.index(ETAOIN[p]))
    return s
        
def getBestScore(msg):
    scores = {}
    #print(">>", msg)
    for n in range(128):        
        xored = xor(bytes.fromhex(msg), (chr(n) * len(msg)).encode())
        try:
            s = getScore(xored.decode('ascii'))
            scores[n] = s
            #print(sorted_count)
        except UnicodeError:
            pass
    if len(scores) == 0:
        return 0, 999
    k, s = min(scores.items(), key=operator.itemgetter(1))
    k = chr(k)
    #print(xor(bytes.fromhex(line), (k * len(line)).encode()).decode('utf-8'))
    return k, s
            
with open(sys.argv[1], 'r') as f:
    for line in f:
        if line[-1] == '\n':
            line = line[:-1]
        k, s = getBestScore(line)
        print(s)
        if s < min_s:
            min_s = s
            min_k = k
            min_m = line

print(min_m, min_k)
msg = xor(bytes.fromhex(min_m), (min_k * len(min_m)).encode())
print(msg)
print(msg.decode('utf-8'))

#k = chr(min(scores.items(), key=operator.itemgetter(1))[0])
#print(xor(bytes.fromhex(line), (k * len(line)).encode()).decode('utf-8'))
