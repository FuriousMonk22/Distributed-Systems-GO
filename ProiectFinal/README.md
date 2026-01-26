# Lamport Mutual Exclusion; Go + Python (Timeline Visualization)

This project implements **Lamport’s mutual exclusion algorithm** using **logical clocks** in Go and provides a **Python script** that automatically runs the Go program and generates a **timeline graph** of process execution using `matplotlib`.

---

## 📁 Project Structure


 `graph.py` and `main.go` are expected to be in the same folder (`ProiectFinal/`), unless you modify the path in `graph.py`.

---

## 🧩 Requirements

### Go
- Go **1.18+**
```bash
go version
```

### Python
```bash
pip install matplotlib
cd ProiectFinal 
python3 graph.py
```

---

## 🧾 Input Parameters
- You will be prompted to enter simulation parameters. Press Enter to keep the default values.
```bash
Numar procese (p_count) [5]:
Clock initial pentru pid 0 [0]:
Clock initial pentru pid 1 [0]:
Clock initial pentru pid 2 [0]:
Clock initial pentru pid 3 [0]:
Clock initial pentru pid 4 [0]:
Poll interval (ms) [25]:
CS duration (ms) [200]:
Main runtime (sec) [5]:
```

---

## 🔄 Execution Flow
- graph.py reads the input parameters.
- it runs:
```bash
go run main.go
```
- main.go executes Lamport’s mutual exclusion algorithm and prints logs like:
```bash
pid 1 request
pid 1 entered CS
pid 1 left CS
```
- graph.py parses the output and generates a timeline diagram.

---

## 📊 Output Files

After running the simulation, the following files are generated:

- **`lamport_timeline.png`**  
  A timeline diagram that visualizes the execution of the Lamport mutual exclusion algorithm.
  - X-axis represents the **logical order of events**
  - Y-axis represents the **processes** (`P0`, `P1`, `P2`, …)
  - Each process is shown using a **unique color**
  - Critical sections are displayed as **non-overlapping segments**, proving mutual exclusion

- **`go_output.txt`**  
  Contains the full standard output produced by `main.go`.
  - Useful for debugging
  - Can be inspected to verify request, enter, and leave events
  - Used internally by `graph.py` to build the timeline

---

## ⚙️ Modifying the Path to `main.go`

If `graph.py` cannot locate `main.go`, you need to update the path used to run the Go program.

Open **`graph.py`** and locate the following function:

```python
def run_go_and_capture(input_str: str):
    subprocess.run(
        ["go", "run", "./ProiectFinal/main.go"],
        ...
    )
```