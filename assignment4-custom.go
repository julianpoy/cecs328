package main

import (
	"fmt"
	"math"
	"os"
	"bufio"
	"strconv"
	"io"
	"sort"
)

func check(err error) {
    if err != nil {
        panic(err)
    }
}

var f io.Writer;
func openLog() {
	var err error;
    f, err = os.Create("./assignment4.log");
    check(err)
}

func writeLog(log string) {
	w := bufio.NewWriter(f)

    _, err := fmt.Fprintf(w, "%v\n", log)
    check(err)
    w.Flush()
}

type victim struct {
	id int
	posX int
	posY int
	tod int
	scheduledPickup bool
	saved bool
}

type ambulance struct {
	posX int
	posY int
	occupants []*victim
}

type hospital struct {
	name string
	posX int
	posY int
	ambulances []*ambulance // Length of this is also parked capacity
	aParked []*ambulance // Ambulences in parked
	aWaiting []*ambulance // Ambulences waiting to park
}

type event struct {
    time int
    name string
    vict *victim
    ambl *ambulance
    hosp *hospital
}

func (e *event) print() {
	fmt.Println(e.time);
}

//Returns < 0 if this object is less than obj
//Returns > 0 if this object is greater than obj
//Returns 0 if these objects are equal.
func (e *event) compare(comp *event) int {
    if (e.time < comp.time) {
    	return -1;
    } else if (e.time > comp.time) {
    	return 1;
    }
    return 0;
}

// ======== EVENT HEAP CLASS ==========
type eventHeap struct {
	array [400]*event
	load int
}

func (h *eventHeap) percolateDown(start int) {
	child := 0;
	tmp := h.array[start];
	hole := 0;

	for hole = start; hole * 2 <= h.load; hole = child {
		child = hole * 2;
		if (child != h.load && h.array[child + 1].compare(h.array[child]) < 0) {
			child++;
		}
		if (h.array[child].compare(tmp) < 0) {
			h.array[hole] = h.array[child];
		} else {
			break;
		}
	}
	h.array[hole] = tmp;
}

func (h *eventHeap) pop() *event {
	data := h.array[1];
	h.array[1] = h.array[h.load];
	h.load--;
	h.percolateDown(1);
	return data;
}

func (h *eventHeap) insert(obj *event) {
	//percolate up
	hole := h.load+1; //number of elements in heap
	for ; hole > 1 && obj.compare(h.array[hole/2]) < 0; hole /= 2 {
		h.array[hole] = h.array[hole/2];
	}
	h.array[hole] = obj;
	h.load++;
}

// ====== END EVENT HEAP CLASS ======

func idxOfAmbulance(ambl *ambulance, ambls []*ambulance) int {
    for p, v := range ambls {
        if (v == ambl) {
            return p
        }
    }
    return -1
}

var evtHeap eventHeap;

func mainMenu() {

}

var hospitals []*hospital;
var victims []*victim;

func calcDistance(posX int, posY int, posX2 int, posY2 int) int {
	dist := (posX - posX2) + (posY - posY2);
	return int(math.Abs(float64(dist)));
}

// Returns hospital closest and distance to that hospital
func findClosestHospital(posX int, posY int) (*hospital, int) {
	var closestHospital *hospital;
	var bestDist int;
	run := false;
	for i:=0; i<len(hospitals); i++ {
		var hospitalDistance = calcDistance(posX, posY, hospitals[i].posX, hospitals[i].posY);
		if (!run || bestDist > hospitalDistance) {
			closestHospital = hospitals[i];
			bestDist = hospitalDistance;
		}
		run = true;
	}
	return closestHospital, bestDist;
}

// Rescue time from current location (assuming space, etc)
func calcRescueTime(vict *victim, posX int, posY int) (*hospital, int) {
	victimDistance := calcDistance(vict.posX, vict.posY, posX, posY);

	hosp, hospitalDistance := findClosestHospital(vict.posX, vict.posY);

	// Current time, distance to victim, loading time, distance to hospital, unloading time
	time := clock + victimDistance + 3 + hospitalDistance + 1;

	return hosp, time;
}

func checkVictimElegibility(ambl *ambulance, vict *victim) (bool) {
	// Save processor time by checking the following conditions
	if (vict.scheduledPickup || vict.tod <= clock) {
		return false;
	}

	_, rescueTime := calcRescueTime(vict, ambl.posX, ambl.posY);

	if (rescueTime > vict.tod) {
		return false;
	} else if (len(ambl.occupants) > 0 && ambl.occupants[0].tod > rescueTime) {
		writeLog("Disregarded victim to not kill current victim in truck")
		return false;
	}
	return true;
}

func findVictim(ambl *ambulance) (*victim, int, bool) {
	var myVictim *victim;
	var foundElegibleVictim bool;
	for i:=0; i<len(victims); i++ {
		var isElegible bool;
		isElegible = checkVictimElegibility(ambl, victims[i]);
		if (isElegible) {
			myVictim = victims[i];
			foundElegibleVictim = true;
			break;
		}
	}
	if (!foundElegibleVictim) {
		return nil, 0, false;
	} else {
		victimArrivalTime := calcDistance(ambl.posX, ambl.posY, myVictim.posX, myVictim.posY);
		return myVictim, victimArrivalTime, true;
	}
}

var clock int;
func handleEvent(e *event) {
	switch e.name {
	case "LeaveHospital":
		writeLog("Leave Hospital Event " + strconv.Itoa(e.time));

		i := idxOfAmbulance(e.ambl, e.hosp.aParked);
		if (i > -1) {
			writeLog("Ambulance unparking");
			e.hosp.aParked = append(e.hosp.aParked[:i], e.hosp.aParked[i+1:]...)

			if (len(e.hosp.aWaiting) > 0) {
				writeLog("Ambulance from waitlist given a spot");
				tmp := e.hosp.aWaiting[0];
				i = 0;
				e.hosp.aWaiting = append(e.hosp.aParked[:i], e.hosp.aParked[i+1:]...)
				e.hosp.aParked = append(e.hosp.aParked, tmp);
				evt := &event{ time: e.time+1, name: "UnloadVictim", ambl: e.ambl, hosp: e.hosp };
				evtHeap.insert(evt);
			}
		}

		myVictim, arrivalT, foundElegibleVictim := findVictim(e.ambl);
		if (!foundElegibleVictim) {
			writeLog("Hospital ambulance found no victim to pick up.")
		} else {
			myVictim.scheduledPickup = true;
			evt := &event{ time: arrivalT, name: "LoadVictim", ambl: e.ambl, vict: myVictim };
			evtHeap.insert(evt);
		}
	case "RescueVictim":
		writeLog("RescueVictim Event " + strconv.Itoa(e.time));
		e.ambl.occupants = append(e.ambl.occupants, e.vict);
		e.ambl.posX = e.vict.posX;
		e.ambl.posY = e.vict.posY;

		headToHospital := false;
		if (len(e.ambl.occupants) < 2) {
			// Try to pick up another victim
			myVictim, arrivalT, foundElegibleVictim := findVictim(e.ambl);
			if (foundElegibleVictim) {
				writeLog("Loaded a second victim")
				myVictim.scheduledPickup = true;
				evt := &event{ time: arrivalT, name: "LoadVictim", ambl: e.ambl, vict: myVictim };
				evtHeap.insert(evt);
			} else {
				headToHospital = true;
			}
		} else {
			headToHospital = true;
		}
		if (headToHospital) {
			hosp, dist := findClosestHospital(e.ambl.posX, e.ambl.posY);
			evt := &event{ time: e.time+dist, name: "ArriveHospital", ambl: e.ambl, hosp: hosp };
			evtHeap.insert(evt);
		}
	case "LoadVictim":
		writeLog("LoadVictim Event " + strconv.Itoa(e.time) + " VictimID: " + strconv.Itoa(e.vict.id));
		e.ambl.posX = e.vict.posX;
		e.ambl.posY = e.vict.posY;
		evt := &event{ time: e.time+3, name: "RescueVictim", ambl: e.ambl, vict: e.vict };
		evtHeap.insert(evt);
	case "ArriveHospital":
		// NOTE: add ambulance to garage queue
		writeLog("ArriveHospital Event " + strconv.Itoa(e.time))
		e.ambl.posX = e.hosp.posX;
		e.ambl.posY = e.hosp.posY;

		if (len(e.hosp.ambulances) > len(e.hosp.aParked)) {
			// Parking spot open!
			e.hosp.aParked = append(e.hosp.aParked, e.ambl);
			evt := &event{ time: e.time+1, name: "UnloadVictim", ambl: e.ambl, hosp: e.hosp };
			evtHeap.insert(evt);
		} else {
			// Join the waitlist for an open spot
			e.hosp.aWaiting = append(e.hosp.aWaiting, e.ambl);
		}
	case "UnloadVictim":
		writeLog("UnloadVictim Event " + strconv.Itoa(e.time));
		for i:=0; i<len(e.ambl.occupants); i++ {
			if (e.ambl.occupants[i].tod < e.time) {
				writeLog("Victim Died " + strconv.Itoa(e.ambl.occupants[i].id))
			} else {
				writeLog("Victim Saved ID: " + strconv.Itoa(e.ambl.occupants[i].id))
				e.ambl.occupants[i].saved = true;
			}
		}
		e.ambl.occupants = e.ambl.occupants[:0];

		evt := &event{ time: e.time, name: "LeaveHospital", ambl: e.ambl, hosp: e.hosp };
		evtHeap.insert(evt);
	case "VictimDeath":
		writeLog("VictimDeath Event " + strconv.Itoa(e.time))
	default:
		writeLog("ERROR - EVENT TYPE NOT RECOGNIZED: " + e.name);
	}
}

func schedule_leaveHospital(ambl *ambulance, hosp *hospital, t int){
	evt := &event{ time: t, name: "LeaveHospital", ambl: ambl, hosp: hosp };
	evtHeap.insert(evt);
}

func finalReport() {
	victimCount := len(victims);
	alive := 0;
	dead := 0;
	w := bufio.NewWriter(f)

    _, err := fmt.Fprintf(w, "|%10s|%10s|%10s|%10s|\n", "Victim ID", "Position", "S-Time", "Saved")
    check(err)
	for i:=0; i<victimCount; i++ {
		victim := victims[i];
		if (victim.saved) {
			alive++;
		} else {
			dead++;
		}
		victPos := "(" + strconv.Itoa(victim.posX) + "," + strconv.Itoa(victim.posY) + ")"
		_, err = fmt.Fprintf(w, "|%10d|%10s|%10d|%10s|\n", victim.id, victPos, victim.tod, strconv.FormatBool(victim.saved))
    	check(err)
	}
	w.Flush()
	writeLog("Saved: " + strconv.Itoa(alive));
	writeLog("Perished: " + strconv.Itoa(dead));
}

func main() {
	openLog();

	evtHeap.load = 0;
	clock = 0;

	// Initialize victims
	victims = []*victim {
		&victim{ id: 0, posX: 50, posY: 55, tod: 35 },
		&victim{ id: 1, posX: 48, posY: 64, tod: 42 },
		&victim{ id: 2, posX: 49, posY: 53, tod: 32 },
		&victim{ id: 3, posX: 53, posY: 56, tod: 39 },
		&victim{ id: 4, posX: 53, posY: 48, tod: 31 },
		&victim{ id: 5, posX: 51, posY: 47, tod: 28 },
		&victim{ id: 6, posX: 52, posY: 51, tod: 33 },
		&victim{ id: 7, posX: 52, posY: 50, tod: 32 },
		&victim{ id: 8, posX: 52, posY: 60, tod: 42 },
		&victim{ id: 9, posX: 47, posY: 65, tod: 42 },
		&victim{ id: 10, posX: 57, posY: 54, tod: 31 },
		&victim{ id: 11, posX: 69, posY: 50, tod: 39 },
		&victim{ id: 12, posX: 57, posY: 57, tod: 34 },
		&victim{ id: 13, posX: 56, posY: 58, tod: 34 },
		&victim{ id: 14, posX: 64, posY: 50, tod: 34 },
		&victim{ id: 15, posX: 62, posY: 51, tod: 33 },
		&victim{ id: 16, posX: 56, posY: 56, tod: 32 },
		&victim{ id: 17, posX: 63, posY: 61, tod: 44 },
		&victim{ id: 18, posX: 60, posY: 51, tod: 31 },
		&victim{ id: 19, posX: 58, posY: 53, tod: 31 },
		&victim{ id: 20, posX: 57, posY: 72, tod: 39 },
		&victim{ id: 21, posX: 66, posY: 60, tod: 36 },
		&victim{ id: 22, posX: 77, posY: 56, tod: 43 },
		&victim{ id: 23, posX: 57, posY: 62, tod: 29 },
		&victim{ id: 24, posX: 65, posY: 65, tod: 40 },
		&victim{ id: 25, posX: 58, posY: 69, tod: 37 },
		&victim{ id: 26, posX: 61, posY: 56, tod: 27 },
		&victim{ id: 27, posX: 65, posY: 57, tod: 32 },
		&victim{ id: 28, posX: 63, posY: 70, tod: 43 },
		&victim{ id: 29, posX: 65, posY: 56, tod: 31 }};

	sort.Slice(victims, func(i, j int) bool {
	  return victims[i].tod < victims[j].tod
	})

	// Initialize ambulences and hospitals
	ambulances := []*ambulance {
		&ambulance{ posX: 45, posY: 45 },
		&ambulance{ posX: 45, posY: 45 },
		&ambulance{ posX: 45, posY: 45 },
		&ambulance{ posX: 45, posY: 45 }};
	hospitals = append(hospitals, &hospital{ name: "Austerlitz", posX: 45, posY: 45, ambulances: ambulances })

	ambulances = []*ambulance {
		&ambulance{ posX: 50, posY: 50 },
		&ambulance{ posX: 50, posY: 50 },
		&ambulance{ posX: 50, posY: 50 }};
	hospitals = append(hospitals, &hospital{ name: "Pasteur", posX: 50, posY: 50, ambulances: ambulances })

	ambulances = []*ambulance {
		&ambulance{ posX: 55, posY: 55 },
		&ambulance{ posX: 55, posY: 55 },
		&ambulance{ posX: 55, posY: 55 }};
	hospitals = append(hospitals, &hospital{ name: "De Gaulle", posX: 55, posY: 55, ambulances: ambulances })

	// Load all ambulances into the event heap for initialization
	for i:=0; i<len(hospitals); i++ {
		for j:=0; j<len(hospitals[i].ambulances); j++ {
			schedule_leaveHospital(hospitals[i].ambulances[j], hospitals[i], 0);
		}
	}

	// While there are still events to process, process events.
	for ; evtHeap.load > 0; {
		e := evtHeap.pop();
		handleEvent(e);
		clock = e.time;
	}

	finalReport();
}
