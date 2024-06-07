"use client";

import { useRouter } from "next/navigation";
import React from "react";
import { roomIdValidator } from "./validator/roomid";
import {incorrectRoomId} from "./validator/incorrectroom";
import { create,join } from "./state/websocket/websocketslice";
import { useAppDispatch, useAppSelector, useAppStore } from "./lib/hooks";

function Page() {
  const router = useRouter();
  const dispatch=useAppDispatch()
  const connection=useAppSelector(state=>state.socket.socket)

  function connect() {
    const socket = new WebSocket("ws://127.0.0.1:5000/create");
    dispatch(create(socket))
    socket.onopen = (event) => {
      console.log("connection established");
    };
    socket.onmessage = (event) => {
      const roomId = JSON.parse(event.data);
      // console.log(roomId);
      if (roomIdValidator(roomId)) {
        // console.log(socket)
        router.push(`/game/${roomId.Id}`);
      }
    };
  }

  async function handler(e: any) {
    e.preventDefault();
    const socket = new WebSocket("ws://127.0.0.1:5000/join");
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
        router.push(`/game/${res.Id}`);
          // console.log(socket)
      }else{
        console.log("incorrect room id")
      }
    };
  }
  
  return (
    <div className="h-screen w-sreen flex justify-center">
      <div className="flex flex-col pt-40 gap-y-10">
        <button onClick={() => connect()}>create session</button>
        <form onSubmit={handler}>
          <input type="text" id="roomId" placeholder="Room Id" />
          <button type="submit">join session</button>
        </form>
      </div>
    </div>
  );
}

export default Page;
