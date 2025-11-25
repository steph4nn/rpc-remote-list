package remotelist

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

func getLastLogId() int64 {
	logFile := "../pkg/db/logs/operations.log"

	file, err := os.Open(logFile)
	if err != nil {
		return 0
	}
	defer file.Close()

	var lastID int64 = 0
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		if strings.TrimSpace(line) == "" {
			continue
		}

		var entry LogEntry
		if err := json.Unmarshal([]byte(line), &entry); err != nil {
			continue
		}

		if entry.Id > lastID {
			lastID = entry.Id
		}
	}

	return lastID
}

func getState() map[int][]int {

	snapshot, lastLogId := getLatestSnapshot()

	if snapshot != nil {
		fmt.Printf("Snapshot carregado. Último log aplicado: %d\n", lastLogId)
		ApplyLogsOnSnapshot(snapshot, lastLogId)
		return snapshot
	}

	data, err := os.ReadFile("../pkg/db/db.json")
	if err != nil || len(data) == 0 {
		fmt.Println("Iniciando com estado vazio")
		return make(map[int][]int)
	}

	var list map[int][]int
	if err := json.Unmarshal(data, &list); err != nil {
		return make(map[int][]int)
	}

	fmt.Println("Estado carregado de db.json")
	return list
}

func getLatestSnapshot() (map[int][]int, int64) {
	dir := "../pkg/db/snapshots/"
	files, err := os.ReadDir(dir)
	if err != nil || len(files) == 0 {
		return nil, 0
	}

	var latestFile os.DirEntry
	var latestTimestamp int64

	for _, f := range files {
		var ts int64
		n, err := fmt.Sscanf(f.Name(), "snapshot_%d.json", &ts)
		if err == nil && n == 1 && ts > latestTimestamp {
			latestTimestamp = ts
			latestFile = f
		}
	}

	data, err := os.ReadFile(filepath.Join(dir, latestFile.Name()))
	if err != nil {
		return nil, 0
	}

	var snap Snapshot
	if err := json.Unmarshal(data, &snap); err != nil {
		return nil, 0
	}

	return snap.List, snap.LastLogApplied
}

func SaveState(l *RemoteList) error {
	data, err := json.MarshalIndent(l.list, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal data: %w", err)
	}
	if err := os.WriteFile("../pkg/db/db.json", data, 0644); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}
	return nil
}

func ApplyLogsOnSnapshot(list map[int][]int, lastLogId int64) {
	logFile := "../pkg/db/logs/operations.log"

	file, err := os.Open(logFile)
	if err != nil {
		fmt.Println("Nenhum log para aplicar")
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		var entry LogEntry

		if err := json.Unmarshal([]byte(line), &entry); err != nil {
			log.Printf("Erro ao parsear log: %v\n", err)
			continue
		}

		if entry.Id <= lastLogId {
			continue
		}

		switch entry.Operation {
		case "APPEND":
			list[entry.ListId] = append(list[entry.ListId], entry.Value)
		case "REMOVE":
			if len(list[entry.ListId]) > 0 {
				list[entry.ListId] = list[entry.ListId][:len(list[entry.ListId])-1]
				if len(list[entry.ListId]) == 0 {
					delete(list, entry.ListId)
				}
			}
		}
	}
}

func AppendLog(logId int64, args AppendArgs) {
	logEntry := LogEntry{
		Id:        logId,
		Operation: "APPEND",
		ListId:    args.ListId,
		Value:     args.Value,
	}
	WriteLogEntry(logEntry)
}

func RemoveLog(logID int64, args RemoveArgs, removedValue int) {
	logEntry := LogEntry{
		Id:        logID,
		Operation: "REMOVE",
		ListId:    args.ListId,
		Value:     removedValue,
	}
	WriteLogEntry(logEntry)
}

func WriteLogEntry(entry LogEntry) {
	logDir := "../pkg/db/logs"
	logFile := filepath.Join(logDir, "operations.log")

	file, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		log.Printf("Erro ao abrir arquivo de log: %v\n", err)
		return
	}
	defer file.Close()

	data, err := json.Marshal(entry)
	if err != nil {
		log.Printf("Erro ao serializar log: %v\n", err)
		return
	}

	file.Write(data)
	file.WriteString("\n")
}

func SaveSnapshot(l *RemoteList, lastLog int64) error {
	timestamp := time.Now().Unix()
	snapshot := Snapshot{
		Timestamp:      timestamp,
		LastLogApplied: lastLog,
		List:           l.list,
	}

	data, err := json.MarshalIndent(snapshot, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal snapshot: %w", err)
	}

	fileName := fmt.Sprintf("../pkg/db/snapshots/snapshot_%d.json", timestamp)
	if err := os.WriteFile(fileName, data, 0644); err != nil {
		return fmt.Errorf("failed to write snapshot file: %w", err)
	}
	return nil
}

func DeleteOldSnapshots() {
	dir := "../pkg/db/snapshots/"
	files, err := os.ReadDir(dir)
	if err != nil || len(files) <= 2 {
		return
	}

	var snaps []snapFile
	for _, f := range files {
		var ts int64
		n, err := fmt.Sscanf(f.Name(), "snapshot_%d.json", &ts)
		if err == nil && n == 1 {
			snaps = append(snaps, snapFile{f.Name(), ts})
		}
	}

	sort.Slice(snaps, func(i, j int) bool {
		return snaps[i].timestamp > snaps[j].timestamp
	})

	for _, s := range snaps[2:] {
		_ = os.Remove(dir + s.name)
	}
}

func StartSnapshotRoutine(l *RemoteList) {
	go func() {
		for {
			time.Sleep(15 * time.Second)
			l.mu.Lock()
			listCopy := make(map[int][]int)
			for k, v := range l.list {
				listCopy[k] = make([]int, len(v))
				copy(listCopy[k], v)
			}
			currentLogId := l.logCounter
			SaveSnapshot(l, currentLogId)
			l.mu.Unlock()
			DeleteOldSnapshots()
		}
	}()
}
