package blockchain

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	"github.com/hyperledger/fabric-sdk-go/pkg/gateway"
	"github.com/tuyenngduc/certificate-management-system/backend/internal/models"
)

type FabricConfig struct {
	ChannelName   string
	ChaincodeName string
	WalletPath    string
	CCPPath       string
	Identity      string
}

func NewFabricConfigFromEnv() *FabricConfig {
	return &FabricConfig{
		ChannelName:   getEnv("FABRIC_CHANNEL", "mychannel"),
		ChaincodeName: getEnv("FABRIC_CHAINCODE", "certificate"),
		WalletPath:    getEnv("FABRIC_WALLET_PATH", "./wallet"),
		CCPPath:       getEnv("FABRIC_CCP_PATH", "./connection-org1.yaml"),
		Identity:      getEnv("FABRIC_IDENTITY", "admin"),
	}

}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

type FabricClient struct {
	cfg *FabricConfig
}

func NewFabricClient(cfg *FabricConfig) *FabricClient {
	return &FabricClient{cfg: cfg}
}

func (fc *FabricClient) IssueCertificate(cert any) (string, error) {
	fmt.Println("Step: Opening wallet:", fc.cfg.WalletPath)
	fmt.Println("Step: Checking identity:", fc.cfg.Identity)
	fmt.Println("Step: Connecting gateway with CCP:", fc.cfg.CCPPath)
	fmt.Println("Step: Getting network:", fc.cfg.ChannelName)
	fmt.Println("Step: Getting contract:", fc.cfg.ChaincodeName)

	wallet, err := gateway.NewFileSystemWallet(fc.cfg.WalletPath)
	if err != nil {
		return "", fmt.Errorf("không tạo được wallet: %v", err)
	}
	if !wallet.Exists(fc.cfg.Identity) {
		return "", fmt.Errorf("không tìm thấy identity %s trong wallet", fc.cfg.Identity)
	}

	gw, err := gateway.Connect(
		gateway.WithConfig(config.FromFile(filepath.Clean(fc.cfg.CCPPath))),
		gateway.WithIdentity(wallet, fc.cfg.Identity),
	)
	if err != nil {
		return "", fmt.Errorf("kết nối gateway thất bại: %v", err)
	}
	defer gw.Close()

	network, err := gw.GetNetwork(fc.cfg.ChannelName)
	if err != nil {
		return "", fmt.Errorf("không lấy được network: %v", err)
	}

	contract := network.GetContract(fc.cfg.ChaincodeName)
	certBytes, err := json.Marshal(cert)
	if err != nil {
		return "", fmt.Errorf("không thể marshal certificate: %v", err)
	}

	endorsingPeers := []string{
		"peer0.org1.example.com:7051",
		"peer0.org2.example.com:9051",
	}

	tx, err := contract.CreateTransaction("IssueCertificate",
		gateway.WithEndorsingPeers(endorsingPeers...),
	)
	fmt.Println("Step: Creating transaction with endorsingPeers:", endorsingPeers)

	if err != nil {
		return "", fmt.Errorf("không tạo được transaction: %v", err)
	}

	result, err := tx.Submit(string(certBytes))
	fmt.Printf("cert bytes = %s\n", string(certBytes))
	if err != nil {
		return "", fmt.Errorf("invoke chaincode lỗi: %v", err)
	}

	txID := string(result)
	fmt.Printf("transaction ID from chaincode = %s\n", txID)

	return txID, nil
}

func (fc *FabricClient) GetCertificateByID(certificateID string) (*models.CertificateOnChain, error) {
	wallet, err := gateway.NewFileSystemWallet(fc.cfg.WalletPath)
	if err != nil {
		return nil, fmt.Errorf("cannot create wallet: %v", err)
	}
	if !wallet.Exists(fc.cfg.Identity) {
		return nil, fmt.Errorf("identity %s does not exist in wallet", fc.cfg.Identity)
	}

	gw, err := gateway.Connect(
		gateway.WithConfig(config.FromFile(filepath.Clean(fc.cfg.CCPPath))),
		gateway.WithIdentity(wallet, fc.cfg.Identity),
	)
	if err != nil {
		return nil, fmt.Errorf("gateway connect failed: %v", err)
	}
	defer gw.Close()

	network, err := gw.GetNetwork(fc.cfg.ChannelName)
	if err != nil {
		return nil, fmt.Errorf("cannot get network: %v", err)
	}

	contract := network.GetContract(fc.cfg.ChaincodeName)

	result, err := contract.EvaluateTransaction("ReadCertificate", certificateID)
	if err != nil {
		return nil, fmt.Errorf("EvaluateTransaction error: %v", err)
	}

	var cert models.CertificateOnChain
	if err := json.Unmarshal(result, &cert); err != nil {
		return nil, fmt.Errorf("unmarshal error: %v", err)
	}
	return &cert, nil
}

func (fc *FabricClient) UpdateCertificate(cert any) error {
	wallet, err := gateway.NewFileSystemWallet(fc.cfg.WalletPath)
	if err != nil {
		return fmt.Errorf("cannot create wallet: %v", err)
	}
	if !wallet.Exists(fc.cfg.Identity) {
		return fmt.Errorf("identity %s does not exist in wallet", fc.cfg.Identity)
	}

	gw, err := gateway.Connect(
		gateway.WithConfig(config.FromFile(filepath.Clean(fc.cfg.CCPPath))),
		gateway.WithIdentity(wallet, fc.cfg.Identity),
	)
	if err != nil {
		return fmt.Errorf("gateway connect failed: %v", err)
	}
	defer gw.Close()

	network, err := gw.GetNetwork(fc.cfg.ChannelName)
	if err != nil {
		return fmt.Errorf("cannot get network: %v", err)
	}

	contract := network.GetContract(fc.cfg.ChaincodeName)

	certBytes, err := json.Marshal(cert)
	if err != nil {
		return fmt.Errorf("cannot marshal certificate: %v", err)
	}

	endorsingPeers := []string{
		"peer0.org1.example.com:7051",
		"peer0.org2.example.com:9051",
	}

	tx, err := contract.CreateTransaction("UpdateCertificate",
		gateway.WithEndorsingPeers(endorsingPeers...),
	)
	if err != nil {
		return fmt.Errorf("failed to create transaction: %v", err)
	}

	_, err = tx.Submit(string(certBytes))
	if err != nil {
		return fmt.Errorf("chaincode invoke (UpdateCertificate) failed: %v", err)
	}

	return nil
}
