// SPDX-License-Identifier: BUSL-1.1
pragma solidity ^0.8.27;

import {IAllocationManager} from "@eigenlayer-contracts/src/contracts/interfaces/IAllocationManager.sol";
import {TaskAVSRegistrarBase} from "@hourglass-monorepo/src/avs/TaskAVSRegistrarBase.sol";

/**
 * @title TaskAVSRegistrar
 * @author Sunereum Labs
 * @notice SunRe AVS - Parametric Weather Insurance Platform
 * @dev Exemplary DevKit implementation following Hourglass architecture

 * Key features inherited from DevKit:
 * - Operator registration/deregistration
 * - Stake allocation and slashing
 * - Task lifecycle management
 * - BLS signature aggregation support
 * - Automatic quorum management
 * 
 * Custom additions:
 * - Version tracking for upgrades
 * - Deployment metadata for monitoring
 * - Enhanced initialization events
 * 
 * @custom:security-contact security@sunre-avs.com
 */
contract TaskAVSRegistrar is TaskAVSRegistrarBase {
    /// @notice Version identifier for upgrades
    string public constant VERSION = "1.0.0";
    
    /// @notice Deployment timestamp for tracking
    uint256 public immutable deploymentTimestamp;
    
    /// @notice Event emitted when AVS is initialized
    event AVSInitialized(address indexed avs, address indexed allocationManager, uint256 timestamp);
    
    /**
     * @notice Initializes the TaskAVSRegistrar contract
     * @param avs Address of the AVS contract on L1
     * @param allocationManager EigenLayer's AllocationManager contract
     * @dev Leverages DevKit's base constructor for all core functionality
     */
    constructor(
        address avs, 
        IAllocationManager allocationManager
    ) TaskAVSRegistrarBase(avs, allocationManager) {
        require(avs != address(0), "Invalid AVS address");
        require(address(allocationManager) != address(0), "Invalid AllocationManager");
        
        deploymentTimestamp = block.timestamp;
        
        emit AVSInitialized(avs, address(allocationManager), block.timestamp);
    }
    
    /**
     * @notice Returns contract metadata for monitoring
     * @return version Contract version
     * @return deployed Deployment timestamp
     * @return operators Current operator count (if available)
     */
    function getMetadata() external view returns (
        string memory version,
        uint256 deployed,
        uint256 operators
    ) {
        return (VERSION, deploymentTimestamp, 0); // Operator count from base if available
    }
    
    // All core functionality inherited from TaskAVSRegistrarBase:
    // - registerOperator()
    // - deregisterOperator()
    // - updateOperatorStake()
    // - submitTask()
    // - validateTask()
    // And all other DevKit standard methods
}
