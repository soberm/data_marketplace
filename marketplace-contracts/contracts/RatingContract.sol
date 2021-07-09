pragma solidity >=0.4.21 <0.6.0;

import "./DeviceContract.sol";
import "./TradingContract.sol";

contract RatingContract {

    mapping(address => bool) private settlements;

    DeviceContract private deviceContract;
    TradingContract private tradingContract;

    constructor(address _deviceContractAddress)
    public
    {
        deviceContract = DeviceContract(_deviceContractAddress);
    }

    function handleSettlement(address _device)
    external
    {
        require(settlements[msg.sender], 'not authorized to rate device');
        uint rating = deviceContract.findRatingByAddress(_device);
        if (rating >= 100) {
            return;
        }
        deviceContract.updateRating(_device, rating + 1);
    }

    function handleDispute(address _device)
    external
    {
        require(settlements[msg.sender], 'not authorized to rate device');
        uint rating = deviceContract.findRatingByAddress(_device);
        if (rating <= 5) {
            deviceContract.updateRating(_device, 0);
        } else {
            deviceContract.updateRating(_device, rating - 5);
        }
    }

    function addSettlementContract(address _addr)
    external
    {
        require(msg.sender == address(tradingContract), 'not authorized to add settlement');
        settlements[_addr] = true;
    }

    function setTradingContract(address _addr)
    external
    {
        require(address(tradingContract) == address(0), 'trading contract already set');
        tradingContract = TradingContract(_addr);
    }

}
