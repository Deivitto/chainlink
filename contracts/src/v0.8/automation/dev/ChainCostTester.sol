// SPDX-License-Identifier: MIT
pragma solidity 0.8.19;

import {Metis_GasPriceOracle} from "./../../vendor/@metis/contracts/L2/predeploys/Metis_GasPriceOracle.sol";
import {IScrollL1GasPriceOracle} from "./../../vendor/@scroll-tech/contracts/src/L2/predeploys/IScrollL1GasPriceOracle.sol";
import {ArbGasInfo} from "../../vendor/@arbitrum/nitro-contracts/src/precompiles/ArbGasInfo.sol";

contract ChainCostTester {
    event BlockLog(uint256 block, uint256 timestamp, bytes32 hash);
    event GasDetails(uint256 nativeGasUsed, uint256 l1Cost, uint256 l1DataLength, uint256 l1BaseFee, uint256 nativeGasPriceFromOracle, uint256 txGasPrice);

    Metis_GasPriceOracle public immutable METIS_ORACLE = Metis_GasPriceOracle(0x420000000000000000000000000000000000000F);
    IScrollL1GasPriceOracle public immutable SCROLL_ORACLE = IScrollL1GasPriceOracle(0x5300000000000000000000000000000000000002);
    ArbGasInfo internal constant ARB_NITRO_ORACLE = ArbGasInfo(0x000000000000000000000000000000000000006C);
    bytes public data;
    uint256 public sum;
    bytes public arbPadding;
    bytes public scrollPadding;
    bytes public metisPadding;
    bytes public gnosisPadding;
    bytes public zkSyncPadding;
    bytes public zkEVMPadding;
    bytes public celoPadding;
    bytes public ethSepoliaPadding;

    function setScrollPadding(bytes calldata _padding) public {
        scrollPadding = _padding;
    }

    function setMetisPadding(bytes calldata _padding) public {
        metisPadding = _padding;
    }

    function setGnosisPadding(bytes calldata _padding) public {
        gnosisPadding = _padding;
    }

    function setZKSyncPadding(bytes calldata _padding) public {
        zkSyncPadding = _padding;
    }

    function setZKEVMPadding(bytes calldata _padding) public {
        zkEVMPadding = _padding;
    }

    function setCeloPadding(bytes calldata _padding) public {
        celoPadding = _padding;
    }

    function computeAndUpdateData(bytes calldata _data, uint256 _n) public {
        uint256 g1 = gasleft();
        emit BlockLog(block.number, block.timestamp, blockhash(block.number - 1));

        // some computation
        uint256 _sum = 0;
        for (uint256 i = 0; i < _n; i++) {
            _sum += i;
        }
        sum = _sum;

        // store bytes
        data = _data;

        uint256 l1Cost;
        uint256 l1BaseFee;
        uint256 nativeGasPriceFromOracle;
        bytes memory totalBytes;
        if (isArbitrum()) {
            totalBytes = bytes.concat(_data, arbPadding);
            l1Cost = ARB_NITRO_ORACLE.getCurrentTxL1GasFees();
            l1BaseFee = ARB_NITRO_ORACLE.getL1GasPriceEstimate();
            nativeGasPriceFromOracle = 0;
        } else if (isMetis()) {
            totalBytes = bytes.concat(_data, metisPadding);
            l1Cost = METIS_ORACLE.getL1Fee(totalBytes);
            l1BaseFee = METIS_ORACLE.l1BaseFee();
            nativeGasPriceFromOracle = METIS_ORACLE.gasPrice();
        } else if (isScroll()) {
            totalBytes = bytes.concat(_data, scrollPadding);
            l1Cost = SCROLL_ORACLE.getL1Fee(totalBytes);
            l1BaseFee = SCROLL_ORACLE.l1BaseFee();
            nativeGasPriceFromOracle = 0; // the oracle does not expose this
        } else if (isGnosis()) {
            totalBytes = bytes.concat(_data, gnosisPadding);
            l1Cost = 0;
            l1BaseFee = 0;
            nativeGasPriceFromOracle = 0;
        } else if (isPolygonZKEVM()) {
            totalBytes = bytes.concat(_data, zkEVMPadding);
            l1Cost = 0;
            l1BaseFee = 0;
            nativeGasPriceFromOracle = 0;
        } else if (isZKSync()) {
            totalBytes = bytes.concat(_data, zkSyncPadding);
            l1Cost = 0;
            l1BaseFee = 0;
            nativeGasPriceFromOracle = 0;
        } else if (isCelo()) {
            totalBytes = bytes.concat(_data, celoPadding);
            l1Cost = 0;
            l1BaseFee = 0;
            nativeGasPriceFromOracle = 0;
        } else if (isETHSepolia()) {
            totalBytes = bytes.concat(_data, ethSepoliaPadding);
            l1Cost = 0;
            l1BaseFee = 0;
            nativeGasPriceFromOracle = 0;
        }

        uint256 g2 = gasleft();
        emit GasDetails(g1 - g2, l1Cost, totalBytes.length, l1BaseFee, nativeGasPriceFromOracle, tx.gasprice);
    }

//    function updateData(bytes calldata _data) public {
//        uint256 g1 = gasLeft();
//        emit BlockLog(block.number, block.timestamp, blockhash(block.number));
//
//        // store bytes
//        data = _data;
//
//        uint256 g2 = gasLeft();
//        emit GasDetails(g2 - g1, );
//    }
//
//    function compute(uint256 _n) public {
//        uint256 g1 = gasLeft();
//        emit BlockLog(block.number, block.timestamp, blockhash(block.number));
//
//        // some computation
//        uint256 _sum = 0;
//        for (uint256 i = 0; i < _n; i++) {
//            _sum += i;
//        }
//        sum = _sum;
//
//        uint256 g2 = gasLeft();
//        emit GasDetails(g2 - g1, );
//    }

    function isScroll() public view returns (bool) {
        // Scroll Sepolia or Scroll mainnet
        return block.chainid == 534351 || block.chainid == 534352;
    }

    function isMetis() public view returns (bool) {
        // Metis Goerli or Metis Andromeda mainnet
        return block.chainid == 599 || block.chainid == 1088;
    }

    function isGnosis() public view returns (bool) {
        // Gnosis Chiado or Gnosis mainnet
        return block.chainid == 10200 || block.chainid == 100;
    }

    function isZKSync() public view returns (bool) {
        // zkSync Sepolia or zkSync mainnet
        return block.chainid == 300 || block.chainid == 324;
    }

    function isPolygonZKEVM() public view returns (bool) {
        // zkEVM testnet or zkEVM mainnet
        return block.chainid == 1442 || block.chainid == 1101;
    }

    function isCelo() public view returns (bool) {
        // celo Alfajores testnet or celo mainnet
        return block.chainid == 44787 || block.chainid == 42220;
    }

    function isArbitrum() public view returns (bool) {
        // Arb Sepolia testnet or Arb mainnet
        return block.chainid == 421614 || block.chainid == 42161;
    }

    function isETHSepolia() public view returns (bool) {
        return block.chainid == 11155111;
    }
}
//
//interface Metis_GasPriceOracle {
//    /*************
//     * Variables *
//     *************/
//    // Current L2 gas price
//    uint256 public gasPrice;
//    // Current L1 base fee
//    uint256 public l1BaseFee;
//    // Amortized cost of batch submission per transaction
//    uint256 public overhead;
//    // Value to scale the fee up by
//    uint256 public scalar;
//    // Number of decimals of the scalar
//    uint256 public decimals;
//
//    // minimum gas to bridge the asset back to l1
//    uint256 public minErc20BridgeCost;
//
//    /**
//     * Computes the L1 portion of the fee
//     * based on the size of the RLP encoded tx
//     * and the current l1BaseFee
//     * @param _data Unsigned RLP encoded tx, 6 elements
//     * @return L1 fee that should be paid for the tx
//     */
//    function getL1Fee(bytes memory _data) public view returns (uint256) {
//        uint256 l1GasUsed = getL1GasUsed(_data);
//        uint256 l1Fee = l1GasUsed * l1BaseFee;
//        uint256 divisor = 10**decimals;
//        uint256 unscaled = l1Fee * scalar;
//        uint256 scaled = unscaled / divisor;
//        return scaled;
//    }
//
//    // solhint-disable max-line-length
//    /**
//     * Computes the amount of L1 gas used for a transaction
//     * The overhead represents the per batch gas overhead of
//     * posting both transaction and state roots to L1 given larger
//     * batch sizes.
//     * 4 gas for 0 byte
//     * https://github.com/ethereum/go-ethereum/blob/9ada4a2e2c415e6b0b51c50e901336872e028872/params/protocol_params.go#L33
//     * 16 gas for non zero byte
//     * https://github.com/ethereum/go-ethereum/blob/9ada4a2e2c415e6b0b51c50e901336872e028872/params/protocol_params.go#L87
//     * This will need to be updated if calldata gas prices change
//     * Account for the transaction being unsigned
//     * Padding is added to account for lack of signature on transaction
//     * 1 byte for RLP V prefix
//     * 1 byte for V
//     * 1 byte for RLP R prefix
//     * 32 bytes for R
//     * 1 byte for RLP S prefix
//     * 32 bytes for S
//     * Total: 68 bytes of padding
//     * @param _data Unsigned RLP encoded tx, 6 elements
//     * @return Amount of L1 gas used for a transaction
//     */
//    // solhint-enable max-line-length
//    function getL1GasUsed(bytes memory _data) public view returns (uint256) {
//        uint256 total = 0;
//        for (uint256 i = 0; i < _data.length; i++) {
//            if (_data[i] == 0) {
//                total += 4;
//            } else {
//                total += 16;
//            }
//        }
//        uint256 unsigned = total + overhead;
//        return unsigned + (68 * 16);
//    }
//}
//
//interface IScrollL1GasPriceOracle {
//    /**********
//     * Events *
//     **********/
//
//    /// @notice Emitted when current fee overhead is updated.
//    /// @param overhead The current fee overhead updated.
//    event OverheadUpdated(uint256 overhead);
//
//    /// @notice Emitted when current fee scalar is updated.
//    /// @param scalar The current fee scalar updated.
//    event ScalarUpdated(uint256 scalar);
//
//    /// @notice Emitted when current l1 base fee is updated.
//    /// @param l1BaseFee The current l1 base fee updated.
//    event L1BaseFeeUpdated(uint256 l1BaseFee);
//
//    /*************************
//     * Public View Functions *
//     *************************/
//
//    /// @notice Return the current l1 fee overhead.
//    function overhead() external view returns (uint256);
//
//    /// @notice Return the current l1 fee scalar.
//    function scalar() external view returns (uint256);
//
//    /// @notice Return the latest known l1 base fee.
//    function l1BaseFee() external view returns (uint256);
//
//    /// @notice Computes the L1 portion of the fee based on the size of the rlp encoded input
//    ///         transaction, the current L1 base fee, and the various dynamic parameters.
//    /// @param data Unsigned fully RLP-encoded transaction to get the L1 fee for.
//    /// @return L1 fee that should be paid for the tx
//    function getL1Fee(bytes memory data) external view returns (uint256);
//
//    /// @notice Computes the amount of L1 gas used for a transaction. Adds the overhead which
//    ///         represents the per-transaction gas overhead of posting the transaction and state
//    ///         roots to L1. Adds 74 bytes of padding to account for the fact that the input does
//    ///         not have a signature.
//    /// @param data Unsigned fully RLP-encoded transaction to get the L1 gas for.
//    /// @return Amount of L1 gas used to publish the transaction.
//    function getL1GasUsed(bytes memory data) external view returns (uint256);
//
//    /*****************************
//     * Public Mutating Functions *
//     *****************************/
//
//    /// @notice Allows whitelisted caller to modify the l1 base fee.
//    /// @param _l1BaseFee New l1 base fee.
//    function setL1BaseFee(uint256 _l1BaseFee) external;
//}