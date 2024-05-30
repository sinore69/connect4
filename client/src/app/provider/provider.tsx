"use client"
import { connection } from "../state/connection"
import { Provider } from "react-redux"

export function ReduxProvider({children}:{children:React.ReactNode}){
    return <Provider store={connection}>{children}</Provider>
}