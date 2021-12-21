# wallet-SDK

ComingChat substrate wallet SDK

## Struct

### Wallet

#### **Wallet 工具类方法（静态方法）**

* ##### 获取助记词：

  ```java
  String mnemonicString = Wallet.genMnemonic();
  ```

  throw err	



* ##### 创建钱包:

  1. 从助记词或私钥

    ```java
    Wallet wallet = Wallet.newWallet(seedOrPhrase);
    ```

    | 参数         | 类型   | 描述           | 获取方式 |
    | ------------ | ------ | -------------- | -------- |
    | seedOrPhrase | string | 助记词或者私钥 |          |
    | network      | int    | ss58           |          |

    throw err

  

  2. 从keystore
  
    ```java
    Wallet wallet = Wallet.NewWalletFromKeyStore(keyStoreJson, password);
    ```
  
    | 参数         | 类型   | 描述     | 获取方式 |
    | ------------ | ------ | -------- | -------- |
    | keyStoreJson | string | keystore |          |
    | password     | string | 密码     |          |
  
    throw err



* ##### 地址转公钥

  ```java
  String publicKey = Wallet.addressToPublicKey(address);
  ```

  | 参数    | 类型   | 描述 | 获取方式 |
  | ------- | ------ | ---- | -------- |
  | address | string | 地址 |          |



* ##### 公钥转地址

    ```java
    String address = Wallet.publicKeyToAddress(publicKey, network);
    ```

    | 参数      | 类型   | 描述             | 获取方式 |
    | --------- | ------ | ---------------- | -------- |
    | publicKey | string | 公钥16进制字符串 |          |
    | network   | int    | ss58             |          |



#### **Wallet类方法**

  * ##### 获取公钥：

    ```java
    String publicKeyHex = wallet.getPublicKeyHex();
    byte[] publicKey = wallet.getPublicKey();
    ```


  * ##### 获取私钥：

    ```java
    String privateKey = wallet.getPrivateKeyHex();
    ```

    throw err



  * ##### 签名:

    ```java
    byte[] signature = wallet.sign(message, password);
    ```

    | 参数     | 类型   | 描述                                            | 获取方式 |
    | -------- | ------ | ----------------------------------------------- | -------- |
    | message  | []byte | 签名内容                                        |          |
    | password | string | keystore 密码，如果是助记词或者私钥创建则随便填 |          |

    throw err
    
    
    
  * ##### 检查Keystore密码

    ```java
    bool isCorrect = wallet.checkPassword(password);
    ```

    | 参数     | 类型   | 描述 | 获取方式 |
    | -------- | ------ | ---- | -------- |
    | password | string |      |          |
    

throw err



  * ##### 获取地址:

    ```java
    String address = wallet.getAddress(network);
    ```

    | 参数    | 类型 | 描述 | 获取方式 |
    | ------- | ---- | ---- | -------- |
    | network | int  | ss58 |          |

    throw err



### Tx

  * ##### 创建Tx

    ```java
    Tx tx = Tx.NewTx(metadataString);
    ```

    | 参数           | 类型   | 描述                   | 获取方式 |
    | -------------- | ------ | ---------------------- | -------- |
    | metadataString | string | metadata的16进制string | 接口获取 |

    throw err
    
  * ##### 从Hex创建Transaction

    ```
    Transaction t = tx.newTxFromHex(isChainX, txHex)
    ```

    | 参数     | 类型   | 描述           | 获取方式             |
    | -------- | ------ | -------------- | -------------------- |
    | isChainX | bool   | 是否是chainx   |                      |
    | txHex    | string | tx的16进制数据 | 从Wallet Connect获取 |

    throw err



  * ##### 创建Balance转账

    ```java
    Transaction t = tx.newBalanceTransferTx(dest, amount);
    ```

    | 参数   | 类型   | 描述        | 获取方式 |
    | ------ | ------ | ----------- | -------- |
    | dest   | string | 对方address | 输入     |
    | amount | long   | 金额        | 用户输入 |
    
    throw err



  * ##### 创建ChainX转账

    ```java
    Transaction t = tx.newChainXBalanceTransferTx(dest, amount);
    ```

    参数同上

    throw err



  * ##### 创建NFT转账

    ```java
    Transaction t = tx.newComingNftTransferTx(dest, cid);
    ```

    | 参数 | 类型   | 描述        | 获取方式 |
    | ---- | ------ | ----------- | -------- |
    | dest | string | 对方address | 输入     |
    | cid  | long   | 金额        | 用户输入 |
    
    throw err



  * ##### 创建XBTC转账

    ```java
    Transaction t = tx.newXAssetsTransferTx(dest, amount);
    ```

    参数同上

    throw err



  * ##### 创建mini门限转账

    ```java
    Transaction t = tx.NewThreshold(thresholdPublicKey, destAddress, aggSignature, aggPublicKey, controlBlock, message, scriptHash, transferAmount, blockNumber);
    ```

    | 参数               | 类型   | 描述                         | 获取方式                  |
    | :----------------- | ------ | ---------------------------- | ------------------------- |
    | thresholdPublicKey | string | 门限钱包公钥                 | generateThresholdPubkey() |
    | destAddress        | string | 对方地址                     |                           |
    | aggSignature       | string | 聚合签名                     | getAggSignature()         |
    | aggPublicKey       | string | 聚合公钥                     | getAggPubkey()            |
    | controlBlock       | string | controlBlock                 | generateControlBlock()    |
    | message            | string | 576520617265206c6567696f6e21 | 写死                      |
    | scriptHash         | string | scriptHash                   | 接口获取                  |
    | transferAmount     | int64  | 转账金额                     | 用户输入                  |
    | blockNumber        | int32  | blockNumber                  | 与scriptHash同接口获取    |
    
    throw err





### Transaction

  * ##### 获取未签名Tx

    ```java
    String unSignTx = t.getUnSignTx();
    ```
    throw err
    此方法用于预估矿工费使用



  * ##### 获取需要签名内容

    ```java
    byte[] signData = t.getSignData(genesisHashString, nonce, specVersion, transVersion);
    ```

    | 参数              | 类型   | 描述        | 获取方式 |
    | ----------------- | ------ | ----------- | -------- |
    | genesisHashString | string | genesisHash | 接口获取 |
    | nonce             | int64  | nonce       | 接口获取 |
    | specVersion       | int32  |             | 接口获取 |
    | transVersion      | int32  |             | 接口获取 |

    throw err

  * ##### 获取Tx

    ```java
    String sendTx = t.getTx(signerPublicKey, signatureData);
    ```

    | 参数            | 类型   | 描述       | 获取方式                                |
    | --------------- | ------ | ---------- | --------------------------------------- |
    | signerPublicKey | byte[] | 签名者公钥 | [获取公钥](#获取公钥)的getPublicKey方法 |
    | signatureData   | byte[] | 签名结果   |                                         |

    throw err





## 交易构造流程

1. 导入助记词或私钥，生成钱包对象

2. 调用[创建钱包](#创建钱包)获得钱包对象wallet用于签名

3. 调用接口获取构造参数(nonce, specVersion, transVersion )

4. Metadata的获取根据接口判断

5. 调用[创建Tx](#创建Tx)获得tx对象用于构造交易

6. 通过tx对象创建交易t

7. 用钱包Sign方法对交易t的内容([获取需要签名内容](#获取需要签名内容))签名

8. 将签名方公钥([获取公钥](#获取公钥))和签名内容放入GetTx的方法中

9. 交易t调用[获取Tx](#获取Tx)获得SendTx

```java
try {
  String mnemonicString = Wallet.genMnemonic();
	Wallet wallet = Wallet.newWallet(mnemonicString);
  Tx tx = Tx.NewTx(metadataString);
  String dest = "5UczqUVGsoQpZnBCZkDtxvLxJ42KnUfaGTzPkQmZeAAug4s9";
  String genesisHashString = "0xfb58f83706a065ced8f658fafaba97e6e49b772287e332077c499784184eda9f";
  long amount = 100000;
  Transaction t = tx.newBalanceTransferTx(dest, amount);
  byte[] signedMessage = wallet.sign(t.getSignData(genesisHashString, 0, 115, 1));
  String sendTx = t.getTx(wallet.getPublicKey(), signedMessage);
}catch  (Exception e){
  
}
```

