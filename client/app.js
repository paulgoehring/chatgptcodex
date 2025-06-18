const e = React.createElement;

function Board() {
  const [board, setBoard] = React.useState([]);
  const [selection, setSelection] = React.useState(null);

  const fetchBoard = () => {
    fetch('http://localhost:8080/board')
      .then(res => res.json())
      .then(setBoard);
  };

  React.useEffect(fetchBoard, []);

  const handleClick = (row, col) => {
    if (!selection) {
      setSelection({ row, col });
      return;
    }
    const from = String.fromCharCode('a'.charCodeAt(0) + selection.col) + (8 - selection.row);
    const to = String.fromCharCode('a'.charCodeAt(0) + col) + (8 - row);
    fetch('http://localhost:8080/move', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ from, to })
    }).then(fetchBoard);
    setSelection(null);
  };

  return e('div', { className: 'board' },
    board.map((row, r) =>
      row.split('').map((cell, c) =>
        e('div', {
          key: r + '-' + c,
          className: 'cell ' + ((r + c) % 2 === 0 ? 'white' : 'black'),
          onClick: () => handleClick(r, c)
        }, cell === '.' ? '' : cell)
      )
    )
  );
}

ReactDOM.createRoot(document.getElementById('root')).render(e(Board));
