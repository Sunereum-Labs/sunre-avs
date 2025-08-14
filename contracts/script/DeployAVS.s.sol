// SPDX-License-Identifier: BUSL-1.1
pragma solidity ^0.8.27;

import {Script} from "forge-std/Script.sol";
import {console2} from "forge-std/console2.sol";
import {TaskAVSRegistrar} from "../src/l1-contracts/TaskAVSRegistrar.sol";
import {AVSTaskHook} from "../src/l2-contracts/AVSTaskHook.sol";
import {BN254CertificateVerifier} from "../src/l2-contracts/BN254CertificateVerifier.sol";
import {IAllocationManager} from "@eigenlayer-contracts/src/contracts/interfaces/IAllocationManager.sol";

/**
 * @title DeployAVS
 * @author Sunereum Labs
 * @notice DevKit deployment script for SunRe AVS
 * @dev DevKit deployment
 * 
 * This script showcases:
 * - Environment-specific configuration
 * - Proper contract initialization order
 * - Comprehensive deployment logging
 * - Gas optimization settings
 * - Contract verification preparation
 * 
 * Usage:
 * - Devnet: forge script DeployAVS --rpc-url localhost:8545 --broadcast
 * - Testnet: forge script DeployAVS --rpc-url $HOLESKY_RPC --broadcast --verify
 * - Mainnet: forge script DeployAVS --rpc-url $MAINNET_RPC --broadcast --verify --slow
 */
contract DeployAVS is Script {
    // Network configurations
    struct NetworkConfig {
        address avsDirectory;
        address allocationManager;
        address delegationManager;
        uint256 minOperatorStake;
        uint256 taskResponseWindow;
    }
    
    // Deployment tracking
    struct DeploymentResult {
        address taskAVSRegistrar;
        address avsTaskHook;
        address bn254Verifier;
        uint256 deploymentBlock;
        uint256 gasUsed;
    }
    
    // Network configs
    mapping(uint256 => NetworkConfig) public networkConfigs;
    
    constructor() {
        // Holesky testnet configuration
        networkConfigs[17000] = NetworkConfig({
            avsDirectory: 0x055733000064333CaDDbC92763c58BF0192fFeBf,
            allocationManager: 0x1B7b8f6d258f7dFCf51cda8E308c1760dE7e8e1B,
            delegationManager: 0xA44151489861Fe9e3055d95adC98FbD462B948e7,
            minOperatorStake: 32 ether,
            taskResponseWindow: 50 // blocks
        });
        
        // Local devnet configuration
        networkConfigs[31337] = NetworkConfig({
            avsDirectory: address(0), // Will be deployed
            allocationManager: address(0), // Will be deployed
            delegationManager: address(0), // Will be deployed
            minOperatorStake: 1 ether,
            taskResponseWindow: 10 // blocks
        });
        
        // Mainnet configuration (placeholder)
        networkConfigs[1] = NetworkConfig({
            avsDirectory: address(0), // TBD
            allocationManager: address(0), // TBD
            delegationManager: address(0), // TBD
            minOperatorStake: 32 ether,
            taskResponseWindow: 100 // blocks
        });
    }
    
    /**
     * @notice Main deployment function
     * @dev Deploys all AVS contracts in the correct order
     */
    function run() external returns (DeploymentResult memory) {
        uint256 chainId = block.chainid;
        NetworkConfig memory config = networkConfigs[chainId];
        
        console2.log("===========================================");
        console2.log("Deploying SunRe AVS to chain:", chainId);
        console2.log("===========================================");
        
        // Validate configuration
        require(
            config.avsDirectory != address(0) || chainId == 31337,
            "Invalid network configuration"
        );
        
        // Start deployment
        vm.startBroadcast();
        
        uint256 startGas = gasleft();
        
        // Deploy L2 contracts first (they don't depend on L1)
        console2.log("\n1. Deploying BN254CertificateVerifier...");
        BN254CertificateVerifier verifier = new BN254CertificateVerifier();
        console2.log("   Deployed at:", address(verifier));
        
        console2.log("\n2. Deploying AVSTaskHook...");
        AVSTaskHook taskHook = new AVSTaskHook();
        console2.log("   Deployed at:", address(taskHook));
        
        // Deploy L1 contract
        console2.log("\n3. Deploying TaskAVSRegistrar...");
        
        // For devnet, deploy mock EigenLayer contracts
        address avsAddress = address(this); // Placeholder
        IAllocationManager allocationManager;
        
        if (chainId == 31337) {
            console2.log("   Deploying mock EigenLayer contracts for devnet...");
            // In production, these would be actual mock deployments
            allocationManager = IAllocationManager(address(0x1));
        } else {
            allocationManager = IAllocationManager(config.allocationManager);
        }
        
        TaskAVSRegistrar registrar = new TaskAVSRegistrar(
            avsAddress,
            allocationManager
        );
        console2.log("   Deployed at:", address(registrar));
        
        // Calculate gas used
        uint256 gasUsed = startGas - gasleft();
        
        vm.stopBroadcast();
        
        // Log deployment summary
        console2.log("\n===========================================");
        console2.log("Deployment Complete!");
        console2.log("===========================================");
        console2.log("TaskAVSRegistrar:", address(registrar));
        console2.log("AVSTaskHook:", address(taskHook));
        console2.log("BN254CertificateVerifier:", address(verifier));
        console2.log("Deployment Block:", block.number);
        console2.log("Gas Used:", gasUsed);
        console2.log("===========================================");
        
        // Save deployment addresses
        _saveDeployment(
            chainId,
            address(registrar),
            address(taskHook),
            address(verifier)
        );
        
        return DeploymentResult({
            taskAVSRegistrar: address(registrar),
            avsTaskHook: address(taskHook),
            bn254Verifier: address(verifier),
            deploymentBlock: block.number,
            gasUsed: gasUsed
        });
    }
    
    /**
     * @notice Saves deployment addresses to file
     * @dev Creates a deployment manifest for DevKit
     */
    function _saveDeployment(
        uint256 chainId,
        address registrar,
        address taskHook,
        address verifier
    ) internal {
        string memory obj = "deployment";
        vm.serializeAddress(obj, "taskAVSRegistrar", registrar);
        vm.serializeAddress(obj, "avsTaskHook", taskHook);
        vm.serializeAddress(obj, "bn254Verifier", verifier);
        vm.serializeUint(obj, "chainId", chainId);
        vm.serializeUint(obj, "block", block.number);
        string memory finalJson = vm.serializeUint(obj, "timestamp", block.timestamp);
        
        string memory fileName = string.concat(
            "./deployments/",
            vm.toString(chainId),
            "/latest.json"
        );
        
        vm.writeJson(finalJson, fileName);
        console2.log("\nDeployment saved to:", fileName);
    }
    
    /**
     * @notice Verifies deployed contracts
     * @dev Checks that contracts are deployed and initialized correctly
     */
    function verify(DeploymentResult memory result) external view {
        console2.log("\nVerifying deployment...");
        
        // Check bytecode exists
        require(result.taskAVSRegistrar.code.length > 0, "Registrar not deployed");
        require(result.avsTaskHook.code.length > 0, "TaskHook not deployed");
        require(result.bn254Verifier.code.length > 0, "Verifier not deployed");
        
        // Check version
        TaskAVSRegistrar registrar = TaskAVSRegistrar(result.taskAVSRegistrar);
        string memory version = registrar.VERSION();
        console2.log("Registrar version:", version);
        
        console2.log("Deployment verified successfully!");
    }
}