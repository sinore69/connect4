import { configureStore } from '@reduxjs/toolkit'
import socketSliceReducer from "./websocket/websocketslice"
export const connection= configureStore({
  reducer: {
    socket:socketSliceReducer
  },
  middleware: (getDefaultMiddleware) =>
    getDefaultMiddleware({
      serializableCheck: false,
    }),
  }) 
export type RootState=ReturnType<typeof connection.getState>
export type AppDispatch = typeof connection.dispatch