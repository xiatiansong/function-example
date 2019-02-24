#include <stdio.h>

typedef unsigned char		uint8_t;
typedef unsigned short int	uint16_t;
typedef unsigned int		uint32_t;
typedef unsigned long long int	uint64_t;

#define MAX_CHARARISTOR_UNIT 		12
#define MAX_ZIP_UNIT				1500
#define MAX_IR_RX_LEN				(MAX_ZIP_UNIT*2)

typedef enum
{
	false,
	true
}bool;

typedef struct
{
	uint8_t 	u8Header;
	uint16_t 	u16BrandID;
	uint32_t 	u32RemoteID;
	uint16_t 	u16Freq;
	uint32_t 	u32IrCommand;
	uint8_t 	u8ChararistorNum;
	uint16_t	u16ZipDataNum;
	uint8_t 	u8Version;
	uint8_t 	u8HeaderChkSum;

	uint16_t 	au16ChararistorUnit[MAX_CHARARISTOR_UNIT];
	uint8_t 	au8ZipData[MAX_ZIP_UNIT];
}tsLumi_IrProtocol;

typedef struct
{
	uint16_t u16Len;
	uint16_t au16DataBuf[MAX_IR_RX_LEN];
}tsIrRxParam_t;

static bool Decode_Rx_Ir_Bits(uint32_t u32StandardV, uint32_t u32RealV, uint16_t u16DeltPer)
{
	uint32_t u32Delt = 0;

	if(u32StandardV > u32RealV)
		return false;

	u32Delt = u32RealV - u32StandardV;

	if(u32Delt < u16DeltPer)
		return true;
	else
		return false;
}

static bool Ir_Rx_RawData_Preprocess(tsIrRxParam_t *psIrRxParam, tsLumi_IrProtocol *pPreFmt)
{
	uint8_t i2, u8Index=0;
	uint16_t i1,u16DeltValue, u16Count,u16TempMin =0;
	uint16_t au16Average[MAX_CHARARISTOR_UNIT];
	uint32_t u32Sum=0;
	if((psIrRxParam == NULL) || (pPreFmt == NULL) || (psIrRxParam->u16Len < 10))
		return false;

	u16DeltValue = 100;
	for(u8Index=0; u8Index<MAX_CHARARISTOR_UNIT; u8Index++)
  	{
  		//Find out the minimus signal
  		u16TempMin = 0xFFFF;
	 	for(i1=0; i1<psIrRxParam->u16Len; i1++)
	  	{
	  		for(i2=0; i2<pPreFmt->u8ChararistorNum; i2++)
	  		{
	  			if(psIrRxParam->au16DataBuf[i1] == pPreFmt->au16ChararistorUnit[i2])
					break;
	  		}
			if(i2 != pPreFmt->u8ChararistorNum)
				continue;

	  		if(psIrRxParam->au16DataBuf[i1] < u16TempMin)
	  			u16TempMin = psIrRxParam->au16DataBuf[i1];
	 	}


		if(u16TempMin == 0xFFFF)
			break;
		if(pPreFmt->u8ChararistorNum == 0)
		{
			u16DeltValue = u16TempMin/2;
		}
		else
		{
			u16DeltValue = u16TempMin/3;
		}
		//calulate average value
		u32Sum = 0;
		u16Count = 0;
		for(i1=0; i1<psIrRxParam->u16Len; i1++)
	  	{
	  		if(Decode_Rx_Ir_Bits (u16TempMin, psIrRxParam->au16DataBuf[i1],u16DeltValue))
	  		{
	  			u32Sum += psIrRxParam->au16DataBuf[i1];
				psIrRxParam->au16DataBuf[i1] = u16TempMin;
				u16Count++;
	  		}
	 	}
		au16Average[pPreFmt->u8ChararistorNum] = u32Sum	/u16Count;
		pPreFmt->au16ChararistorUnit[pPreFmt->u8ChararistorNum] = u16TempMin;

		//printf("u4ChararistorNum[%d][%d][%d][%d]: %d\n", pPreFmt->u8ChararistorNum,u16TempMin, u16Count, u16DeltValue, au16Average[pPreFmt->u8ChararistorNum]);

		pPreFmt->u8ChararistorNum++;
	}

	//Preprocess use the average value
	for(i1=0; i1<psIrRxParam->u16Len; i1++)
	{
		for(i2=0; i2<pPreFmt->u8ChararistorNum; i2++)
	  	{
	  		if(psIrRxParam->au16DataBuf[i1] == pPreFmt->au16ChararistorUnit[i2])
			{
				psIrRxParam->au16DataBuf[i1] = au16Average[i2];
				break;
	  		}
	  	}
	}
	for(i2=0; i2<pPreFmt->u8ChararistorNum; i2++)
		pPreFmt->au16ChararistorUnit[i2] = au16Average[i2];

	//pPreFmt->u8ChararistorNum == MAX_CHARARISTOR_UNIT, error ir
	if(MAX_CHARARISTOR_UNIT == pPreFmt->u8ChararistorNum)
		return false;
	else
		return true;
}

static bool Ir_Rx_Preprocess_Zip1(tsIrRxParam_t *psIrRxParam, tsLumi_IrProtocol *pPreFmt)
{
	uint16_t i1=0;
	uint8_t u8Index = 0;

	if((psIrRxParam == NULL) || (pPreFmt == NULL) || (psIrRxParam->u16Len < 20))
		return false;

	pPreFmt->u8Version = 2;
	//Zip Now
	pPreFmt->u16ZipDataNum = psIrRxParam->u16Len/2;
	for(i1=0; i1<pPreFmt->u16ZipDataNum; i1++)
	{
		for(u8Index=0; u8Index<pPreFmt->u8ChararistorNum; u8Index++)
		{
			if(psIrRxParam->au16DataBuf[i1*2] == pPreFmt->au16ChararistorUnit[u8Index])
				break;
		}
		if(u8Index == pPreFmt->u8ChararistorNum)
			printf("Error: didn't find out the chararistor = %d \n", i1*2);

		pPreFmt->au8ZipData[i1] = u8Index << 4;
		for(u8Index=0; u8Index<pPreFmt->u8ChararistorNum; u8Index++)
		{
			if(psIrRxParam->au16DataBuf[i1*2+1] == pPreFmt->au16ChararistorUnit[u8Index])
				break;
		}
		if(u8Index == pPreFmt->u8ChararistorNum)
			printf("Error: didn't find out the chararistor = %d \n", i1*2+1);
		pPreFmt->au8ZipData[i1] |= u8Index;
	}

#if 0
	wmprintf("au8ZipUnit[%d]: ", pPreFmt->u16ZipDataNum);
	uint32_t u32Sum = 0;
	for(i1=0; i1<pPreFmt->u16ZipDataNum; i1++)
	{
		wmprintf(" %02x ",  pPreFmt->au8ZipData[i1]);
		u32Sum += pPreFmt->au8ZipData[i1];
	}
	wmprintf(" %d\r\n", u32Sum);
#endif
	return true;
}


int main(void)
{
	int i;
	uint16_t au16DataBuf[MAX_IR_RX_LEN]={3484,1690,442,1248,442,442,442,442,442,3484,1690,442,1248,442,442,442,442,442,3484,1690,442,1248,442,442,442,442};
	tsIrRxParam_t psIrRxParam={0};
	tsLumi_IrProtocol pPreFmt={0};

	for(i=0;0!=au16DataBuf[i];i++)
	{
		psIrRxParam.au16DataBuf[i] = au16DataBuf[i];
	}

	psIrRxParam.u16Len = i;
	if(true == Ir_Rx_RawData_Preprocess(&psIrRxParam, &pPreFmt))
		Ir_Rx_Preprocess_Zip1(&psIrRxParam, &pPreFmt);

	for(i=0;i<pPreFmt.u8ChararistorNum;i++)
	{
		printf("Chararistor[%d]=%d\n", i, pPreFmt.au16ChararistorUnit[i]);
	}

	for(i=0;i<pPreFmt.u16ZipDataNum;i++)
	{
		printf("%.2x", pPreFmt.au8ZipData[i]);
	}
	printf("\n");

	return 0;
}