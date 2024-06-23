"use client";

import React, { useEffect, useState } from "react";
import { isInitialState } from "@/app/validator/gamestate";
import { useAppSelector, useAppDispatch } from "@/app/lib/hooks";
import { isMessage } from "@/app/validator/message";
import { replace } from "@/app/state/websocket/messageslice";
import { roomIdValidator } from "@/app/validator/roomid";
type board = {
  Board: number[][];
  MoveCount: number;
  Disable: boolean;
  LastMove: {
    RowIndex: number;
    ColIndex: number;
  };
};

function Page({ params }: { params: { id: string; value: string } }) {
  const socket = useAppSelector((state) => state.socket.socket);
  const message = useAppSelector((state) => state.Message.Text);
  const dispatch = useAppDispatch();
  const [disable, setdisable] = useState<boolean>(true);
  const [moveCount, setMoveCount] = useState<number>(0);
  const [banner, setbanner] = useState(params.value === "t" ? true : false);
  const [board, setboard] = useState<number[][]>([
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
  ]);
  useEffect(() => {
    socket.onmessage = (event) => {
      const res = JSON.parse(event.data);
      if (isInitialState(res)) {
        setdisable(res.Disable);
        setbanner(false);
        dispatch(replace(res.Disable ? "your opponents turn" : "your turn"));
      } else if (isMessage(res)) {
        dispatch(replace(res.Text));
      } else if (roomIdValidator(res)) {
        console.log(res.Id);
      } else {
        const newboard = res.Board;
        setboard([...newboard]);
        setMoveCount(res.moveCount);
        setdisable(res.Disable);
        dispatch(replace(res.Disable ? "your opponents turn" : "your turn"));
      }
    };
  }, [socket, board]);

  function handleClick(rowIndex: number, colIndex: number) {
    const data: board = {
      Board: board,
      MoveCount: moveCount,
      Disable: disable,
      LastMove: {
        RowIndex: rowIndex,
        ColIndex: colIndex,
      },
    };
    socket.send(JSON.stringify(data));
  }

  return (
    <>
      <div className="h-screen w-screen pt-28 flex justify-center bg-white">
        <div>
          {board.map((row: number[], rowIndex: number) => (
            <div key={rowIndex} className="flex">
              {row.map((col: number, colIndex: number) => (
                <button
                  disabled={disable}
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
        <div className="flex flex-col w-80">
          <div className="pl-32">{message}</div>
          {banner && (
            <div className="pl-32 pt-20">
              share this code:{params.id}
            </div>
          )}
        </div>
      </div>
    </>
  );
}

export default Page;
