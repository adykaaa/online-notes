import { createContext, useReducer } from "react";

export const NoteContext = createContext();

export const noteTitleReducer = (state,action) => {
    switch (action.type) {
        case 'TITLE_UPDATE':
            return {text: action.payload}
        default:
            return state
    }
}

export const noteTextReducer = (state,action) => {
    switch (action.type) {
        case 'TEXT_UPDATE':
            return {text: action.payload}
        default:
            return state
    }
}


export const NoteContextProvider =({ children })=> {
    const [title, dispatchTitle] = useReducer(noteTitleReducer, {
        title: ""
    })
    const [text, dispatchText] = useReducer(noteTextReducer, {
        text: ""
    })

    return (
        <NoteContext.Provider value={{ ...title, ...text,dispatchTitle,dispatchText}}>
            { children }
        </NoteContext.Provider>
    )
}