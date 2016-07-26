package main

import (	
			"os"			
         "io/ioutil"
			s "strings"
			"fmt"
			"bufio"
			"strconv"
		 )


type gate struct { 
	output_channel string 
	gate_type string
	input_chan_1 chan bool
	input_chan_2 chan bool
	output_channel_2 string
}

type output_thing struct {
	name string
	value bool
}

func check(e error) {
    if e != nil {
        panic(e)
    }
}

func main() {
	argsWithoutProg := os.Args[1]
	adder, err := ioutil.ReadFile(argsWithoutProg)
   check(err)

	adder_contents := s.Split(string(adder), "\n")

	channels := make(map[string][]chan bool)
	outputs := make(map[string] chan bool)
	var gates []gate
	var clock chan bool = nil

   for i := 0; i < len(adder_contents); i++ {
		eachline := s.Split(string(adder_contents[i]), " ")
			if eachline[0] == "xor" {
				inchan1 := make(chan bool, 1)
				inchan2 := make(chan bool, 1)
				channels[eachline[1]] = append(channels[eachline[1]], inchan1)
				channels[eachline[2]] = append(channels[eachline[2]], inchan2)
				outchan := channels[eachline[4]]
				if s.Contains(eachline[4], "output") {
    				final_outchan := make(chan bool, 1)
    				channels[eachline[4]] = append(outchan, final_outchan)
    				outputs[eachline[4]] = final_outchan
				}
				gates = append(gates, gate{eachline[4], "xor", inchan1, inchan2, ""})	
			} else if eachline[0] == "and" {
				inchan1and := make(chan bool, 1)
				inchan2and := make(chan bool, 1)
				channels[eachline[1]] = append(channels[eachline[1]], inchan1and)
				channels[eachline[2]] = append(channels[eachline[2]], inchan2and)
				outchanand := channels[eachline[4]]
				if s.Contains(eachline[4], "output") {
					final_outchanand := make(chan bool, 1)
					channels[eachline[4]] = append(outchanand, final_outchanand)
					outputs[eachline[4]] = final_outchanand
				}
				gates = append(gates, gate{eachline[4], "and", inchan1and, inchan2and, ""})			
			} else if eachline[0] == "or" {
				inchan1or := make(chan bool, 1)
				inchan2or := make(chan bool, 1)
				channels[eachline[1]] = append(channels[eachline[1]], inchan1or)
				channels[eachline[2]] = append(channels[eachline[2]], inchan2or)
				outchanor := channels[eachline[4]]
				if s.Contains(eachline[4], "output") {
					final_outchanor := make(chan bool, 1)	
					channels[eachline[4]] = append(outchanor, final_outchanor)
					outputs[eachline[4]] = final_outchanor
				}
				gates = append(gates, gate{eachline[4], "or", inchan1or, inchan2or, ""})
			} else if eachline[0] == "not" {
				in_not := make(chan bool, 1)
				channels[eachline[1]] = append(channels[eachline[1]], in_not)
				out_not := channels[eachline[4]]
				if s.Contains(eachline[4], "output") {
					final_outnot := make(chan bool, 1)
					channels[eachline[4]] = append(out_not, final_outnot)
					outputs[eachline[4]] = final_outnot
				}	
				gates = append(gates, gate{eachline[4], "not", in_not, nil, ""})
			} else if eachline[0] == "nand" {
				in_nand_1 := make(chan bool, 1)
				in_nand_2 := make(chan bool, 1)
				channels[eachline[1]] = append(channels[eachline[1]], in_nand_1)
				channels[eachline[2]] = append(channels[eachline[2]], in_nand_2)
				out_nand := channels[eachline[4]]
				if s.Contains(eachline[4], "output") {
					final_outnand := make(chan bool, 1)
					channels[eachline[4]] = append(out_nand, final_outnand)
					outputs[eachline[4]] = final_outnand
				}
				gates = append(gates, gate{eachline[4], "nand", in_nand_1, in_nand_2, ""})
			} else if eachline[0] == "nor" {
				in_nor_1 := make(chan bool, 1)
				in_nor_2 := make(chan bool, 1)
				channels[eachline[1]] = append(channels[eachline[1]], in_nor_1)
				channels[eachline[2]] = append(channels[eachline[2]], in_nor_2)
				out_nor := channels[eachline[4]]
				if s.Contains(eachline[4], "output") {
					final_outnor := make(chan bool, 1)
					channels[eachline[4]] = append(out_nor, final_outnor)
					outputs[eachline[4]] = final_outnor
				}
				gates = append(gates, gate{eachline[4], "nor", in_nor_1, in_nor_2, ""})
			} else if eachline[0] == "dflipflop" {
				d_in := make(chan bool, 1)
				cl_in := make(chan bool, 1)
				channels[eachline[1]] = append(channels[eachline[1]], d_in)
				channels[eachline[2]] = append(channels[eachline[2]], cl_in)
				//out_q := channels[eachline[4]]
				//out_notq := channels[eachline[5]]
				if s.Contains(eachline[4], "output") {
					final_outq := make(chan bool, 1)
					outputs[eachline[4]] = final_outq
				}
				if s.Contains(eachline[5], "output") {
					final_outnotq := make(chan bool, 1)
					outputs[eachline[5]] = final_outnotq
				}
				gates = append(gates, gate{eachline[4], "dflipflop", d_in, cl_in, eachline[5]})
			} else if eachline[0] == "clock" {
				clock = make(chan bool)
			}
	}

var pulses_per_sim int
if clock != nil {
	reader := bufio.NewReader(os.Stdin)
		//fmt.Println("Number of clock pulses per second: ")
		//text, _ := reader.ReadString('\n')
		//pulses_per_sec, err := strconv.Atoi(text)
		fmt.Println("num pulses")
		text2, _ := reader.ReadString('\n')
		pulses_per_sim, _ = strconv.Atoi(text2[:len(text2) - 1])
		//sec_per_pulse := 1 / pulses_per_sec
}
	for key := range gates {
		if gates[key].gate_type == "xor" {
			go xor(gates[key].input_chan_1, gates[key].input_chan_2, channels[gates[key].output_channel], gates[key].output_channel)
		} else if gates[key].gate_type == "and" {
			go and(gates[key].input_chan_1, gates[key].input_chan_2,channels[gates[key].output_channel], gates[key].output_channel)
		} else if gates[key].gate_type == "or" {
			go or(gates[key].input_chan_1, gates[key].input_chan_2,channels[gates[key].output_channel], gates[key].output_channel)
		} else if gates[key].gate_type == "dflipflop" {
			go dflipflop(gates[key].input_chan_1, gates[key].input_chan_2,
channels[gates[key].output_channel], channels[gates[key].output_channel_2], pulses_per_sim,
outputs[gates[key].output_channel], outputs[gates[key].output_channel_2])
		}
	}

	for key := range channels {
		if s.Contains(key, "extern") {
			reader := bufio.NewReader(os.Stdin)
			fmt.Print("Enter boolean value for: " + key + " ")
			text, _ := reader.ReadString('\n')
			
			if text == "true\n" {
				send_chans := channels[key]
				for i := 0; i < len(send_chans); i++ {
					select {
        				case send_chans[i] <- true:
        				default: 
    				}
				}
			} else if text == "false\n" {
				send_chans := channels[key]
				for i := 0; i < len(send_chans); i++ {
					select {
						case send_chans[i] <- false:
						default:
					}
				}
			}
		}	
	}

	if clock != nil {
		clock_channels := channels["clock"]
		for i := 0; i < pulses_per_sim; i++ {
			for j := 0; j < len(clock_channels); j++ {
				select {
					case clock_channels[j] <- true:
					default:
				}
			}
			//time.Sleep(time.Second * sec_per_pulse)
		}
		}


	for key := range outputs {
		output := <- outputs[key]
		if output == true {
			fmt.Println(key + ": true")
		} else if output == false {
			fmt.Println(key + ": false")
		}
	}

}

func and(in1and chan bool, in2and chan bool, outchanand []chan bool, outname string) {
	firsttime := true
	gotone := false
	gottwo := false
	in1 := false
	in2 := false
	oldvalue1 := false
	oldvalue2 := false

	for true {
		select {
			case in1 = <- in1and: gotone = true
			default:
		}
		select {
			case in2 = <- in2and: gottwo = true
			default:
		}

		if (oldvalue1 != in1 || oldvalue2 != in2 || firsttime == true) &&  (gotone == true && gottwo == true) {
			firsttime = false
			result := in1 && in2
			for i := 0; i < len(outchanand); i++ {
   			select {
      			case outchanand[i] <- result:
      			default: 
    			}
			}
		oldvalue1 = in1
		oldvalue2 = in2
		}
	}	
}

func or(in1or chan bool, in2or chan bool, outchanor []chan bool, outname string) {
	firsttime := true
	gotone := false
	gottwo := false
	in1 := false
	in2 := false
	oldvalue1 := false
	oldvalue2 := false

	for true {
		select {
			case in1 = <- in1or: gotone = true
			default:
		}
		select {
  			case in2 = <- in2or: gottwo = true
			default:
		}

		if (oldvalue1 != in1 || oldvalue2 != in2 || firsttime == true) &&  (gotone == true && gottwo == true) {
			firsttime = false
			result := in1 || in2

			for i := 0; i < len(outchanor); i++ {
				select {
					case outchanor[i] <- result:
					default:
				}
			}
			oldvalue1 = in1
			oldvalue2 = in2
		}	
	}
}

func not(in1not chan bool, outnot []chan bool) {
	firsttime := true
	gotone := false
	in1 := false
	oldvalue1 := false


	for true {
		select {
			case in1 = <- in1not: gotone = true
			default:
		}
		
		if (oldvalue1 != in1 || firsttime == true) &&  gotone == true  {
			firsttime = false
			result := !in1

			for i := 0; i < len(outnot); i++ {
				select {
					case outnot[i] <- result:
					default:
				}
			}
			oldvalue1 = in1
		}
	}
}

func xor(in1xor chan bool, in2xor chan bool, outchan []chan bool, outname string) {
	firsttime := true
	gotone := false
	gottwo := false
	in1 := false
	in2 := false
	oldvalue1 := false
	oldvalue2 := false

	for true {
		select {
			case in1 = <- in1xor: gotone = true
			default:
		}
		select {
			case in2 = <- in2xor: gottwo = true
			default:
		}

		if (oldvalue1 != in1 || oldvalue2 != in2 || firsttime == true) &&  (gotone == true && gottwo == true) {
			firsttime = false
			result := ((in1 && !in2) || (!in1 && in2))
		
			for i := 0; i < len(outchan); i++ {
				select {
					case outchan[i] <- result:
					default: 
				}
			}
			oldvalue1 = in1
			oldvalue2 = in2
		}
	}
}

func nand(in1nand chan bool, in2nand chan bool, outnand []chan bool) {
	firsttime := true
	gotone := false
	gottwo := false
	in1 := false
	in2 := false
	oldvalue1 := false
	oldvalue2 := false

	for true {
		select {
			case in1 = <- in1nand: gotone = true
			default:
		}
		select {
			case in2 = <- in2nand: gottwo = true
			default:
		}

		if (oldvalue1 != in1 || oldvalue2 != in2 || firsttime == true) &&  (gotone == true && gottwo == true) {
			firsttime = false
			result := !(in1 && in2)

			for i := 0; i < len(outnand); i++ {
				select {
					case outnand[i] <- result:
					default:
				}
			}
			oldvalue1 = in1
			oldvalue2 = in2
		}
	}
}

func nor(in1nor chan bool, in2nor chan bool, outnor []chan bool) {
	firsttime := true
	gotone := false
	gottwo := false
	in1 := false
	in2 := false
	oldvalue1 := false
	oldvalue2 := false

	for true {
      select {
			case in1 = <- in1nor: gotone = true
			default:
		}
		select {
			case in2 = <- in2nor: gottwo = true
			default:
		}

		if (oldvalue1 != in1 || oldvalue2 != in2 || firsttime == true) &&  (gotone == true && gottwo == true) {
			firsttime = false
			result := !(in1 || in2)

			for i := 0; i < len(outnor); i++ {
				select {
					case outnor[i] <- result:
					default:
				}
			}
			oldvalue1 = in1
			oldvalue2 = in2
		}
	}
}

func dflipflop(d_chan chan bool, clock chan bool, q []chan bool, notq []chan bool, pulses_per_sim int, final_output_q chan bool, final_output_qnot chan bool) {
	d := false
	cl := false
	counter := 0	

	for counter < pulses_per_sim {
		select {
			case d = <- d_chan: 
			default:
		}	
		select {
			case cl = <- clock:
			default:
		}

		if cl == true {
			counter++
			fmt.Println("cl==true")
			for i := 0; i < len(q); i++ {
				select {
					case q[i] <- d: fmt.Println("sending to q")
					default:
				}
			}
		
			for i := 0; i < len(notq); i++ {
				select {
					case notq[i] <- !d: fmt.Println("sending to not q")
					default:
				}
			}
			cl = false
		}
	}
	
		if final_output_q != nil {
			select {
				case final_output_q <- d:
				default:
			}
		}
		if final_output_qnot != nil {
			select {
				case final_output_qnot <- !d:
				default:
			}	
		}
}	
