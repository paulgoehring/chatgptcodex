const e = React.createElement;

function Board() {
  const [board, setBoard] = React.useState([]);
  const [gameOver, setGameOver] = React.useState(false);
  const [dragFrom, setDragFrom] = React.useState(null);

  const pieceIcons = {
    K: "\u2654",
    Q: "\u2655",
    R: "\u2656",
    B: "\u2657",
    N: "\u2658",
    P: "\u2659",
    k: "\u265A",
    q: "\u265B",
    r: "\u265C",
    b: "\u265D",
    n: "\u265E",
    p: "\u265F",
  };

  const pieceClass = (p) =>
    p >= "A" && p <= "Z" ? "white-piece" : "black-piece";

  const newGame = () => {
    fetch("http://localhost:8080/newgame", { method: "POST" }).then(fetchBoard);
  };

  const fetchBoard = () => {
    fetch("http://localhost:8080/board")
      .then((res) => res.json())
      .then((data) => {
        setBoard(data.board);
        if (data.gameOver && !gameOver) {
          if (confirm("Game over. Start new game?")) {
            newGame();
          } else {
            setGameOver(true);
          }
        } else {
          setGameOver(data.gameOver);
        }
      });
  };

  React.useEffect(fetchBoard, []);

  const onDragStart = (r, c) => () => {
    setDragFrom({ r, c });
  };

  const onDrop = (r, c) => () => {
    if (!dragFrom) return;
    const from =
      String.fromCharCode("a".charCodeAt(0) + dragFrom.c) + (8 - dragFrom.r);
    const to = String.fromCharCode("a".charCodeAt(0) + c) + (8 - r);
    fetch("http://localhost:8080/move", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ from, to }),
    })
      .then((res) => res.json())
      .then((data) => {
        fetchBoard();
        if (data.gameOver && !gameOver) {
          if (confirm("Game over. Start new game?")) {
            newGame();
          } else {
            setGameOver(true);
          }
        }
      });
    setDragFrom(null);
  };

  const allowDrop = (e) => e.preventDefault();

  return e(
    "div",
    { className: "board" },
    board.map((row, r) =>
      row.split("").map((cell, c) =>
        e(
          "div",
          {
            key: r + "-" + c,
            className: "cell " + ((r + c) % 2 === 0 ? "white" : "black"),
            onDragOver: allowDrop,
            onDrop: onDrop(r, c),
          },
          cell === "."
            ? ""
            : e(
                "span",
                {
                  draggable: true,
                  className: pieceClass(cell),
                  onDragStart: onDragStart(r, c),
                },
                pieceIcons[cell] || cell,
              ),
        ),
      ),
    ),
  );
}

ReactDOM.createRoot(document.getElementById("root")).render(e(Board));
