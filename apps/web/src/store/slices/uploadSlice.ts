import { createAsyncThunk, createSlice } from '@reduxjs/toolkit'

export type UploadResult = {
  knowledgeSourceId: string
  ingestionJobId: string
  documentId: string
  status: string
  fileName: string
  createdAt: string
}

export type IngestionJob = {
  id: string
  knowledgeSourceId: string
  status: string
  errorMessage?: string
  startedAt?: string
  completedAt?: string
  createdAt: string
  updatedAt: string
}

type UploadState = {
  uploading: boolean
  error?: string
  lastUpload?: UploadResult
  job?: IngestionJob
}

const initialState: UploadState = {
  uploading: false,
  error: undefined,
  lastUpload: undefined,
  job: undefined,
}

export const uploadDocumentToKnowledgeBase = createAsyncThunk<
  UploadResult,
  { knowledgeBaseId: string; file: File }
>('upload/uploadDocumentToKnowledgeBase', async ({ knowledgeBaseId, file }) => {
  const form = new FormData()
  form.append('file', file)

  const res = await fetch(`/api/v1/knowledge-bases/${knowledgeBaseId}/sources`, {
    method: 'POST',
    body: form,
  })
  if (!res.ok) throw new Error(`Upload failed: ${res.status}`)
  const json = await res.json()
  return json.data
})

export const fetchIngestionJob = createAsyncThunk<
  IngestionJob,
  { ingestionJobId: string }
>('upload/fetchIngestionJob', async ({ ingestionJobId }) => {
  const res = await fetch(`/api/v1/ingestion-jobs/${ingestionJobId}`)
  if (!res.ok) throw new Error(`Fetch job failed: ${res.status}`)
  const json = await res.json()
  return json.data
})

const slice = createSlice({
  name: 'upload',
  initialState,
  reducers: {
    clearUpload(state) {
      state.lastUpload = undefined
      state.job = undefined
      state.error = undefined
      state.uploading = false
    },
  },
  extraReducers: (builder) => {
    builder
      .addCase(uploadDocumentToKnowledgeBase.pending, (state) => {
        state.uploading = true
        state.error = undefined
      })
      .addCase(uploadDocumentToKnowledgeBase.fulfilled, (state, action) => {
        state.uploading = false
        state.lastUpload = action.payload
      })
      .addCase(uploadDocumentToKnowledgeBase.rejected, (state, action) => {
        state.uploading = false
        state.error = action.error.message
      })
      .addCase(fetchIngestionJob.fulfilled, (state, action) => {
        state.job = action.payload
      })
  },
})

export const { clearUpload } = slice.actions
export const uploadReducer = slice.reducer

