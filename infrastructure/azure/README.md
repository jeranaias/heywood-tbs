# Azure Infrastructure — Heywood Phase 2

Bicep templates for deploying Heywood Phase 2 resources on Azure Government (IL5).

## Prerequisites

1. **Azure Government subscription** with IL5 authorization
2. **Azure CLI** installed with government cloud configured:
   ```bash
   az cloud set --name AzureUSGovernment
   az login
   ```
3. **IATT approved** before deploying to production (see `docs/authorization/iatt-application-draft.md`)
4. **Resource group** created in `usgovvirginia` or `usgovarizona`

## Resources Deployed

| Resource | SKU | Purpose | Monthly Cost |
|----------|-----|---------|-------------|
| Azure OpenAI | S0 | GPT-4o + embeddings | $200-500 |
| Azure AI Search | S1 | RAG vector store with CMK | $250 |
| Azure App Service | B1 | Custom connector + anonymization | $55 |
| Azure Key Vault | Standard | Secrets + CMK keys | $10 |
| Log Analytics | PerGB | Audit logging + monitoring | $25 |
| Application Insights | — | App performance monitoring | Included |
| Managed Identity | — | Service-to-service auth | Free |
| Budget Alert | — | Cost monitoring at 75/90/100% | Free |

## Deployment

### 1. Update parameters

Edit `parameters.dev.json`:
- Set `adminObjectId` to your Azure AD Object ID
- Verify `location` matches your authorized region

```bash
# Find your Object ID
az ad signed-in-user show --query id -o tsv
```

### 2. Deploy

```bash
# Create resource group
az group create \
  --name rg-usmc-tbs-heywood-dev \
  --location usgovvirginia \
  --tags project=heywood classification=CUI impactLevel=IL5

# Deploy infrastructure
az deployment group create \
  --resource-group rg-usmc-tbs-heywood-dev \
  --template-file main.bicep \
  --parameters @parameters.dev.json \
  --name heywood-phase2-deploy
```

### 3. Verify deployment

```bash
# List deployed resources
az resource list \
  --resource-group rg-usmc-tbs-heywood-dev \
  --output table

# Verify OpenAI endpoint
az cognitiveservices account show \
  --name usmc-tbs-heywood-dev-oai-ugv \
  --resource-group rg-usmc-tbs-heywood-dev \
  --query properties.endpoint
```

## Security Notes

- **No API keys** — all service-to-service auth uses Managed Identity
- **No public endpoints** — Azure OpenAI and AI Search have public access disabled
- **CMK encryption** — AI Search uses Customer-Managed Key from Key Vault
- **Purge protection** — Key Vault has purge protection enabled (cannot be disabled once set)
- **TLS 1.2+** — enforced on all endpoints
- **Audit logging** — all API calls logged to Log Analytics (90-day retention)
- **Budget alerts** — notifications at 75%, 90%, and 100% of $1,500/month ceiling

## Naming Convention

Resources follow DoD cloud naming standards:

```
{org}-{project}-{environment}-{resource-type}-{region-short}
usmc-tbs-heywood-dev-oai-ugv
```

| Abbreviation | Meaning |
|---|---|
| `oai` | Azure OpenAI |
| `srch` | Azure AI Search |
| `app` | App Service |
| `kv` | Key Vault |
| `log` | Log Analytics |
| `ai` | Application Insights |
| `asp` | App Service Plan |
| `id` | Managed Identity |
| `ugv` | US Gov Virginia |
| `uga` | US Gov Arizona |

## Post-Deployment Steps

1. Configure Key Vault network rules to allow App Service subnet
2. Create AI Search index for TBS doctrine (see `schemas/search/` — TBD Phase 2)
3. Deploy App Service code (custom connector + anonymization layer)
4. Register custom connector in Power Platform
5. Configure Azure Entra ID app registration for Power App auth
6. Run IATT security tests (see `docs/authorization/iatt-application-draft.md` Section 7)
