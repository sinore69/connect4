"use client";
import React from "react";
function page() {
  function connect() {
    const socket=new WebSocket("ws://127.0.0.1:5000/echo")
    socket.onopen=(event)=>{
      const data={
        message:"Here's some text that the server is urgently awaiting!",
        num:1999
      }
      console.log(data)
      socket.send(JSON.stringify(data));
      
    }
  }
  return (
    <div className="h-screen w-sreen flex justify-center">
      <div className="flex flex-col pt-40 gap-y-10">
        <button onClick={() => connect()}>create session</button>
        <button>join session</button>
      </div>
    </div>
  );
}

export default page;
