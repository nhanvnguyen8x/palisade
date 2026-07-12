import { createAsyncThunk, createSlice } from '@reduxjs/toolkit'

export type ChatContext = {
  chunkID: string
  text: string
  score: number
}

export type ChatResponse = {
  answer: string
  contexts: ChatContext[]
}

type ChatState = {
  sending: boolean
  error?: string
  response?: ChatResponse
}

const initialState: ChatState = {
  sending: false,
  error: undefined,
  response: undefined,
}

export const askAI = createAsyncThunk<
  ChatResponse,
  { knowledgeBaseId: string; question: string; topK?: number }
>('chat/askAI', async ({ knowledgeBaseId, question, topK }) => {
  const res = await fetch('/api/v1/chat', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ knowledgeBaseId, question, topK: topK ?? 5 }),
  })
  if (!res.ok) {
    const text = await res.text()
    throw new Error(`Chat failed: ${res.status} ${text}`)
  }
  const json = await res.json()
  return json.data
})

const slice = createSlice({
  name: 'chat',
  initialState,
  reducers: {
    clearChat(state) {
      state.response = undefined
      state.error = undefined
      state.sending = false
    },
  },
  extraReducers: (builder) => {
    builder
      .addCase(askAI.pending, (state) => {
        state.sending = true
        state.error = undefined
        state.response = undefined
      })
      .addCase(askAI.fulfilled, (state, action) => {
        state.sending = false
        state.response = action.payload
      })
      .addCase(askAI.rejected, (state, action) => {
        state.sending = false
        state.error = action.error.message
      })
  },
})

export const { clearChat } = slice.actions
export const chatReducer = slice.reducer

