#! /usr/bin/python3

#import urllib.request, sys, urllib.urlencode
import sys, time, threading
from httplib2 import Http
from random import shuffle
import socket, queue
import numpy

q = queue.Queue()

alphabet = ""
alphabet = "abcdefghijklmnopqrstuvwxyz"
alphabet += alphabet.upper()
alphabet += "0123456789"
alphabet = list(alphabet)
N_THREADS = 10

key_base = ""

timing = {}
threads = {}

input = sys.stdin.readline()[:-1]

url = "http://54.83.207.90:4242/?input=" + input
args = "/?input=" + input
ip = "54.83.207.90"
port = 4242

def test_key(key):
    h = Http()
    k = key_base + key
    #print("testing key:", key)
    headers = {'Connection': 'close'}
    data = "input=" + input + "&key=" + k + "&submit=Submit+Query"
    start = time.time()
    resp, content = h.request(url, "POST", data, headers=headers)
    end = time.time()
    delta = end - start
    timing[k[-1]] += delta
    #print("finishing:", key)

def test_key2():
    while True:
        n = q.get()
        k = "A"*n
        #k = key_base + key
        #if len(k) < 4:
        #    k += "AAA"
        s = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
        s.connect((ip, port))
        data = "input=" + input + "&key=" + k + "&submit=Submit+Query"
        msg = "POST "+args+" HTTP/1.1\r\n" +\
              "Content-Length: " + str(len(data)) + "\r\n\r\n" +\
              data
        
        msg = msg.encode("ASCII")
        #print(msg)
        #MSGLEN = len(msg)
        #totalsent = 0
        #while totalsent < MSGLEN:
        #    sent = s.send(msg[totalsent:])
        #    if sent == 0:
        #        raise RuntimeError("socket connection broken")
        #    totalsent += sent
        
        #while True:
        #    chunk = s.recv(1500)
        #    if chunk == b'':
        #        raise RuntimeError("socket connection broken")
        #    msg += chunk
        start = time.time()
        sent = s.send(msg)
        msg = s.recv(1500)
        msg += s.recv(1500)
        end = time.time()
        delta = end - start
        timing[n].append(delta)
        s.close()

        q.task_done()
        #print(msg.decode("ASCII"))
    
    
    
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
        
for n in range(0, 8): #16
    for i in range(32):
        timing[i] = [] 
    for i in range(0, 128): #10
        max_delta = 0
        char = ''
        shuffle(alphabet)
        for i in range(32):
            q.put(i)
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
        mean_timing[k] = numpy.mean(mean_timing[k])

    max_t = 0
    max_k = ''
    for k in mean_timing:
        if mean_timing[k] > max_t:
            max_t = mean_timing[k]
            max_k = k

    k = max_k
    print(k)
    break

    #print(timing)
    key_base += k
    print(key_base)
    if check_key(key_base):
        print(key_base)
        break
