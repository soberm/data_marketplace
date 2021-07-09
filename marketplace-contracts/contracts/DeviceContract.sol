pragma solidity >=0.4.21 <0.6.0;

import "./UserContract.sol";
import "./ProductContract.sol";
import "./NegotiationContract.sol";
import "./TradingContract.sol";
import "./RatingContract.sol";

contract DeviceContract {

    event CreatedDevice(
        address indexed addr
    );

    event UpdatedDevice(
        address indexed addr
    );

    event RemovedDevice(
        address indexed addr
    );

    struct Device {
        address addr;
        uint index;
        address user;
        string name;
        string description;
        bytes publicKey;
        uint rating;
        bool deleted;
        uint[] products;
        uint[] negotiationRequests;
        uint[] negotiations;
        uint[] tradingRequests;
        uint[] trades;
    }

    mapping(address => Device) private devices;
    address[] private index;

    UserContract private userContract;
    ProductContract private productContract;
    NegotiationContract private negotiationContract;
    TradingContract private tradingContract;
    RatingContract private ratingContract;

    constructor(address _userContractAddress)
    public
    {
        userContract = UserContract(_userContractAddress);
    }

    function create(
        address _addr,
        string calldata _name,
        string calldata _description,
        bytes calldata _publicKey
    )
    external
    returns (address addr)
    {
        require(!existsByAddress(_addr), 'device already exists');
        require(userContract.existsByAddressAndDeleted(msg.sender, false), 'user does not exist');

        Device storage device = devices[_addr];

        device.addr = _addr;
        device.user = msg.sender;
        device.name = _name;
        device.description = _description;
        device.publicKey = _publicKey;
        device.index = index.push(device.addr) - 1;

        userContract.addDevice(msg.sender, device.addr);

        emit CreatedDevice(device.addr);

        return device.addr;
    }

    function update(
        address _addr,
        string calldata _name,
        string calldata _description,
        bytes calldata _publicKey
    )
    external
    returns (address addr)
    {
        require(existsByAddressAndDeleted(_addr, false), 'device does not exist');

        Device storage device = devices[_addr];

        require(msg.sender == device.user, 'not authorized to update device');

        device.name = _name;
        device.description = _description;
        device.publicKey = _publicKey;

        emit UpdatedDevice(device.addr);

        return device.addr;
    }

    function updateRating(
        address _addr,
        uint _rating
    )
    external
    {
        require(existsByAddressAndDeleted(_addr, false), 'device does not exist');
        require(msg.sender == address(ratingContract), 'not authorized to update rating');

        Device storage device = devices[_addr];
        device.rating = _rating;
    }

    function remove(address _addr)
    external
    {
        require(existsByAddressAndDeleted(_addr, false), 'device does not exist');

        Device storage device = devices[_addr];

        require(msg.sender == device.user || msg.sender == address(userContract), 'not authorized to remove device');

        device.deleted = true;

        removeProducts(device.products);

        emit RemovedDevice(device.addr);
    }

    function addProduct(address _address, uint _product)
    external
    {
        require(msg.sender == address(productContract), 'not authorized to add product to device');

        Device storage device = devices[_address];
        device.products.push(_product);
    }

    function addNegotiationRequest(address _addr, uint _request)
    external
    {
        require(msg.sender == address(negotiationContract), 'not authorized to add negotiation request to device');

        Device storage device = devices[_addr];
        device.negotiationRequests.push(_request);
    }

    function addNegotiation(address _addr, uint _negotiation)
    external
    {
        require(msg.sender == address(negotiationContract), 'not authorized to add negotiation to device');

        Device storage device = devices[_addr];
        device.negotiations.push(_negotiation);
    }

    function addTradingRequest(address _addr, uint _request)
    external
    {
        require(msg.sender == address(tradingContract), 'not authorized to add trading request to device');

        Device storage device = devices[_addr];
        device.tradingRequests.push(_request);
    }

    function addTrade(address _addr, uint _trade)
    external
    {
        require(msg.sender == address(tradingContract), 'not authorized to add trade to device');

        Device storage device = devices[_addr];
        device.trades.push(_trade);
    }

    function removeNegotiationRequest(address _addr, uint _request)
    external
    {
        require(msg.sender == address(negotiationContract), 'not authorized to remove negotiation request from device');

        Device storage device = devices[_addr];
        removeRequest(_request, device.negotiationRequests);
    }

    function removeTradingRequest(address _addr, uint _request)
    external
    {
        require(msg.sender == address(tradingContract), 'not authorized to remove trading request from device');

        Device storage device = devices[_addr];
        removeRequest(_request, device.tradingRequests);
    }

    function setProductContract(address _addr)
    external
    {
        require(address(productContract) == address(0), 'product contract already set');
        productContract = ProductContract(_addr);
    }

    function setNegotiationContract(address _addr)
    external
    {
        require(address(negotiationContract) == address(0), 'negotiation contract already set');
        negotiationContract = NegotiationContract(_addr);
    }

    function setTradingContract(address _addr)
    external
    {
        require(address(tradingContract) == address(0), 'trading contract already set');
        tradingContract = TradingContract(_addr);
    }

    function setRatingContract(address _addr)
    external
    {
        require(address(ratingContract) == address(0), 'rating contract already set');
        ratingContract = RatingContract(_addr);
    }

    function findByIndex(uint _index)
    external
    view
    returns (
        address addr,
        address user,
        string memory name,
        string memory description,
        bytes memory publicKey,
        uint rating,
        bool deleted
    )
    {
        require(_index >= 0 && _index < index.length, 'device does not exist');

        Device memory device = devices[index[_index]];

        addr = device.addr;
        user = device.user;
        name = device.name;
        description = device.description;
        publicKey = device.publicKey;
        rating = device.rating;
        deleted = device.deleted;
    }

    function findByAddress(address _addr)
    external
    view
    returns (
        address addr,
        address user,
        string memory name,
        string memory description,
        bytes memory publicKey,
        uint rating,
        bool deleted
    )
    {
        require(existsByAddress(_addr), 'device does not exist');

        Device memory device = devices[_addr];

        addr = device.addr;
        user = device.user;
        name = device.name;
        description = device.description;
        publicKey = device.publicKey;
        rating = device.rating;
        deleted = device.deleted;
    }

    function findProductsByAddress(address _addr)
    external
    view
    returns (uint[] memory products)
    {
        require(existsByAddress(_addr), 'device does not exist');

        Device memory device = devices[_addr];
        products = device.products;
    }

    function findNegotiationRequestsByAddress(address _addr)
    external
    view
    returns (uint[] memory negotiationRequests)
    {
        require(existsByAddress(_addr), 'device does not exist');

        Device memory device = devices[_addr];
        negotiationRequests = device.negotiationRequests;
    }

    function findNegotiationsByAddress(address _addr)
    external
    view
    returns (uint[] memory negotiations)
    {
        require(existsByAddress(_addr), 'device does not exist');

        Device memory device = devices[_addr];
        negotiations = device.negotiations;
    }

    function findTradingRequestsByAddress(address _addr)
    external
    view
    returns (uint[] memory tradingRequests)
    {
        require(existsByAddress(_addr), 'device does not exist');

        Device memory device = devices[_addr];
        tradingRequests = device.tradingRequests;
    }

    function findTradesByAddress(address _addr)
    external
    view
    returns (uint[] memory trades)
    {
        require(existsByAddress(_addr), 'device does not exist');

        Device memory device = devices[_addr];
        trades = device.trades;
    }

    function findRatingByAddress(address _address)
    external
    view
    returns (uint rating)
    {
        require(existsByAddress(_address), 'device does not exist');
        rating = devices[_address].rating;
    }

    function isOwnedByUser(address _addr, address _user)
    external
    view
    returns (bool isOwned)
    {
        require(existsByAddress(_addr), 'device does not exist');
        return devices[_addr].user == _user;
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
        return (index[devices[_addr].index] == _addr);
    }

    function existsByAddressAndDeleted(address _addr, bool _deleted)
    public
    view
    returns (bool exists)
    {
        return existsByAddress(_addr) && devices[_addr].deleted == _deleted;
    }

    function removeProducts(uint[] storage products)
    private
    {
        for (uint i = 0; i < products.length; i++) {
            productContract.remove(products[i]);
        }
    }

    function removeRequest(uint _request, uint[] storage requests)
    private
    {
        for (uint i = 0; i < requests.length; i++) {
            if (requests[i] == _request) {
                requests[i] = requests[requests.length - 1];
                requests.length--;
                return;
            }
        }
    }

}
