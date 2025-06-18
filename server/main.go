package main

import (
	"encoding/json"
	"net/http"
	"sync"
)

type Board struct {
	mu   sync.Mutex
	grid [8][8]rune
}

func NewBoard() *Board {
	b := &Board{}
	setup := []string{
		"rnbqkbnr",
		"pppppppp",
		"........",
		"........",
		"........",
		"........",
		"PPPPPPPP",
		"RNBQKBNR",
	}
	for i := 0; i < 8; i++ {
		for j := 0; j < 8; j++ {
			b.grid[i][j] = rune(setup[i][j])
		}
	}
	return b
}

// square converts notation like e2 to array indices
func square(pos string) (int, int, bool) {
	if len(pos) != 2 {
		return 0, 0, false
	}
	file := int(pos[0] - 'a')
	rank := int('8' - pos[1])
	if file < 0 || file > 7 || rank < 0 || rank > 7 {
		return 0, 0, false
	}
	return rank, file, true
}

func (b *Board) Move(from, to string) bool {
	b.mu.Lock()
	defer b.mu.Unlock()
	fr, fc, ok1 := square(from)
	tr, tc, ok2 := square(to)
	if !ok1 || !ok2 {
		return false
	}
	piece := b.grid[fr][fc]
	if piece == '.' {
		return false
	}
	b.grid[tr][tc] = piece
	b.grid[fr][fc] = '.'
	return true
}

func (b *Board) State() []string {
	b.mu.Lock()
	defer b.mu.Unlock()
	rows := make([]string, 8)
	for i := 0; i < 8; i++ {
		row := make([]rune, 8)
		for j := 0; j < 8; j++ {
			row[j] = b.grid[i][j]
		}
		rows[i] = string(row)
	}
	return rows
}

func main() {
	board := NewBoard()

	http.HandleFunc("/board", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		json.NewEncoder(w).Encode(board.State())
	})

	http.HandleFunc("/move", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		if r.Method == http.MethodOptions {
			w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
			return
		}
		var req struct {
			From string `json:"from"`
			To   string `json:"to"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if board.Move(req.From, req.To) {
			json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
		} else {
			json.NewEncoder(w).Encode(map[string]string{"status": "invalid"})
		}
	})

	http.ListenAndServe(":8080", nil)
}
