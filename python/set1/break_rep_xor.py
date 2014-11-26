#! /usr/bin/env python3

import sys, binascii, base64, operator

LETTERS = 'ABCDEFGHIJKLMNOPQRSTUVWXYZ 1234567890\',.?!:;-"$#=*<>/'
#ETAOIN =  'ETAOIN SHRDLCUMWFGYPBVKJXQZ1234567890\',.?!:;-"$#=*<>/'
ETAOIN = """ etaoinsrhldcumfgpyw\nb,.vk-"_'x)(;0j1q=2:z/*!?$35>{}49[]867\+|&<%@#^`~"""
ETAOIN = ETAOIN.upper()

def xor(a, b):
    """Xor two byte strings together"""
    x = b''
    for n in range(len(a)):
        x += bytes([a[n] ^ b[n]])
    return x

def rep_xor(data, key):
    """Xor a byte string with a repetition of key"""
    rep_key = (key * int(len(data) / len(key) + 1)).encode()
    return xor(data, rep_key)

def ham_dist(a, b):
    """Number of differing bits between byte strings"""
    xored = xor(a, b)
    count = 0
    for n in xored:
        count += bin(n).count("1")
    return count

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
        xored = xor(msg, (chr(n) * len(msg)).encode())
        try:
            s = getScore(xored.decode('utf-8'))
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

#a = "this is a test"
#b = "wokka wokka!!!"
#print(ham_dist(a.encode(), b.encode()))

MAX_KEYSIZE = 40
MIN_KEYSIZE = 2

avg_n = 4
len_tries = 4

with open (sys.argv[1], 'r') as f:
    data_b64 = f.read()
    data = base64.b64decode(data_b64)
    dists = {}
    for ks in range(MIN_KEYSIZE, MAX_KEYSIZE + 1):
        dist = 0
        #print("ks: ", ks)
        for p in range(avg_n):
            first = data[ks * p:ks * p + ks]
            second = data[ks * p + ks:ks * p + ks * 2]
            #print(ham_dist(first, second) / ks)
            dist += ham_dist(first, second)
        dist = dist / (avg_n * ks)
        dists[ks] = dist
        print("mean: ", dist)
    #print(dists)
    sorted_dists = sorted(dists.items(), key=operator.itemgetter(1))
    #print(sorted_dists)
    min_s = 999
    min_k = ""
    for n in range(len_tries):
        ks = sorted_dists[n][0]
        blocks = {}
        for q in range(ks):
            blocks[q] = b''
        for p in range(int(len(data) / ks)):
            d = data[ks * p: ks * (p + 1)]
            #print("chunk: ", binascii.hexlify(d))
            for q in range(ks):
                blocks[q] += bytes([d[q]])
            #print(blocks)
        #print(blocks[0][:16])
        key = ""
        score_avg = 0
        for q in range(ks):
            k, s = getBestScore(blocks[q])
            key += k
            score_avg += s
        score_avg /= ks
        print(score_avg, key)
        if score_avg < min_s:
            min_s = score_avg
            min_k = key
print("KEY:", min_k)
#min_k = "Terminator X: Bring the noise"
print(rep_xor(data, min_k).decode('utf-8'))


