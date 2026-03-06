import { useState, useEffect, useCallback } from 'react'
import {
  Settings, Database, FileSpreadsheet, Cloud, Server,
  CheckCircle2, AlertCircle, Loader2, Upload, Shield,
  Brain, Search, Mail, RefreshCw,
} from 'lucide-react'
import { api } from '../lib/api'
import type { AppSettings, SystemInfo } from '../lib/types'

type DataSourceType = 'json' | 'excel' | 'sharepoint' | 'cosmos' | 'postgres' | 'sqlserver'

const DATA_SOURCES: { type: DataSourceType; label: string; icon: typeof Database; desc: string }[] = [
  { type: 'json', label: 'JSON Files', icon: Server, desc: 'Default demo data — bundled in Docker image' },
  { type: 'excel', label: 'Excel Upload', icon: FileSpreadsheet, desc: 'Upload .xlsx rosters and spreadsheets' },
  { type: 'sharepoint', label: 'SharePoint', icon: Cloud, desc: 'Connect to SharePoint lists via Microsoft Graph' },
  { type: 'cosmos', label: 'Cosmos DB', icon: Database, desc: 'Azure Cosmos DB for cloud production' },
  { type: 'postgres', label: 'PostgreSQL', icon: Database, desc: 'Self-hosted PostgreSQL database' },
  { type: 'sqlserver', label: 'Azure SQL', icon: Database, desc: 'Azure SQL / SQL Server database' },
]

export function SettingsPage() {
  const [settings, setSettings] = useState<AppSettings | null>(null)
  const [systemInfo, setSystemInfo] = useState<SystemInfo | null>(null)
  const [loading, setLoading] = useState(true)
  const [saving, setSaving] = useState(false)
  const [saveMsg, setSaveMsg] = useState('')
  const [testResult, setTestResult] = useState<{ status: string; message: string } | null>(null)
  const [testing, setTesting] = useState(false)
  const [uploadStatus, setUploadStatus] = useState('')

  const loadData = useCallback(async () => {
    try {
      setLoading(true)
      const [s, info] = await Promise.all([
        api.getSettings(),
        api.getSystemInfo(),
      ])
      setSettings(s)
      setSystemInfo(info)
    } catch (err) {
      console.error('Failed to load settings:', err)
    } finally {
      setLoading(false)
    }
  }, [])

  useEffect(() => { loadData() }, [loadData])

  const handleSave = async () => {
    if (!settings) return
    try {
      setSaving(true)
      setSaveMsg('')
      const result = await api.updateSettings(settings)
      setSaveMsg(result.note || 'Settings saved')
    } catch (err) {
      setSaveMsg('Failed to save settings')
    } finally {
      setSaving(false)
    }
  }

  const handleTestConnection = async () => {
    if (!settings) return
    try {
      setTesting(true)
      setTestResult(null)
      const params: { type: string; connectionString?: string; tenantId?: string; clientId?: string; clientSecret?: string; siteUrl?: string } = { type: settings.dataSource.type }
      if (settings.dataSource.type === 'sharepoint') {
        params.tenantId = settings.dataSource.sharepoint.tenantId
        params.clientId = settings.dataSource.sharepoint.clientId
        params.clientSecret = settings.dataSource.sharepoint.clientSecret
        params.siteUrl = settings.dataSource.sharepoint.siteUrl
      } else if (['cosmos', 'postgres', 'sqlserver'].includes(settings.dataSource.type)) {
        params.connectionString = settings.dataSource.database.connectionString
      }
      const result = await api.testConnection(params)
      setTestResult(result)
    } catch {
      setTestResult({ status: 'error', message: 'Connection test failed' })
    } finally {
      setTesting(false)
    }
  }

  const handleFileUpload = async (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0]
    if (!file) return
    try {
      setUploadStatus('Uploading...')
      const result = await api.uploadFile(file)
      setUploadStatus(`Uploaded: ${result.filename} (${(result.size / 1024).toFixed(1)} KB)`)
      if (settings) {
        setSettings({
          ...settings,
          dataSource: { ...settings.dataSource, type: 'excel', excelPath: result.path },
        })
      }
    } catch {
      setUploadStatus('Upload failed')
    }
  }

  const updateDataSourceType = (type: DataSourceType) => {
    if (!settings) return
    setSettings({
      ...settings,
      dataSource: { ...settings.dataSource, type },
    })
    setTestResult(null)
  }

  if (loading) {
    return (
      <div className="flex items-center justify-center h-64">
        <Loader2 className="w-8 h-8 animate-spin text-slate-400" />
      </div>
    )
  }

  if (!settings || !systemInfo) {
    return (
      <div className="text-center py-12 text-slate-500">
        Failed to load settings. Ensure you are logged in as XO or Staff.
      </div>
    )
  }

  return (
    <div className="max-w-4xl mx-auto space-y-6">
      {/* Header */}
      <div className="flex items-center justify-between">
        <div className="flex items-center gap-3">
          <Settings className="w-6 h-6 text-slate-700" />
          <h1 className="text-2xl font-bold text-slate-900">Settings</h1>
        </div>
        <button
          onClick={handleSave}
          disabled={saving}
          className="px-4 py-2 bg-[var(--color-navy)] text-white rounded-lg hover:opacity-90 disabled:opacity-50 flex items-center gap-2"
        >
          {saving ? <Loader2 className="w-4 h-4 animate-spin" /> : null}
          Save Changes
        </button>
      </div>

      {saveMsg && (
        <div className={`p-3 rounded-lg text-sm ${saveMsg.includes('Failed') ? 'bg-red-50 text-red-700' : 'bg-green-50 text-green-700'}`}>
          {saveMsg}
        </div>
      )}

      {/* System Info */}
      <section className="bg-white rounded-xl border border-slate-200 p-6">
        <h2 className="text-lg font-semibold text-slate-900 mb-4 flex items-center gap-2">
          <Shield className="w-5 h-5" />
          System Information
        </h2>
        <div className="grid grid-cols-2 sm:grid-cols-4 gap-4">
          <InfoCard label="Version" value={systemInfo.version} />
          <InfoCard label="Auth Mode" value={systemInfo.authMode === 'demo' ? 'Demo (Role Picker)' : 'CAC/PKI'} />
          <InfoCard label="Students Loaded" value={String(systemInfo.studentCount)} />
          <InfoCard label="Data Source" value={systemInfo.dataSource.toUpperCase()} />
        </div>
      </section>

      {/* AI Configuration */}
      <section className="bg-white rounded-xl border border-slate-200 p-6">
        <h2 className="text-lg font-semibold text-slate-900 mb-4 flex items-center gap-2">
          <Brain className="w-5 h-5" />
          AI Configuration
        </h2>
        <div className="grid grid-cols-1 sm:grid-cols-3 gap-4">
          <div>
            <div className="text-xs text-slate-500 uppercase tracking-wider mb-1">Status</div>
            <div className="flex items-center gap-2">
              {systemInfo.ai.status === 'not configured' ? (
                <AlertCircle className="w-4 h-4 text-yellow-500" />
              ) : (
                <CheckCircle2 className="w-4 h-4 text-green-500" />
              )}
              <span className="text-sm font-medium">{systemInfo.ai.status}</span>
            </div>
          </div>
          <div>
            <div className="text-xs text-slate-500 uppercase tracking-wider mb-1">Model</div>
            <div className="text-sm font-medium">{systemInfo.ai.model}</div>
          </div>
          <div>
            <div className="text-xs text-slate-500 uppercase tracking-wider mb-1">API Key</div>
            <div className="text-sm font-medium font-mono">{systemInfo.ai.keyHint || 'Not set'}</div>
          </div>
        </div>

        <div className="mt-4 pt-4 border-t border-slate-100">
          <div className="flex items-center gap-2 mb-2">
            <Search className="w-4 h-4 text-slate-400" />
            <span className="text-sm font-medium text-slate-700">Web Search (SearXNG)</span>
          </div>
          <div className="flex items-center gap-3">
            <input
              type="text"
              value={settings.ai.searxngUrl}
              onChange={e => setSettings({ ...settings, ai: { ...settings.ai, searxngUrl: e.target.value } })}
              className="flex-1 px-3 py-2 border border-slate-300 rounded-lg text-sm"
              placeholder="http://localhost:8888"
            />
            <StatusDot status={systemInfo.searxng.status === 'configured' ? 'ok' : 'error'} />
          </div>
        </div>
      </section>

      {/* Data Source */}
      <section className="bg-white rounded-xl border border-slate-200 p-6">
        <h2 className="text-lg font-semibold text-slate-900 mb-4 flex items-center gap-2">
          <Database className="w-5 h-5" />
          Data Source
        </h2>

        {/* Source selector cards */}
        <div className="grid grid-cols-2 sm:grid-cols-3 gap-3 mb-6">
          {DATA_SOURCES.map(ds => (
            <button
              key={ds.type}
              onClick={() => updateDataSourceType(ds.type)}
              className={`p-3 rounded-lg border text-left transition-all ${
                settings.dataSource.type === ds.type
                  ? 'border-[var(--color-navy)] bg-blue-50 ring-1 ring-[var(--color-navy)]'
                  : 'border-slate-200 hover:border-slate-300 hover:bg-slate-50'
              }`}
            >
              <ds.icon className={`w-5 h-5 mb-1 ${settings.dataSource.type === ds.type ? 'text-[var(--color-navy)]' : 'text-slate-400'}`} />
              <div className="text-sm font-medium">{ds.label}</div>
              <div className="text-xs text-slate-500 mt-0.5">{ds.desc}</div>
            </button>
          ))}
        </div>

        {/* Source-specific config */}
        {settings.dataSource.type === 'excel' && (
          <div className="space-y-3">
            <div className="border-2 border-dashed border-slate-300 rounded-lg p-6 text-center">
              <Upload className="w-8 h-8 text-slate-400 mx-auto mb-2" />
              <p className="text-sm text-slate-600 mb-2">Drop your Excel (.xlsx) or CSV file here</p>
              <input
                type="file"
                accept=".xlsx,.csv"
                onChange={handleFileUpload}
                className="text-sm"
              />
              {uploadStatus && <p className="text-xs mt-2 text-slate-500">{uploadStatus}</p>}
            </div>
            {settings.dataSource.excelPath && (
              <div className="text-sm text-slate-600">
                Active file: <code className="bg-slate-100 px-1 rounded">{settings.dataSource.excelPath}</code>
              </div>
            )}
          </div>
        )}

        {settings.dataSource.type === 'sharepoint' && (
          <div className="space-y-3">
            <div className="grid grid-cols-2 gap-3">
              <SettingsInput
                label="Tenant ID"
                value={settings.dataSource.sharepoint.tenantId}
                onChange={v => setSettings({
                  ...settings,
                  dataSource: { ...settings.dataSource, sharepoint: { ...settings.dataSource.sharepoint, tenantId: v } },
                })}
              />
              <SettingsInput
                label="Client ID"
                value={settings.dataSource.sharepoint.clientId}
                onChange={v => setSettings({
                  ...settings,
                  dataSource: { ...settings.dataSource, sharepoint: { ...settings.dataSource.sharepoint, clientId: v } },
                })}
              />
            </div>
            <SettingsInput
              label="Client Secret"
              value={settings.dataSource.sharepoint.clientSecret}
              onChange={v => setSettings({
                ...settings,
                dataSource: { ...settings.dataSource, sharepoint: { ...settings.dataSource.sharepoint, clientSecret: v } },
              })}
              type="password"
            />
            <SettingsInput
              label="Site URL"
              value={settings.dataSource.sharepoint.siteUrl}
              onChange={v => setSettings({
                ...settings,
                dataSource: { ...settings.dataSource, sharepoint: { ...settings.dataSource.sharepoint, siteUrl: v } },
              })}
              placeholder="https://yourtenant.sharepoint.com/sites/TBS"
            />
            <div>
              <label className="text-xs text-slate-500 uppercase tracking-wider mb-1 block">Cloud Environment</label>
              <select
                value={settings.dataSource.sharepoint.cloud}
                onChange={e => setSettings({
                  ...settings,
                  dataSource: { ...settings.dataSource, sharepoint: { ...settings.dataSource.sharepoint, cloud: e.target.value } },
                })}
                className="px-3 py-2 border border-slate-300 rounded-lg text-sm w-full"
              >
                <option value="commercial">Commercial</option>
                <option value="gcc-high">GCC High</option>
                <option value="dod">DoD</option>
              </select>
            </div>
          </div>
        )}

        {['cosmos', 'postgres', 'sqlserver'].includes(settings.dataSource.type) && (
          <div className="space-y-3">
            {settings.dataSource.type !== 'cosmos' && (
              <div>
                <label className="text-xs text-slate-500 uppercase tracking-wider mb-1 block">Database Type</label>
                <select
                  value={settings.dataSource.database.type || settings.dataSource.type}
                  onChange={e => setSettings({
                    ...settings,
                    dataSource: { ...settings.dataSource, type: e.target.value as DataSourceType, database: { ...settings.dataSource.database, type: e.target.value } },
                  })}
                  className="px-3 py-2 border border-slate-300 rounded-lg text-sm w-full"
                >
                  <option value="postgres">PostgreSQL</option>
                  <option value="sqlserver">Azure SQL / SQL Server</option>
                  <option value="cosmos">Azure Cosmos DB</option>
                </select>
              </div>
            )}
            <SettingsInput
              label="Connection String"
              value={settings.dataSource.database.connectionString}
              onChange={v => setSettings({
                ...settings,
                dataSource: { ...settings.dataSource, database: { ...settings.dataSource.database, connectionString: v } },
              })}
              type="password"
              placeholder="Host=...; Database=...; User=...; Password=..."
            />
          </div>
        )}

        {/* Test Connection button */}
        {settings.dataSource.type !== 'json' && (
          <div className="mt-4 flex items-center gap-3">
            <button
              onClick={handleTestConnection}
              disabled={testing}
              className="px-4 py-2 border border-slate-300 rounded-lg text-sm hover:bg-slate-50 flex items-center gap-2"
            >
              {testing ? <Loader2 className="w-4 h-4 animate-spin" /> : <RefreshCw className="w-4 h-4" />}
              Test Connection
            </button>
            {testResult && (
              <div className={`flex items-center gap-1.5 text-sm ${testResult.status === 'ok' ? 'text-green-600' : testResult.status === 'error' ? 'text-red-600' : 'text-yellow-600'}`}>
                {testResult.status === 'ok' ? <CheckCircle2 className="w-4 h-4" /> : <AlertCircle className="w-4 h-4" />}
                {testResult.message}
              </div>
            )}
          </div>
        )}
      </section>

      {/* Outlook Integration */}
      <section className="bg-white rounded-xl border border-slate-200 p-6">
        <h2 className="text-lg font-semibold text-slate-900 mb-4 flex items-center gap-2">
          <Mail className="w-5 h-5" />
          Outlook Integration
        </h2>

        <div className="flex items-center gap-3 mb-4">
          <label className="flex items-center gap-2 cursor-pointer">
            <input
              type="checkbox"
              checked={settings.outlook.enabled}
              onChange={e => setSettings({ ...settings, outlook: { ...settings.outlook, enabled: e.target.checked } })}
              className="w-4 h-4 rounded border-slate-300"
            />
            <span className="text-sm font-medium">Enable Outlook mail and calendar sync</span>
          </label>
        </div>

        {settings.outlook.enabled && (
          <div className="space-y-3">
            <div className="grid grid-cols-2 gap-3">
              <SettingsInput
                label="Tenant ID"
                value={settings.outlook.tenantId}
                onChange={v => setSettings({ ...settings, outlook: { ...settings.outlook, tenantId: v } })}
              />
              <SettingsInput
                label="Client ID"
                value={settings.outlook.clientId}
                onChange={v => setSettings({ ...settings, outlook: { ...settings.outlook, clientId: v } })}
              />
            </div>
            <SettingsInput
              label="Client Secret"
              value={settings.outlook.clientSecret}
              onChange={v => setSettings({ ...settings, outlook: { ...settings.outlook, clientSecret: v } })}
              type="password"
            />
            <div className="grid grid-cols-2 gap-3">
              <div>
                <label className="text-xs text-slate-500 uppercase tracking-wider mb-1 block">Cloud Environment</label>
                <select
                  value={settings.outlook.cloud}
                  onChange={e => setSettings({ ...settings, outlook: { ...settings.outlook, cloud: e.target.value } })}
                  className="px-3 py-2 border border-slate-300 rounded-lg text-sm w-full"
                >
                  <option value="commercial">Commercial</option>
                  <option value="gcc-high">GCC High</option>
                  <option value="dod">DoD</option>
                </select>
              </div>
              <SettingsInput
                label="Sync Interval (minutes)"
                value={String(settings.outlook.syncIntervalMinutes)}
                onChange={v => setSettings({ ...settings, outlook: { ...settings.outlook, syncIntervalMinutes: parseInt(v) || 5 } })}
                type="number"
              />
            </div>
            <SettingsInput
              label="Master Calendar ID (optional)"
              value={settings.outlook.masterCalendarId}
              onChange={v => setSettings({ ...settings, outlook: { ...settings.outlook, masterCalendarId: v } })}
              placeholder="For TBS-wide shared calendar events"
            />
          </div>
        )}
      </section>
    </div>
  )
}

// ---- Subcomponents ----

function InfoCard({ label, value }: { label: string; value: string }) {
  return (
    <div className="bg-slate-50 rounded-lg p-3">
      <div className="text-xs text-slate-500 uppercase tracking-wider">{label}</div>
      <div className="text-sm font-semibold text-slate-900 mt-1">{value}</div>
    </div>
  )
}

function SettingsInput({ label, value, onChange, type = 'text', placeholder }: {
  label: string
  value: string
  onChange: (v: string) => void
  type?: string
  placeholder?: string
}) {
  return (
    <div>
      <label className="text-xs text-slate-500 uppercase tracking-wider mb-1 block">{label}</label>
      <input
        type={type}
        value={value}
        onChange={e => onChange(e.target.value)}
        placeholder={placeholder}
        className="w-full px-3 py-2 border border-slate-300 rounded-lg text-sm"
      />
    </div>
  )
}

function StatusDot({ status }: { status: 'ok' | 'error' | 'pending' }) {
  const colors = {
    ok: 'bg-green-500',
    error: 'bg-red-500',
    pending: 'bg-yellow-500',
  }
  return <div className={`w-2.5 h-2.5 rounded-full ${colors[status]}`} />
}
