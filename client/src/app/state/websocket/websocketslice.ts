import { PayloadAction, createSlice } from "@reduxjs/toolkit"

interface Socket{
    socket:WebSocket
}

const initialState:Socket={
    socket:WebSocket as unknown as WebSocket
}

const SocketSlice=createSlice({
    name:"socket",
    initialState,
    reducers:{
        create:(state,actions:PayloadAction<WebSocket>)=>{
            state.socket=actions.payload
        },
        join:(state,actions:PayloadAction<WebSocket>)=>{
            state.socket=actions.payload
        }
    }
})
export const {create,join}=SocketSlice.actions
export default SocketSlice.reducer
