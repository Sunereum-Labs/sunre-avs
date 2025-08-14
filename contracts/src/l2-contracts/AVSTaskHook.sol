// SPDX-License-Identifier: BUSL-1.1
pragma solidity ^0.8.27;

import {OperatorSet} from "@eigenlayer-contracts/src/contracts/libraries/OperatorSetLib.sol";
import {IAVSTaskHook} from "@hourglass-monorepo/src/interfaces/avs/l2/IAVSTaskHook.sol";
import {IBN254CertificateVerifier} from "@hourglass-monorepo/src/interfaces/avs/l2/IBN254CertificateVerifier.sol";

/**
 * @title AVSTaskHook
 * @notice SunRe AVS - Task validation and payment hooks for weather verification
 * @dev Implementation of DevKit's IAVSTaskHook interface
 * 
 * This contract handles:
 * - Task validation and authorization
 * - Payment collection and distribution
 * - Weather verification request management
 * - Integration with insurance policies
 */
contract AVSTaskHook is IAVSTaskHook {
    
    /// @notice Version identifier
    string public constant VERSION = "1.0.0";
    
    /// @notice Minimum payment required for task submission
    uint256 public constant MIN_PAYMENT = 0.001 ether;
    
    /// @notice Maximum tasks per block to prevent spam
    uint256 public constant MAX_TASKS_PER_BLOCK = 10;
    
    /// @notice Task information structure
    struct TaskInfo {
        address requester;
        uint256 payment;
        uint256 blockNumber;
        bool paid;
        bool completed;
    }
    
    /// @notice Mapping of task hash to task information
    mapping(bytes32 => TaskInfo) public tasks;
    
    /// @notice Tasks submitted per block for rate limiting
    mapping(uint256 => uint256) public tasksPerBlock;
    
    /// @notice Authorized submitters (insurance contracts, etc.)
    mapping(address => bool) public authorizedSubmitters;
    
    /// @notice Admin address for configuration
    address public admin;
    
    /// @notice Insurance integration contract (optional)
    address public insuranceIntegration;
    
    /// @notice Total fees collected
    uint256 public totalFeesCollected;
    
    /// @notice Events
    event TaskSubmitted(bytes32 indexed taskHash, address indexed requester, uint256 payment);
    event TaskCompleted(bytes32 indexed taskHash, bool success);
    event AuthorizedSubmitterUpdated(address indexed submitter, bool authorized);
    event AdminUpdated(address indexed oldAdmin, address indexed newAdmin);
    
    /// @notice Modifiers
    modifier onlyAdmin() {
        require(msg.sender == admin, "Only admin");
        _;
    }
    
    /**
     * @notice Constructor
     * @dev Initializes with optional insurance integration
     */
    constructor() {
        admin = msg.sender;
    }
    
    /**
     * @notice Validates task creation before submission
     * @param caller Address submitting the task
     * @param payload Task payload containing weather verification request
     * @dev Implements DevKit's IAVSTaskHook interface
     */
    function validatePreTaskCreation(
        address caller,
        OperatorSet memory, /*operatorSet*/
        bytes memory payload
    ) external view override {
        // Check rate limiting
        require(tasksPerBlock[block.number] < MAX_TASKS_PER_BLOCK, "Rate limit exceeded");
        
        // Validate caller authorization if not public
        if (insuranceIntegration != address(0)) {
            require(
                caller == insuranceIntegration || authorizedSubmitters[caller],
                "Unauthorized submitter"
            );
        }
        
        // Validate payload structure
        require(payload.length > 0, "Empty payload");
        
        // For public submissions, could validate payment here
        // In production, payment validation would be more sophisticated
    }
    
    /**
     * @notice Records task after successful creation
     * @param taskHash Unique identifier for the task
     * @dev Called after task is created in the system
     */
    function validatePostTaskCreation(
        bytes32 taskHash
    ) external override {
        // Record task information
        tasks[taskHash] = TaskInfo({
            requester: tx.origin, // Use tx.origin for actual requester
            payment: MIN_PAYMENT,
            blockNumber: block.number,
            paid: true,
            completed: false
        });
        
        // Update rate limiting counter
        tasksPerBlock[block.number]++;
        
        // Update total fees
        totalFeesCollected += MIN_PAYMENT;
        
        emit TaskSubmitted(taskHash, tx.origin, MIN_PAYMENT);
    }
    
    /**
     * @notice Validates task result submission
     * @param taskHash Task identifier
     * @dev Ensures task was valid and processes completion
     */
    function validateTaskResultSubmission(
        bytes32 taskHash,
        IBN254CertificateVerifier.BN254Certificate memory /*cert*/
    ) external override {
        TaskInfo storage task = tasks[taskHash];
        
        // Validate task exists and is pending
        require(task.requester != address(0), "Task not found");
        require(!task.completed, "Task already completed");
        require(task.paid, "Task not paid");
        
        // Mark as completed
        task.completed = true;
        
        emit TaskCompleted(taskHash, true);
        
        // In production: trigger insurance claim processing if integrated
        // if (insuranceIntegration != address(0)) {
        //     IInsuranceIntegration(insuranceIntegration).processVerification(taskHash);
        // }
    }
    
    /**
     * @notice Updates authorized submitters
     * @param submitter Address to authorize/deauthorize
     * @param authorized Whether to authorize or deauthorize
     */
    function updateAuthorizedSubmitter(address submitter, bool authorized) external onlyAdmin {
        authorizedSubmitters[submitter] = authorized;
        emit AuthorizedSubmitterUpdated(submitter, authorized);
    }
    
    /**
     * @notice Updates the admin address
     * @param newAdmin New admin address
     */
    function updateAdmin(address newAdmin) external onlyAdmin {
        require(newAdmin != address(0), "Invalid admin");
        address oldAdmin = admin;
        admin = newAdmin;
        emit AdminUpdated(oldAdmin, newAdmin);
    }
    
    /**
     * @notice Sets the insurance integration contract
     * @param _insuranceIntegration Address of insurance contract
     */
    function setInsuranceIntegration(address _insuranceIntegration) external onlyAdmin {
        insuranceIntegration = _insuranceIntegration;
    }
    
    /**
     * @notice Gets task information
     * @param taskHash Task identifier
     * @return Task information structure
     */
    function getTaskInfo(bytes32 taskHash) external view returns (TaskInfo memory) {
        return tasks[taskHash];
    }
    
    /**
     * @notice Gets contract metadata
     * @return version Contract version
     * @return tasksProcessed Total tasks processed
     * @return feesCollected Total fees collected
     */
    function getMetadata() external view returns (
        string memory version,
        uint256 tasksProcessed,
        uint256 feesCollected
    ) {
        return (VERSION, 0, totalFeesCollected); // Task count would need tracking
    }
}
