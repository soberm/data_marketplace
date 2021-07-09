pragma solidity >=0.4.21 <0.6.0;

import "./UserContract.sol";
import "./TradingContract.sol";

contract BrokerContract {

    enum Location {
        BR, // Brazil
        EUNE, // Europe Nordic & East
        EUW, // Europe West
        LAN, // Latin America North
        LAS, // Latin America South
        NA, // North America
        OCE, // Oceania
        RU, // Russia
        TR, // Turkey
        JP, // Japan
        PH, // The Philippines
        SG, // Singapore, Malaysia, and Indonesia
        TW, // Taiwan, Hong Kong, and Macao
        VN, // Vietnam
        TH, // Thailand
        KR, // Republic of Korea
        CN      // People's Republic of China
    }

    event CreatedBroker(
        address indexed addr
    );

    event UpdatedBroker(
        address indexed addr
    );

    event RemovedBroker(
        address indexed addr
    );

    struct Broker {
        address addr;
        uint index;
        address user;
        string name;
        string hostAddr;
        Location location;
        bool deleted;
        uint[] trades;
    }

    mapping(address => Broker) private brokers;
    address[] private index;

    UserContract private userContract;
    TradingContract private tradingContract;

    constructor(address _userContractAddress)
    public
    {
        userContract = UserContract(_userContractAddress);
    }

    function create(
        address _addr,
        string calldata _name,
        string calldata _hostAddr,
        Location _location
    )
    external
    returns (address addr)
    {
        require(!existsByAddress(_addr), 'broker already exists');
        require(userContract.existsByAddressAndDeleted(msg.sender, false), 'user does not exist');

        Broker storage broker = brokers[_addr];

        broker.addr = _addr;
        broker.user = msg.sender;
        broker.name = _name;
        broker.hostAddr = _hostAddr;
        broker.location = _location;
        broker.index = index.push(broker.addr) - 1;

        userContract.addBroker(msg.sender, broker.addr);

        emit CreatedBroker(broker.addr);

        return broker.addr;
    }

    function update(
        address _addr,
        string calldata _name,
        string calldata _hostAddr,
        Location _location
    )
    external
    returns (address addr)
    {
        require(existsByAddressAndDeleted(_addr, false), 'broker does not exist');

        Broker storage broker = brokers[_addr];

        require(msg.sender == broker.user, 'not authorized to update broker');

        broker.name = _name;
        broker.hostAddr = _hostAddr;
        broker.location = _location;

        emit UpdatedBroker(broker.addr);

        return broker.addr;
    }

    function remove(address _addr)
    external
    {
        require(existsByAddressAndDeleted(_addr, false), 'broker does not exist');

        Broker storage broker = brokers[_addr];

        require(msg.sender == broker.user || msg.sender == address(userContract), 'not authorized to remove broker');

        broker.deleted = true;

        emit RemovedBroker(broker.addr);
    }

    function addTrade(address _addr, uint _trade)
    external
    {
        require(msg.sender == address(tradingContract), 'not authorized to add trade to broker');

        Broker storage broker = brokers[_addr];
        broker.trades.push(_trade);
    }

    function setTradingContract(address _addr)
    external
    {
        require(address(tradingContract) == address(0), 'trading contract already set');
        tradingContract = TradingContract(_addr);
    }

    function findByIndex(uint _index)
    external
    view
    returns (
        address addr,
        address user,
        string memory name,
        string memory hostAddr,
        Location location,
        uint[] memory trades,
        bool deleted
    )
    {
        require(_index >= 0 && _index < index.length, 'broker does not exist');

        Broker memory broker = brokers[index[_index]];

        addr = broker.addr;
        user = broker.user;
        name = broker.name;
        hostAddr = broker.hostAddr;
        location = broker.location;
        trades = broker.trades;
        deleted = broker.deleted;
    }

    function findByAddress(address _addr)
    external
    view
    returns (
        address addr,
        address user,
        string memory name,
        string memory hostAddr,
        Location location,
        uint[] memory trades,
        bool deleted
    )
    {
        require(existsByAddress(_addr), 'broker does not exist');

        Broker memory broker = brokers[_addr];

        addr = broker.addr;
        user = broker.user;
        name = broker.name;
        hostAddr = broker.hostAddr;
        location = broker.location;
        trades = broker.trades;
        deleted = broker.deleted;
    }

    function count()
    external
    view
    returns (uint counter)
    {
        return index.length;
    }

    function countryAbbr(Location _location)
    external
    pure
    returns (string memory abbr)
    {
        if (_location == Location.BR) return "BR";
        if (_location == Location.EUNE) return "EUNE";
        if (_location == Location.EUW) return "EUW";
        if (_location == Location.LAN) return "LAN";
        if (_location == Location.LAS) return "LAS";
        if (_location == Location.NA) return "NA";
        if (_location == Location.OCE) return "OCE";
        if (_location == Location.RU) return "RU";
        if (_location == Location.TR) return "TR";
        if (_location == Location.JP) return "JP";
        if (_location == Location.PH) return "PH";
        if (_location == Location.SG) return "SG";
        if (_location == Location.TW) return "TW";
        if (_location == Location.VN) return "VN";
        if (_location == Location.TH) return "TH";
        if (_location == Location.KR) return "KR";
        if (_location == Location.CN) return "CN";
        revert("location does not exist");
    }

    function countryCode(string calldata abbr)
    external
    pure
    returns (Location code)
    {
        if (keccak256(abi.encodePacked(abbr)) == keccak256("BR")) return Location.BR;
        if (keccak256(abi.encodePacked(abbr)) == keccak256("EUNE")) return Location.EUNE;
        if (keccak256(abi.encodePacked(abbr)) == keccak256("EUW")) return Location.EUW;
        if (keccak256(abi.encodePacked(abbr)) == keccak256("LAN")) return Location.LAN;
        if (keccak256(abi.encodePacked(abbr)) == keccak256("LAS")) return Location.LAS;
        if (keccak256(abi.encodePacked(abbr)) == keccak256("NA")) return Location.NA;
        if (keccak256(abi.encodePacked(abbr)) == keccak256("OCE")) return Location.OCE;
        if (keccak256(abi.encodePacked(abbr)) == keccak256("RU")) return Location.RU;
        if (keccak256(abi.encodePacked(abbr)) == keccak256("TR")) return Location.TR;
        if (keccak256(abi.encodePacked(abbr)) == keccak256("JP")) return Location.JP;
        if (keccak256(abi.encodePacked(abbr)) == keccak256("PH")) return Location.PH;
        if (keccak256(abi.encodePacked(abbr)) == keccak256("SG")) return Location.SG;
        if (keccak256(abi.encodePacked(abbr)) == keccak256("TW")) return Location.TW;
        if (keccak256(abi.encodePacked(abbr)) == keccak256("VN")) return Location.VN;
        if (keccak256(abi.encodePacked(abbr)) == keccak256("TH")) return Location.TH;
        if (keccak256(abi.encodePacked(abbr)) == keccak256("KR")) return Location.KR;
        if (keccak256(abi.encodePacked(abbr)) == keccak256("CN")) return Location.CN;
        revert("location does not exist");
    }

    function existsByAddress(address _addr)
    public
    view
    returns (bool isContained)
    {
        if (index.length == 0) return false;
        return (index[brokers[_addr].index] == _addr);
    }

    function existsByAddressAndDeleted(address _addr, bool _deleted)
    public
    view
    returns (bool exists)
    {
        return existsByAddress(_addr) && brokers[_addr].deleted == _deleted;
    }

}
