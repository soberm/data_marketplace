pragma solidity >=0.4.21 <0.6.0;

import "./DeviceContract.sol";
import "./BrokerContract.sol";

contract UserContract {

    event CreatedUser(
        address indexed addr
    );

    event UpdatedUser(
        address indexed addr
    );

    event RemovedUser(
        address indexed addr
    );

    struct User {
        address addr;
        uint index;
        string firstName;
        string lastName;
        string company;
        string email;
        bool deleted;
        address[] devices;
        address[] brokers;
    }

    mapping(address => User) private users;
    address[] private index;

    DeviceContract private deviceContract;
    BrokerContract private brokerContract;

    function create(
        string calldata _firstName,
        string calldata _lastName,
        string calldata _company,
        string calldata _email
    )
    external
    returns (address addr)
    {
        require(!existsByAddress(msg.sender), 'user already exists');

        User storage user = users[msg.sender];

        user.addr = msg.sender;
        user.firstName = _firstName;
        user.lastName = _lastName;
        user.company = _company;
        user.email = _email;

        user.index = index.push(msg.sender) - 1;

        emit CreatedUser(user.addr);

        return user.addr;
    }

    function update(
        string calldata _firstName,
        string calldata _lastName,
        string calldata _company,
        string calldata _email
    )
    external
    returns (address addr)
    {
        require(existsByAddressAndDeleted(msg.sender, false), 'user does not exist');

        User storage user = users[msg.sender];

        user.firstName = _firstName;
        user.lastName = _lastName;
        user.company = _company;
        user.email = _email;

        emit UpdatedUser(user.addr);

        return user.addr;
    }

    function remove()
    external
    {
        require(existsByAddressAndDeleted(msg.sender, false), 'user does not exist');

        User storage user = users[msg.sender];
        user.deleted = true;

        removeDevices(user.devices);
        removeBrokers(user.brokers);

        emit RemovedUser(user.addr);
    }

    function addDevice(address _addr, address _device)
    external
    {
        require(msg.sender == address(deviceContract), 'not authorized to add device to user');

        User storage user = users[_addr];
        user.devices.push(_device);
    }

    function addBroker(address _addr, address _broker)
    external
    {
        require(msg.sender == address(brokerContract), 'not authorized to add broker to user');

        User storage user = users[_addr];
        user.brokers.push(_broker);
    }

    function setDeviceContract(address _addr)
    external
    {
        require(address(deviceContract) == address(0), 'device contract already set');
        deviceContract = DeviceContract(_addr);
    }

    function setBrokerContract(address _addr)
    external
    {
        require(address(brokerContract) == address(0), 'broker contract already set');
        brokerContract = BrokerContract(_addr);
    }

    function findByIndex(uint _index)
    external
    view
    returns (
        address addr,
        string memory firstName,
        string memory lastName,
        string memory company,
        string memory email,
        bool deleted,
        address[] memory devices,
        address[] memory brokers
    )
    {
        require(_index >= 0 && _index < index.length, 'user does not exist');

        User memory user = users[index[_index]];

        addr = user.addr;
        firstName = user.firstName;
        lastName = user.lastName;
        company = user.company;
        email = user.email;
        deleted = user.deleted;
        devices = user.devices;
        brokers = user.brokers;
    }

    function findByAddress(address _addr)
    external
    view
    returns (
        address addr,
        string memory firstName,
        string memory lastName,
        string memory company,
        string memory email,
        bool deleted,
        address[] memory devices,
        address[] memory brokers
    )
    {
        require(existsByAddress(_addr), 'user does not exist');

        User memory user = users[_addr];

        addr = user.addr;
        firstName = user.firstName;
        lastName = user.lastName;
        company = user.company;
        email = user.email;
        deleted = user.deleted;
        devices = user.devices;
        brokers = user.brokers;
    }

    function count()
    external
    view
    returns (uint counter)
    {
        return index.length;
    }

    function existsByAddress(address _addr)
    public
    view
    returns (bool exists)
    {
        if (index.length == 0) return false;
        return (index[users[_addr].index] == _addr);
    }

    function existsByAddressAndDeleted(address _addr, bool _deleted)
    public
    view
    returns (bool exists)
    {
        return existsByAddress(_addr) && users[_addr].deleted == _deleted;
    }

    function removeDevices(address[] storage devices)
    private
    {
        for (uint i = 0; i < devices.length; i++) {
            deviceContract.remove(devices[i]);
        }
    }

    function removeBrokers(address[] storage brokers)
    private
    {
        for (uint i = 0; i < brokers.length; i++) {
            brokerContract.remove(brokers[i]);
        }
    }

}