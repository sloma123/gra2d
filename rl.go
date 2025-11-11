package main

import (
	"math"
	"math/rand"
)

type simState struct {
	x, y      float64
	vy        float64
	isJumping bool
}

func resetEpisode() simState {
	return simState{x: 0, y: groundY}
}

func stepSim(s simState, action int) (simState, float64, bool) {
	speed := 0.0
	switch action {
	case 0:
		speed = -stepSpeed
	case 1:
		speed = stepSpeed
	case 2:
		if !s.isJumping {
			s.vy = jumpImpulse
			s.isJumping = true
		}
	case 3:
		speed = 0
	}

	oldDist := distance(s.x, coinX)

	s.x += speed
	s.vy += gravity
	s.y += s.vy

	if s.y >= groundY {
		s.y = groundY
		s.vy = 0
		s.isJumping = false
	}

	if s.x < 0 {
		s.x = 0
	}
	if s.x > float64(screenW-playerW) {
		s.x = float64(screenW - playerW)
	}

	newDist := distance(s.x, coinX)
	r := (oldDist - newDist) / 10.0
	r -= 0.01

	// PRZEGRANA
	if checkCollision(s.x, s.y, playerW, playerH, obstacleX, obstacleY, 30, 30) {
		r -= 100
		return s, r, true
	}

	// WYGRANA
	if checkCollision(s.x, s.y, playerW, playerH, coinX, coinY, 30, 30) {
		r += 100
		return s, r, true
	}

	return s, r, false
}

func (a *Agent) Train(episodes int, maxSteps int) {
	for ep := 0; ep < episodes; ep++ {
		sim := resetEpisode()
		s := a.stateKey(sim.x, sim.isJumping)
		a.epsilon = 0.2 + 0.8*math.Exp(-float64(ep)/float64(episodes/4))
		for t := 0; t < maxSteps; t++ {
			a.ensureState(s)
			aIdx := a.chooseAction(s)
			next, reward, done := stepSim(sim, aIdx)
			sNext := a.stateKey(next.x, next.isJumping)
			a.ensureState(sNext)
			_, bestNext := a.bestAction(sNext)
			q := a.Q[s]
			td := reward + a.gamma*bestNext - q[aIdx]
			q[aIdx] += a.alpha * td
			a.Q[s] = q
			sim = next
			s = sNext
			if done {
				break
			}
		}
	}
}

func (a *Agent) chooseAction(s int) int {
	if rand.Float64() < a.epsilon {
		return rand.Intn(4)
	}
	ba, _ := a.bestAction(s)
	return ba
}

func (a *Agent) stateKey(x float64, isJumping bool) int {
	xb := int(x / a.xBinSize)
	xb = limit(xb, 0, a.xBins-1)
	j := 0
	if isJumping {
		j = 1
	}
	return xb*2 + j
}

func (a *Agent) ensureState(s int) {
	if _, ok := a.Q[s]; !ok {
		a.Q[s] = [4]float64{}
	}
}

func (a *Agent) bestAction(s int) (int, float64) {
	a.ensureState(s)
	q := a.Q[s]
	bestA, bestQ := 0, q[0]
	for i := 1; i < 4; i++ {
		if q[i] > bestQ {
			bestQ = q[i]
			bestA = i
		}
	}
	return bestA, bestQ
}
