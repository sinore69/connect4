export default function Home() {
  const board = [
    [0, 1, 0, 0, 0, 0, 0, 0, 0, 0],
    [0, 0, 0, 0, 0, 0, 0, 0, 0, 0],
    [0, 0, 0, 0, 0, 0, 0, 0, 0, 0],
    [0, 0, 0, 0, 0, 0, 0, 0, 0, 0],
    [0, 0, 0, 2, 0, 2, 0, 0, 0, 0],
    [0, 0, 0, 0, 0, 0, 0, 0, 0, 0],
    [0, 0, 0, 0, 0, 0, 0, 0, 0, 0],
    [0, 0, 0, 0, 0, 0, 0, 0, 0, 0],
    [0, 0, 0, 0, 0, 0, 0, 0, 0, 0],
  ];
  return (
    <>
      <div className="h-screen w-screen pt-28 flex justify-center bg-whitw">
        <div>
          {board.map((row: number[],rowIndex:number) => (
            <div key={rowIndex} className="flex">
              {row.map((col: number,colIndex:number) => (
                <div key={colIndex} className="h-14 w-14 text-white border-2 flex justify-center pt-1.5">
                  <div
                    className={`border-1 rounded-full h-10 w-10 ${
                      col === 0
                        ? "bg-white"
                        : col === 1
                        ? "bg-red-400"
                        : "bg-yellow-400"
                    }`}
                  ></div>
                </div>
              ))}
            </div>
          ))}
        </div>
      </div>
    </>
  );
}
