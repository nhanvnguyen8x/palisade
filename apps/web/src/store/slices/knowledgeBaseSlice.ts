import { createAsyncThunk, createSlice } from '@reduxjs/toolkit'

export type KnowledgeBase = {
  id: string
  workspaceId: string
  name: string
  description: string
  createdAt: string
  updatedAt: string
}

type KBState = {
  orgId: string
  workspaceId: string
  selectedKnowledgeBase?: KnowledgeBase
  bases: KnowledgeBase[]
  loading: boolean
  error?: string
}

const initialState: KBState = {
  orgId: '',
  workspaceId: '',
  selectedKnowledgeBase: undefined,
  bases: [],
  loading: false,
  error: undefined,
}

export const createOrganization = createAsyncThunk<
  { id: string; name: string },
  { name: string }
>('kb/createOrganization', async ({ name }) => {
  const res = await fetch('/api/v1/organizations', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ name }),
  })
  if (!res.ok) throw new Error(`Create org failed: ${res.status}`)
  const json = await res.json()
  return json.data
})

export const createWorkspace = createAsyncThunk<
  { id: string; organizationId: string; name: string },
  { organizationId: string; name: string }
>('kb/createWorkspace', async ({ organizationId, name }) => {
  const res = await fetch(`/api/v1/organizations/${organizationId}/workspaces`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ name }),
  })
  if (!res.ok) throw new Error(`Create workspace failed: ${res.status}`)
  const json = await res.json()
  return json.data
})

export const createKnowledgeBase = createAsyncThunk<
  KnowledgeBase,
  { workspaceId: string; name: string; description: string }
>('kb/createKnowledgeBase', async ({ workspaceId, name, description }) => {
  const res = await fetch(`/api/v1/workspaces/${workspaceId}/knowledge-bases`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ name, description }),
  })
  if (!res.ok) throw new Error(`Create knowledge base failed: ${res.status}`)
  const json = await res.json()
  return json.data
})

export const listKnowledgeBasesByWorkspace = createAsyncThunk<
  KnowledgeBase[],
  { workspaceId: string }
>('kb/listKnowledgeBasesByWorkspace', async ({ workspaceId }) => {
  const res = await fetch(`/api/v1/workspaces/${workspaceId}/knowledge-bases`)
  if (!res.ok) throw new Error(`List KB failed: ${res.status}`)
  const json = await res.json()
  return json.data
})

const slice = createSlice({
  name: 'knowledgeBase',
  initialState,
  reducers: {
    setSelectedKnowledgeBase(state, action) {
      state.selectedKnowledgeBase = state.bases.find((b) => b.id === action.payload)
    },
    setOrgId(state, action) {
      state.orgId = action.payload
    },
    setWorkspaceId(state, action) {
      state.workspaceId = action.payload
    },
  },
  extraReducers: (builder) => {
    builder
      .addCase(listKnowledgeBasesByWorkspace.pending, (state) => {
        state.loading = true
        state.error = undefined
      })
      .addCase(listKnowledgeBasesByWorkspace.fulfilled, (state, action) => {
        state.loading = false
        state.bases = action.payload
        state.selectedKnowledgeBase = action.payload[0]
      })
      .addCase(listKnowledgeBasesByWorkspace.rejected, (state, action) => {
        state.loading = false
        state.error = action.error.message
      })
      .addCase(createKnowledgeBase.fulfilled, (state, action) => {
        state.bases = [action.payload, ...state.bases]
        state.selectedKnowledgeBase = action.payload
      })
  },
})

export const { setSelectedKnowledgeBase, setOrgId, setWorkspaceId } = slice.actions

export const knowledgeBaseReducer = slice.reducer

