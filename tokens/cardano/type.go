package cardano

type TransactionChainingMap struct {
	InputKey  UtxoKey   `json:"inputKey"`
	AssetsMap AssetsMap `json:"assetsMap"`
}

type TransactionChainingKey struct {
	SpentUtxoMap              map[UtxoKey]bool      `json:"SpentUtxoMap"`
	SpentUtxoListGropByTxHash map[string]*[]UtxoKey `json:"SpentUtxoListGropByTxHash"`
}

type Tip struct {
	Slot  uint64 `json:"slot"`
	Block uint64 `json:"block"`
	Epoch uint64 `json:"epoch"`
	Hash  string `json:"hash"`
}

type UtxoKey struct {
	TxHash  string `json:"txHash"`
	TxIndex uint64 `json:"txIndex"`
}

type OutputsResult struct {
	Outputs []Output `json:"utxos"`
}

type TransactionResult struct {
	Transactions []Transaction `json:"transactions"`
}

type Transaction struct {
	Block         Block      `json:"block"`
	Hash          string     `json:"hash"`
	Metadata      []Metadata `json:"metadata"`
	Outputs       []Output   `json:"outputs"`
	ValidContract bool       `json:"validContract"`
}

type Output struct {
	Address string  `json:"address"`
	Index   uint64  `json:"index"`
	Tokens  []Token `json:"tokens"`
	Value   string  `json:"value"`
	TxHash  string  `json:"txHash"`
}

type Token struct {
	Asset    Asset  `json:"asset"`
	Quantity string `json:"quantity"`
}

type Asset struct {
	PolicyId  string `json:"policyId"`
	AssetName string `json:"assetName"`
}

type Block struct {
	EpochNo uint64 `json:"epochNo"`
	Number  uint64 `json:"number"`
	SlotNo  uint64 `json:"slotNo"`
}

type Metadata struct {
	Key   string        `json:"key"`
	Value MetadataValue `json:"value"`
}

type MetadataValue struct {
	Bind      string `json:"bind,omitempty"`
	ToChainId string `json:"toChainId,omitempty"`
}

type AssetsMap map[string]string

type RawTransaction struct {
	Fee     string               `json:"fee"`
	TxIns   []UtxoKey            `json:"txIns"`
	TxOuts  map[string]AssetsMap `json:"txOuts"`
	Mint    AssetsMap            `json:"mint"`
	TxIndex uint64               `json:"txIndex"`
	OutFile string               `json:"outFile"`
}

func (*RawTransaction) ProtoMessage() {}

type SignedTransaction struct {
	FilePath  string    `json:"filePath"`
	TxIns     []UtxoKey `json:"txIns"`
	TxHash    string    `json:"txHash"`
	TxIndex   uint64    `json:"txIndex"`
	AssetsMap AssetsMap `json:"assetsMap"`
}
