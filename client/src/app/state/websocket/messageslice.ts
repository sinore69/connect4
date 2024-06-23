import { PayloadAction, createSlice } from "@reduxjs/toolkit"

interface Message{
    Text:string
}

const initialState:Message={
    Text:""
}

const MessageSlice=createSlice({
    name:"Message",
    initialState,
    reducers:{
replace:(state,actions:PayloadAction<string>)=>{
    state.Text=actions.payload
}
}
})
export const {replace}=MessageSlice.actions
export default MessageSlice.reducer
