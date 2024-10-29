package wallbox

var wallboxStatusCodes = []string{
	"Ready",
	"Charging",
	"Connected waiting car",
	"Connected waiting schedule",
	"Paused",
	"Schedule end",
	"Locked",
	"Error",
	"Connected waiting current assignation",
	"Unconfigured power sharing",
	"Queue by power boost",
	"Discharging",
	"Connected waiting admin auth for mid",
	"Connected mid safety margin exceeded",
	"OCPP unavailable",
	"OCPP charge finishing",
	"OCPP reserved",
	"Updating",
	"Queue by eco smart",
}

var stateOverrides = map[int]int{
	0xA1: 0,
	0xA2: 9,
	0xA3: 14,
	0xA4: 15,
	0xA6: 17,
	0xB1: 3,
	0xB2: 4,
	0xB3: 3,
	0xB4: 2,
	0xB5: 2,
	0xB6: 4,
	0xB7: 8,
	0xB8: 8,
	0xB9: 10,
	0xBA: 10,
	0xBB: 12,
	0xBC: 13,
	0xBD: 18,
	0xC1: 1,
	0xC2: 1,
	0xC3: 11,
	0xC4: 11,
	0xD1: 6,
	0xD2: 6,
}

var stateMachineStates = map[int]string{
	0xE:  "Error",
	0xF:  "Unviable",
	0xA1: "Ready",
	0xA2: "PS Unconfig",
	0xA3: "Unavailable",
	0xA4: "Finish",
	0xA5: "Reserved",
	0xA6: "Updating",
	0xB1: "Connected 1", // Make new session?
	0xB2: "Connected 2",
	0xB3: "Connected 3", // Waiting schedule ?
	0xB4: "Connected 4",
	0xB5: "Connected 5", // Connected waiting car ?
	0xB6: "Connected 6", // Paused
	0xB7: "Waiting 1",
	0xB8: "Waiting 2",
	0xB9: "Waiting 3",
	0xBA: "Waiting 4",
	0xBB: "Mid 1",
	0xBC: "Mid 2",
	0xBD: "Waiting eco power",
	0xC1: "Charging 1",
	0xC2: "Charging 2",
	0xC3: "Discharging 1",
	0xC4: "Discharging 2",
	0xD1: "Lock",
	0xD2: "Wait Unlock",
}

type ControlPilotState struct {
	Status  string
	Detail  string
	Message string
}

var ControlPilotStates = map[int]ControlPilotState{
	0xE: {
		Status:  "E",
		Detail:  "E",
		Message: "Error",
	},
	0xF: {
		Status:  "F",
		Detail:  "F",
		Message: "Failure",
	},
	0xA1: {
		Status:  "A",
		Detail:  "A1",
		Message: "Ready (S1 at 12V, car not connected)", // S1 at 12V, car not connected
	},
	0xA2: {
		Status:  "A",
		Detail:  "A2",
		Message: "Ready",
	},
	0xB1: {
		Status:  "B",
		Detail:  "B1",
		Message: "Connected (S1 at 9V, car connected not allowed charge)", // S1 at 9V, car connected not allowed charge
	},
	0xB2: {
		Status:  "B",
		Detail:  "B2",
		Message: "Connected (S1 at Oscillator, car connected allowed charge)", // S1 at Oscillator, car connected allowed charge
	},
	0xC1: {
		Status:  "C",
		Detail:  "C1",
		Message: "Charging",
	},
	0xC2: {
		Status:  "C",
		Detail:  "C2",
		Message: "Charging (S2 closed))", // S2 closed
	},
}
