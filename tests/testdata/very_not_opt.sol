// SPDX-License-Identifier: GPL-3.0
pragma solidity ^0.8.0;

contract VeryUnoptimisedContract {
    uint128 public id;
    address public addr;
    uint256 public balance;
    string public favouritePet;
    bool public isActive;

    struct Employee {
        uint256 id;
        uint32 salary;
        uint256 age;
        bool isActive;
        address addr;
        string role;
        uint16 department;
    }

    function NotPureNotModFunc(uint256[] numbers) external returns (uint256) {
        uint256 sum = 0;
        for (uint256 i = 0; i < numbers.length; ++i) {
            sum += numbers[i];
        }
        return sum;
    }

    function PureNoModFunc(uint256[] numbers) external pure returns (uint256) {
        uint256 sum = 0;
        for (uint256 i = 0; i < numbers.length; ++i) {
            sum += numbers[i];
        }
        return sum;
    }
}
