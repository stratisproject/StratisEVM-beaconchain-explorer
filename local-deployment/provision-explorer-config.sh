#! /bin/bash
# CL_PORT=$(kurtosis enclave inspect my-testnet | grep 4000/tcp | tr -s ' ' | cut -d " " -f 6 | sed -e 's/http\:\/\/127.0.0.1\://' | head -n 1)
# echo "CL Node port is $CL_PORT"

# EL_PORT=$(kurtosis enclave inspect my-testnet | grep 8545/tcp | tr -s ' ' | cut -d " " -f 5 | sed -e 's/127.0.0.1\://' | head -n 1)
# echo "EL Node port is $EL_PORT"

# REDIS_PORT=$(kurtosis enclave inspect my-testnet | grep 6379/tcp | tr -s ' ' | cut -d " " -f 6 | sed -e 's/tcp\:\/\/127.0.0.1\://' | head -n 1)
# echo "Redis port is $REDIS_PORT"

# POSTGRES_PORT=$(kurtosis enclave inspect my-testnet | grep 5432/tcp | tr -s ' ' | cut -d " " -f 6 | sed -e 's/postgresql\:\/\/127.0.0.1\://' | head -n 1)
# echo "Postgres port is $POSTGRES_PORT"

# LBT_PORT=$(kurtosis enclave inspect my-testnet | grep 9000/tcp | tr -s ' ' | cut -d " " -f 6 | sed -e 's/tcp\:\/\/127.0.0.1\://' | tail -n 1)
# echo "Little bigtable port is $LBT_PORT"

CL_PORT=3500
EL_PORT=8545
REDIS_PORT=6379
POSTGRES_PORT=5432
LBT_PORT=8086

cat <<EOF > .env
CL_PORT=$CL_PORT
EL_PORT=$EL_PORT
REDIS_PORT=$REDIS_PORT
REDIS_SESSIONS_PORT=$REDIS_SESSIONS_PORT
POSTGRES_PORT=$POSTGRES_PORT
LBT_PORT=$LBT_PORT
EOF

touch elconfig.json
cat >elconfig.json <<EOL
{
  "chainId": 205205,
  "homesteadBlock": 0,
  "daoForkSupport": true,
  "eip150Block": 0,
  "eip155Block": 0,
  "eip158Block": 0,
  "byzantiumBlock": 0,
  "constantinopleBlock": 0,
  "petersburgBlock": 0,
  "istanbulBlock": 0,
  "muirGlacierBlock": 0,
  "berlinBlock": 0,
  "londonBlock": 0,
  "arrowGlacierBlock": 0,
  "grayGlacierBlock": 0,
  "shanghaiTime": 1705497931,
  "cancunTime": 1705497931,
  "terminalTotalDifficulty": 0,
  "terminalTotalDifficultyPassed": true
}
EOL

touch config.yml

cat >config.yml <<EOL
chain:
  clConfigPath: 'node'
  elConfigPath: 'local-deployment/elconfig.json'
readerDatabase:
  name: explorer
  host: 127.0.0.1
  port: "$POSTGRES_PORT"
  user: postgres
  password: "test"
writerDatabase:
  name: explorer
  host: 127.0.0.1
  port: "$POSTGRES_PORT"
  user: postgres
  password: "test"
bigtable:
  project: explorer
  instance: explorer
  emulator: true
  emulatorPort: $LBT_PORT
eth1ErigonEndpoint: 'http://127.0.0.1:$EL_PORT'
eth1GethEndpoint: 'http://127.0.0.1:$EL_PORT'
redisCacheEndpoint: '127.0.0.1:$REDIS_PORT'
redisSessionStoreEndpoint: '127.0.0.1:$REDIS_SESSIONS_PORT'
tieredCacheProvider: 'redis'
frontend:
  siteSSL: false
  siteDomain: "localhost:8080"
  siteName: 'Open Source Ethereum (ETH) Testnet Explorer' # Name of the site, displayed in the title tag
  siteSubtitle: "Showing a local testnet."
  server:
    host: '0.0.0.0' # Address to listen on
    port: '8080' # Port to listen on
  readerDatabase:
    name: explorer
    host: 127.0.0.1
    port: "$POSTGRES_PORT"
    user: postgres
    password: "test"
  writerDatabase:
    name: explorer
    host: 127.0.0.1
    port: "$POSTGRES_PORT"
    user: postgres
    password: "test"
  sessionSecret: "11111111111111111111111111111111"
  jwtSigningSecret: "1111111111111111111111111111111111111111111111111111111111111111"
  jwtIssuer: "localhost"
  jwtValidityInMinutes: 30
  maxMailsPerEmailPerDay: 10
  mail:
    mailgun:
      sender: no-reply@localhost
      domain: mg.localhost
      privateKey: "key-11111111111111111111111111111111"
  csrfAuthKey: '1111111111111111111111111111111111111111111111111111111111111111'
  legal:
    termsOfServiceUrl: "tos.pdf"
    privacyPolicyUrl: "privacy.pdf"
    imprintTemplate: '{{ define "js" }}{{ end }}{{ define "css" }}{{ end }}{{ define "content" }}Imprint{{ end }}'
  stripe:
    sapphire: price_sapphire
    emerald: price_emerald
    diamond: price_diamond
  ratelimitUpdateInterval: 1s

indexer:
  # fullIndexOnStartup: false # Perform a one time full db index on startup
  # indexMissingEpochsOnStartup: true # Check for missing epochs and export them after startup
  node:
    host: 127.0.0.1
    port: '$CL_PORT'
    type: prysm
  eth1DepositContractFirstBlock: 0
EOL

echo "generated config written to config.yml"

bash init-schema.sh
