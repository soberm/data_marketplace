pragma solidity >=0.4.21 <0.6.0;

import "./DeviceContract.sol";
import "./NegotiationContract.sol";
import "./TradingContract.sol";

contract ProductContract {

    event CreatedProduct(
        uint id,
        address indexed user
    );

    event UpdatedProduct(
        uint indexed id,
        address indexed user
    );

    event RemovedProduct(
        uint indexed id,
        address indexed user
    );

    struct Product {
        uint id;
        uint index;
        address device;
        string name;
        string description;
        string dataType;
        uint frequency;
        uint cost;
        bool deleted;
        uint[] negotiationRequests;
        uint[] negotiations;
        uint[] tradingRequests;
        uint[] trades;
    }

    DeviceContract private deviceContract;
    NegotiationContract private negotiationContract;
    TradingContract private tradingContract;

    mapping(uint => Product) private products;
    uint[] private index;
    uint private idCounter;

    constructor(address _deviceContractAddress)
    public
    {
        deviceContract = DeviceContract(_deviceContractAddress);
    }

    function create(
        address _device,
        string calldata _name,
        string calldata _description,
        string calldata _dataType,
        uint _frequency,
        uint _cost
    )
    external
    returns (uint id)
    {
        require(deviceContract.existsByAddress(_device), 'device does not exist');
        require(deviceContract.isOwnedByUser(_device, msg.sender), 'not authorized to create product');

        uint nextId = idCounter++;

        Product storage product = products[nextId];

        product.id = nextId;
        product.device = _device;
        product.name = _name;
        product.description = _description;
        product.dataType = _dataType;
        product.frequency = _frequency;
        product.cost = _cost;
        product.index = index.push(product.id) - 1;

        deviceContract.addProduct(_device, product.id);

        emit CreatedProduct(product.id, msg.sender);

        return product.id;
    }

    function update(
        uint _id,
        string calldata _name,
        string calldata _description,
        string calldata _dataType,
        uint _frequency,
        uint _cost
    )
    external
    returns (uint id)
    {
        require(existsByIdAndDeleted(_id, false), 'product does not exist');

        Product storage product = products[_id];

        require(deviceContract.isOwnedByUser(product.device, msg.sender), 'not authorized to update product');

        product.name = _name;
        product.description = _description;
        product.dataType = _dataType;
        product.frequency = _frequency;
        product.cost = _cost;

        emit UpdatedProduct(product.id, msg.sender);

        return product.id;
    }

    function remove(uint _id)
    external
    {
        require(existsByIdAndDeleted(_id, false), 'product does not exist');

        Product storage product = products[_id];

        require(
            deviceContract.isOwnedByUser(product.device, msg.sender) || msg.sender == address(deviceContract),
            'not authorized to delete product'
        );

        product.deleted = true;

        emit RemovedProduct(product.id, msg.sender);
    }

    function addNegotiationRequest(uint _id, uint _request)
    external
    {
        require(msg.sender == address(negotiationContract), 'not authorized to add negotiation request to product');

        Product storage product = products[_id];
        product.negotiationRequests.push(_request);
    }

    function addTradingRequest(uint _id, uint _request)
    external
    {
        require(msg.sender == address(tradingContract), 'not authorized to add trading request to product');

        Product storage product = products[_id];
        product.tradingRequests.push(_request);
    }

    function addTrade(uint _id, uint _trade)
    external
    {
        require(msg.sender == address(tradingContract), 'not authorized to add trade to product');

        Product storage product = products[_id];
        product.trades.push(_trade);
    }

    function addNegotiation(uint _id, uint _negotiation)
    external
    {
        require(msg.sender == address(negotiationContract), 'not authorized to add negotiation to product');

        Product storage product = products[_id];
        product.negotiations.push(_negotiation);
    }

    function removeTradingRequest(uint _id, uint _request)
    public
    {
        require(msg.sender == address(tradingContract), 'not authorized to remove trading request from product');

        Product storage product = products[_id];
        removeRequest(_request, product.tradingRequests);
    }

    function removeNegotiationRequest(uint _id, uint _request)
    public
    {
        require(
            msg.sender == address(negotiationContract),
            'not authorized to remove negotiation request from product'
        );

        Product storage product = products[_id];
        removeRequest(_request, product.negotiationRequests);
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

    function findByIndex(uint _index)
    external
    view
    returns (
        uint id,
        address device,
        string memory name,
        string memory description,
        string memory dataType,
        uint frequency,
        uint cost,
        bool deleted
    )
    {
        require(_index >= 0 && _index < index.length, 'product does not exist');

        Product memory product = products[index[_index]];

        id = product.id;
        device = product.device;
        name = product.name;
        description = product.description;
        dataType = product.dataType;
        frequency = product.frequency;
        cost = product.cost;
        deleted = product.deleted;
    }

    function findById(uint _id)
    external
    view
    returns (
        uint id,
        address device,
        string memory name,
        string memory description,
        string memory dataType,
        uint frequency,
        uint cost,
        bool deleted
    )
    {
        require(existsById(_id), 'product does not exist');

        Product memory product = products[_id];

        id = product.id;
        device = product.device;
        name = product.name;
        description = product.description;
        dataType = product.dataType;
        frequency = product.frequency;
        cost = product.cost;
        deleted = product.deleted;
    }

    function findNegotiationRequestsById(uint _id)
    external
    view
    returns (uint[] memory negotiationRequests)
    {
        require(existsById(_id), 'product does not exist');

        Product memory product = products[_id];
        negotiationRequests = product.negotiationRequests;
    }

    function findNegotiationsById(uint _id)
    external
    view
    returns (uint[] memory negotiations)
    {
        require(existsById(_id), 'product does not exist');

        Product memory product = products[_id];
        negotiations = product.negotiations;
    }

    function findTradingRequestsById(uint _id)
    external
    view
    returns (uint[] memory tradingRequests)
    {
        require(existsById(_id), 'product does not exist');

        Product memory product = products[_id];
        tradingRequests = product.tradingRequests;
    }

    function findTradesById(uint _id)
    external
    view
    returns (uint[] memory trades)
    {
        require(existsById(_id), 'product does not exist');

        Product memory product = products[_id];
        trades = product.trades;
    }

    function findCostById(uint _id)
    external
    view
    returns (uint cost)
    {
        require(existsById(_id), 'product does not exist');
        cost = products[_id].cost;
    }

    function findFrequencyById(uint _id)
    external
    view
    returns (uint frequency)
    {
        require(existsById(_id), 'product does not exist');
        frequency = products[_id].frequency;
    }

    function isOwnedByDevice(uint _id, address _device)
    external
    view
    returns (bool isOwned)
    {
        require(existsById(_id), 'product does not exist');
        return products[_id].device == _device;
    }

    function count()
    external
    view
    returns (uint counter)
    {
        return index.length;
    }

    function existsById(uint _id)
    public
    view
    returns (bool isContained)
    {
        if (index.length == 0) return false;
        return (index[products[_id].index] == _id);
    }

    function existsByIdAndDeleted(uint _id, bool _deleted)
    public
    view
    returns (bool exists)
    {
        return existsById(_id) && products[_id].deleted == _deleted;
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
