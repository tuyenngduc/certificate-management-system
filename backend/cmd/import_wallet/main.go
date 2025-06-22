package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/hyperledger/fabric-sdk-go/pkg/gateway"
)

func main() {
	credPath := "/home/tuyenngduc/go/src/github.com/tuyenngduc/fabric-samples/test-network/organizations/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp"

	// Đọc file cert
	certPath := filepath.Join(credPath, "signcerts", "cert.pem")
	cert, err := os.ReadFile(certPath)
	if err != nil {
		panic(fmt.Errorf("lỗi đọc cert: %w", err))
	}

	// Đọc file private key (chỉ có 1 file trong keystore)
	keyDir := filepath.Join(credPath, "keystore")
	keyFiles, err := os.ReadDir(keyDir)
	if err != nil {
		panic(fmt.Errorf("lỗi đọc thư mục keystore: %w", err))
	}
	if len(keyFiles) == 0 {
		panic("không tìm thấy file private key trong keystore")
	}
	keyPath := filepath.Join(keyDir, keyFiles[0].Name())
	key, err := os.ReadFile(keyPath)
	if err != nil {
		panic(fmt.Errorf("lỗi đọc private key: %w", err))
	}

	// Tạo ví và import identity
	wallet, err := gateway.NewFileSystemWallet("wallet")
	if err != nil {
		panic(fmt.Errorf("lỗi tạo wallet: %w", err))
	}

	identity := gateway.NewX509Identity("Org1MSP", string(cert), string(key))
	if err = wallet.Put("admin", identity); err != nil {
		panic(fmt.Errorf("lỗi import identity: %w", err))
	}

	fmt.Println("Đã import admin Org1 vào ví 'wallet'")
}
