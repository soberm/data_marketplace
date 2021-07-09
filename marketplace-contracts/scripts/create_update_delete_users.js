const UserContract = artifacts.require('UserContract');

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

        for (let i = 0; i < accounts.length; i++) {
            let result = await userContract.create(
                user.firstName, user.lastName, user.company, user.email, {from: accounts[i]}
            );
            result = await userContract.update(
                user.firstName, user.lastName, user.company, user.email, {from: accounts[Math.floor(Math.random() * i)]}
            );
        }

        for (let i = 0; i < accounts.length; i++) {
            let result = await userContract.remove({from: accounts[i]});
            console.log(i, ",", result.receipt.gasUsed)
        }
    }
    catch (e) {
        console.log(e)
    }
};
