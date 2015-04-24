#include <stdio.h>
#include <iostream>
#include <string.h>
#include <openssl/aes.h>
#include <openssl/evp.h>

using namespace std;

int32_t Base64Encode(const char * msg, int encoded_length, char * result){
    return EVP_EncodeBlock((unsigned char*)result, (const unsigned char*)msg, encoded_length);
}

//int32_t Base64Decode(const char *encoded, int encoded_length, char *decoded) {
int32_t Base64Decode(const char * msg, int encoded_length, char * result){
    return EVP_DecodeBlock((unsigned char*)result, (const unsigned char*)msg, encoded_length);
}


int main() {
    unsigned char key[32];      // 秘钥
    unsigned char iv[32];       // 向量
    //EVP_CIPHER_CTX encrypt_ctx;
    EVP_CIPHER_CTX decrypt_ctx; // 环境

    int i, nrounds = 1;
    const unsigned char *key_data = (unsigned char *)"kdwxdsp";
    
    i = EVP_BytesToKey(EVP_aes_256_cbc(), EVP_md5(), NULL, key_data,
                       strlen((char *)key_data), nrounds, key, iv);
    if (i != 32)
    {
        cerr << "set key fail" << endl;
        return -1;
    }

    // 初始化加密解密环境
    //EVP_CIPHER_CTX_init(&encrypt_ctx);
    //EVP_EncryptInit_ex(&encrypt_ctx, EVP_aes_256_cbc(), NULL, key, iv);
    EVP_CIPHER_CTX_init(&decrypt_ctx);
    //EVP_DecryptInit_ex(&decrypt_ctx, EVP_aes_256_cbc(), NULL, key, iv);

    //=========================================================================
    // 开始解密
    
    const unsigned char *cipher_b64_text = (const unsigned char *)"Vz3XJBIVoIM0UcTuiA49JQ==";  // base64 密文
    //unsigned char *cipher_text = new unsigned char[strlen((char *)cipher_b64_text)](); // 开辟密文空间
    char *cipher_text = new char[strlen((char *)cipher_b64_text)](); // 开辟密文空间

    int cipher_text_len = Base64Decode((char *)cipher_b64_text, strlen((const char *)cipher_b64_text), cipher_text);

    unsigned char plaintext[cipher_text_len + AES_BLOCK_SIZE]{};
    int plaintext_len = 0;
    int tmp_len = 0;
    int len = 0;
    // 
    EVP_DecryptInit_ex(&decrypt_ctx, EVP_aes_256_cbc(), NULL, key, iv);
    EVP_DecryptUpdate(&decrypt_ctx, plaintext, &tmp_len, (unsigned char *)cipher_text, cipher_text_len);
    plaintext_len = tmp_len;

    cout << "plaintext_len1:" << plaintext_len << endl;
    
    EVP_DecryptFinal_ex(&decrypt_ctx, plaintext + tmp_len, &len);
    plaintext_len += len;

    cout << "plaintext_len2:" << plaintext_len << endl;

    cout << plaintext << endl;
}
