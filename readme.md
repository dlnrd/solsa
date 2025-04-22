# solsa
Command-line Static Analyzer Tool for Solidity Smart Contracts

## Building
```
go build
```

## Testing
```
go test ./tests
```

## Usage
### Input Flag (Required)
```
solsa -i <path-to-file>
```

```
solsa -i <path-to-dir>
```

### Output (Optional)
```
solsa -i <path> -o <output-dir>
```

### Silent (Optional)

```
solsa -i <path> -o <output-dir> -s
```

<details>
<summary><h2>Optimisations</h2></summary>
Below are the gas-inefficient patterns that solsa identifies and refactors.
<details><summary><h3>Calldata Optimisation</h3></summary>
<h4> What is the optimisation? </h4>
In Solidity, memory and calldata are different types of data locations used to store variables and data. They determine how data is stored, accessed, and how much gas is consumed. Memory is the default location (although can be explicitly specified), it is a temporary, mutable (can read and write) and more expensive in gas compared to calldata. Calldata is immutable, read-only and reserved for external function's input parameters (variables on the blockchain). If the external parameter doesn't get modified then gas can be saved by storing the parameter in calldata instead of being copied to memory. Pure functions don't allow on-chain variables to be modified so any parameters should be stored in calldata to save gas.
<h4>Unoptimised Contract</h4>

```solidity
// SPDX-License-Identifier: GPL-3.0
pragma solidity ^0.8.0;

contract NoCalldataUsage {
    // No modification of `numbers` so can optimise into calldata
    function NotPureNotModFunc(uint256[] numbers) external returns (uint256) {
        uint256 sum = 0;
        for (uint256 i = 0; i < numbers.length; ++i) {
            sum += numbers[i];
        } 
        return sum;
    }

    function NotPure_ModFunc(uint256[] numbers) external returns (uint256) {
        uint256 sum = 0;
        for (uint256 i = 0; i < numbers.length; ++i) {
            sum += numbers[i];
            numbers[i] = 0; // Can't be optimised due to this assignment
        }
        return sum;
    }

    // Pure function should be optimised
    function PureNoModFunc(uint256[] numbers) external pure returns (uint256) {
        uint256 sum = 0;
        for (uint256 i = 0; i < numbers.length; ++i) {
            sum += numbers[i];
        }
        return sum;
    }
}
```

<h4>Optimised Contract</h4>

```solidity
// SPDX-License-Identifier: GPL-3.0
pragma solidity ^0.8.0;

contract NoCalldataUsage {
    // Has been optimised as external parameter doesn't get modified
    function NotPure_NotModFunc(uint256[] calldata numbers) external returns (uint256) {
        uint256 sum = 0;
        for (uint256 i = 0; i < numbers.length; ++i) {
            sum += numbers[i];
        }
        return sum;
    }

    // Can't be optimised due to external parameter gets modified
    function NotPure_ModFunc(uint256[] numbers) external returns (uint256) {
        uint256 sum = 0;
        for (uint256 i = 0; i < numbers.length; ++i) {
            sum += numbers[i];
            numbers[i] = 0;
        }
        return sum;
    }

    // Pure functions can't modify external variables so all parameters should be in calldata
    function PureNo_ModFunc(uint256[] calldata numbers) external pure returns (uint256) {
        uint256 sum = 0;
        for (uint256 i = 0; i < numbers.length; ++i) {
            sum += numbers[i];
        }
        return sum;
    }


}
```

</details>

<details><summary><h3>State Varaible Optmisation</h3></summary>
<h4> What is the optimisation? </h4>
This optimisation involves arraging state variables to minimize storage cost and therefore reducing gas usage. In Solidity, each storage slot is 256 bits (32 bytes) so by using a bin-packing algorithm, we can sort each state variable into the smallest number of storage slots, reducing the number of gas usage which is especially important for contracts with a large number of state variables.
<h4>Unoptimised Contract</h4>

```solidity
// SPDX-License-Identifier: GPL-3.0
pragma solidity ^0.8.0;

contract BasicStorage {
    uint128 public id; // 128 bits (slot 1)
    address public addr; // 160 bits (slot 2)
    uint256 public balance; // 256 bits (slot 3)
    string public favouritePet; // 256 bits (slot 4)
    bool public isActive; // 8 bits (slot 5)
}
```

<h4>Optimised Contract</h4>

```solidity
// SPDX-License-Identifier: GPL-3.0
pragma solidity ^0.8.0;

contract BasicStorage {
    uint256 public balance; // 256 bits (slot 1)
    string public favouritePet; // 256 bits (slot 2)
    address public addr; // 160 bits (slot 3)
    bool public isActive; // 8 bits (slot 3)
    uint128 public id; // 128 bits (slot 4)
}
```

</details>

<details><summary><h3>Struct Varaible Optmisation</h3></summary>
<h4> What is the optimisation? </h4>
This optimisation involves arraging struct variables to minimize storage cost and therefore reducing gas usage. In Solidity, each storage slot is 256 bits (32 bytes) so by using a bin-packing algorithm, we can sort each struct variable into the smallest number of storage slots, reducing the number of gas usage which is especially important for contracts with a large number of struct variables.
<h4>Unoptimised Contract</h4>

```solidity
// SPDX-License-Identifier: GPL-3.0
pragma solidity ^0.8.0;

contract EmployeeInfo {
    struct Employee {
        uint256 id; // 256 bits (slot 1)
        uint32 salary; // 32 bits (slot 2)
        uint256 age; // 256 bits (slot 3)
        bool isActive; // 8 bits (slot 4)
        address addr; // 160 bits (slot 4)
        string role; // 256 bits (slot 5)
        uint16 department; // 16 bits (slot 6)
  }
}
```

<h4>Optimised Contract</h4>

```solidity
// SPDX-License-Identifier: GPL-3.0
pragma solidity ^0.8.0;

contract EmployeeInfo {
    struct Employee {
        uint256 id; // 256 bits (slot 1)
        uint256 age; // 256 bits (slot 2)
        string role; // 256 bits (slot 3)
        address addr; // 160 bits (slot 4)
        uint32 salary; // 32 bits (slot 4)
        uint16 department; // 16 bits (slot 4)
        bool isActive; // 8 bits (slot 4)
  }
}
```

</details>

</details>
