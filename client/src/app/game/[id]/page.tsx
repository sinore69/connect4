"use client";
import { RootState } from "@/app/state/connection";
import React, { useEffect, useState } from "react";
import { useDispatch, useSelector } from "react-redux";
type board={
  Board:number[][],
  MoveCount:number,
LastMove:{
  RowIndex:number,
  ColIndex:number
}
}
function page({ params }: { params: { id: string } }) {
  const socket = useSelector((state: RootState) => state.socket.socket);
  const dispatch = useDispatch();
  const [disable,setdisable]=useState<boolean>(false)
  const [moveCount,setMoveCount]=useState<number>(0)
  const[board,setboard]=useState<number[][]>([
    [0, 0, 0, 0, 0, 0, 0, 0, 0, 0],
    [0, 0, 0, 0, 0, 0, 0, 0, 0, 0],
    [0, 0, 0, 0, 0, 0, 0, 0, 0, 0],
    [0, 0, 0, 0, 0, 0, 0, 0, 0, 0],
    [0, 0, 0, 0, 0, 0, 0, 0, 0, 0],
    [0, 0, 0, 0, 0, 0, 0, 0, 0, 0],
    [0, 0, 0, 0, 0, 0, 0, 0, 0, 0],
    [0, 0, 0, 0, 0, 0, 0, 0, 0, 0],
    [0, 0, 0, 0, 0, 0, 0, 0, 0, 0],
    [0, 0, 0, 0, 0, 0, 0, 0, 0, 0],
  ])
  useEffect(()=>{
    socket.onmessage=(event)=>{
      const res=JSON.parse(event.data) 
      const newboard=res.Board
      setboard([...newboard])
      setMoveCount(res.moveCount)
    }
  },[board])
  function handleClick(rowIndex: number, colIndex: number) {
    const data:board = {
      Board: board,
      MoveCount:moveCount,
      LastMove: {
        RowIndex: rowIndex,
        ColIndex: colIndex,
      },
    };
    socket.send(JSON.stringify(data));
  }
  return (
    <>
      <div className="h-screen w-screen pt-28 flex justify-center bg-whitw">
        <div>
          {board.map((row: number[], rowIndex: number) => (
            <div key={rowIndex} className="flex">
              {row.map((col: number, colIndex: number) => (
                <button
                disabled={disable || col===0?false:true}
                onClick={() => handleClick(rowIndex, colIndex)}
                  key={colIndex}
                  className="h-14 w-14 text-white border-2 flex justify-center pt-1.5"
                >
                  <div
                    className={`border-1 rounded-full h-10 w-10 ${
                      col === 0
                        ? "bg-white"
                        : col === 1
                        ? "bg-red-400"
                        : "bg-yellow-400"
                    }`}
                  ></div>
                </button>
              ))}
            </div>
          ))}
        </div>
      </div>
    </>
  );
}

export default page;
