// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

import {BaseTest} from "./BaseTest.t.sol";
import {CapabilityRegistry} from "../CapabilityRegistry.sol";

contract CapabilityRegistry_RemoveNodesTest is BaseTest {
  event NodeRemoved(bytes32 p2pId);

  uint32 private constant TEST_NODE_OPERATOR_ONE_ID = 1;
  uint256 private constant TEST_NODE_OPERATOR_TWO_ID = 2;
  bytes32 private constant INVALID_P2P_ID = bytes32("fake-p2p");

  function setUp() public override {
    BaseTest.setUp();
    changePrank(ADMIN);
    CapabilityRegistry.Capability[] memory capabilities = new CapabilityRegistry.Capability[](2);
    capabilities[0] = s_basicCapability;
    capabilities[1] = s_capabilityWithConfigurationContract;

    s_capabilityRegistry.addNodeOperators(_getNodeOperators());
    s_capabilityRegistry.addCapabilities(capabilities);

    CapabilityRegistry.NodeInfo[] memory nodes = new CapabilityRegistry.NodeInfo[](1);
    bytes32[] memory hashedCapabilityIds = new bytes32[](2);
    hashedCapabilityIds[0] = s_basicHashedCapabilityId;
    hashedCapabilityIds[1] = s_capabilityWithConfigurationContractId;

    nodes[0] = CapabilityRegistry.NodeInfo({
      nodeOperatorId: TEST_NODE_OPERATOR_ONE_ID,
      p2pId: P2P_ID,
      signer: NODE_OPERATOR_ONE_SIGNER_ADDRESS,
      hashedCapabilityIds: hashedCapabilityIds
    });

    changePrank(NODE_OPERATOR_ONE_ADMIN);

    s_capabilityRegistry.addNodes(nodes);
  }

  function test_RevertWhen_CalledByNonNodeOperatorAdminAndNonOwner() public {
    changePrank(STRANGER);
    bytes32[] memory nodes = new bytes32[](1);
    nodes[0] = P2P_ID;

    vm.expectRevert(CapabilityRegistry.AccessForbidden.selector);
    s_capabilityRegistry.removeNodes(nodes);
  }

  function test_RevertWhen_NodeDoesNotExist() public {
    changePrank(NODE_OPERATOR_ONE_ADMIN);
    bytes32[] memory nodes = new bytes32[](1);
    nodes[0] = INVALID_P2P_ID;

    vm.expectRevert(abi.encodeWithSelector(CapabilityRegistry.InvalidNodeP2PId.selector, INVALID_P2P_ID));
    s_capabilityRegistry.removeNodes(nodes);
  }

  function test_RevertWhen_P2PIDEmpty() public {
    changePrank(NODE_OPERATOR_ONE_ADMIN);
    bytes32[] memory nodes = new bytes32[](1);
    nodes[0] = bytes32("");

    vm.expectRevert(abi.encodeWithSelector(CapabilityRegistry.InvalidNodeP2PId.selector, bytes32("")));
    s_capabilityRegistry.removeNodes(nodes);
  }

  function test_RemovesNode() public {
    changePrank(NODE_OPERATOR_ONE_ADMIN);

    bytes32[] memory nodes = new bytes32[](1);
    nodes[0] = P2P_ID;

    vm.expectEmit(address(s_capabilityRegistry));
    emit NodeRemoved(P2P_ID);
    s_capabilityRegistry.removeNodes(nodes);

    (CapabilityRegistry.NodeInfo memory node, uint32 configCount) = s_capabilityRegistry.getNode(P2P_ID);
    assertEq(node.nodeOperatorId, 0);
    assertEq(node.p2pId, bytes32(""));
    assertEq(node.signer, bytes32(""));
    assertEq(node.hashedCapabilityIds.length, 0);
    assertEq(configCount, 0);
  }

  function test_CanAddNodeWithSameSignerAddressAfterRemoving() public {
    changePrank(NODE_OPERATOR_ONE_ADMIN);

    bytes32[] memory nodes = new bytes32[](1);
    nodes[0] = P2P_ID;

    s_capabilityRegistry.removeNodes(nodes);

    CapabilityRegistry.NodeInfo[] memory NodeInfo = new CapabilityRegistry.NodeInfo[](1);
    bytes32[] memory hashedCapabilityIds = new bytes32[](2);
    hashedCapabilityIds[0] = s_basicHashedCapabilityId;
    hashedCapabilityIds[1] = s_capabilityWithConfigurationContractId;

    NodeInfo[0] = CapabilityRegistry.NodeInfo({
      nodeOperatorId: TEST_NODE_OPERATOR_ONE_ID,
      p2pId: P2P_ID,
      signer: NODE_OPERATOR_ONE_SIGNER_ADDRESS,
      hashedCapabilityIds: hashedCapabilityIds
    });

    s_capabilityRegistry.addNodes(NodeInfo);

    (CapabilityRegistry.NodeInfo memory node, uint32 configCount) = s_capabilityRegistry.getNode(P2P_ID);
    assertEq(node.nodeOperatorId, TEST_NODE_OPERATOR_ONE_ID);
    assertEq(node.p2pId, P2P_ID);
    assertEq(node.hashedCapabilityIds.length, 2);
    assertEq(node.hashedCapabilityIds[0], s_basicHashedCapabilityId);
    assertEq(node.hashedCapabilityIds[1], s_capabilityWithConfigurationContractId);
    assertEq(configCount, 1);
  }

  function test_OwnerCanRemoveNodes() public {
    changePrank(ADMIN);

    bytes32[] memory nodes = new bytes32[](1);
    nodes[0] = P2P_ID;

    vm.expectEmit(address(s_capabilityRegistry));
    emit NodeRemoved(P2P_ID);
    s_capabilityRegistry.removeNodes(nodes);

    (CapabilityRegistry.NodeInfo memory node, uint32 configCount) = s_capabilityRegistry.getNode(P2P_ID);
    assertEq(node.nodeOperatorId, 0);
    assertEq(node.p2pId, bytes32(""));
    assertEq(node.signer, bytes32(""));
    assertEq(node.hashedCapabilityIds.length, 0);
    assertEq(configCount, 0);
  }
}
