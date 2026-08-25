// SPDX-License-Identifier: MIT
pragma solidity >=0.6.12 <0.9.0;

contract MyCounter {

    /**
     * 全局数字，用于计数
     */
    uint256 internal counter;

    constructor() {
        counter = 0;
    }

    /**
     * 这是我的测试计数器
     * 自增后，返回自增后的数字
   */
    function increase() public returns (uint256) {
        return ++counter;
    }

    /**
     * 返回最新的 counter
     */
    function getCounter() public view returns (uint256){
        return counter;
    }
}