#ifndef __SIMPLE_AES_H__
#define __SIMPLE_AES_H__
#include <string>
#include <openssl/aes.h>
#include <openssl/evp.h>

class SimpleAES
{
public:
    SimpleAES(){
        /**
        for(int i = 0; i < 32; i++){
            iv[i] = i;
        }
        **/
    };
    ~SimpleAES(){};

public:
    int Init(const std::string &key);
    std::string Encrypt(const std::string &plaintext);
    std::string Decrypt(const std::string &ciphertext);

private:
    EVP_CIPHER_CTX encrypt_ctx;
    EVP_CIPHER_CTX decrypt_ctx;

    unsigned char key[32];
    unsigned char iv[32];
};

#endif

