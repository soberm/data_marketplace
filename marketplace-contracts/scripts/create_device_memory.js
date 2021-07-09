const UserContract = artifacts.require('UserContract');
const DeviceContract = artifacts.require('DeviceContract');

module.exports = async function(callback) {
    let user = {
        firstName: 'michael',
        lastName: 'sober',
        company: 'tu wien',
        email: 'e1326403@student.tuwien.ac.at',
    };

    try {
        let accounts = await web3.eth.getAccounts();
        let userContract = await UserContract.deployed();
        let deviceContract = await DeviceContract.deployed();

        await userContract.create(
            user.firstName, user.lastName, user.company, user.email, {from: accounts[0]}
        );

        let description = "asdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasasasaa";
        for(let k = 1; k < 1024; k++){
            let deviceID = web3.utils.randomHex(20);
            description += "a";
            let result = await deviceContract.create(
                deviceID, "Device", "asdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasasasaa", web3.utils.hexToBytes(web3.utils.randomHex(64)), {from: accounts[0]}
            );
            result = await deviceContract.update(
                deviceID, "Device", description, web3.utils.hexToBytes(web3.utils.randomHex(64)), {from: accounts[0]}
            );
            console.log(k, ",", result.receipt.gasUsed)
        }

    }
    catch (e) {
        console.log(e)
    }
};
