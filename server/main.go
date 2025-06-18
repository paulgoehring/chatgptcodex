package main

import (
	"encoding/json"
	"net/http"
	"sync"
)

type Board struct {
	mu       sync.Mutex
	grid     [8][8]rune
	turn     int
	gameOver bool
}

func NewBoard() *Board {
	b := &Board{turn: 1, gameOver: false}
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

func color(piece rune) int {
	if piece >= 'A' && piece <= 'Z' {
		return 1 // white
	}
	if piece >= 'a' && piece <= 'z' {
		return -1 // black
	}
	return 0
}

func sign(x int) int {
	switch {
	case x > 0:
		return 1
	case x < 0:
		return -1
	}
	return 0
}

func (b *Board) Reset() {
	nb := NewBoard()
	b.mu.Lock()
	defer b.mu.Unlock()
	*b = *nb
}

func (b *Board) clearPath(fr, fc, tr, tc int) bool {
	dr := sign(tr - fr)
	dc := sign(tc - fc)
	r, c := fr+dr, fc+dc
	for r != tr || c != tc {
		if b.grid[r][c] != '.' {
			return false
		}
		r += dr
		c += dc
	}
	return true
}

func (b *Board) validMove(fr, fc, tr, tc int, piece rune) bool {
	if fr == tr && fc == tc {
		return false
	}
	dest := b.grid[tr][tc]
	if dest != '.' && color(dest) == color(piece) {
		return false
	}

	dr := tr - fr
	dc := tc - fc

	switch piece {
	case 'P':
		if dr == -1 && dc == 0 && dest == '.' {
			return true
		}
		if fr == 6 && dr == -2 && dc == 0 && dest == '.' && b.grid[fr-1][fc] == '.' {
			return true
		}
		if dr == -1 && (dc == -1 || dc == 1) && dest != '.' && color(dest) == -1 {
			return true
		}
	case 'p':
		if dr == 1 && dc == 0 && dest == '.' {
			return true
		}
		if fr == 1 && dr == 2 && dc == 0 && dest == '.' && b.grid[fr+1][fc] == '.' {
			return true
		}
		if dr == 1 && (dc == -1 || dc == 1) && dest != '.' && color(dest) == 1 {
			return true
		}
	case 'R', 'r':
		if dr == 0 || dc == 0 {
			if b.clearPath(fr, fc, tr, tc) {
				return true
			}
		}
	case 'B', 'b':
		if abs(dr) == abs(dc) {
			if b.clearPath(fr, fc, tr, tc) {
				return true
			}
		}
	case 'Q', 'q':
		if dr == 0 || dc == 0 || abs(dr) == abs(dc) {
			if b.clearPath(fr, fc, tr, tc) {
				return true
			}
		}
	case 'N', 'n':
		if (abs(dr) == 2 && abs(dc) == 1) || (abs(dr) == 1 && abs(dc) == 2) {
			return true
		}
	case 'K', 'k':
		if abs(dr) <= 1 && abs(dc) <= 1 {
			return true
		}
	}
	return false
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func (b *Board) Move(from, to string) bool {
	b.mu.Lock()
	defer b.mu.Unlock()
	if b.gameOver {
		return false
	}
	fr, fc, ok1 := square(from)
	tr, tc, ok2 := square(to)
	if !ok1 || !ok2 {
		return false
	}
	piece := b.grid[fr][fc]
	if piece == '.' {
		return false
	}
	if color(piece) != b.turn {
		return false
	}
	if !b.validMove(fr, fc, tr, tc, piece) {
		return false
	}
	dest := b.grid[tr][tc]
	b.grid[tr][tc] = piece
	b.grid[fr][fc] = '.'
	if dest == 'k' || dest == 'K' {
		b.gameOver = true
	} else {
		b.turn = -b.turn
	}
	return true
}

type State struct {
	Board    []string `json:"board"`
	GameOver bool     `json:"gameOver"`
}

func (b *Board) State() State {
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
	return State{Board: rows, GameOver: b.gameOver}
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
			json.NewEncoder(w).Encode(map[string]interface{}{"status": "ok", "gameOver": board.gameOver})
		} else {
			json.NewEncoder(w).Encode(map[string]interface{}{"status": "invalid", "gameOver": board.gameOver})
		}
	})

	http.HandleFunc("/newgame", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		if r.Method == http.MethodOptions {
			w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
			return
		}
		if r.Method == http.MethodPost {
			board.Reset()
			json.NewEncoder(w).Encode(board.State())
		}
	})

	http.ListenAndServe(":8080", nil)
}
