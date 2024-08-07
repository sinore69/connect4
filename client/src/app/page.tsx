"use client";

import { useRouter } from "next/navigation";
import React from "react";
import { roomIdValidator } from "./validator/roomid";
import {incorrectRoomId} from "./validator/incorrectroom";
import { create,join } from "./state/websocket/websocketslice";
import { useAppDispatch, useAppSelector } from "./lib/hooks";
import { useState } from "react";
import { replace } from "./state/websocket/messageslice";

function Page() {
  const router = useRouter();
  const dispatch=useAppDispatch()
  const connection=useAppSelector(state=>state.socket.socket)
  const message=useAppSelector((state)=>state.Message.Text)
  const[error,seterror]=useState("");

  function connect() {
    const socket = new WebSocket(`ws://${process.env.NEXT_PUBLIC_IP}:5000/create`);
    dispatch(create(socket))
    socket.onopen = (event) => {
      console.log("connection established");
    };
    socket.onmessage = (event) => {
      const roomId = JSON.parse(event.data);
      // console.log(roomId);
      if (roomIdValidator(roomId)) {
        // console.log(socket)
        dispatch(replace("Waiting for a player to join"))
        router.push(`/game/${roomId.Id}/t`);
      }
    };
  }

  async function handler(e: any) {
    e.preventDefault();
    const socket = new WebSocket(`ws://${process.env.NEXT_PUBLIC_IP}:5000/join`);
    dispatch(join(socket))
    socket.onopen = (event) => {
      const data = {
        id: Number(e.target.roomId.value),
      };
      socket.send(JSON.stringify(data));
    };
    socket.onmessage = (event) => {
      const res=JSON.parse(event.data)
      if(!incorrectRoomId(res)){
        router.push(`/game/${res.Id}/f`);
          // console.log(socket)
      }else{
        console.log(res.Message)
        seterror(res.Message)
      }
    };
  }
  
  return (
    <div className="h-screen w-sreen flex justify-center">
      <div className="flex flex-col pt-40 gap-y-10">
        <button onClick={() => connect()}>create room</button>
        <form onSubmit={handler}>
          <input type="text" id="roomId" placeholder="Room Id" />
          <button type="submit">join room</button>
        </form>
        <div>{error}</div>
      </div>
    </div>
  );
}

export default Page;
