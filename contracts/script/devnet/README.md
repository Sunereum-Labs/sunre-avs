# AVS Contracts Script Setup

This README provides step-by-step instructions for setting up and deploying the AVS contracts using the scripts in this directory. These commands are intended for local development and testing using a forked Ethereum mainnet environment.

## Prerequisites

- [Foundry](https://book.getfoundry.sh/) installed
- [Anvil](https://book.getfoundry.sh/anvil/) (comes with Foundry)
- `make` installed
- Access to a mainnet RPC URL (e.g., from [QuikNode](https://quiknode.io/), [Alchemy](https://www.alchemy.com/), or [Infura](https://infura.io/))

## Step-by-Step Setup

### 1. Start a Forked Mainnet Node

Start a local Anvil node forked from Ethereum mainnet at a specific block (This includes the Eigenlayer Protocol already deployed):

```sh
anvil --fork-url <MAINNET_RPC_URL>
```

This will run a local node at `127.0.0.1:8545`.

### 2. Deploy TaskMailbox Contract

In a new terminal, deploy the TaskMailbox contract to your local fork:

```sh
make deploy-task-mailbox \
  RPC_URL="127.0.0.1:8545" \
  ENVIRONMENT="local"
```

### 3. Deploy AVS L1 Contracts

Deploy the AVS L1 contracts, specifying the AVS address and the Allocation Manager address:

```sh
make deploy-avs-l1-contracts \
  RPC_URL="127.0.0.1:8545" \
  ENVIRONMENT="local" \
  AVS_ADDRESS='0x70997970C51812dc3A010C7d01b50e0d17dc79C8' \
  ALLOCATION_MANAGER_ADDRESS='0x948a420b8CC1d6BFd0B6087C2E7c344a2CD0bc39'
```

### 4. Setup AVS L1

Register the AVS L1 contracts with the EigenLayer core protocol:

```sh
make setup-avs-l1 \
  RPC_URL="127.0.0.1:8545" \
  ENVIRONMENT="local" \
  ALLOCATION_MANAGER_ADDRESS='0x948a420b8CC1d6BFd0B6087C2E7c344a2CD0bc39' \
  METADATA_URI="TestAVS" \
  AGGREGATOR_OPERATOR_SET_ID=0 \
  AGGREGATOR_STRATEGIES='["0xaCB55C530Acdb2849e6d4f36992Cd8c9D50ED8F7","0x93c4b944D05dfe6df7645A86cd2206016c51564D"]' \
  EXECUTOR_OPERATOR_SET_ID=1 \
  EXECUTOR_STRATEGIES='["0xaCB55C530Acdb2849e6d4f36992Cd8c9D50ED8F7","0x93c4b944D05dfe6df7645A86cd2206016c51564D"]'
```

### 5. Deploy AVS L2 Contracts

Deploy the AVS L2 contracts:

```sh
make deploy-avs-l2-contracts \
  RPC_URL="127.0.0.1:8545" \
  ENVIRONMENT="local"
```

### 6. Setup AVS Task Mailbox Config

Configure the TaskMailbox with the required addresses:

```sh
make setup-avs-task-mailbox-config \
  RPC_URL="127.0.0.1:8545" \
  ENVIRONMENT="local" \
  AGGREGATOR_OPERATOR_SET_ID=0 \
  EXECUTOR_OPERATOR_SET_ID=1 \
  TASK_SLA=60
```

### 7. Create Task

Create a Task on L2 Mailbox

```sh
make create-task \
  RPC_URL="127.0.0.1:8545" \
  ENVIRONMENT="local" \
  AVS_ADDRESS='0x70997970C51812dc3A010C7d01b50e0d17dc79C8' \
  EXECUTOR_OPERATOR_SET_ID=1 \
  PAYLOAD='0x0000000000000000000000000000000000000000000000000000000000000005'
```

## Environment Variables

Some scripts require environment variables, such as `PRIVATE_KEY_AVS`, to be set. Refer to the `.env.example` file in the parent directory for required variables and copy it as `.env`. Get the private keys from the Anvil bootup message.

## Script Reference

- `DeployTaskMailbox.s.sol`: Deploys the TaskMailbox contract
- `DeployAVSL1Contracts.s.sol`: Deploys AVS L1 contracts
- `SetupAVSL1.s.sol`: Registers AVS L1 contracts
- `DeployAVSL2Contracts.s.sol`: Deploys AVS L2 contracts
- `SetupAVSTaskMailboxConfig.s.sol`: Configures the TaskMailbox with AVS and verifier addresses

## Additional Notes

- Ensure your local node is running before executing deployment or setup commands.
- The provided addresses are for local development and testing only.
- For more details on each script, review the source files in this directory. 