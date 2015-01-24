package main

import (
	"./mt19937"
	"fmt"
	"math/rand"
	"time"
)

var fake_time uint32

func waitFakeTime(t uint32) {
	fake_time += t
}

func getFakeTime() uint32 {
	return fake_time
}

func initFakeTime(t uint32) {
	fake_time = t
}

func genRandom() (val uint32) {
	m := mt19937.New()

	waitFakeTime(uint32(rand.Int31n(10000-40) + 40))
	fmt.Println("Real seed:", getFakeTime())
	m.Init(getFakeTime())
	val = m.ExtractNumber()
	waitFakeTime(uint32(rand.Int31n(10000-40) + 40))
	return val
}

func crackSeed(val uint32) (seed uint32) {
	t := getFakeTime()
	fmt.Println("Current time:", t)
	m := mt19937.New()
	for {
		//fmt.Println("Trying for ", t)
		m.Init(t)
		if m.ExtractNumber() == val {
			return t
		}
		t--
	}
}

func main() {
	rand.Seed(time.Now().UnixNano())
	initFakeTime(uint32(time.Now().Unix()))
	val := genRandom()
	fmt.Println("Random value:", val)
	seed := crackSeed(val)
	fmt.Println("Seed found:", seed)
	m := mt19937.New()
	m.Init(seed)
	fmt.Println("Test -> seed:", seed, "generates:", m.ExtractNumber())
}
