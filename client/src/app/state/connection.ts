import { configureStore } from '@reduxjs/toolkit'
import socketSliceReducer from "./websocket/websocketslice"
import MessageSliceReducer from './websocket/messageslice'
  export const connection=()=>{
    return configureStore({
      reducer:{
        socket:socketSliceReducer,
        Message:MessageSliceReducer
      },
      middleware: (getDefaultMiddleware) =>
    getDefaultMiddleware({
      serializableCheck: false,
    }),
    })
  }
export type Store=ReturnType<typeof connection>
export type RootState=ReturnType<Store['getState']>
export type AppDispatch =Store['dispatch']