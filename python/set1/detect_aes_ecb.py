#! /usr/bin/env python3 

import sys, base64

def unique_set(seq):
    seen = {}
    pos = 0
    for item in seq:
        if item not in seen:
            seen[item] = True
        else:
            return False
    return True

c_size = 16 * 2
    
filename = sys.argv[1]
with open (filename, 'r') as f:
    for line in f:
        if line[-1] == '\n':
            line = line[:-1]
        blocks = []
        for n in range(int(len(line) / c_size)):
            blocks += [line[c_size * n: c_size * (n + 1)]]
        #print(blocks)
        #sys.exit(0)
        if not unique_set(blocks):
            print(line)
            
