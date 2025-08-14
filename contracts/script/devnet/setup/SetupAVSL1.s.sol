// SPDX-License-Identifier: BUSL-1.1
pragma solidity ^0.8.27;

import {Script, console} from "forge-std/Script.sol";
import {stdJson} from "forge-std/StdJson.sol";
import {
    IAllocationManager,
    IAllocationManagerTypes
} from "@eigenlayer-contracts/src/contracts/interfaces/IAllocationManager.sol";
import {IAVSRegistrar} from "@eigenlayer-contracts/src/contracts/interfaces/IAVSRegistrar.sol";
import {IStrategy} from "@eigenlayer-contracts/src/contracts/interfaces/IStrategy.sol";

contract SetupAVSL1 is Script {
    using stdJson for string;

    function run(
        string memory environment,
        address allocationManager,
        string memory metadataURI,
        uint32 aggregatorOperatorSetId,
        address[] memory aggregatorStrategies,
        uint32 executorOperatorSetId,
        address[] memory executorStrategies
    ) public {
        // Load config and get addresses
        address taskAVSRegistrar = _readConfigAddress(environment, "taskAVSRegistrar");
        console.log("Task AVS Registrar:", taskAVSRegistrar);

        // Load the private key from the environment variable
        uint256 avsPrivateKey = vm.envUint("PRIVATE_KEY_AVS");
        address avs = vm.addr(avsPrivateKey);

        vm.startBroadcast(avsPrivateKey);
        console.log("AVS address:", avs);

        // 1. Update the AVS metadata URI
        IAllocationManager(allocationManager).updateAVSMetadataURI(avs, metadataURI);
        console.log("AVS metadata URI updated:", metadataURI);

        // 2. Set the AVS Registrar
        IAllocationManager(allocationManager).setAVSRegistrar(avs, IAVSRegistrar(taskAVSRegistrar));
        console.log("AVS Registrar set:", address(IAllocationManager(allocationManager).getAVSRegistrar(avs)));

        // 3. Create the operator sets
        IStrategy[] memory aggregatorStrategiesArray = new IStrategy[](aggregatorStrategies.length);
        for (uint256 i = 0; i < aggregatorStrategies.length; i++) {
            aggregatorStrategiesArray[i] = IStrategy(aggregatorStrategies[i]);
        }
        IStrategy[] memory executorStrategiesArray = new IStrategy[](executorStrategies.length);
        for (uint256 i = 0; i < executorStrategies.length; i++) {
            executorStrategiesArray[i] = IStrategy(executorStrategies[i]);
        }
        IAllocationManagerTypes.CreateSetParams[] memory createOperatorSetParams =
            new IAllocationManagerTypes.CreateSetParams[](2);
        createOperatorSetParams[0] = IAllocationManagerTypes.CreateSetParams({
            operatorSetId: aggregatorOperatorSetId,
            strategies: aggregatorStrategiesArray
        });
        createOperatorSetParams[1] = IAllocationManagerTypes.CreateSetParams({
            operatorSetId: executorOperatorSetId,
            strategies: executorStrategiesArray
        });
        IAllocationManager(allocationManager).createOperatorSets(avs, createOperatorSetParams);
        console.log("Operator sets created: ", IAllocationManager(allocationManager).getOperatorSetCount(avs));

        vm.stopBroadcast();
    }

    function _readConfigAddress(string memory environment, string memory key) internal view returns (address) {
        // Load the output file
        string memory avsL1ConfigFile = string.concat("script/", environment, "/output/deploy_avs_l1_output.json");
        string memory avsL1Config = vm.readFile(avsL1ConfigFile);

        // Parse the address
        return stdJson.readAddress(avsL1Config, string.concat(".addresses.", key));
    }
}
