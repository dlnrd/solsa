// SPDX-License-Identifier: GPL-3.0
pragma solidity ^0.8.0;

contract CalldataUnoptimisable {
    function NotPure_ModFunc(uint256[] numbers) external returns (uint256) {
        uint256 sum = 0;
        for (uint256 i = 0; i < numbers.length; ++i) {
            sum += numbers[i];
            numbers[i] = 0;
        }
        return sum;
    }
}
