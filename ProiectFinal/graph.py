import subprocess
import re
import sys
import matplotlib.pyplot as plt


def read_int(prompt: str, default: int) -> int:
    s = input(f"{prompt} [{default}]: ").strip()
    if s == "":
        return default
    try:
        return int(s)
    except ValueError:
        print(f"Valoare invalidă, folosesc default={default}.")
        return default


def run_go_and_capture(input_str: str) -> str:
    try:
        res = subprocess.run(
            ["go", "run", "./ProiectFinal/main.go"],
            input=input_str,
            capture_output=True,
            text=True,
            check=True,
        )
        return res.stdout
    except FileNotFoundError:
        print("Nu găsesc comanda `go`. Instalează Go sau asigură-te că e în PATH.")
        sys.exit(1)
    except subprocess.CalledProcessError as e:
        print("Eroare la rularea `go run main.go`.")
        print("STDOUT:\n", e.stdout)
        print("STDERR:\n", e.stderr)
        sys.exit(1)


def parse_events(go_output: str):
    # Acceptă doar liniile care încep cu "pid <n> ..."
    rx_req = re.compile(r"^pid\s+(\d+)\s+request\s*$")
    rx_ent = re.compile(r"^pid\s+(\d+)\s+entered\s+CS\s*$")
    rx_left = re.compile(r"^pid\s+(\d+)\s+left\s+CS\s*$")

    events = []
    t = 0
    for line in go_output.splitlines():
        line = line.strip()
        m = rx_req.match(line)
        if m:
            t += 1
            events.append((int(m.group(1)), "REQUEST", t))
            continue
        m = rx_ent.match(line)
        if m:
            t += 1
            events.append((int(m.group(1)), "ENTER", t))
            continue
        m = rx_left.match(line)
        if m:
            t += 1
            events.append((int(m.group(1)), "LEAVE", t))
            continue

    return events, t


def plot_timeline(events, t_max, out_png="./ProiectFinal/lamport_timeline.png"):
    if not events:
        print("N-am găsit evenimente în output (linii de forma: 'pid X request/entered CS/left CS').")
        return

    pids = sorted(set(pid for pid, _, _ in events))
    pid_to_y = {pid: i for i, pid in enumerate(pids)}

    # --- ADĂUGAT: culoare unică per proces ---
    cmap = plt.cm.get_cmap("tab20", len(pids))  # N culori distincte
    pid_color = {pid: cmap(i) for i, pid in enumerate(pids)}
    # ----------------------------------------

    # Construim intervalele de CS (ENTER -> LEAVE) pe pid
    enter_time = {}
    cs_intervals = []
    for pid, etype, t in events:
        if etype == "ENTER":
            enter_time[pid] = t
        elif etype == "LEAVE" and pid in enter_time:
            cs_intervals.append((pid, enter_time[pid], t))
            del enter_time[pid]

    # Figura (scalare în funcție de numărul de procese)
    height = max(6, min(30, 0.25 * len(pids) + 4))
    plt.figure(figsize=(14, height))

    # Linii orizontale pentru procese
    for pid in pids:
        y = pid_to_y[pid]
        plt.hlines(y, xmin=0, xmax=t_max + 1, linestyles="dotted", alpha=0.3)

    # Evenimente (markere) — aceeași culoare per pid
    for pid, etype, t in events:
        y = pid_to_y[pid]
        c = pid_color[pid]

        if etype == "REQUEST":
            plt.plot(t, y, "o", color=c)
            plt.text(t, y + 0.12, "REQ", fontsize=7, ha="center", color=c)
        elif etype == "ENTER":
            plt.plot(t, y, "s", color=c)
            plt.text(t, y + 0.12, "ENTER", fontsize=7, ha="center", color=c)
        elif etype == "LEAVE":
            plt.plot(t, y, "s", color=c)
            plt.text(t, y + 0.12, "LEAVE", fontsize=7, ha="center", color=c)

    # Intervalele CS ca segmente groase — aceeași culoare per pid
    for pid, t1, t2 in cs_intervals:
        y = pid_to_y[pid]
        plt.hlines(y, t1, t2, linewidth=4, color=pid_color[pid])

    plt.yticks(range(len(pids)), [f"P{pid}" for pid in pids])
    plt.xlabel("Ordinea evenimentelor (timp logic = ordinea liniilor relevante din output)")
    plt.title("Lamport Mutual Exclusion – Timeline")
    plt.grid(axis="x", linestyle="--", alpha=0.3)

    plt.tight_layout()
    plt.savefig(out_png, dpi=200)
    plt.show()
    print(f"Salvat: {out_png}")



def main():
    # Citești parametrii (aceeași ordine pe care o așteaptă main.go)
    print("Introdu parametrii (Enter = default).")

    p_count = read_int("Numar procese (p_count)", 5)

    clocks = []
    for i in range(p_count):
        clocks.append(read_int(f"Clock initial pentru pid {i}", 0))

    poll_ms = read_int("Poll interval (ms)", 25)
    cs_ms = read_int("CS duration (ms)", 200)
    runtime_sec = read_int("Main runtime (sec)", 5)

    # Construim stdin pentru Go exact în ordinea prompturilor
    input_lines = [str(p_count)]
    input_lines += [str(c) for c in clocks]
    input_lines += [str(poll_ms), str(cs_ms), str(runtime_sec)]
    go_stdin = "\n".join(input_lines) + "\n"

    # Rulează Go și capturează output-ul
    go_out = run_go_and_capture(go_stdin)

    # (opțional) salvează output-ul pentru debugging
    with open("./ProiectFinal/go_output.txt", "w", encoding="utf-8") as f:
        f.write(go_out)

    events, t_max = parse_events(go_out)
    plot_timeline(events, t_max, out_png="./ProiectFinal/lamport_timeline.png")


if __name__ == "__main__":
    main()
