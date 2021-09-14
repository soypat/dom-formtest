package main

import "gonum.org/v1/gonum/spatial/r3"

// Parameters contains some dynamic values that
// may be used to tune the control system.
type Parameters struct {
	HW Hardware
	// Actuator values
	ActuatorZeros, Gains []float64
	// IMU noise variance
	VarAccel, VarGyro   r3.Vec
	ConvGyro, ConvAccel [3]linearConv

	// Sampling period [s]
	Ts float64
	// Complementary filter configuration
	FilterGyro, FilterAccel filterconfig

	// Smoothing parameters

	// Smoothing/estimation trust in system, must be contained in [0,1)
	SysTrust      float64 // Smoothing/estimation trust in system, must be contained in [0,1)
	SmoothSamples int
}

type linearConv struct {
	// slope and bias (offset)
	Slope float64
	Bias  int
}

func (l *linearConv) Convert(v int) float64 { return l.Slope * float64(v+l.Bias) }

// Values pertaining to construction of rocket.
type Hardware struct {
	// IMU position with respect to the CoG
	IMUPosition r3.Vec
	IMUAttitude r3.Vec
	// Max gimbal actuation angle measured from vertical [rad]
	GimbalMaxAngle r3.Vec
}

type filterconfig struct {
	Cutoff    float64
	Bandwidth float64
}
