// Heywood TBS — Phase 2 Azure Government Infrastructure
// Deployment target: Azure Government (usgovvirginia or usgovarizona)
// Classification: IL5 (CUI)
//
// Resources:
//   - Azure OpenAI Service (GPT-4o + text-embedding-ada-002)
//   - Azure AI Search (S1 tier with CMK)
//   - Azure App Service (B1 — custom connector)
//   - Azure Key Vault (secrets + CMK keys)
//   - Azure Monitor (Log Analytics workspace)
//   - Managed Identity (for service-to-service auth)

targetScope = 'resourceGroup'

// ============================================================================
// Parameters
// ============================================================================

@description('Environment name (dev, test, prod)')
@allowed(['dev', 'test', 'prod'])
param environment string = 'dev'

@description('Azure Government region')
@allowed(['usgovvirginia', 'usgovarizona'])
param location string = 'usgovvirginia'

@description('Project name used in resource naming')
param projectName string = 'heywood'

@description('Organization identifier for naming')
param orgPrefix string = 'usmc-tbs'

@description('Azure AD Object ID of the system administrator')
param adminObjectId string

@description('Azure AD Tenant ID')
param tenantId string = subscription().tenantId

@description('App Service Plan SKU')
@allowed(['B1', 'B2', 'S1'])
param appServiceSku string = 'B1'

@description('AI Search SKU')
@allowed(['basic', 'standard'])
param searchSku string = 'standard'

@description('Tags applied to all resources')
param tags object = {
  project: 'heywood'
  environment: environment
  classification: 'CUI'
  impactLevel: 'IL5'
  owner: 'tbs-s3'
  costCenter: 'tbs-training'
  managedBy: 'ssgt-morgan'
}

// ============================================================================
// Variables — DoD Cloud Naming Convention
// ============================================================================

// Pattern: {org}-{project}-{environment}-{resource-type}-{region-short}
var regionShort = location == 'usgovvirginia' ? 'ugv' : 'uga'
var namePrefix = '${orgPrefix}-${projectName}-${environment}'
var nameSuffix = regionShort

// Resource names
var keyVaultName = replace('${namePrefix}-kv-${nameSuffix}', '-', '')  // Key Vault has strict naming
var logAnalyticsName = '${namePrefix}-log-${nameSuffix}'
var appInsightsName = '${namePrefix}-ai-${nameSuffix}'
var openAiName = '${namePrefix}-oai-${nameSuffix}'
var searchName = replace('${namePrefix}-srch-${nameSuffix}', '-', '')  // Search has strict naming
var appServicePlanName = '${namePrefix}-asp-${nameSuffix}'
var appServiceName = '${namePrefix}-app-${nameSuffix}'
var managedIdentityName = '${namePrefix}-id-${nameSuffix}'

// ============================================================================
// Managed Identity
// ============================================================================

resource managedIdentity 'Microsoft.ManagedIdentity/userAssignedManagedIdentities@2023-01-31' = {
  name: managedIdentityName
  location: location
  tags: tags
}

// ============================================================================
// Log Analytics + Application Insights
// ============================================================================

resource logAnalytics 'Microsoft.OperationalInsights/workspaces@2022-10-01' = {
  name: logAnalyticsName
  location: location
  tags: tags
  properties: {
    sku: {
      name: 'PerGB2018'
    }
    retentionInDays: 90
    features: {
      enableLogAccessUsingOnlyResourcePermissions: true
    }
  }
}

resource appInsights 'Microsoft.Insights/components@2020-02-02' = {
  name: appInsightsName
  location: location
  tags: tags
  kind: 'web'
  properties: {
    Application_Type: 'web'
    WorkspaceResourceId: logAnalytics.id
    RetentionInDays: 90
  }
}

// ============================================================================
// Key Vault (secrets + CMK for AI Search)
// ============================================================================

resource keyVault 'Microsoft.KeyVault/vaults@2023-07-01' = {
  name: keyVaultName
  location: location
  tags: tags
  properties: {
    tenantId: tenantId
    sku: {
      family: 'A'
      name: 'standard'
    }
    enabledForDeployment: false
    enabledForDiskEncryption: false
    enabledForTemplateDeployment: false
    enableSoftDelete: true
    softDeleteRetentionInDays: 90
    enablePurgeProtection: true  // Required for CMK — cannot be disabled once enabled
    enableRbacAuthorization: true
    networkAcls: {
      defaultAction: 'Deny'
      bypass: 'AzureServices'
      ipRules: []
      virtualNetworkRules: []
    }
  }
}

// Key Vault access for managed identity
resource kvRoleAssignment 'Microsoft.Authorization/roleAssignments@2022-04-01' = {
  name: guid(keyVault.id, managedIdentity.id, 'Key Vault Secrets User')
  scope: keyVault
  properties: {
    roleDefinitionId: subscriptionResourceId('Microsoft.Authorization/roleDefinitions', '4633458b-17de-408a-b874-0445c86b69e6') // Key Vault Secrets User
    principalId: managedIdentity.properties.principalId
    principalType: 'ServicePrincipal'
  }
}

// Key Vault access for admin
resource kvAdminRoleAssignment 'Microsoft.Authorization/roleAssignments@2022-04-01' = {
  name: guid(keyVault.id, adminObjectId, 'Key Vault Administrator')
  scope: keyVault
  properties: {
    roleDefinitionId: subscriptionResourceId('Microsoft.Authorization/roleDefinitions', '00482a5a-887f-4fb3-b363-3b7fe8e74483') // Key Vault Administrator
    principalId: adminObjectId
    principalType: 'User'
  }
}

// CMK key for AI Search encryption
resource cmkKey 'Microsoft.KeyVault/vaults/keys@2023-07-01' = {
  parent: keyVault
  name: 'heywood-search-cmk'
  properties: {
    kty: 'RSA'
    keySize: 2048
    keyOps: [
      'encrypt'
      'decrypt'
      'wrapKey'
      'unwrapKey'
    ]
    attributes: {
      enabled: true
    }
    rotationPolicy: {
      attributes: {
        expiryTime: 'P1Y'  // 1 year rotation
      }
      lifetimeActions: [
        {
          action: {
            type: 'rotate'
          }
          trigger: {
            timeBeforeExpiry: 'P30D'  // Rotate 30 days before expiry
          }
        }
        {
          action: {
            type: 'notify'
          }
          trigger: {
            timeBeforeExpiry: 'P60D'  // Notify 60 days before expiry
          }
        }
      ]
    }
  }
}

// ============================================================================
// Azure OpenAI Service
// ============================================================================

resource openAi 'Microsoft.CognitiveServices/accounts@2024-04-01-preview' = {
  name: openAiName
  location: location
  tags: tags
  kind: 'OpenAI'
  sku: {
    name: 'S0'
  }
  identity: {
    type: 'UserAssigned'
    userAssignedIdentities: {
      '${managedIdentity.id}': {}
    }
  }
  properties: {
    customSubDomainName: openAiName
    publicNetworkAccess: 'Disabled'
    networkAcls: {
      defaultAction: 'Deny'
    }
    disableLocalAuth: true  // Force Azure AD auth only — no API keys
  }
}

// GPT-4o deployment
resource gpt4oDeployment 'Microsoft.CognitiveServices/accounts/deployments@2024-04-01-preview' = {
  parent: openAi
  name: 'gpt-4o'
  sku: {
    name: 'Standard'
    capacity: 30  // 30K tokens per minute — adjust based on usage
  }
  properties: {
    model: {
      format: 'OpenAI'
      name: 'gpt-4o'
      version: '2024-11-20'
    }
    raiPolicyName: 'Microsoft.DefaultV2'
  }
}

// Embedding model deployment
resource embeddingDeployment 'Microsoft.CognitiveServices/accounts/deployments@2024-04-01-preview' = {
  parent: openAi
  name: 'text-embedding-ada-002'
  sku: {
    name: 'Standard'
    capacity: 120  // 120K tokens per minute for embeddings
  }
  properties: {
    model: {
      format: 'OpenAI'
      name: 'text-embedding-ada-002'
      version: '2'
    }
  }
  dependsOn: [gpt4oDeployment]  // Sequential deployment required
}

// OpenAI Contributor role for managed identity
resource openAiRoleAssignment 'Microsoft.Authorization/roleAssignments@2022-04-01' = {
  name: guid(openAi.id, managedIdentity.id, 'Cognitive Services OpenAI Contributor')
  scope: openAi
  properties: {
    roleDefinitionId: subscriptionResourceId('Microsoft.Authorization/roleDefinitions', 'a001fd3d-188f-4b5d-821b-7da978bf7442') // Cognitive Services OpenAI Contributor
    principalId: managedIdentity.properties.principalId
    principalType: 'ServicePrincipal'
  }
}

// Diagnostic settings for OpenAI
resource openAiDiagnostics 'Microsoft.Insights/diagnosticSettings@2021-05-01-preview' = {
  name: 'heywood-oai-diagnostics'
  scope: openAi
  properties: {
    workspaceId: logAnalytics.id
    logs: [
      {
        category: 'Audit'
        enabled: true
        retentionPolicy: { enabled: true, days: 90 }
      }
      {
        category: 'RequestResponse'
        enabled: true
        retentionPolicy: { enabled: true, days: 90 }
      }
    ]
    metrics: [
      {
        category: 'AllMetrics'
        enabled: true
        retentionPolicy: { enabled: true, days: 90 }
      }
    ]
  }
}

// ============================================================================
// Azure AI Search (S1 + CMK)
// ============================================================================

resource search 'Microsoft.Search/searchServices@2024-03-01-preview' = {
  name: searchName
  location: location
  tags: tags
  sku: {
    name: searchSku
  }
  identity: {
    type: 'UserAssigned'
    userAssignedIdentities: {
      '${managedIdentity.id}': {}
    }
  }
  properties: {
    hostingMode: 'default'
    partitionCount: 1
    replicaCount: 1
    publicNetworkAccess: 'disabled'
    encryptionWithCmk: {
      enforcement: 'Enabled'
    }
    authOptions: {
      aadOrApiKey: {
        aadAuthFailureMode: 'http401WithBearerChallenge'
      }
    }
    disableLocalAuth: true  // Force Azure AD auth only
  }
}

// Search Service Contributor role for managed identity
resource searchRoleAssignment 'Microsoft.Authorization/roleAssignments@2022-04-01' = {
  name: guid(search.id, managedIdentity.id, 'Search Service Contributor')
  scope: search
  properties: {
    roleDefinitionId: subscriptionResourceId('Microsoft.Authorization/roleDefinitions', '7ca78c08-252a-4471-8644-bb5ff32d4ba0') // Search Service Contributor
    principalId: managedIdentity.properties.principalId
    principalType: 'ServicePrincipal'
  }
}

// Search Index Data Contributor role for managed identity
resource searchIndexRoleAssignment 'Microsoft.Authorization/roleAssignments@2022-04-01' = {
  name: guid(search.id, managedIdentity.id, 'Search Index Data Contributor')
  scope: search
  properties: {
    roleDefinitionId: subscriptionResourceId('Microsoft.Authorization/roleDefinitions', '8bbe4f3e-f3f2-4e53-a0d3-1f20962e8e4a') // Search Index Data Contributor
    principalId: managedIdentity.properties.principalId
    principalType: 'ServicePrincipal'
  }
}

// ============================================================================
// App Service (Custom Connector + Anonymization Layer)
// ============================================================================

resource appServicePlan 'Microsoft.Web/serverfarms@2023-12-01' = {
  name: appServicePlanName
  location: location
  tags: tags
  kind: 'linux'
  sku: {
    name: appServiceSku
  }
  properties: {
    reserved: true  // Linux
  }
}

resource appService 'Microsoft.Web/sites@2023-12-01' = {
  name: appServiceName
  location: location
  tags: tags
  kind: 'app,linux'
  identity: {
    type: 'UserAssigned'
    userAssignedIdentities: {
      '${managedIdentity.id}': {}
    }
  }
  properties: {
    serverFarmId: appServicePlan.id
    httpsOnly: true  // No HTTP allowed
    siteConfig: {
      linuxFxVersion: 'NODE|20-lts'
      minTlsVersion: '1.2'
      ftpsState: 'Disabled'
      http20Enabled: true
      alwaysOn: true
      appSettings: [
        {
          name: 'AZURE_OPENAI_ENDPOINT'
          value: openAi.properties.endpoint
        }
        {
          name: 'AZURE_OPENAI_DEPLOYMENT_GPT'
          value: gpt4oDeployment.name
        }
        {
          name: 'AZURE_OPENAI_DEPLOYMENT_EMBEDDING'
          value: embeddingDeployment.name
        }
        {
          name: 'AZURE_SEARCH_ENDPOINT'
          value: 'https://${search.name}.search.azure.us'
        }
        {
          name: 'AZURE_SEARCH_INDEX_NAME'
          value: 'heywood-doctrine-index'
        }
        {
          name: 'AZURE_CLIENT_ID'
          value: managedIdentity.properties.clientId
        }
        {
          name: 'APPLICATIONINSIGHTS_CONNECTION_STRING'
          value: appInsights.properties.ConnectionString
        }
        {
          name: 'KEY_VAULT_URI'
          value: keyVault.properties.vaultUri
        }
        {
          name: 'NODE_ENV'
          value: environment == 'prod' ? 'production' : 'development'
        }
      ]
    }
  }
}

// App Service diagnostic settings
resource appServiceDiagnostics 'Microsoft.Insights/diagnosticSettings@2021-05-01-preview' = {
  name: 'heywood-app-diagnostics'
  scope: appService
  properties: {
    workspaceId: logAnalytics.id
    logs: [
      {
        category: 'AppServiceHTTPLogs'
        enabled: true
        retentionPolicy: { enabled: true, days: 90 }
      }
      {
        category: 'AppServiceConsoleLogs'
        enabled: true
        retentionPolicy: { enabled: true, days: 90 }
      }
      {
        category: 'AppServiceAuditLogs'
        enabled: true
        retentionPolicy: { enabled: true, days: 90 }
      }
    ]
    metrics: [
      {
        category: 'AllMetrics'
        enabled: true
        retentionPolicy: { enabled: true, days: 90 }
      }
    ]
  }
}

// ============================================================================
// Budget Alert
// ============================================================================

resource budget 'Microsoft.Consumption/budgets@2023-11-01' = {
  name: '${namePrefix}-budget'
  properties: {
    category: 'Cost'
    amount: 1500  // $1,500/month budget ceiling
    timeGrain: 'Monthly'
    timePeriod: {
      startDate: '2026-06-01'  // Phase 2 start
      endDate: '2026-12-01'    // Phase 2 + buffer
    }
    notifications: {
      seventyFivePercent: {
        enabled: true
        operator: 'GreaterThanOrEqualTo'
        threshold: 75
        contactEmails: ['heywood-admin@usmc.mil']  // Replace with actual
        thresholdType: 'Actual'
      }
      ninetyPercent: {
        enabled: true
        operator: 'GreaterThanOrEqualTo'
        threshold: 90
        contactEmails: ['heywood-admin@usmc.mil']
        thresholdType: 'Actual'
      }
      oneHundredPercent: {
        enabled: true
        operator: 'GreaterThanOrEqualTo'
        threshold: 100
        contactEmails: ['heywood-admin@usmc.mil']
        thresholdType: 'Actual'
      }
    }
  }
}

// ============================================================================
// Outputs
// ============================================================================

output resourceGroupName string = resourceGroup().name
output managedIdentityClientId string = managedIdentity.properties.clientId
output managedIdentityPrincipalId string = managedIdentity.properties.principalId
output keyVaultUri string = keyVault.properties.vaultUri
output openAiEndpoint string = openAi.properties.endpoint
output searchEndpoint string = 'https://${search.name}.search.azure.us'
output appServiceUrl string = 'https://${appService.properties.defaultHostName}'
output appInsightsInstrumentationKey string = appInsights.properties.InstrumentationKey
output logAnalyticsWorkspaceId string = logAnalytics.id
