import { configureStore } from '@reduxjs/toolkit'
import { knowledgeBaseReducer } from './slices/knowledgeBaseSlice'
import { uploadReducer } from './slices/uploadSlice'
import { chatReducer } from './slices/chatSlice'

export const store = configureStore({
  reducer: {
    knowledgeBase: knowledgeBaseReducer,
    upload: uploadReducer,
    chat: chatReducer,
  },
})

export type RootState = ReturnType<typeof store.getState>
export type AppDispatch = typeof store.dispatch

