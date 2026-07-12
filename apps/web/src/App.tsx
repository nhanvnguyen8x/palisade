import React, { useMemo, useState } from 'react'
import { useAppDispatch, useAppSelector } from './store/hooks'
import {
  createKnowledgeBase,
  createOrganization,
  createWorkspace,
  setOrgId,
  setSelectedKnowledgeBase,
  setWorkspaceId,
} from './store/slices/knowledgeBaseSlice'
import { uploadDocumentToKnowledgeBase, fetchIngestionJob } from './store/slices/uploadSlice'
import { askAI } from './store/slices/chatSlice'

function TabButton({ active, onClick, children }: { active: boolean; onClick: () => void; children: React.ReactNode }) {
  return (
    <button
      onClick={onClick}
      style={{
        padding: '10px 14px',
        borderRadius: 8,
        border: '1px solid #ddd',
        background: active ? '#111827' : 'white',
        color: active ? 'white' : '#111827',
        cursor: 'pointer',
      }}
    >
      {children}
    </button>
  )
}

export function App() {
  const dispatch = useAppDispatch()
  const kbState = useAppSelector((s) => s.knowledgeBase)
  const uploadState = useAppSelector((s) => s.upload)
  const chatState = useAppSelector((s) => s.chat)

  const [tab, setTab] = useState<'upload' | 'chat'>('upload')

  const [orgName, setOrgName] = useState('')
  const [workspaceName, setWorkspaceName] = useState('')
  const [kbName, setKbName] = useState('')
  const [kbDescription, setKbDescription] = useState('')

  const [selectedKBIdOverride, setSelectedKBIdOverride] = useState<string>('')
  const effectiveKBId = useMemo(() => {
    return selectedKBIdOverride || kbState.selectedKnowledgeBase?.id || ''
  }, [selectedKBIdOverride, kbState.selectedKnowledgeBase?.id])

  const [uploadFile, setUploadFile] = useState<File | null>(null)

  const [chatQuestion, setChatQuestion] = useState('')

  const canUpload = !!effectiveKBId && !!uploadFile && !uploadState.uploading

  return (
    <div style={{ maxWidth: 1100, margin: '0 auto', padding: 20, fontFamily: 'system-ui, -apple-system, Segoe UI, Roboto, Arial' }}>
      <h1 style={{ marginTop: 0 }}>Palisade</h1>

      <div style={{ display: 'flex', gap: 10, marginBottom: 18 }}>
        <TabButton active={tab === 'upload'} onClick={() => setTab('upload')}>
          Upload
        </TabButton>
        <TabButton active={tab === 'chat'} onClick={() => setTab('chat')}>
          Chat
        </TabButton>
      </div>

      {tab === 'upload' && (
        <div style={{ display: 'grid', gridTemplateColumns: '1fr 1fr', gap: 18 }}>
          <div style={{ border: '1px solid #e5e7eb', borderRadius: 12, padding: 16 }}>
            <h3 style={{ marginTop: 0 }}>1) Create Organization</h3>
            <div style={{ display: 'flex', flexDirection: 'column', gap: 10 }}>
              <input placeholder="Organization name" value={orgName} onChange={(e) => setOrgName(e.target.value)} />
              <button
                onClick={async () => {
                  const r = await dispatch(createOrganization({ name: orgName })).unwrap()
                  dispatch(setOrgId(r.id))
                }}
                disabled={!orgName}
              >
                Create Org
              </button>
              {kbState.orgId && <div>Org ID: <code>{kbState.orgId}</code></div>}
            </div>
          </div>

          <div style={{ border: '1px solid #e5e7eb', borderRadius: 12, padding: 16 }}>
            <h3 style={{ marginTop: 0 }}>2) Create Workspace</h3>
            <div style={{ display: 'flex', flexDirection: 'column', gap: 10 }}>
              <input placeholder="Workspace name" value={workspaceName} onChange={(e) => setWorkspaceName(e.target.value)} disabled={!kbState.orgId} />
              <button
                onClick={async () => {
                  const r = await dispatch(
                    createWorkspace({ organizationId: kbState.orgId, name: workspaceName }),
                  ).unwrap()
                  dispatch(setWorkspaceId(r.id))
                }}
                disabled={!kbState.orgId || !workspaceName}
              >
                Create Workspace
              </button>
              {kbState.workspaceId && <div>Workspace ID: <code>{kbState.workspaceId}</code></div>}
            </div>
          </div>

          <div style={{ border: '1px solid #e5e7eb', borderRadius: 12, padding: 16, gridColumn: '1 / -1' }}>
            <h3 style={{ marginTop: 0 }}>3) Create Knowledge Base + Upload PDF</h3>
            <div style={{ display: 'grid', gridTemplateColumns: '1fr 1fr', gap: 12 }}>
              <div style={{ display: 'flex', flexDirection: 'column', gap: 8 }}>
                <input placeholder="Knowledge base name" value={kbName} onChange={(e) => setKbName(e.target.value)} disabled={!kbState.workspaceId} />
                <input placeholder="Description" value={kbDescription} onChange={(e) => setKbDescription(e.target.value)} disabled={!kbState.workspaceId} />
                <button
                  onClick={async () => {
                    const r = await dispatch(
                      createKnowledgeBase({
                        workspaceId: kbState.workspaceId,
                        name: kbName,
                        description: kbDescription,
                      }),
                    ).unwrap()
                    dispatch(setSelectedKnowledgeBase(r.id))
                    setSelectedKBIdOverride(r.id)
                  }}
                  disabled={!kbState.workspaceId || !kbName}
                >
                  Create KB
                </button>
              </div>

              <div style={{ display: 'flex', flexDirection: 'column', gap: 8 }}>
                <div>Selected KB ID</div>
                <input
                  placeholder="(optional) override KB id"
                  value={selectedKBIdOverride}
                  onChange={(e) => setSelectedKBIdOverride(e.target.value)}
                />

                <input
                  type="file"
                  accept="application/pdf"
                  onChange={(e) => {
                    const f = e.target.files?.[0] ?? null
                    setUploadFile(f)
                  }}
                />
                <button
                  disabled={!canUpload}
                  onClick={async () => {
                    if (!uploadFile) return
                    const r = await dispatch(
                      uploadDocumentToKnowledgeBase({
                        knowledgeBaseId: effectiveKBId,
                        file: uploadFile,
                      }),
                    ).unwrap()
                    // fetch job right away
                    dispatch(fetchIngestionJob({ ingestionJobId: r.ingestionJobId })).catch(() => {})
                  }}
                >
                  Upload PDF
                </button>

                {uploadState.error && <div style={{ color: '#b91c1c' }}>Error: {uploadState.error}</div>}
              </div>
            </div>

            <div style={{ marginTop: 14 }}>
              {uploadState.lastUpload && (
                <div>
                  Upload accepted. Ingestion job: <code>{uploadState.lastUpload.ingestionJobId}</code> status: <b>{uploadState.lastUpload.status}</b>
                </div>
              )}
              {uploadState.job && (
                <div style={{ marginTop: 10 }}>
                  Current job status: <b>{uploadState.job.status}</b>
                  {uploadState.job.errorMessage && <div style={{ color: '#b91c1c' }}>Error: {uploadState.job.errorMessage}</div>}
                </div>
              )}
              {uploadState.lastUpload && (
                <button
                  style={{ marginTop: 10 }}
                  onClick={() => dispatch(fetchIngestionJob({ ingestionJobId: uploadState.lastUpload!.ingestionJobId }))}
                >
                  Refresh Job
                </button>
              )}
            </div>
          </div>
        </div>
      )}

      {tab === 'chat' && (
        <div style={{ border: '1px solid #e5e7eb', borderRadius: 12, padding: 16 }}>
          <h3 style={{ marginTop: 0 }}>Chat (knowledge-base scoped)</h3>
          <div style={{ display: 'flex', flexDirection: 'column', gap: 10, maxWidth: 820 }}>
            <div>
              Knowledge Base ID: <code>{effectiveKBId || '-'}</code>
            </div>
            <input
              placeholder="Ask a question about uploaded documents..."
              value={chatQuestion}
              onChange={(e) => setChatQuestion(e.target.value)}
              disabled={!effectiveKBId}
            />
            <button
              disabled={!effectiveKBId || !chatQuestion || chatState.sending}
              onClick={async () => {
                await dispatch(
                  askAI({
                    knowledgeBaseId: effectiveKBId,
                    question: chatQuestion,
                    topK: 5,
                  }),
                ).unwrap()
              }}
            >
              {chatState.sending ? 'Thinking...' : 'Send'}
            </button>

            {chatState.error && <div style={{ color: '#b91c1c' }}>Error: {chatState.error}</div>}

            {chatState.response && (
              <div>
                <h4>Answer</h4>
                <pre style={{ whiteSpace: 'pre-wrap', background: '#f9fafb', padding: 12, borderRadius: 10 }}>{chatState.response.answer}</pre>

                <h4>Top contexts</h4>
                {chatState.response.contexts.slice(0, 3).map((c) => (
                  <div key={c.chunkID} style={{ marginBottom: 12, border: '1px solid #e5e7eb', borderRadius: 10, padding: 10 }}>
                    <div style={{ display: 'flex', justifyContent: 'space-between' }}>
                      <code>{c.chunkID}</code>
                      <span>score: {c.score.toFixed(4)}</span>
                    </div>
                    <pre style={{ whiteSpace: 'pre-wrap', marginTop: 8 }}>{c.text}</pre>
                  </div>
                ))}
              </div>
            )}
          </div>
        </div>
      )}
    </div>
  )
}

