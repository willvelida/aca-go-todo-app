@description('The location to deploy our resources to. Default is location of resource group')
param location string = resourceGroup().location

@description('The name of our application.')
param applicationName string = uniqueString(resourceGroup().id)

@description('The name of the Azure Container Registry')
param containerRegistryName string = 'acr${applicationName}'

@description('The name of the Key Vault')
param keyVaultName string = 'kv${applicationName}'

@description('The name of the Log Analytics workspace')
param logAnalyticsWorkspaceName string = 'law${applicationName}'

@description('The name of the App Insights workspace')
param appInsightsName string = 'appins${applicationName}'

@description('The name of the Container App Environment')
param containerEnvironmentName string = 'env${applicationName}'

@description('The name of the Cosmos DB account')
param cosmosDbAccountName string = 'db${applicationName}'

@description('The database name')
param databaseName string = 'TodoDB'

@description('The container name')
param containerName string = 'todos'

@description('The latest deployment timestamp')
param lastDeployed string = utcNow('d')

var frontendName = 'todo-frontend'
var backendName = 'todo-backend'
var tags = {
  ApplicationName: 'Todo-App'
  Language: 'Golang'
  Environment: 'Production'
  LastDeployed: lastDeployed
}

resource keyVault 'Microsoft.KeyVault/vaults@2022-07-01' = {
  name: keyVaultName
  location: location
  tags: tags
  properties: {
    sku: {
      family: 'A'
      name: 'standard'
    }
    tenantId: tenant().tenantId
    enabledForDeployment: true
    enabledForTemplateDeployment: true
    enableSoftDelete: false
    accessPolicies: [
    ]
  }
}

resource appInsights 'Microsoft.Insights/components@2020-02-02' = {
  name: appInsightsName
  tags: tags
  location: location
  kind: 'web'
  properties: {
    Application_Type: 'web'
    WorkspaceResourceId: logAnalytics.outputs.logAnalyticsId
  }
}

module logAnalytics 'modules/log-analytics.bicep' = {
  name: 'log-analytics'
  params: {
    tags: tags
    keyVaultName: keyVault.name
    location: location
    logAnalyticsWorkspaceName: logAnalyticsWorkspaceName
  }
}

module containerRegistry 'modules/container-registry.bicep' = {
  name: 'acr'
  params: {
    tags: tags
    containerRegistryName: containerRegistryName
    keyVaultName: keyVault.name
    location: location
  }
}

module cosmos 'modules/cosmos-db.bicep' = {
  name: 'cosmos'
  params: {
    accountName: cosmosDbAccountName 
    containerName: containerName
    databaseName: databaseName
    location: location
    tags: tags
  }
}

module env 'modules/container-app-environment.bicep' = {
  name: 'container-app-env'
  params: {
    containerEnvironmentName: containerEnvironmentName
    location: location
    logAnalyticsCustomerId: logAnalytics.outputs.customerId 
    logAnalyticsSharedKey: keyVault.getSecret('log-analytics-shared-key')
    tags: tags
  }
}

module restApi 'modules/http-container-app.bicep' = {
  name: 'rest-api'
  params: {
    acrPasswordSecret: keyVault.getSecret('acr-primary-password') 
    acrServerName: containerRegistry.outputs.loginServer
    acrUsername: keyVault.getSecret('acr-username')
    containerAppEnvId: env.outputs.containerAppEnvId
    containerAppName: backendName
    containerImage: 'mcr.microsoft.com/azuredocs/containerapps-helloworld:latest'
    cpuCore: '0.5'
    isExternal: false
    location: location
    memorySize: '1.0'
    tags: tags
  }
}

module frontend 'modules/http-container-app.bicep' = {
  name: 'front-end'
  params: {
    acrPasswordSecret: keyVault.getSecret('acr-primary-password')
    acrServerName: containerRegistry.outputs.loginServer
    acrUsername: keyVault.getSecret('acr-username')
    containerAppEnvId: env.outputs.containerAppEnvId
    containerAppName: frontendName
    containerImage: 'mcr.microsoft.com/azuredocs/containerapps-helloworld:latest'
    cpuCore: '0.5'
    isExternal: true
    location: location
    memorySize: '1.0'
    tags: tags
  }
}
