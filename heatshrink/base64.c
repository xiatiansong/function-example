#include <string.h>
#include <stdlib.h>
#include <stdio.h>
#include "base64.h"

static const unsigned char base64[64] = 
{
    'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H',
    'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P',
    'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X',
    'Y', 'Z', 'a', 'b', 'c', 'd', 'e', 'f',
    'g', 'h', 'i', 'j', 'k', 'l', 'm', 'n',
    'o', 'p', 'q', 'r', 's', 't', 'u', 'v',
    'w', 'x', 'y', 'z', '0', '1', '2', '3',
    '4', '5', '6', '7', '8', '9', '+', '/'
};

static const unsigned char base64_back[128] = 
{
    255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255,
    255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 
    255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 62, 255, 255, 255, 63,
    52, 53, 54, 55, 56, 57, 58, 59, 60, 61, 255, 255, 255,  0, 255, 255,
    255,  0,  1,  2,  3,  4,  5,  6,  7,  8,  9, 10, 11, 12, 13, 14, 
    15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 255, 255, 255, 255, 255, 
    255, 26, 27, 28, 29, 30, 31, 32, 33, 34, 35, 36, 37, 38, 39, 40, 
    41, 42, 43, 44, 45, 46, 47, 48, 49, 50, 51, 255, 255, 255, 255, 255,
};


static void base64_encrypt(const unsigned char * pbuf, unsigned char * cbuf)
{
    int temp = ((pbuf[0]) << 16) | ((pbuf[1]) << 8) | ((pbuf[2]) << 0);
    int i = 0;


    for(i = 0; i<4; i++){
        int index = (temp >> (18-i*6)) & 0x3F;
        cbuf[i] = base64[index];	
    }
}

static void base64_decrypt(const unsigned char * cbuf, unsigned char * pbuf)
{
    int temp = 0;
    int i = 0;
    for(i = 0; i<4; i++){
        temp |= ((base64_back[cbuf[i]]) << (18-6*i));
    }
    for(i = 0; i<3; i++){
        pbuf[i] = (temp >> (16-i*8)) & 0XFF;
    }
}

static void show_base64(const unsigned char * buf, int len)
{
	int i = 0;
    for(i=0; i<len; i++){
        printf("%c",buf[i]);
    }
    printf("\n");
}

void base64_encrypt_text(const unsigned char * pbuf,int plen)
{
	int clen = (plen % 3) ? (plen/3 + 1) : (plen/3);

	//printf("\r\nbase64_encrypt_text 0,plen=%d,  %x %x %x\r\n",plen,pbuf[0],pbuf[1],pbuf[2]);
	unsigned char * buf = (unsigned char *)malloc(clen*3);
	unsigned char * cbuf = (unsigned char *)malloc(clen * 4);
	if(NULL == cbuf || NULL == buf){
		printf("\r\nbase64_encrypt_text error\r\n");
		return;
	}
	memset(cbuf, 0, clen*4);
	memset(buf, 0, clen*3);
	memcpy(buf, pbuf, plen);


	int i = 0;
	for(i = 0; i < clen; i++){
		base64_encrypt(&buf[i*3], &cbuf[i*4]);
	}

	if(plen % 3 == 2){
		cbuf[clen*4 - 1] = '=';
	}
	else if(plen % 3 == 1){
		cbuf[clen*4 - 1] = cbuf[clen*4 - 2] = '=';
	}

	printf("\r\nbase64_encrypt_text  ==%s!!!!\r\n",cbuf);
	show_base64(cbuf, clen*4);
	free(buf);
	free(cbuf);
}

void base64_decrypt_text(const unsigned char * cbuf,int clen)
{
	int plen = clen/4;

	unsigned char * pbuf = (unsigned char *)malloc(plen*3);

	if(NULL == pbuf){
		return;
	}
	memset(pbuf, 0, plen*3);

	int i = 0;
	for(i = 0; i < plen; i++){
		base64_decrypt(&cbuf[i*4], &pbuf[i*3]);
	}
	show_base64(pbuf, plen*3);
	free(pbuf);
}


void vBase64DecryptData(const unsigned char * p_buf,int i32_len,unsigned char *p_out_buf,int *p_out_buf_len)
{
	int plen = i32_len/4;
	int out_len = 0;
	unsigned char * temp_buf = (unsigned char *)malloc(plen*3);

	if(NULL == temp_buf){
		printf("\r\nvBase64DecryptData error\r\n");
		return;
	}
	memset(temp_buf, 0, plen*3);

	int i = 0;
	int k = 0;
	for(i = 0; i < plen; i++){
		k = i;
		base64_decrypt(&p_buf[i*4], &temp_buf[i*3]);
	}

	out_len = plen*3;

	if ('=' == p_buf[k*4 + 2])
		out_len = plen*3 - 2;
	else if ('=' == p_buf[k*4 + 3])
		out_len = plen*3 - 1;		
	
	memcpy(p_out_buf,temp_buf,out_len);
	*p_out_buf_len = out_len;
	free(temp_buf);
}


void vBase64EncryptData(const unsigned char * p_buf,int i32_len,unsigned char *p_out_buf,int *p_out_buf_len)
{
	int clen = (i32_len % 3) ? (i32_len/3 + 1) : (i32_len/3);


	unsigned char * buf = (unsigned char *)malloc(clen*3);
	unsigned char * cbuf = (unsigned char *)malloc(clen * 4);
	if(NULL == cbuf || NULL == buf){
		printf("\r\nvBase64EncryptData error\r\n");
		return;
	}
	memset(cbuf, 0, clen*4);
	memset(buf, 0, clen*3);
	memcpy(buf, p_buf, i32_len);

	int i = 0;
	for(i = 0; i < clen; i++){
		base64_encrypt(&buf[i*3], &cbuf[i*4]);
	}

	if(i32_len % 3 == 2){
		cbuf[clen*4 - 1] = '=';
	}
	else if(i32_len % 3 == 1){
		cbuf[clen*4 - 1] = cbuf[clen*4 - 2] = '=';
	}

	memcpy(p_out_buf,cbuf,clen*4);
	*p_out_buf_len = clen*4;

	free(buf);
	free(cbuf);
}