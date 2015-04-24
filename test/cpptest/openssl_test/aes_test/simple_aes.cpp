#include <iostream>
#include <string.h>
#include "simple_aes.h"
//#include "encoding/iBase64.h"

using namespace std;


int32_t Base64Encode(const char * msg, int encoded_length, char * result){
    return EVP_EncodeBlock((unsigned char*)result, (const unsigned char*)msg, encoded_length);
}

//int32_t Base64Decode(const char *encoded, int encoded_length, char *decoded) {
int32_t Base64Decode(const char * msg, int encoded_length, char * result){
    return EVP_DecodeBlock((unsigned char*)result, (const unsigned char*)msg, encoded_length);
}

//Calculates the length of a decoded string
size_t CalcDecodeLength(const char* b64input)
{
    size_t len = strlen(b64input), padding = 0;

    if (b64input[len-1] == '=' && b64input[len-2] == '=') //last two chars are =
        padding = 2;
    else if (b64input[len-1] == '=') //last char is =
        padding = 1;

    return (len*3)/4 - padding;
}

int SimpleAES::Init(const std::string &key_data_str)
{
    int i, nrounds = 1;
    unsigned char *key_data = new unsigned char[key_data_str.size() + 1]();
    std::copy(key_data_str.begin(), key_data_str.end(), key_data);
    key_data[key_data_str.size()] = '\0';

    i = EVP_BytesToKey(EVP_aes_256_cbc(), EVP_md5(), NULL, key_data,
                       key_data_str.length(), nrounds, key, iv);
    if (i != 32)
    {
        delete [] key_data;
        return -1;
    }

    EVP_CIPHER_CTX_init(&encrypt_ctx);
    EVP_EncryptInit_ex(&encrypt_ctx, EVP_aes_256_cbc(), NULL, key, iv);
    EVP_CIPHER_CTX_init(&decrypt_ctx);
    EVP_DecryptInit_ex(&decrypt_ctx, EVP_aes_256_cbc(), NULL, key, iv);

    delete [] key_data;

    return 0;
}

std::string SimpleAES::Encrypt(const std::string &plaintext)
{
    int len = plaintext.size() + 1;
    unsigned char *plaintext_data = new unsigned char[len]();
    std::copy(plaintext.begin(), plaintext.end(), plaintext_data);
    plaintext_data[plaintext.size()] = '\0';

    len = plaintext.size();

    int c_len = len + AES_BLOCK_SIZE, f_len = 0;
    unsigned char *ciphertext = new unsigned char[c_len]();

    EVP_EncryptInit_ex(&encrypt_ctx, EVP_aes_256_cbc(), NULL, key, iv);
    EVP_EncryptUpdate(&encrypt_ctx, ciphertext, &c_len, plaintext_data, len);
    EVP_EncryptFinal_ex(&encrypt_ctx, ciphertext + c_len, &f_len);

    len = c_len + f_len;

    char res[len]{};
    //Comm::EncodeBase64(ciphertext, cipherstr, len);
    Base64Encode((const char *)ciphertext, len, res);
    std::string cipherstr(res);

    delete [] ciphertext;
    delete [] plaintext_data;

    return cipherstr;
}

std::string SimpleAES::Decrypt(const std::string &cipherstr_b64)
{
    unsigned char *ciphertext = new unsigned char[CalcDecodeLength(cipherstr_b64.c_str())]();

    unsigned char *cipherstr_b64_data = new unsigned char[cipherstr_b64.size() + 1]();
    std::copy(cipherstr_b64.begin(), cipherstr_b64.end(), cipherstr_b64_data);
    cipherstr_b64_data[cipherstr_b64.size()] = '\0';

    //int ciphertext_len = Comm::DecodeBase64(cipherstr_b64_data, ciphertext, cipherstr_b64.size());
    //int ciphertext_len = Base64Decode((const char *)cipherstr_b64_data, cipherstr_b64.size(), ciphertext);
    //int ciphertext_len = Base64Decode((const char *)ciphertext, cipherstr_b64.size(), cipherstr_b64_data);
    // 原始串通过Base64 decode 到ciphertext 中
    int ciphertext_len = Base64Decode((const char *)cipherstr_b64_data, cipherstr_b64.size(), (char *)ciphertext);

    int len = 0, plaintext_len = 0;
    unsigned char *plaintext = new unsigned char[cipherstr_b64.size()]();
    memset(plaintext, 0, cipherstr_b64.size());
    //unsigned char plaintext[cipherstr_b64.size() + AES_BLOCK_SIZE]{};
    //unsigned char plaintext[ciphertext_len + AES_BLOCK_SIZE]{};

    EVP_DecryptInit_ex(&decrypt_ctx, EVP_aes_256_cbc(), NULL, key, iv);
    EVP_DecryptUpdate(&decrypt_ctx, plaintext, &len, ciphertext, ciphertext_len);
    plaintext_len = len;

    cout << "plaintext_len1:" << plaintext_len << endl;
    //fprintf(stdout, "%s\n", plaintext);

    EVP_DecryptFinal_ex(&decrypt_ctx, plaintext + len, &len);
    plaintext_len += len;

    cout << "plaintext_len2:" << plaintext_len << endl;
    std::string plaintext_str(reinterpret_cast<const char*>(plaintext), plaintext_len);
    //std::string plaintext_str((char *)plaintext);

    delete [] ciphertext;
    delete [] cipherstr_b64_data;
    delete [] plaintext;

    return plaintext_str;
}




int main(int argc, char** argv) {

    //cout << "cipher_block_size:" << AES_BLOCK_SIZE << endl;
    //string key("1234");
    string key("kdwxdsp");
    //string msg0("1000");
    //string msg1("Vz3XJBIVoIM0UcTuiA49JQ==");
    
    SimpleAES en_de_er;
    en_de_er.Init(key);
    
    string ciphertext = en_de_er.Encrypt(argv[1]);

    cout << argv[1] << " encrypt to:" << ciphertext << endl;
    cout << ciphertext << " decrypt to:" << en_de_er.Decrypt(ciphertext) << endl;

}

