// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

contract UnoptimisedStruct {
    struct Employee {
        uint256 id;
        uint32 salary;
        uint256 age;
        bool isActive;
        address addr;
        string role;
        uint16 department;
    }
}

