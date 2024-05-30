"use client"
import { RootState } from "@/app/state/connection";
import React from "react";
import { UseSelector, useDispatch, useSelector } from "react-redux";
function page({ params }: { params: { id: string } }) {
  const socket = useSelector((state: RootState) => state.socket.socket);
  const dispatch = useDispatch();
  const board = [
    [0, 0, 0, 0, 0, 0, 0, 0, 0, 0],
    [0, 0, 0, 0, 0, 0, 0, 0, 0, 0],
    [0, 0, 0, 0, 0, 0, 0, 0, 0, 0],
    [0, 0, 0, 0, 0, 0, 0, 0, 0, 0],
    [0, 0, 0, 0, 0, 0, 0, 0, 0, 0],
    [0, 0, 0, 0, 0, 0, 0, 0, 0, 0],
    [0, 0, 0, 0, 0, 0, 0, 0, 0, 0],
    [0, 0, 0, 0, 0, 0, 0, 0, 0, 0],
    [0, 0, 0, 0, 0, 0, 0, 0, 0, 0],
  ];
  console.log(socket)
  return (
    <>
      <div className="h-screen w-screen pt-28 flex justify-center bg-whitw">
        <div>
          {board.map((row: number[], rowIndex: number) => (
            <div key={rowIndex} className="flex">
              {row.map((col: number, colIndex: number) => (
                <div
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
                </div>
              ))}
            </div>
          ))}
        </div>
      </div>
    </>
  );
}

export default page;
