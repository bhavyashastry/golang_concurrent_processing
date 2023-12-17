package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"sync"
	"time"
)

type RequestPayload struct {
	ToSort [][]int `json:"to_sort"`
}

type ResponsePayload struct {
	SortedArrays [][]int `json:"sorted_arrays"`
	TimeNs       int64   `json:"time_ns"`
}

func processSequential(w http.ResponseWriter, r *http.Request) {
	// Your processSequential implementation
	var payload RequestPayload
	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	startTime := time.Now()

	sortedArrays := make([][]int, len(payload.ToSort))
	for i, arr := range payload.ToSort {
		sortedArr := make([]int, len(arr))
		copy(sortedArr, arr)
		sort.Ints(sortedArr)
		sortedArrays[i] = sortedArr
	}

	timeTaken := time.Since(startTime).Nanoseconds()

	response := ResponsePayload{
		SortedArrays: sortedArrays,
		TimeNs:       timeTaken,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func processConcurrent(w http.ResponseWriter, r *http.Request) {
	// Your processConcurrent implementation
	var payload RequestPayload
	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	startTime := time.Now()

	var wg sync.WaitGroup
	wg.Add(len(payload.ToSort))
	mutex := sync.Mutex{}
	sortedArrays := make([][]int, len(payload.ToSort))

	for i, arr := range payload.ToSort {
		go func(index int, array []int) {
			defer wg.Done()

			sortedArr := make([]int, len(array))
			copy(sortedArr, array)
			sort.Ints(sortedArr)

			mutex.Lock()
			sortedArrays[index] = sortedArr
			mutex.Unlock()
		}(i, arr)
	}

	wg.Wait()

	timeTaken := time.Since(startTime).Nanoseconds()

	response := ResponsePayload{
		SortedArrays: sortedArrays,
		TimeNs:       timeTaken,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func testScript() {
	url := "http://localhost:8000/process-single"

	payload := map[string]interface{}{
		"to_sort": [][]int{
			{3, 1, 2},
			{6, 4, 5},
			{9, 7, 8},
		},
	}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		fmt.Println("Error encoding JSON:", err)
		return
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonPayload))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return
	}
	defer resp.Body.Close()

	fmt.Println("Response Status:", resp.Status)
	fmt.Println("Response Body:")
	buf := new(bytes.Buffer)
	_, _ = buf.ReadFrom(resp.Body)
	fmt.Println(buf.String())
}

func main() {
	// Starting the HTTP server
	http.HandleFunc("/process-single", processSequential)
	http.HandleFunc("/process-concurrent", processConcurrent)

	// Starting the test script
	go testScript()

	// Listening on port 8000
	if err := http.ListenAndServe(":8000", nil); err != nil {
		fmt.Println("Server error:", err)
	}
}
