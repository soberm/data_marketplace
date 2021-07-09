module.exports = async function(callback) {
    var accounts = await web3.eth.getAccounts();
    var transactionConfig = {
        from: accounts[0],
        to: "0x10991E218011e0f4ef9555EaF7d66b8aB5Ef24AB",
        value: web3.utils.toWei("0.10", "ether"),
    };
    web3.eth.sendTransaction(transactionConfig);
};
