import { useState, useEffect, useCallback } from 'react'
import {
  Settings, Database, FileSpreadsheet, Cloud, Server,
  CheckCircle2, AlertCircle, Loader2, Upload, Shield,
  Brain, Search, Mail, RefreshCw, HelpCircle, ChevronDown,
  ChevronRight, Info, Zap, BookOpen,
} from 'lucide-react'
import { api } from '../lib/api'
import type { AppSettings, SystemInfo } from '../lib/types'

type DataSourceType = 'json' | 'excel' | 'sharepoint' | 'cosmos' | 'postgres' | 'sqlserver'

export function SettingsPage() {
  const [settings, setSettings] = useState<AppSettings | null>(null)
  const [systemInfo, setSystemInfo] = useState<SystemInfo | null>(null)
  const [loading, setLoading] = useState(true)
  const [saving, setSaving] = useState(false)
  const [saveMsg, setSaveMsg] = useState('')
  const [testResult, setTestResult] = useState<{ status: string; message: string } | null>(null)
  const [testing, setTesting] = useState(false)
  const [uploadStatus, setUploadStatus] = useState('')
  const [showAdvancedData, setShowAdvancedData] = useState(false)
  const [showAdvancedOutlook, setShowAdvancedOutlook] = useState(false)
  const [expandedHelp, setExpandedHelp] = useState<string | null>(null)

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
      setSaveMsg(result.note || 'Settings saved successfully')
    } catch {
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

  const toggleHelp = (id: string) => {
    setExpandedHelp(expandedHelp === id ? null : id)
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

  // Calculate setup progress
  const steps = [
    { label: 'Student data loaded', done: systemInfo.studentCount > 0 },
    { label: 'AI assistant configured', done: systemInfo.ai.status !== 'not configured' },
    { label: 'Outlook connected', done: settings.outlook.enabled && !!settings.outlook.tenantId },
  ]
  const doneCount = steps.filter(s => s.done).length

  return (
    <div className="max-w-4xl mx-auto space-y-6">
      {/* Header */}
      <div className="flex items-center justify-between">
        <div className="flex items-center gap-3">
          <Settings className="w-6 h-6 text-slate-700" />
          <div>
            <h1 className="text-2xl font-bold text-slate-900">Setup & Settings</h1>
            <p className="text-sm text-slate-500">Configure Heywood for your unit</p>
          </div>
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

      {/* Setup Progress */}
      <section className="bg-gradient-to-r from-[var(--color-navy)] to-slate-700 rounded-xl p-6 text-white">
        <div className="flex items-center gap-3 mb-4">
          <Zap className="w-5 h-5" />
          <h2 className="text-lg font-semibold">Setup Progress</h2>
          <span className="ml-auto text-sm opacity-80">{doneCount} of {steps.length} complete</span>
        </div>
        <div className="flex gap-2 mb-3">
          {steps.map((_, i) => (
            <div key={i} className={`flex-1 h-2 rounded-full ${i < doneCount ? 'bg-green-400' : 'bg-white/20'}`} />
          ))}
        </div>
        <div className="space-y-2">
          {steps.map((step, i) => (
            <div key={i} className="flex items-center gap-2 text-sm">
              {step.done ? (
                <CheckCircle2 className="w-4 h-4 text-green-400 flex-shrink-0" />
              ) : (
                <div className="w-4 h-4 rounded-full border-2 border-white/40 flex-shrink-0" />
              )}
              <span className={step.done ? 'opacity-70 line-through' : ''}>{step.label}</span>
            </div>
          ))}
        </div>
      </section>

      {/* System Status - Compact */}
      <section className="bg-white rounded-xl border border-slate-200 p-5">
        <div className="flex items-center gap-2 mb-3">
          <Shield className="w-4 h-4 text-slate-500" />
          <h2 className="text-sm font-semibold text-slate-700 uppercase tracking-wider">Current Status</h2>
        </div>
        <div className="grid grid-cols-2 sm:grid-cols-4 gap-3">
          <StatusCard label="Students" value={String(systemInfo.studentCount)} status={systemInfo.studentCount > 0 ? 'ok' : 'warn'} />
          <StatusCard label="Data Source" value={systemInfo.dataSource === 'json' ? 'Demo Data' : systemInfo.dataSource.toUpperCase()} status="ok" />
          <StatusCard label="AI" value={systemInfo.ai.status === 'not configured' ? 'Not Set Up' : 'Active'} status={systemInfo.ai.status === 'not configured' ? 'warn' : 'ok'} />
          <StatusCard label="Auth" value={systemInfo.authMode === 'demo' ? 'Demo Mode' : 'CAC/PKI'} status="ok" />
        </div>
      </section>

      {/* STEP 1: Student Data */}
      <section className="bg-white rounded-xl border border-slate-200 overflow-hidden">
        <div className="px-6 py-4 bg-slate-50 border-b border-slate-200">
          <div className="flex items-center gap-3">
            <div className="w-7 h-7 rounded-full bg-[var(--color-navy)] text-white flex items-center justify-center text-sm font-bold">1</div>
            <div>
              <h2 className="text-base font-semibold text-slate-900">Student Data</h2>
              <p className="text-xs text-slate-500">Where does Heywood get your student roster?</p>
            </div>
            <HelpButton id="data-help" expanded={expandedHelp} onToggle={toggleHelp} />
          </div>
          {expandedHelp === 'data-help' && (
            <div className="mt-3 p-3 bg-blue-50 rounded-lg text-xs text-blue-800 leading-relaxed">
              <BookOpen className="w-4 h-4 inline mr-1" />
              <strong>How this works:</strong> Heywood needs your student roster to track performance, flag at-risk students, and support counseling.
              Most units start by uploading their existing Excel spreadsheet — Heywood will read the column headers and map them automatically.
              You can also connect directly to SharePoint or a database if your unit uses those.
            </div>
          )}
        </div>

        <div className="p-6">
          {/* Simple choices */}
          <div className="grid grid-cols-1 sm:grid-cols-2 gap-3 mb-4">
            <DataSourceCard
              icon={Server}
              title="Demo Data"
              desc="200 sample students for testing and evaluation"
              selected={settings.dataSource.type === 'json'}
              onClick={() => updateDataSourceType('json')}
              badge={systemInfo.dataSource === 'json' ? 'Active' : undefined}
            />
            <DataSourceCard
              icon={FileSpreadsheet}
              title="Upload Excel Roster"
              desc="Upload your unit's .xlsx or .csv spreadsheet"
              selected={settings.dataSource.type === 'excel'}
              onClick={() => updateDataSourceType('excel')}
              recommended
            />
          </div>

          {/* Excel Upload */}
          {settings.dataSource.type === 'excel' && (
            <div className="space-y-4 mt-4">
              <div className="border-2 border-dashed border-blue-200 bg-blue-50/50 rounded-lg p-6 text-center">
                <Upload className="w-8 h-8 text-blue-400 mx-auto mb-2" />
                <p className="text-sm text-slate-700 font-medium mb-1">Upload your roster spreadsheet</p>
                <p className="text-xs text-slate-500 mb-3">Supports .xlsx and .csv files. Heywood will automatically detect your column headers.</p>
                <input
                  type="file"
                  accept=".xlsx,.csv"
                  onChange={handleFileUpload}
                  className="text-sm"
                />
                {uploadStatus && (
                  <p className={`text-xs mt-2 ${uploadStatus.includes('failed') ? 'text-red-600' : 'text-green-600'}`}>
                    {uploadStatus}
                  </p>
                )}
              </div>
              {settings.dataSource.excelPath && (
                <div className="flex items-center gap-2 text-sm text-green-700 bg-green-50 px-3 py-2 rounded-lg">
                  <CheckCircle2 className="w-4 h-4" />
                  File loaded: {settings.dataSource.excelPath.split('/').pop()}
                </div>
              )}
              <p className="text-xs text-slate-500">
                After uploading, Heywood maps your columns automatically (e.g., "Last Name", "EDIPI", "Platoon").
                You can adjust any mappings that don't match.
              </p>
            </div>
          )}

          {/* Advanced data sources toggle */}
          <button
            onClick={() => setShowAdvancedData(!showAdvancedData)}
            className="flex items-center gap-2 text-sm text-slate-500 hover:text-slate-700 mt-4"
          >
            {showAdvancedData ? <ChevronDown className="w-4 h-4" /> : <ChevronRight className="w-4 h-4" />}
            Advanced: SharePoint, Database connections
          </button>

          {showAdvancedData && (
            <div className="mt-3 space-y-3 pl-4 border-l-2 border-slate-200">
              <p className="text-xs text-slate-500">
                These options are for units with IT support that want live data sync. Most units should use Excel upload instead.
              </p>
              <div className="grid grid-cols-1 sm:grid-cols-2 gap-3">
                <DataSourceCard
                  icon={Cloud}
                  title="SharePoint Lists"
                  desc="Connect to your unit's SharePoint site"
                  selected={settings.dataSource.type === 'sharepoint'}
                  onClick={() => updateDataSourceType('sharepoint')}
                  small
                />
                <DataSourceCard
                  icon={Database}
                  title="Database"
                  desc="Cosmos DB, PostgreSQL, or SQL Server"
                  selected={['cosmos', 'postgres', 'sqlserver'].includes(settings.dataSource.type)}
                  onClick={() => updateDataSourceType('postgres')}
                  small
                />
              </div>

              {/* SharePoint config */}
              {settings.dataSource.type === 'sharepoint' && (
                <div className="space-y-3 bg-slate-50 p-4 rounded-lg">
                  <p className="text-xs text-slate-600">
                    Your S-6 or IT shop can provide these values. They come from an Azure App Registration.
                  </p>
                  <div className="grid grid-cols-2 gap-3">
                    <SettingsInput
                      label="Tenant ID"
                      value={settings.dataSource.sharepoint.tenantId}
                      onChange={v => setSettings({
                        ...settings,
                        dataSource: { ...settings.dataSource, sharepoint: { ...settings.dataSource.sharepoint, tenantId: v } },
                      })}
                      help="From Azure Portal > App Registrations"
                    />
                    <SettingsInput
                      label="Client ID"
                      value={settings.dataSource.sharepoint.clientId}
                      onChange={v => setSettings({
                        ...settings,
                        dataSource: { ...settings.dataSource, sharepoint: { ...settings.dataSource.sharepoint, clientId: v } },
                      })}
                      help="Also called Application ID"
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
                    help="Generated in App Registration > Certificates & Secrets"
                  />
                  <SettingsInput
                    label="SharePoint Site URL"
                    value={settings.dataSource.sharepoint.siteUrl}
                    onChange={v => setSettings({
                      ...settings,
                      dataSource: { ...settings.dataSource, sharepoint: { ...settings.dataSource.sharepoint, siteUrl: v } },
                    })}
                    placeholder="https://yourtenant.sharepoint.com/sites/TBS"
                    help="The URL of your SharePoint site that has your data lists"
                  />
                  <div>
                    <label className="text-xs text-slate-500 mb-1 block">Network</label>
                    <select
                      value={settings.dataSource.sharepoint.cloud}
                      onChange={e => setSettings({
                        ...settings,
                        dataSource: { ...settings.dataSource, sharepoint: { ...settings.dataSource.sharepoint, cloud: e.target.value } },
                      })}
                      className="px-3 py-2 border border-slate-300 rounded-lg text-sm w-full"
                    >
                      <option value="commercial">Commercial (most units)</option>
                      <option value="gcc-high">GCC High (MCEN / IL5)</option>
                      <option value="dod">DoD Cloud</option>
                    </select>
                  </div>
                </div>
              )}

              {/* Database config */}
              {['cosmos', 'postgres', 'sqlserver'].includes(settings.dataSource.type) && (
                <div className="space-y-3 bg-slate-50 p-4 rounded-lg">
                  <p className="text-xs text-slate-600">
                    Your database admin or cloud team can provide the connection string.
                  </p>
                  <div>
                    <label className="text-xs text-slate-500 mb-1 block">Database Type</label>
                    <select
                      value={settings.dataSource.database.type || settings.dataSource.type}
                      onChange={e => setSettings({
                        ...settings,
                        dataSource: { ...settings.dataSource, type: e.target.value as DataSourceType, database: { ...settings.dataSource.database, type: e.target.value } },
                      })}
                      className="px-3 py-2 border border-slate-300 rounded-lg text-sm w-full"
                    >
                      <option value="cosmos">Azure Cosmos DB</option>
                      <option value="postgres">PostgreSQL</option>
                      <option value="sqlserver">Azure SQL / SQL Server</option>
                    </select>
                  </div>
                  <SettingsInput
                    label="Connection String"
                    value={settings.dataSource.database.connectionString}
                    onChange={v => setSettings({
                      ...settings,
                      dataSource: { ...settings.dataSource, database: { ...settings.dataSource.database, connectionString: v } },
                    })}
                    type="password"
                    help="The full connection string from your database provider"
                  />
                </div>
              )}

              {/* Test Connection */}
              {settings.dataSource.type !== 'json' && settings.dataSource.type !== 'excel' && (
                <div className="flex items-center gap-3">
                  <button
                    onClick={handleTestConnection}
                    disabled={testing}
                    className="px-4 py-2 border border-slate-300 rounded-lg text-sm hover:bg-slate-50 flex items-center gap-2"
                  >
                    {testing ? <Loader2 className="w-4 h-4 animate-spin" /> : <RefreshCw className="w-4 h-4" />}
                    Test Connection
                  </button>
                  {testResult && (
                    <div className={`flex items-center gap-1.5 text-sm ${testResult.status === 'ok' ? 'text-green-600' : 'text-red-600'}`}>
                      {testResult.status === 'ok' ? <CheckCircle2 className="w-4 h-4" /> : <AlertCircle className="w-4 h-4" />}
                      {testResult.message}
                    </div>
                  )}
                </div>
              )}
            </div>
          )}
        </div>
      </section>

      {/* STEP 2: AI Assistant */}
      <section className="bg-white rounded-xl border border-slate-200 overflow-hidden">
        <div className="px-6 py-4 bg-slate-50 border-b border-slate-200">
          <div className="flex items-center gap-3">
            <div className="w-7 h-7 rounded-full bg-[var(--color-navy)] text-white flex items-center justify-center text-sm font-bold">2</div>
            <div>
              <h2 className="text-base font-semibold text-slate-900">AI Assistant</h2>
              <p className="text-xs text-slate-500">Power Heywood's chat and analysis features</p>
            </div>
            <HelpButton id="ai-help" expanded={expandedHelp} onToggle={toggleHelp} />
          </div>
          {expandedHelp === 'ai-help' && (
            <div className="mt-3 p-3 bg-blue-50 rounded-lg text-xs text-blue-800 leading-relaxed">
              <BookOpen className="w-4 h-4 inline mr-1" />
              <strong>How this works:</strong> Heywood uses an AI model (like GPT-4) to understand your questions and generate analysis.
              You need an API key from OpenAI or Azure OpenAI. Your S-6 or IT shop can set this up.
              Without it, Heywood still works for viewing data, but the chat assistant will use placeholder responses.
            </div>
          )}
        </div>

        <div className="p-6">
          <div className="flex items-center gap-4">
            <div className={`flex items-center gap-2 px-3 py-2 rounded-lg text-sm font-medium ${
              systemInfo.ai.status === 'not configured'
                ? 'bg-yellow-50 text-yellow-800'
                : 'bg-green-50 text-green-800'
            }`}>
              {systemInfo.ai.status === 'not configured' ? (
                <AlertCircle className="w-4 h-4" />
              ) : (
                <CheckCircle2 className="w-4 h-4" />
              )}
              {systemInfo.ai.status === 'not configured' ? 'Not configured — using placeholder responses' : `Active — ${systemInfo.ai.model}`}
            </div>
          </div>

          <div className="mt-4 text-xs text-slate-500">
            <Info className="w-3.5 h-3.5 inline mr-1" />
            AI is configured via environment variables on the server (OPENAI_API_KEY or AZURE_OPENAI_*). Ask your IT admin to set these.
          </div>

          {/* Web Search */}
          <div className="mt-4 pt-4 border-t border-slate-100">
            <div className="flex items-center gap-2 mb-2">
              <Search className="w-4 h-4 text-slate-400" />
              <span className="text-sm font-medium text-slate-700">Web Search</span>
              <span className="text-xs text-slate-400">(optional)</span>
            </div>
            <p className="text-xs text-slate-500 mb-2">
              Lets Heywood search the web for doctrine references, regulations, and current info.
            </p>
            <div className="flex items-center gap-3">
              <input
                type="text"
                value={settings.ai.searxngUrl}
                onChange={e => setSettings({ ...settings, ai: { ...settings.ai, searxngUrl: e.target.value } })}
                className="flex-1 px-3 py-2 border border-slate-300 rounded-lg text-sm"
                placeholder="http://localhost:8888 (SearXNG URL)"
              />
              <StatusDot status={systemInfo.searxng.status === 'configured' ? 'ok' : 'off'} />
            </div>
          </div>
        </div>
      </section>

      {/* STEP 3: Outlook & Calendar */}
      <section className="bg-white rounded-xl border border-slate-200 overflow-hidden">
        <div className="px-6 py-4 bg-slate-50 border-b border-slate-200">
          <div className="flex items-center gap-3">
            <div className="w-7 h-7 rounded-full bg-[var(--color-navy)] text-white flex items-center justify-center text-sm font-bold">3</div>
            <div>
              <h2 className="text-base font-semibold text-slate-900">Outlook Mail & Calendar</h2>
              <p className="text-xs text-slate-500">Sync your Outlook calendar and mail into Heywood</p>
            </div>
            <HelpButton id="outlook-help" expanded={expandedHelp} onToggle={toggleHelp} />
          </div>
          {expandedHelp === 'outlook-help' && (
            <div className="mt-3 p-3 bg-blue-50 rounded-lg text-xs text-blue-800 leading-relaxed">
              <BookOpen className="w-4 h-4 inline mr-1" />
              <strong>How this works:</strong> When connected, Heywood shows your Outlook calendar events alongside the TBS training schedule,
              and displays recent mail in the Calendar page sidebar. Heywood can also answer questions like "What's on my schedule today?"
              in chat. Your S-6 needs to create an Azure App Registration with Calendar.Read and Mail.Read permissions.
              In demo mode, Heywood uses sample calendar events so you can see how it looks.
            </div>
          )}
        </div>

        <div className="p-6">
          <div className="flex items-center gap-3 mb-4">
            <label className="relative inline-flex items-center cursor-pointer">
              <input
                type="checkbox"
                checked={settings.outlook.enabled}
                onChange={e => setSettings({ ...settings, outlook: { ...settings.outlook, enabled: e.target.checked } })}
                className="sr-only peer"
              />
              <div className="w-11 h-6 bg-slate-200 peer-focus:outline-none rounded-full peer peer-checked:after:translate-x-full rtl:peer-checked:after:-translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-[2px] after:start-[2px] after:bg-white after:border-slate-300 after:border after:rounded-full after:h-5 after:w-5 after:transition-all peer-checked:bg-green-500" />
            </label>
            <span className="text-sm font-medium text-slate-700">
              {settings.outlook.enabled ? 'Outlook sync enabled' : 'Outlook sync disabled'}
            </span>
          </div>

          {!settings.outlook.enabled && (
            <div className="text-sm text-slate-500 bg-slate-50 p-4 rounded-lg">
              Calendar page is using demo data. Enable Outlook sync to show your real calendar and mail.
            </div>
          )}

          {settings.outlook.enabled && (
            <div className="space-y-4">
              <div className={`flex items-center gap-2 px-3 py-2 rounded-lg text-sm ${
                settings.outlook.tenantId
                  ? 'bg-green-50 text-green-800'
                  : 'bg-yellow-50 text-yellow-800'
              }`}>
                {settings.outlook.tenantId ? (
                  <><CheckCircle2 className="w-4 h-4" /> Connected</>
                ) : (
                  <><AlertCircle className="w-4 h-4" /> Credentials needed — ask your S-6 for the Azure App Registration details</>
                )}
              </div>

              <button
                onClick={() => setShowAdvancedOutlook(!showAdvancedOutlook)}
                className="flex items-center gap-2 text-sm text-slate-500 hover:text-slate-700"
              >
                {showAdvancedOutlook ? <ChevronDown className="w-4 h-4" /> : <ChevronRight className="w-4 h-4" />}
                {settings.outlook.tenantId ? 'Edit connection details' : 'Enter connection details'}
              </button>

              {showAdvancedOutlook && (
                <div className="space-y-3 bg-slate-50 p-4 rounded-lg">
                  <p className="text-xs text-slate-600">
                    These values come from an Azure App Registration. Your S-6 or IT admin can provide them.
                  </p>
                  <div className="grid grid-cols-2 gap-3">
                    <SettingsInput
                      label="Tenant ID"
                      value={settings.outlook.tenantId}
                      onChange={v => setSettings({ ...settings, outlook: { ...settings.outlook, tenantId: v } })}
                      help="From Azure Portal > App Registrations > Overview"
                    />
                    <SettingsInput
                      label="Client ID"
                      value={settings.outlook.clientId}
                      onChange={v => setSettings({ ...settings, outlook: { ...settings.outlook, clientId: v } })}
                      help="Also called Application (client) ID"
                    />
                  </div>
                  <SettingsInput
                    label="Client Secret"
                    value={settings.outlook.clientSecret}
                    onChange={v => setSettings({ ...settings, outlook: { ...settings.outlook, clientSecret: v } })}
                    type="password"
                    help="Generated in Certificates & Secrets section"
                  />
                  <div className="grid grid-cols-2 gap-3">
                    <div>
                      <label className="text-xs text-slate-500 mb-1 block">Network</label>
                      <select
                        value={settings.outlook.cloud}
                        onChange={e => setSettings({ ...settings, outlook: { ...settings.outlook, cloud: e.target.value } })}
                        className="px-3 py-2 border border-slate-300 rounded-lg text-sm w-full"
                      >
                        <option value="commercial">Commercial (most units)</option>
                        <option value="gcc-high">GCC High (MCEN / IL5)</option>
                        <option value="dod">DoD Cloud</option>
                      </select>
                    </div>
                    <SettingsInput
                      label="Sync Interval (minutes)"
                      value={String(settings.outlook.syncIntervalMinutes)}
                      onChange={v => setSettings({ ...settings, outlook: { ...settings.outlook, syncIntervalMinutes: parseInt(v) || 5 } })}
                      type="number"
                      help="How often to check for new events (default: 5)"
                    />
                  </div>
                  <SettingsInput
                    label="Shared Calendar ID"
                    value={settings.outlook.masterCalendarId}
                    onChange={v => setSettings({ ...settings, outlook: { ...settings.outlook, masterCalendarId: v } })}
                    placeholder="(optional) For TBS-wide shared calendar"
                    help="If your unit has a shared master calendar, enter its ID here"
                  />
                </div>
              )}
            </div>
          )}
        </div>
      </section>
    </div>
  )
}

// ---- Subcomponents ----

function StatusCard({ label, value, status }: { label: string; value: string; status: 'ok' | 'warn' | 'error' }) {
  const colors = { ok: 'border-green-200 bg-green-50', warn: 'border-yellow-200 bg-yellow-50', error: 'border-red-200 bg-red-50' }
  const dots = { ok: 'bg-green-500', warn: 'bg-yellow-500', error: 'bg-red-500' }
  return (
    <div className={`rounded-lg p-3 border ${colors[status]}`}>
      <div className="flex items-center gap-1.5">
        <div className={`w-2 h-2 rounded-full ${dots[status]}`} />
        <span className="text-xs text-slate-500">{label}</span>
      </div>
      <div className="text-sm font-semibold text-slate-900 mt-1">{value}</div>
    </div>
  )
}

function DataSourceCard({ icon: Icon, title, desc, selected, onClick, badge, recommended, small }: {
  icon: typeof Database
  title: string
  desc: string
  selected: boolean
  onClick: () => void
  badge?: string
  recommended?: boolean
  small?: boolean
}) {
  return (
    <button
      onClick={onClick}
      className={`${small ? 'p-3' : 'p-4'} rounded-lg border text-left transition-all relative ${
        selected
          ? 'border-[var(--color-navy)] bg-blue-50 ring-1 ring-[var(--color-navy)]'
          : 'border-slate-200 hover:border-slate-300 hover:bg-slate-50'
      }`}
    >
      {recommended && !selected && (
        <span className="absolute -top-2 right-2 text-[10px] bg-green-100 text-green-700 px-1.5 py-0.5 rounded-full font-medium">
          Recommended
        </span>
      )}
      {badge && (
        <span className="absolute -top-2 right-2 text-[10px] bg-blue-100 text-blue-700 px-1.5 py-0.5 rounded-full font-medium">
          {badge}
        </span>
      )}
      <Icon className={`w-5 h-5 mb-1.5 ${selected ? 'text-[var(--color-navy)]' : 'text-slate-400'}`} />
      <div className={`${small ? 'text-sm' : 'text-sm'} font-medium`}>{title}</div>
      <div className="text-xs text-slate-500 mt-0.5">{desc}</div>
    </button>
  )
}

function HelpButton({ id, expanded, onToggle }: { id: string; expanded: string | null; onToggle: (id: string) => void }) {
  return (
    <button
      onClick={() => onToggle(id)}
      className={`ml-auto flex items-center gap-1 text-xs px-2 py-1 rounded-lg transition-colors ${
        expanded === id ? 'bg-blue-100 text-blue-700' : 'text-slate-400 hover:text-slate-600 hover:bg-slate-100'
      }`}
    >
      <HelpCircle className="w-3.5 h-3.5" />
      {expanded === id ? 'Hide help' : "What's this?"}
    </button>
  )
}

function SettingsInput({ label, value, onChange, type = 'text', placeholder, help }: {
  label: string
  value: string
  onChange: (v: string) => void
  type?: string
  placeholder?: string
  help?: string
}) {
  return (
    <div>
      <label className="text-xs text-slate-500 mb-1 block">{label}</label>
      <input
        type={type}
        value={value}
        onChange={e => onChange(e.target.value)}
        placeholder={placeholder}
        className="w-full px-3 py-2 border border-slate-300 rounded-lg text-sm"
      />
      {help && <p className="text-[11px] text-slate-400 mt-1">{help}</p>}
    </div>
  )
}

function StatusDot({ status }: { status: 'ok' | 'error' | 'off' }) {
  const colors = { ok: 'bg-green-500', error: 'bg-red-500', off: 'bg-slate-300' }
  return <div className={`w-2.5 h-2.5 rounded-full ${colors[status]}`} />
}
