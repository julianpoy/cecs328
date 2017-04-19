package heap

import "fmt"

type event struct {
    time int
    name string
    victimId int
}

func (e event) print() {
	fmt.Println(e.time);
}

//Returns < 0 if this object is less than obj
//Returns > 0 if this object is greater than obj
//Returns 0 if these objects are equal.
func (e event) compare(comp event) int {
    if (e.time < comp.time) {
    	return -1;
    } else if (e.time > comp.time) {
    	return 1;
    }
    return 0;
}

type heap struct {
	array []event
	load int
}

func (h heap) build() {
	h.array = make([]event, 200);
	fmt.Printf("%v", h.array);
}

func (h heap) percolateDown(start int) {
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

func (h heap) pop() event {
	data := h.array[1];
	h.array[1] = h.array[h.load];
	fmt.Printf("val - ");
	h.array[1].print();
	h.load--;
	h.percolateDown(1);
	return data;
}

func (h heap) insert(obj event) {
	//percolate up
	hole := h.load + 1; //number of elements in heap
	for ; hole > 1 && obj.compare(h.array[hole/2]) < 0; hole /= 2 {
		h.array[hole] = h.array[hole/2];
	}
	h.array[hole] = obj;
	h.load++;
}