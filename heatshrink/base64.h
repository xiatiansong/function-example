#ifndef     __BASE64_H__
#define     __BASE64_H__

#ifdef  __cplusplus
extern "C" {
#endif

void vBase64DecryptData(const unsigned char * p_buf,int i32_len,unsigned char *p_out_buf,int *p_out_buf_len);
void vBase64EncryptData(const unsigned char * p_buf,int i32_len,unsigned char *p_out_buf,int *p_out_buf_len);

#ifdef  __cplusplus
}
#endif

#endif

