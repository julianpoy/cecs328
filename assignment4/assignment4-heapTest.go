package main

import "fmt"
import "bufio"
import "os"
import "strconv"
import "strings"

type intHeap struct {
	array []*int
	load int
}

func compare(i *int, j *int) int {
    if (*i < *j) {
    	return -1;
    } else if (*i > *j) {
    	return 1;
    }
    return 0;
}

func (h *intHeap) print() {
	for i:=1; i<=h.load; i++ {
		fmt.Print(*h.array[i]);
		fmt.Print(" ");
	}
	fmt.Println("");
}

func (h *intHeap) percolateDown(start int) {
	child := 0;
	tmp := h.array[start];
	hole := 0;

	for hole = start; hole * 2 <= h.load; hole = child {
		child = hole * 2;
		if (child != h.load && compare(h.array[child + 1], h.array[child]) < 0) {
			child++;
		}
		if (compare(h.array[child], tmp) < 0) {
			h.array[hole] = h.array[child];
		} else {
			break;
		}
	}
	h.array[hole] = tmp;
}

func (h *intHeap) pop() *int {
	if (h.load == 0) {
		cantFind := -1;
		return &cantFind;
	}
	data := h.array[1];
	h.array[1] = h.array[h.load];
	h.load--;
	h.percolateDown(1);
	return data;
}

func (h *intHeap) insert(obj int) {
	//percolate up
	hole := h.load+1; //number of elements in heap
	for ; hole > 1 && compare(&obj, h.array[hole/2]) < 0; hole /= 2 {
		h.array[hole] = h.array[hole/2];
	}
	h.array[hole] = &obj;
	h.load++;
}

func readIn() string {
	reader := bufio.NewReader(os.Stdin)
	text, _ := reader.ReadString('\n')
	return strings.TrimRight(text, "\r\n");
}

var globalHeap *intHeap;

func build(a []*int){
	globalHeap.array = a;
	// Load is always less than length by 1
	globalHeap.load = len(a)-1;
	for i:=globalHeap.load/2; i>0; i-- {
		globalHeap.percolateDown(i);
	}
}

func Atoi(s string) int {
	myInt, _ := strconv.Atoi(s);
	return myInt;
}

func mainMenu() {
	fmt.Println("1. Create a heap");
	fmt.Println("2. Insert an element")
	fmt.Println("3. Pop an element")

	input := readIn();

	switch input {
	case "1":
		globalHeap = new(intHeap);
		fmt.Println("Do you want to populate that with initial values?");

		input = readIn();

		if (input == "y") {
			fmt.Println("Enter values now hitting enter every time ( to quit): ");
			var vals []*int;
			initialVals := 0;
			vals = append(vals, &initialVals);
			inputting := true;
			for ; inputting; {
				input = readIn();
				if (input != "q") {
					conv := Atoi(input);
					vals = append(vals, &conv);
				} else {
					break;
				}
			}
			build(vals);
			globalHeap.print();
		}
	case "2":
		fmt.Println("Enter your value: ");
		input := readIn();
		globalHeap.insert(Atoi(input));
		globalHeap.print();
	case "3":
		popped := globalHeap.pop();
		fmt.Println("Popped: ", *popped);
		fmt.Println("Heap:");
		globalHeap.print();
	}
	mainMenu();
}

func main() {
	mainMenu();
}
