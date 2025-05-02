package utils

import (
	"context"
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"slv.sh/slv/internal/core/crypto"
	"slv.sh/slv/internal/core/secretkey"
)

var (
	sKey            *crypto.SecretKey
	sKeyInitialized bool
)

func SecretKey() (secretKey *crypto.SecretKey, err error) {
	if sKey == nil && !sKeyInitialized {
		sKeyInitialized = true
		sk, _ := secretkey.Get()
		if clientset, _ := getKubeClientSet(); clientset != nil {
			if sk == nil {
				sk, err = GetSecretKeyFor(clientset, GetCurrentNamespace())
			}
			if err == nil && sk != nil {
				sKey = sk
				var pkEC, pkPQ *crypto.PublicKey
				if pkEC, err = sk.PublicKey(false); err == nil {
					if pkPQ, err = sk.PublicKey(true); err == nil {
						var publicKeyEC, publicKeyPQ string
						if publicKeyEC, err = pkEC.String(); err == nil {
							if publicKeyPQ, err = pkPQ.String(); err == nil {
								err = putPublicKeyToConfigMap(clientset, publicKeyEC, publicKeyPQ)
							}
						}
					}
				}
			}
		}
		if err != nil {
			sKey = nil
		}
	}
	return sKey, err
}

func GetPublicKeyFromK8s(namespace string, pq bool) (string, error) {
	clientset, err := getKubeClientSet()
	if err != nil {
		return "", fmt.Errorf("failed to get k8s clientset: %w", err)
	}
	configMap, err := clientset.CoreV1().ConfigMaps(namespace).Get(context.Background(), slvK8sEnvSecret, metav1.GetOptions{})
	if err != nil {
		return "", err
	}
	configMapKey := publicKeyNameEC
	if pq {
		configMapKey = publicKeyNamePQ
	}
	publicKeyStr, ok := configMap.Data[configMapKey]
	if !ok {
		return "", fmt.Errorf("public key not found in configmap: %s", slvK8sEnvSecret)
	}
	return publicKeyStr, nil
}
