var UserContract = artifacts.require("UserContract");
var DeviceContract = artifacts.require("DeviceContract");
var ProductContract = artifacts.require("ProductContract");
var BrokerContract = artifacts.require("BrokerContract");
var NegotiationContract = artifacts.require("NegotiationContract");
var TradingContract = artifacts.require("TradingContract");
var RatingContract = artifacts.require("RatingContract");

module.exports = function (deployer) {

    let userContractInstance;
    let deviceContractInstance;
    let brokerContractInstance;
    let productContractInstance;
    let negotiationContractInstance;
    let tradingContractInstance;
    let ratingContractInstance;

    deployer.deploy(UserContract)
        .then(function () {
            return deployer.deploy(DeviceContract, UserContract.address);
        })
        .then(function () {
            return UserContract.deployed()
                .then(function (instance) {
                    userContractInstance = instance;
                });
        })
        .then(function () {
            userContractInstance.setDeviceContract(DeviceContract.address);
            return DeviceContract.deployed()
                .then(function (instance) {
                    deviceContractInstance = instance;
                });
        })
        .then(function () {
            return deployer.deploy(ProductContract, DeviceContract.address);
        })
        .then(function () {
            return ProductContract.deployed()
                .then(function (instance) {
                    deviceContractInstance.setProductContract(ProductContract.address);
                    productContractInstance = instance;
                });
        })
        .then(function () {
            return deployer.deploy(BrokerContract, UserContract.address);
        })
        .then(function () {
            userContractInstance.setBrokerContract(BrokerContract.address);
            return BrokerContract.deployed()
                .then(function (instance) {
                    brokerContractInstance = instance;
                });
        })
        .then(function () {
            return deployer.deploy(NegotiationContract, DeviceContract.address, ProductContract.address);
        })
        .then(function () {
            return NegotiationContract.deployed()
                .then(function (instance) {
                    deviceContractInstance.setNegotiationContract(NegotiationContract.address);
                    productContractInstance.setNegotiationContract(NegotiationContract.address);
                    negotiationContractInstance = instance;
                });
        })
        .then(function () {
            return deployer.deploy(
                RatingContract,
                DeviceContract.address
            );
        })
        .then(function () {
            return RatingContract.deployed()
                .then(function (instance) {
                    deviceContractInstance.setRatingContract(RatingContract.address);
                    ratingContractInstance = instance;
                });
        })
        .then(function () {
            return deployer.deploy(
                TradingContract,
                DeviceContract.address,
                ProductContract.address,
                BrokerContract.address,
                NegotiationContract.address,
                RatingContract.address
            );
        })
        .then(function () {
            return TradingContract.deployed()
                .then(function (instance) {
                    deviceContractInstance.setTradingContract(TradingContract.address);
                    brokerContractInstance.setTradingContract(TradingContract.address);
                    productContractInstance.setTradingContract(TradingContract.address);
                    ratingContractInstance.setTradingContract(TradingContract.address);
                    tradingContractInstance = instance;
                });
        });
};
