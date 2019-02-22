#ifndef     __HEATSHRINK_APP_H__
#define     __HEATSHRINK_APP_H__

#ifdef  __cplusplus
extern "C" {
#endif

int  LumiHeatshrinkBase64Encode(char *p_encode_in_buf, int u32_encode_in_len, char *p_encode_out_buf);

int  LumiHeatshrinkBase64Decode(char *p_decode_in_buf, int u32_decode_in_len , char *p_decode_out_buf);

#ifdef  __cplusplus
}
#endif

#endif
