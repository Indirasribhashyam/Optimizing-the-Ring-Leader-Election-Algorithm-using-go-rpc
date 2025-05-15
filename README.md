# Optimizing the Ring Leader Election Algorithm using Go RPC

Leader election is an essential operation in distributed systems. This project optimizes the classic Ring Election Algorithm using Go’s concurrency features like goroutines, atomic operations, and channels. It ensures only one election runs at a time, minimizes overhead, and handles failures gracefully.

## Problem Statement

In traditional ring algorithms, the following problems occur:
- Multiple processes may start elections simultaneously.
- Dead or unresponsive processes can cause delays or infinite waits.
- There is no guaranteed mechanism for all processes to acknowledge the leader.

## Description

This project enhances the Ring Leader Election Algorithm to address common issues like:
- Multiple simultaneous elections
- Deadlocks due to non-responsive processes
- Message overhead
- Lack of leader acknowledgment

We leverage Go's atomic operations, channels, and timeout mechanisms to make the election robust and efficient.

## Key Optimizations

**Prevention of Multiple Simultaneous Elections:**
Uses an atomic boolean (electionRunning) to ensure only one election occurs at a time, reducing unnecessary message overhead.

**Failure Handling and Timeout Mechanism:**
Implements a timeout (time.After(3 * time.Second)) to handle process failures and restart elections if no response is received.

**Reduction of Message Overhead:**
Guarantees a single election process to avoid excessive message traffic.

**Leader Acknowledgment and Robustness:**
Introduces leader acknowledgment channels (Leader channels) to ensure all processes receive the leader ID, even if messages are lost.

**Handling Process Failures During Election:**
Utilizes a randomized "alive" status (IsAlive) to skip elections for downed processes.

**Leader Failure Detection and Election Restart:**
Monitors leader status and restarts the election if no leader is detected within 5 seconds.

**Prioritization of Stronger Nodes:**
The first alive process starts the election, and higher IDs have a better chance of winning, prioritizing stronger nodes.

**Deadlock Prevention:**
Implements timeout-based retries to prevent deadlocks and ensure elections complete successfully.

## Technologies Used

- **Go (Golang)**: Concurrency, channels, atomic operations
- **Randomization**: Simulate process failures
- **Goroutines**: Parallel execution
- **Standard Library Only**: No external dependencies

## Usage


1. **Clone the repository:**

    ```bash
    git clone https://github.com/Indirasribhashyam/Optimizing-the-Ring-Leader-Election-Algorithm-using-go-rpc.git
    ```

2. **Change directory to the project folder:**

    ```bash
    cd Optimizing-the-Ring-Leader-Election-Algorithm-using-go-rpc
    ```

3. **Compile and Run the application:**

    ```bash
    go run main.go
    ```

## Future Work

* Integrate actual Go RPC between processes.
* Create UI-based visualization of the ring and election progress.
* Extend to dynamic rings (process joins/leaves).
* Add metrics dashboard for election time, retries, etc.
