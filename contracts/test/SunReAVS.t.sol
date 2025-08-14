// SPDX-License-Identifier: BUSL-1.1
pragma solidity ^0.8.27;

import {Test} from "forge-std/Test.sol";
import {console2} from "forge-std/console2.sol";
import {TaskAVSRegistrar} from "../src/l1-contracts/TaskAVSRegistrar.sol";
import {AVSTaskHook} from "../src/l2-contracts/AVSTaskHook.sol";
import {BN254CertificateVerifier} from "../src/l2-contracts/BN254CertificateVerifier.sol";
import {IAllocationManager} from "@eigenlayer-contracts/src/contracts/interfaces/IAllocationManager.sol";
import {OperatorSet} from "@eigenlayer-contracts/src/contracts/libraries/OperatorSetLib.sol";
import {IBN254CertificateVerifier} from "@hourglass-monorepo/src/interfaces/avs/l2/IBN254CertificateVerifier.sol";

/**
 * @title SunReAVS Tests
 * @author Sunereum Labs
 * @notice Comprehensive test suite for SunRe AVS
 * @dev Demonstrates DevKit testing best practices
 * 
 * Test coverage includes:
 * - Contract deployment and initialization
 * - Task submission and validation
 * - Operator registration and management
 * - BLS signature verification
 * - Economic model validation
 * - Security checks and edge cases
 */
contract SunReAVSTest is Test {
    // Contracts
    TaskAVSRegistrar public registrar;
    AVSTaskHook public taskHook;
    BN254CertificateVerifier public verifier;
    
    // Test accounts
    address public deployer = address(0x1);
    address public operator1 = address(0x2);
    address public operator2 = address(0x3);
    address public operator3 = address(0x4);
    address public user = address(0x5);
    address public attacker = address(0x666);
    
    // Mock addresses
    address public mockAVS = address(0x100);
    IAllocationManager public mockAllocationManager = IAllocationManager(address(0x101));
    
    // Test constants
    uint256 constant MIN_STAKE = 32 ether;
    uint256 constant TASK_FEE = 0.001 ether;
    uint256 constant QUORUM_THRESHOLD = 67; // 67%
    
    // Events to test
    event AVSInitialized(address indexed avs, address indexed allocationManager, uint256 timestamp);
    event TaskSubmitted(bytes32 indexed taskHash, address indexed requester, uint256 payment);
    event TaskCompleted(bytes32 indexed taskHash, bool success);
    
    /**
     * @notice Test setup
     */
    function setUp() public {
        // Setup test environment
        vm.label(deployer, "Deployer");
        vm.label(operator1, "Operator1");
        vm.label(operator2, "Operator2");
        vm.label(operator3, "Operator3");
        vm.label(user, "User");
        vm.label(attacker, "Attacker");
        
        // Fund accounts
        vm.deal(deployer, 100 ether);
        vm.deal(operator1, 100 ether);
        vm.deal(operator2, 100 ether);
        vm.deal(operator3, 100 ether);
        vm.deal(user, 10 ether);
        vm.deal(attacker, 1 ether);
        
        // Deploy contracts
        vm.startPrank(deployer);
        
        registrar = new TaskAVSRegistrar(mockAVS, mockAllocationManager);
        taskHook = new AVSTaskHook();
        verifier = new BN254CertificateVerifier();
        
        vm.stopPrank();
    }
    
    /**
     * @notice Test contract deployment
     */
    function testDeployment() public {
        assertEq(registrar.VERSION(), "1.0.0", "Invalid registrar version");
        assertGt(registrar.deploymentTimestamp(), 0, "Invalid deployment timestamp");
        
        // Verify initialization event
        vm.expectEmit(true, true, false, true);
        emit AVSInitialized(mockAVS, address(mockAllocationManager), block.timestamp);
        
        new TaskAVSRegistrar(mockAVS, mockAllocationManager);
    }
    
    /**
     * @notice Test invalid deployment parameters
     */
    function testInvalidDeployment() public {
        // Test zero AVS address
        vm.expectRevert("Invalid AVS address");
        new TaskAVSRegistrar(address(0), mockAllocationManager);
        
        // Test zero AllocationManager
        vm.expectRevert("Invalid AllocationManager");
        new TaskAVSRegistrar(mockAVS, IAllocationManager(address(0)));
    }
    
    /**
     * @notice Test task submission
     */
    function testTaskSubmission() public {
        vm.startPrank(user);
        
        // Prepare task payload
        bytes memory payload = abi.encode(
            int256(40), // latitude (NYC)
            int256(-74), // longitude
            block.timestamp,
            "POL-TEST-001"
        );
        
        // Submit task through hook
        OperatorSet memory operatorSet;
        taskHook.validatePreTaskCreation(user, operatorSet, payload);
        
        // Verify task creation
        bytes32 taskHash = keccak256(payload);
        taskHook.validatePostTaskCreation(taskHash);
        
        // Check task info
        AVSTaskHook.TaskInfo memory taskInfo = taskHook.getTaskInfo(taskHash);
        assertEq(taskInfo.requester, user, "Invalid requester");
        assertEq(taskInfo.payment, TASK_FEE, "Invalid payment");
        assertTrue(taskInfo.paid, "Task not marked as paid");
        assertFalse(taskInfo.completed, "Task should not be completed");
        
        vm.stopPrank();
    }
    
    /**
     * @notice Test rate limiting
     */
    function testRateLimiting() public {
        vm.startPrank(user);
        
        bytes memory payload = abi.encode(int256(40), int256(-74), block.timestamp, "POL-TEST");
        OperatorSet memory operatorSet;
        
        // Submit maximum allowed tasks
        for (uint i = 0; i < 10; i++) {
            bytes memory uniquePayload = abi.encode(int256(40), int256(-74), block.timestamp, i);
            taskHook.validatePreTaskCreation(user, operatorSet, uniquePayload);
            taskHook.validatePostTaskCreation(keccak256(uniquePayload));
        }
        
        // 11th task should fail
        vm.expectRevert("Rate limit exceeded");
        taskHook.validatePreTaskCreation(user, operatorSet, payload);
        
        vm.stopPrank();
    }
    
    /**
     * @notice Test BLS certificate verification
     */
    function testBLSVerification() public {
        // Create mock certificate
        IBN254CertificateVerifier.BN254Certificate memory cert;
        // cert.taskHash = keccak256("test"); // Field not available in interface
        cert.nonsignerIndices = new uint32[](0); // All operators signed
        
        // Verify certificate
        uint96[] memory stakes = verifier.verifyCertificate(cert);
        assertGt(stakes.length, 0, "No stakes returned");
        
        // Test proportion verification
        uint16[] memory thresholds = new uint16[](1);
        thresholds[0] = 6700; // 67%
        assertTrue(verifier.verifyCertificateProportion(cert, thresholds), "Proportion check failed");
        
        // Test nominal verification
        uint96[] memory nominalThresholds = new uint96[](1);
        nominalThresholds[0] = 1000 ether;
        assertTrue(verifier.verifyCertificateNominal(cert, nominalThresholds), "Nominal check failed");
    }
    
    /**
     * @notice Test weather consensus validation
     */
    function testWeatherConsensus() public {
        // Test valid consensus
        assertTrue(verifier.validateWeatherConsensus(3, 3000 ether), "Valid consensus rejected");
        
        // Test insufficient operators
        assertFalse(verifier.validateWeatherConsensus(2, 3000 ether), "Insufficient operators accepted");
        
        // Test insufficient stake
        assertFalse(verifier.validateWeatherConsensus(3, 2000 ether), "Insufficient stake accepted");
    }
    
    /**
     * @notice Test task completion
     */
    function testTaskCompletion() public {
        // Setup task
        bytes32 taskHash = keccak256("test-task");
        vm.prank(user);
        taskHook.validatePostTaskCreation(taskHash);
        
        // Complete task
        IBN254CertificateVerifier.BN254Certificate memory cert;
        // cert.taskHash = taskHash; // Field not available in interface
        
        vm.expectEmit(true, false, false, true);
        emit TaskCompleted(taskHash, true);
        
        taskHook.validateTaskResultSubmission(taskHash, cert);
        
        // Verify completion
        AVSTaskHook.TaskInfo memory taskInfo = taskHook.getTaskInfo(taskHash);
        assertTrue(taskInfo.completed, "Task not marked as completed");
    }
    
    /**
     * @notice Test unauthorized access
     */
    function testUnauthorizedAccess() public {
        vm.startPrank(attacker);
        
        // Try to update admin
        vm.expectRevert("Only admin");
        taskHook.updateAdmin(attacker);
        
        // Try to update authorized submitter
        vm.expectRevert("Only admin");
        taskHook.updateAuthorizedSubmitter(attacker, true);
        
        vm.stopPrank();
    }
    
    /**
     * @notice Test admin functions
     */
    function testAdminFunctions() public {
        vm.startPrank(deployer);
        
        // Update authorized submitter
        taskHook.updateAuthorizedSubmitter(operator1, true);
        assertTrue(taskHook.authorizedSubmitters(operator1), "Submitter not authorized");
        
        // Update admin
        taskHook.updateAdmin(operator1);
        assertEq(taskHook.admin(), operator1, "Admin not updated");
        
        vm.stopPrank();
    }
    
    /**
     * @notice Test gas optimization
     */
    function testGasOptimization() public {
        uint256 gasStart = gasleft();
        
        // Deploy contracts
        new TaskAVSRegistrar(mockAVS, mockAllocationManager);
        new AVSTaskHook();
        new BN254CertificateVerifier();
        
        uint256 gasUsed = gasStart - gasleft();
        console2.log("Total deployment gas:", gasUsed);
        
        // Ensure deployment is gas-efficient
        assertLt(gasUsed, 5000000, "Deployment too gas-intensive");
    }
    
    /**
     * @notice Fuzz test for task validation
     */
    function testFuzzTaskValidation(
        uint256 latitude,
        uint256 longitude,
        uint256 timestamp,
        string memory policyId
    ) public {
        // Bound inputs to valid ranges
        latitude = bound(latitude, 0, 180) - 90; // -90 to 90
        longitude = bound(longitude, 0, 360) - 180; // -180 to 180
        timestamp = bound(timestamp, block.timestamp - 86400, block.timestamp + 86400);
        
        bytes memory payload = abi.encode(latitude, longitude, timestamp, policyId);
        OperatorSet memory operatorSet;
        
        // Should not revert for valid inputs
        vm.prank(user);
        taskHook.validatePreTaskCreation(user, operatorSet, payload);
    }
    
    /**
     * @notice Integration test for full task lifecycle
     */
    function testFullTaskLifecycle() public {
        console2.log("Starting full task lifecycle test...");
        
        // 1. Submit task
        vm.startPrank(user);
        bytes memory payload = abi.encode(int256(40), int256(-74), block.timestamp, "POL-INT-001");
        bytes32 taskHash = keccak256(payload);
        
        OperatorSet memory operatorSet;
        taskHook.validatePreTaskCreation(user, operatorSet, payload);
        taskHook.validatePostTaskCreation(taskHash);
        console2.log("Task submitted:", vm.toString(taskHash));
        
        // 2. Operators process task (simulated)
        vm.warp(block.timestamp + 60); // Fast forward 1 minute
        
        // 3. Submit result with BLS certificate
        IBN254CertificateVerifier.BN254Certificate memory cert;
        // cert.taskHash = taskHash; // Field not available in interface
        cert.nonsignerIndices = new uint32[](0); // All signed
        
        taskHook.validateTaskResultSubmission(taskHash, cert);
        console2.log("Task completed successfully");
        
        // 4. Verify final state
        AVSTaskHook.TaskInfo memory taskInfo = taskHook.getTaskInfo(taskHash);
        assertTrue(taskInfo.completed, "Task not completed");
        assertEq(taskInfo.requester, user, "Invalid requester");
        
        console2.log("Full lifecycle test passed!");
        vm.stopPrank();
    }
}